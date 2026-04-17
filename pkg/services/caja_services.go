package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/prunus/pkg/dto"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
	zaplogger "github.com/prunus/pkg/utils/logger"
	"go.uber.org/zap"
)

// ServiceCaja encapsula la lógica de negocio para la gestión de cajas y sesiones.
type ServiceCaja struct {
	store        store.StoreCaja
	usuarioStore store.StoreUsuario
	logger       *zap.Logger
}

// NewServiceCaja crea una nueva instancia del servicio de cajas.
func NewServiceCaja(s store.StoreCaja, u store.StoreUsuario, logger *zap.Logger) *ServiceCaja {
	return &ServiceCaja{
		store:        s,
		usuarioStore: u,
		logger:       logger,
	}
}

func (s *ServiceCaja) Logger() *zap.Logger {
	return s.logger
}

func (s *ServiceCaja) validateCaja(c *models.Caja) error {
	if c.Nombre == "" {
		return errors.New("el nombre de la caja es requerido")
	}
	if c.IDSucursal == uuid.Nil {
		return errors.New("el ID de la sucursal es requerido")
	}
	return nil
}

// GetAllCajas obtiene todas las cajas registradas.
func (s *ServiceCaja) GetAllCajas(ctx context.Context) ([]*models.Caja, error) {
	zaplogger.WithContext(ctx, s.logger).Info("Obteniendo todas las cajas")
	return s.store.GetAllCajas(ctx)
}

// GetCajaByID obtiene una caja por su ID.
func (s *ServiceCaja) GetCajaByID(ctx context.Context, id uuid.UUID) (*models.Caja, error) {
	if id == uuid.Nil {
		return nil, errors.New("el ID de la caja es requerido")
	}
	return s.store.GetCajaByID(ctx, id)
}

// CreateCaja registra una nueva caja física.
func (s *ServiceCaja) CreateCaja(ctx context.Context, c models.Caja) (*models.Caja, error) {
	if err := s.validateCaja(&c); err != nil {
		zaplogger.WithContext(ctx, s.logger).Warn("Validación fallida al crear caja", zap.Error(err))
		return nil, err
	}
	return s.store.CreateCaja(ctx, &c)
}

// UpdateCaja actualiza la información de una caja.
func (s *ServiceCaja) UpdateCaja(ctx context.Context, id uuid.UUID, c models.Caja) (*models.Caja, error) {
	if id == uuid.Nil {
		return nil, errors.New("el ID de la caja es requerido")
	}
	if err := s.validateCaja(&c); err != nil {
		return nil, err
	}
	return s.store.UpdateCaja(ctx, id, &c)
}

// DeleteCaja realiza un borrado lógico de la caja.
func (s *ServiceCaja) DeleteCaja(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("el ID de la caja es requerido")
	}
	return s.store.DeleteCaja(ctx, id)
}

// AbrirSesion inicia un nuevo turno para un cajero en una caja específica.
func (s *ServiceCaja) AbrirSesion(ctx context.Context, cajaID, usuarioID uuid.UUID, montoApertura float64, desglose []dto.DenominacionDTO) (*models.SesionCaja, error) {
	zaplogger.WithContext(ctx, s.logger).Info("Intentando abrir sesión de caja",
		zap.String("id_caja", cajaID.String()),
		zap.String("id_usuario", usuarioID.String()),
	)

	activa, err := s.store.GetSesionActivaByUsuario(ctx, usuarioID)
	if err != nil {
		return nil, err
	}
	if activa != nil {
		return nil, fmt.Errorf("el usuario ya tiene una sesión abierta (ID: %s)", activa.IDSesion)
	}

	sesion := &models.SesionCaja{
		IDCaja:        cajaID,
		IDUsuario:     usuarioID,
		MontoApertura: montoApertura,
	}
	res, err := s.store.CreateSesion(ctx, sesion)
	if err != nil {
		return nil, err
	}

	// Registrar desglose inicial si existe
	if len(desglose) > 0 {
		_ = s.store.RegistrarArqueoDesglose(ctx, res.IDSesion, "APERTURA", desglose)
	}

	if err := s.usuarioStore.UpdateTurnoStatus(ctx, usuarioID, true); err != nil {
		zaplogger.WithContext(ctx, s.logger).Error("Sesión abierta pero falló actualización de estado de usuario", zap.Error(err))
	}

	zaplogger.WithContext(ctx, s.logger).Info("Sesión de caja abierta exitosamente", zap.String("id_sesion", res.IDSesion.String()))
	return res, nil
}

