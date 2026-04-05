package services

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/prunus/pkg/dto"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
)

type ServiceFactura struct {
	store  store.StoreFactura
	logger *slog.Logger
}

func NewServiceFactura(s store.StoreFactura, logger *slog.Logger) *ServiceFactura {
	return &ServiceFactura{
		store:  s,
		logger: logger,
	}
}

func (s *ServiceFactura) CreateFactura(ctx context.Context, f models.Factura, items []*models.DetalleFactura) (*models.Factura, error) {
	return s.store.CreateFactura(ctx, &f, items)
}

func (s *ServiceFactura) RegistrarFacturaCompleta(ctx context.Context, req dto.FacturaCompletaRequest, idUsuario uuid.UUID) (*dto.FacturaResponse, error) {
	return s.store.RegistrarFacturaCompleta(ctx, req, idUsuario)
}

func (s *ServiceFactura) GetFactura(ctx context.Context, id uuid.UUID) (*models.Factura, []*models.DetalleFactura, error) {
	return s.store.GetFacturaByID(ctx, id)
}

func (s *ServiceFactura) GetAllFacturas(ctx context.Context, params dto.PaginationParams) ([]*models.Factura, error) {
	return s.store.GetAllFacturas(ctx, params)
}

func (s *ServiceFactura) GetImpuestos(ctx context.Context) ([]*models.Impuesto, error) {
	return s.store.GetAllImpuestos(ctx)
}

func (s *ServiceFactura) GetFormasPago(ctx context.Context) ([]*models.FormaPago, error) {
	return s.store.GetAllFormasPago(ctx)
}
