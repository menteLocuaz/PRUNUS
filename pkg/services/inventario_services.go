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

func (s *ServiceInventario) RegistrarMovimiento(ctx context.Context, m models.MovimientoInventario) (*models.MovimientoInventario, error) {
	return s.store.RegistrarMovimiento(ctx, &m)
}

func (s *ServiceInventario) RegistrarMovimientoMasivo(ctx context.Context, idSucursal, idUsuario uuid.UUID, tipoMov, referencia string, items []models.MovimientoItem) ([]*models.MovimientoInventario, error) {
	if len(items) == 0 {
		return nil, errors.New("debe proporcionar al menos un item para el movimiento")
	}
	return s.store.RegistrarMovimientoMasivo(ctx, idSucursal, idUsuario, tipoMov, referencia, items)
}

func (s *ServiceInventario) GetMovimientos(ctx context.Context, productoID uuid.UUID) ([]*models.MovimientoInventario, error) {
	return s.store.GetMovimientosByProducto(ctx, productoID)
}

func (s *ServiceInventario) GetAlertasStock(ctx context.Context, sucursalID uuid.UUID) ([]*models.Inventario, error) {
	return s.store.GetAlertasStock(ctx, sucursalID)
}

func (s *ServiceInventario) GetValuacion(ctx context.Context, sucursalID uuid.UUID) (float64, error) {
	return s.store.GetValuacion(ctx, sucursalID)
}
