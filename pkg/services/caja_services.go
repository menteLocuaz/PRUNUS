package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
)

type ServiceCaja struct {
	store  store.StoreCaja
	logger *slog.Logger
}

func NewServiceCaja(s store.StoreCaja, logger *slog.Logger) *ServiceCaja {
	return &ServiceCaja{
		store:  s,
		logger: logger,
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
