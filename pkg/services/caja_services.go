package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
)

type ServiceCaja struct {
	store        store.StoreCaja
	usuarioStore store.StoreUsuario
	logger       *slog.Logger
}

func NewServiceCaja(s store.StoreCaja, u store.StoreUsuario, logger *slog.Logger) *ServiceCaja {
	return &ServiceCaja{
		store:        s,
		usuarioStore: u,
		logger:       logger,
	}
}

func (s *ServiceCaja) GetAllCajas(ctx context.Context) ([]*models.Caja, error) {
	return s.store.GetAllCajas(ctx)
}

func (s *ServiceCaja) GetCajaByID(ctx context.Context, id uuid.UUID) (*models.Caja, error) {
	return s.store.GetCajaByID(ctx, id)
}

func (s *ServiceCaja) CreateCaja(ctx context.Context, c models.Caja) (*models.Caja, error) {
	return s.store.CreateCaja(ctx, &c)
}

func (s *ServiceCaja) UpdateCaja(ctx context.Context, id uuid.UUID, c models.Caja) (*models.Caja, error) {
	return s.store.UpdateCaja(ctx, id, &c)
}

func (s *ServiceCaja) DeleteCaja(ctx context.Context, id uuid.UUID) error {
	return s.store.DeleteCaja(ctx, id)
}

func (s *ServiceCaja) AbrirSesion(ctx context.Context, cajaID, usuarioID uuid.UUID, montoApertura float64) (*models.SesionCaja, error) {
	sesion := &models.SesionCaja{
		IDCaja:        cajaID,
		IDUsuario:     usuarioID,
		MontoApertura: montoApertura,
	}
	return s.store.CreateSesion(ctx, sesion)
}

func (s *ServiceCaja) CerrarSesion(ctx context.Context, sesionID uuid.UUID, montoCierre float64) (*models.SesionCaja, error) {
	ahora := time.Now()
	sesion := &models.SesionCaja{
		MontoCierre: montoCierre,
		FechaCierre: &ahora,
		Estado:      "CERRADA",
	}
	return s.store.UpdateSesion(ctx, sesionID, sesion)
}

func (s *ServiceCaja) RegistrarMovimiento(ctx context.Context, m models.MovimientoCaja) (*models.MovimientoCaja, error) {
	return s.store.CreateMovimiento(ctx, &m)
}

func (s *ServiceCaja) GetMovimientos(ctx context.Context, sesionID uuid.UUID) ([]*models.MovimientoCaja, error) {
	return s.store.GetMovimientosBySesion(ctx, sesionID)
}

// ArqueoYCierre realiza el cierre de jornada con comparativa y desasignación
func (s *ServiceCaja) ArqueoYCierre(ctx context.Context, sesionID, usuarioID uuid.UUID, montoFisico float64) (map[string]interface{}, error) {
	// 1. Obtener ventas en efectivo registradas en sistema para este turno
	ventasSistema, err := s.store.GetVentasEfectivoBySesion(ctx, sesionID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener ventas de sistema: %w", err)
	}

	// 2. Calcular diferencia (Arqueo)
	diferencia := montoFisico - ventasSistema
	resultado := "CUADRADO"
	if diferencia < 0 {
		resultado = "FALTANTE"
	} else if diferencia > 0 {
		resultado = "SOBRANTE"
	}

	// 3. Cerrar la sesión en el store
	ahora := time.Now()
	sesion := &models.SesionCaja{
		MontoCierre: montoFisico,
		FechaCierre: &ahora,
		Estado:      "CERRADA",
	}
	_, err = s.store.UpdateSesion(ctx, sesionID, sesion)
	if err != nil {
		return nil, fmt.Errorf("error al cerrar sesión: %w", err)
	}

	// 4. Desasignar al usuario (Seguridad)
	err = s.usuarioStore.UpdateTurnoStatus(ctx, usuarioID, false)
	if err != nil {
		s.logger.ErrorContext(ctx, "Sesión cerrada pero falló desasignación de usuario", 
			slog.String("id_usuario", usuarioID.String()), 
			slog.Any("error", err),
		)
	}

	s.logger.InfoContext(ctx, "Cierre de jornada completado", 
		slog.String("id_sesion", sesionID.String()),
		slog.Float64("sistema", ventasSistema),
		slog.Float64("fisico", montoFisico),
		slog.Float64("diferencia", diferencia),
		slog.String("resultado", resultado),
	)

	return map[string]interface{}{
		"ventas_sistema": ventasSistema,
		"efectivo_fisico": montoFisico,
		"diferencia":     diferencia,
		"resultado":      resultado,
		"mensaje":        fmt.Sprintf("Cierre completado. Resultado: %s", resultado),
	}, nil
}
