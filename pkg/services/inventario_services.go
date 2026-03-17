package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
)

type ServiceInventario struct {
	store  store.StoreInventario
	logger *slog.Logger
}

func NewServiceInventario(s store.StoreInventario, logger *slog.Logger) *ServiceInventario {
	return &ServiceInventario{
		store:  s,
		logger: logger,
	}
}

func (s *ServiceInventario) GetAllInventario(ctx context.Context) ([]*models.Inventario, error) {
	return s.store.GetAllInventario(ctx)
}

func (s *ServiceInventario) GetInventarioByID(ctx context.Context, id uuid.UUID) (*models.Inventario, error) {
	return s.store.GetInventarioByID(ctx, id)
}

func (s *ServiceInventario) CreateInventario(ctx context.Context, inventario models.Inventario) (*models.Inventario, error) {
	// Verificar si ya existe inventario para ese producto en esa sucursal
	existing, err := s.store.GetInventarioByProductoYSucursal(ctx, inventario.IDProducto, inventario.IDSucursal)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("ya existe un registro de inventario para este producto en esta sucursal")
	}

	return s.store.CreateInventario(ctx, &inventario)
}

func (s *ServiceInventario) UpdateInventario(ctx context.Context, id uuid.UUID, inventario models.Inventario) (*models.Inventario, error) {
	return s.store.UpdateInventario(ctx, id, &inventario)
}

func (s *ServiceInventario) DeleteInventario(ctx context.Context, id uuid.UUID) error {
	return s.store.DeleteInventario(ctx, id)
}
