package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/prunus/pkg/dto"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
	zaplogger "github.com/prunus/pkg/utils/logger"
	"go.uber.org/zap"
)

type ServiceProveedor struct {
	store  store.StoreProveedor
	logger *zap.Logger
}

func NewServiceProveedor(s store.StoreProveedor, logger *zap.Logger) *ServiceProveedor {
	return &ServiceProveedor{
		store:  s,
		logger: logger,
	}
}

func (s *ServiceProveedor) GetAllProveedores(ctx context.Context, params dto.PaginationParams) ([]*models.Proveedor, error) {
	return s.store.GetAllProveedores(ctx, params)
}

func (s *ServiceProveedor) GetProveedorByID(ctx context.Context, id uuid.UUID) (*models.Proveedor, error) {
	if id == uuid.Nil {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de obtener proveedor con ID nulo")
		return nil, errors.New("el ID del proveedor es requerido")
	}
	return s.store.GetProveedorByID(ctx, id)
}

func (s *ServiceProveedor) CreateProveedor(ctx context.Context, proveedor models.Proveedor) (*models.Proveedor, error) {
	if proveedor.RazonSocial == "" {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de creación de proveedor con razón social vacía")
		return nil, errors.New("la razón social del proveedor es obligatoria")
	}
	if proveedor.NitRut == "" {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de creación de proveedor sin NIT/RUT", zap.String("razon_social", proveedor.RazonSocial))
		return nil, errors.New("el NIT/RUT del proveedor es obligatorio")
	}
	if proveedor.IDStatus == uuid.Nil {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de creación de proveedor sin estatus", zap.String("razon_social", proveedor.RazonSocial))
		return nil, errors.New("el ID de estatus es obligatorio")
	}
	return s.store.CreateProveedor(ctx, &proveedor)
}

func (s *ServiceProveedor) UpdateProveedor(ctx context.Context, id uuid.UUID, proveedor models.Proveedor) (*models.Proveedor, error) {
	if id == uuid.Nil {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de actualización de proveedor con ID nulo")
		return nil, errors.New("el ID del proveedor es requerido")
	}
	if proveedor.RazonSocial == "" {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de actualización de proveedor con razón social vacía", zap.String("id", id.String()))
		return nil, errors.New("la razón social del proveedor es obligatoria")
	}
	return s.store.UpdateProveedor(ctx, id, &proveedor)
}

func (s *ServiceProveedor) DeleteProveedor(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de eliminación de proveedor con ID nulo")
		return errors.New("el ID del proveedor es requerido")
	}
	return s.store.DeleteProveedor(ctx, id)
}
