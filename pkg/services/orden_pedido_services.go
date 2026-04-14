package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/prunus/pkg/dto"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
	"go.uber.org/zap"
)

type ServiceOrdenPedido struct {
	store  store.StoreOrdenPedido
	logger *zap.Logger
}

func NewServiceOrdenPedido(s store.StoreOrdenPedido, logger *zap.Logger) *ServiceOrdenPedido {
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

func (s *ServiceOrdenPedido) GetAllOrdenes(ctx context.Context, params dto.PaginationParams) ([]*models.OrdenPedido, error) {
	return s.store.GetAllOrdenes(ctx, params)
}

func (s *ServiceOrdenPedido) UpdateStatus(ctx context.Context, id uuid.UUID, statusID uuid.UUID) error {
	return s.store.UpdateOrdenStatus(ctx, id, statusID)
}
