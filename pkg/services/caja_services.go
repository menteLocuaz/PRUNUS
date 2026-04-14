package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
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
func (s *ServiceCaja) AbrirSesion(ctx context.Context, cajaID, usuarioID uuid.UUID, montoApertura float64) (*models.SesionCaja, error) {
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

// ArqueoYCierre realiza el cierre de jornada comparando el efectivo físico con el registrado en sistema.
func (s *ServiceCaja) ArqueoYCierre(ctx context.Context, sesionID, usuarioID uuid.UUID, montoFisico float64) (map[string]interface{}, error) {
	zaplogger.WithContext(ctx, s.logger).Info("Iniciando arqueo y cierre de caja", zap.String("id_sesion", sesionID.String()))

	sesionActual, err := s.store.GetSesionByID(ctx, sesionID)
	if err != nil {
		return nil, err
	}
	if sesionActual.Estado == "CERRADA" {
		return nil, errors.New("la sesión ya se encuentra cerrada")
	}

	ventasSistema, err := s.store.GetVentasEfectivoBySesion(ctx, sesionID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener ventas de sistema: %w", err)
	}

	diferencia := montoFisico - ventasSistema
	resultado := "CUADRADO"
	if diferencia < -0.01 {
		resultado = "FALTANTE"
	} else if diferencia > 0.01 {
		resultado = "SOBRANTE"
	}

	ahora := time.Now()
	sesionUpdate := &models.SesionCaja{
		MontoCierre: montoFisico,
		FechaCierre: &ahora,
		Estado:      "CERRADA",
	}
	_, err = s.store.UpdateSesion(ctx, sesionID, sesionUpdate)
	if err != nil {
		return nil, fmt.Errorf("error al cerrar sesión: %w", err)
	}

	err = s.usuarioStore.UpdateTurnoStatus(ctx, usuarioID, false)
	if err != nil {
		zaplogger.WithContext(ctx, s.logger).Error("Sesión cerrada pero falló desasignación de usuario",
			zap.String("id_usuario", usuarioID.String()),
			zap.Error(err),
		)
	}

	zaplogger.WithContext(ctx, s.logger).Info("Cierre de jornada completado",
		zap.String("id_sesion", sesionID.String()),
		zap.Float64("sistema", ventasSistema),
		zap.Float64("fisico", montoFisico),
		zap.Float64("diferencia", diferencia),
		zap.String("resultado", resultado),
	)

	return map[string]interface{}{
		"ventas_sistema":  ventasSistema,
		"efectivo_fisico": montoFisico,
		"diferencia":      diferencia,
		"resultado":       resultado,
		"mensaje":         fmt.Sprintf("Cierre completado. Resultado: %s", resultado),
	}, nil
}

// GetMovimientos obtiene el historial de movimientos de una sesión.
func (s *ServiceCaja) GetMovimientos(ctx context.Context, sesionID uuid.UUID) ([]*models.MovimientoCaja, error) {
	if sesionID == uuid.Nil {
		return nil, errors.New("el ID de la sesión es requerido")
	}
	return s.store.GetMovimientosBySesion(ctx, sesionID)
}
