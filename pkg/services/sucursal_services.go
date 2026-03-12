package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
)

type ServiceSucursal struct {
	store  store.StoreSucursal
	logger *slog.Logger
}

func NewServiceSucursal(s store.StoreSucursal, logger *slog.Logger) *ServiceSucursal {
	return &ServiceSucursal{
		store:  s,
		logger: logger,
	}
}

// obtine todas las sucursales
func (s *ServiceSucursal) GetAllSucursales(ctx context.Context) ([]*models.Sucursal, error) {
	return s.store.GetAllSucursales(ctx)
}

// obtien solo una sucursla
func (s *ServiceSucursal) GetSucursalByID(ctx context.Context, id uuid.UUID) (*models.Sucursal, error) {
	if id == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de obtener sucursal con ID nulo")
		return nil, errors.New("el ID de la sucursal es requerido")
	}
	return s.store.GetSucursalByID(ctx, id)
}

// crea sucursal
func (s *ServiceSucursal) CreateSucursal(ctx context.Context, sucursal models.Sucursal) (*models.Sucursal, error) {
	if sucursal.NombreSucursal == "" {
		s.logger.WarnContext(ctx, "Intento de creación de sucursal con nombre vacío")
		return nil, errors.New("falta el nombre de la sucursal")
	}
	if sucursal.IDEmpresa == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de creación de sucursal sin empresa", slog.String("nombre", sucursal.NombreSucursal))
		return nil, errors.New("falta el id de la empresa")
	}
	if sucursal.IDStatus == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de creación de sucursal sin estatus", slog.String("nombre", sucursal.NombreSucursal))
		return nil, errors.New("falta el id del estatus")
	}

	return s.store.CreateSucursal(ctx, &sucursal)
}

// actualizar empresa
func (s *ServiceSucursal) UpdateSucursal(ctx context.Context, id uuid.UUID, sucursal models.Sucursal) (*models.Sucursal, error) {
	if id == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de actualización de sucursal con ID nulo")
		return nil, errors.New("el ID de la sucursal es requerido")
	}
	if sucursal.NombreSucursal == "" {
		s.logger.WarnContext(ctx, "Intento de actualización de sucursal con nombre vacío", slog.String("id", id.String()))
		return nil, errors.New("falta el nombre de la sucursal")
	}
	return s.store.UpdateSucursal(ctx, id, &sucursal)
}

// eliminar empresa
func (s *ServiceSucursal) DeleteSucursal(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de eliminación de sucursal con ID nulo")
		return errors.New("el ID de la sucursal es requerido")
	}
	return s.store.DeleteSucursal(ctx, id)
}
