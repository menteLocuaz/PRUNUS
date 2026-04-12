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
	if proveedor.RazonSocial == "" {
		s.logger.WarnContext(ctx, "Intento de creación de proveedor con razón social vacía")
		return nil, errors.New("la razón social del proveedor es obligatoria")
	}
	if proveedor.NitRut == "" {
		s.logger.WarnContext(ctx, "Intento de creación de proveedor sin NIT/RUT", slog.String("razon_social", proveedor.RazonSocial))
		return nil, errors.New("el NIT/RUT del proveedor es obligatorio")
	}
	if proveedor.IDStatus == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de creación de proveedor sin estatus", slog.String("razon_social", proveedor.RazonSocial))
		return nil, errors.New("el ID de estatus es obligatorio")
	}
	return s.store.CreateProveedor(ctx, &proveedor)
}

func (s *ServiceProveedor) UpdateProveedor(ctx context.Context, id uuid.UUID, proveedor models.Proveedor) (*models.Proveedor, error) {
	if id == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de actualización de proveedor con ID nulo")
		return nil, errors.New("el ID del proveedor es requerido")
	}
	if proveedor.RazonSocial == "" {
		s.logger.WarnContext(ctx, "Intento de actualización de proveedor con razón social vacía", slog.String("id", id.String()))
		return nil, errors.New("la razón social del proveedor es obligatoria")
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