// RegistrarMovimiento registra un ingreso o egreso manual de efectivo en la sesión.
func (s *ServiceCaja) RegistrarMovimiento(ctx context.Context, m models.MovimientoCaja) (*models.MovimientoCaja, error) {
	if m.IDSesion == uuid.Nil {
		return nil, errors.New("el ID de la sesión es requerido")
	}
	if m.Monto <= 0 {
		return nil, errors.New("el monto debe ser mayor a cero")
	}
	return s.store.CreateMovimiento(ctx, &m)
}

// ArqueoYCierre realiza una conciliación financiera completa antes de cerrar la jornada.
func (s *ServiceCaja) ArqueoYCierre(ctx context.Context, req dto.CierreCajaRequest, usuarioID uuid.UUID) (*dto.ResumenCierreDTO, error) {
	zaplogger.WithContext(ctx, s.logger).Info("Iniciando arqueo y cierre avanzado", zap.String("id_sesion", req.IDControlEstacion.String()))

	// 1. Validar sesión
	sesionActual, err := s.store.GetSesionByID(ctx, req.IDControlEstacion)
	if err != nil {
		return nil, err
	}
	if sesionActual.Estado == "CERRADA" {
		return nil, errors.New("la sesión ya se encuentra cerrada")
	}

	// 2. Obtener Resumen Financiero de Sistema (Ventas, Retiros, Gastos)
	resumen, err := s.store.GetResumenFinanciero(ctx, req.IDControlEstacion)
	if err != nil {
		return nil, fmt.Errorf("error al conciliar montos: %w", err)
	}

	// 3. Cruzar con lo declarado físicamente
	resumen.SaldoReal = req.MontoDeclarado
	resumen.Diferencia = resumen.SaldoReal - resumen.SaldoEsperado
	
	resumen.Resultado = "CUADRADO"
	if resumen.Diferencia < -0.01 {
		resumen.Resultado = "FALTANTE"
	} else if resumen.Diferencia > 0.01 {
		resumen.Resultado = "SOBRANTE"
	}

	// 4. Registrar desglose físico de monedas/billetes
	if len(req.Desglose) > 0 {
		err = s.store.RegistrarArqueoDesglose(ctx, req.IDControlEstacion, "CIERRE", req.Desglose)
		if err != nil {
			zaplogger.WithContext(ctx, s.logger).Error("Fallo al registrar desglose físico", zap.Error(err))
		}
	}

	// 5. Cerrar sesión en BD
	ahora := time.Now()
	sesionUpdate := &models.SesionCaja{
		MontoCierre: resumen.SaldoReal,
		FechaCierre: &ahora,
		Estado:      "CERRADA",
	}
	_, err = s.store.UpdateSesion(ctx, req.IDControlEstacion, sesionUpdate)
	if err != nil {
		return nil, fmt.Errorf("error al cerrar sesión: %w", err)
	}

	// 6. Liberar usuario
	_ = s.usuarioStore.UpdateTurnoStatus(ctx, usuarioID, false)

	zaplogger.WithContext(ctx, s.logger).Info("Cierre de jornada avanzado completado",
		zap.String("id_sesion", req.IDControlEstacion.String()),
		zap.Float64("esperado", resumen.SaldoEsperado),
		zap.Float64("real", resumen.SaldoReal),
		zap.String("resultado", resumen.Resultado),
	)

	return resumen, nil
}

// GetMovimientos obtiene el historial de movimientos de una sesión.
func (s *ServiceCaja) GetMovimientos(ctx context.Context, sesionID uuid.UUID) ([]*models.MovimientoCaja, error) {
	if sesionID == uuid.Nil {
		return nil, errors.New("el ID de la sesión es requerido")
	}
	return s.store.GetMovimientosBySesion(ctx, sesionID)
}
