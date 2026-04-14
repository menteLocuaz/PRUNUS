package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
	zaplogger "github.com/prunus/pkg/utils/logger"
	"go.uber.org/zap"
)

type ServiceSucursal struct {
	store  store.StoreSucursal
	logger *zap.Logger
}

func NewServiceSucursal(s store.StoreSucursal, logger *zap.Logger) *ServiceSucursal {
	return &ServiceSucursal{
		store:  s,
		logger: logger,
	}
}

func (s *ServiceSucursal) GetAllSucursales(ctx context.Context) ([]*models.Sucursal, error) {
	return s.store.GetAllSucursales(ctx)
}

func (s *ServiceSucursal) GetSucursalByID(ctx context.Context, id uuid.UUID) (*models.Sucursal, error) {
	if id == uuid.Nil {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de obtener sucursal con ID nulo")
		return nil, errors.New("el ID de la sucursal es requerido")
	}
	return s.store.GetSucursalByID(ctx, id)
}

func (s *ServiceSucursal) CreateSucursal(ctx context.Context, sucursal models.Sucursal) (*models.Sucursal, error) {
	if sucursal.NombreSucursal == "" {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de creación de sucursal con nombre vacío")
		return nil, errors.New("falta el nombre de la sucursal")
	}
	if sucursal.IDEmpresa == uuid.Nil {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de creación de sucursal sin empresa", zap.String("nombre", sucursal.NombreSucursal))
		return nil, errors.New("falta el id de la empresa")
	}
	if sucursal.IDStatus == uuid.Nil {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de creación de sucursal sin estatus", zap.String("nombre", sucursal.NombreSucursal))
		return nil, errors.New("falta el id del estatus")
	}
	return s.store.CreateSucursal(ctx, &sucursal)
}

func (s *ServiceSucursal) UpdateSucursal(ctx context.Context, id uuid.UUID, sucursal models.Sucursal) (*models.Sucursal, error) {
	if id == uuid.Nil {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de actualización de sucursal con ID nulo")
		return nil, errors.New("el ID de la sucursal es requerido")
	}
	if sucursal.NombreSucursal == "" {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de actualización de sucursal con nombre vacío", zap.String("id", id.String()))
		return nil, errors.New("falta el nombre de la sucursal")
	}
	return s.store.UpdateSucursal(ctx, id, &sucursal)
}

func (s *ServiceSucursal) DeleteSucursal(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de eliminación de sucursal con ID nulo")
		return errors.New("el ID de la sucursal es requerido")
	}
	return s.store.DeleteSucursal(ctx, id)
}
