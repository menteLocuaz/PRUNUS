package services

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
)

type ServiceOrdenPedido struct {
	store  store.StoreOrdenPedido
	logger *slog.Logger
}

func NewServiceOrdenPedido(s store.StoreOrdenPedido, logger *slog.Logger) *ServiceOrdenPedido {
	return &ServiceOrdenPedido{
		store:  s,
		logger: logger,
	}
}

func (s *ServiceOrdenPedido) CreateOrden(ctx context.Context, o models.OrdenPedido) (*models.OrdenPedido, error) {
	return s.store.CreateOrden(ctx, &o)
}

func (s *ServiceOrdenPedido) GetOrdenByID(ctx context.Context, id uuid.UUID) (*models.OrdenPedido, error) {
	return s.store.GetOrdenByID(ctx, id)
}

func (s *ServiceOrdenPedido) GetAllOrdenes(ctx context.Context) ([]*models.OrdenPedido, error) {
	return s.store.GetAllOrdenes(ctx)
}

func (s *ServiceOrdenPedido) UpdateStatus(ctx context.Context, id uuid.UUID, statusID uuid.UUID) error {
	return s.store.UpdateOrdenStatus(ctx, id, statusID)
}
