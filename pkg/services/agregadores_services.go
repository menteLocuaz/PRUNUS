package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
	"go.uber.org/zap"
)

type ServiceAgregadores struct {
	store  store.StoreAgregadores
	logger *zap.Logger
}

func NewServiceAgregadores(s store.StoreAgregadores, logger *zap.Logger) *ServiceAgregadores {
	return &ServiceAgregadores{
		store:  s,
		logger: logger,
	}
}

func (s *ServiceAgregadores) GetAllAgregadores(ctx context.Context) ([]*models.Agregador, error) {
	return s.store.GetAllAgregadores(ctx)
}

func (s *ServiceAgregadores) GetAgregadorByID(ctx context.Context, id uuid.UUID) (*models.Agregador, error) {
	return s.store.GetAgregadorByID(ctx, id)
}

func (s *ServiceAgregadores) CreateAgregador(ctx context.Context, a models.Agregador) (*models.Agregador, error) {
	return s.store.CreateAgregador(ctx, &a)
}

func (s *ServiceAgregadores) UpdateAgregador(ctx context.Context, id uuid.UUID, a models.Agregador) (*models.Agregador, error) {
	return s.store.UpdateAgregador(ctx, id, &a)
}

func (s *ServiceAgregadores) DeleteAgregador(ctx context.Context, id uuid.UUID) error {
	return s.store.DeleteAgregador(ctx, id)
}

func (s *ServiceAgregadores) CreateOrdenAgregador(ctx context.Context, o models.OrdenAgregador) (*models.OrdenAgregador, error) {
	return s.store.CreateOrdenAgregador(ctx, &o)
}
