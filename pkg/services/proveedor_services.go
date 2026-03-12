package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
)

type ServiceProveedor struct {
	store  store.StoreProveedor
	logger *slog.Logger
}

func NewServiceProveedor(s store.StoreProveedor, logger *slog.Logger) *ServiceProveedor {
	return &ServiceProveedor{
		store:  s,
		logger: logger,
	}
}

func (s *ServiceProveedor) GetAllProveedores(ctx context.Context) ([]*models.Proveedor, error) {
	return s.store.GetAllProveedores(ctx)
}

func (s *ServiceProveedor) GetProveedorByID(ctx context.Context, id uuid.UUID) (*models.Proveedor, error) {
	if id == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de obtener proveedor con ID nulo")
		return nil, errors.New("el ID del proveedor es requerido")
	}
	return s.store.GetProveedorByID(ctx, id)
}

func (s *ServiceProveedor) CreateProveedor(ctx context.Context, proveedor models.Proveedor) (*models.Proveedor, error) {
	if proveedor.Nombre == "" {
		s.logger.WarnContext(ctx, "Intento de creación de proveedor con nombre vacío")
		return nil, errors.New("falta el nombre del proveedor")
	}
	if proveedor.IDSucursal == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de creación de proveedor sin sucursal", slog.String("nombre", proveedor.Nombre))
		return nil, errors.New("falta el id de la sucursal")
	}
	if proveedor.IDEmpresa == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de creación de proveedor sin empresa", slog.String("nombre", proveedor.Nombre))
		return nil, errors.New("falta el id de la empresa")
	}
	if proveedor.IDStatus == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de creación de proveedor sin estatus", slog.String("nombre", proveedor.Nombre))
		return nil, errors.New("falta el id de estatus")
	}
	return s.store.CreateProveedor(ctx, &proveedor)
}

func (s *ServiceProveedor) UpdateProveedor(ctx context.Context, id uuid.UUID, proveedor models.Proveedor) (*models.Proveedor, error) {
	if id == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de actualización de proveedor con ID nulo")
		return nil, errors.New("el ID del proveedor es requerido")
	}
	if proveedor.Nombre == "" {
		s.logger.WarnContext(ctx, "Intento de actualización de proveedor con nombre vacío", slog.String("id", id.String()))
		return nil, errors.New("falta el nombre del proveedor")
	}
	return s.store.UpdateProveedor(ctx, id, &proveedor)
}

func (s *ServiceProveedor) DeleteProveedor(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de eliminación de proveedor con ID nulo")
		return errors.New("el ID del proveedor es requerido")
	}
	return s.store.DeleteProveedor(ctx, id)
}
