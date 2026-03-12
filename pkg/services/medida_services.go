package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
)

type ServiceUnidad struct {
	store  store.StoreUnidad
	logger *slog.Logger
}

func NewServiceUnidad(s store.StoreUnidad, logger *slog.Logger) *ServiceUnidad {
	return &ServiceUnidad{
		store:  s,
		logger: logger,
	}
}

func (s *ServiceUnidad) GetAllUnidades(ctx context.Context) ([]*models.Unidad, error) {
	return s.store.GetAllUnidades(ctx)
}

func (s *ServiceUnidad) GetUnidadByID(ctx context.Context, id uuid.UUID) (*models.Unidad, error) {
	if id == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de obtener unidad con ID nulo")
		return nil, errors.New("el ID de la unidad es requerido")
	}
	return s.store.GetUnidadByID(ctx, id)
}

func (s *ServiceUnidad) CreateUnidad(ctx context.Context, unidad models.Unidad) (*models.Unidad, error) {
	if unidad.Nombre == "" {
		s.logger.WarnContext(ctx, "Intento de creación de unidad con nombre vacío")
		return nil, errors.New("falta el nombre de la unidad")
	}
	if unidad.IDSucursal == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de creación de unidad sin sucursal", slog.String("nombre", unidad.Nombre))
		return nil, errors.New("falta el id de la sucursal")
	}
	return s.store.CreateUnidad(ctx, &unidad)
}

func (s *ServiceUnidad) UpdateUnidad(ctx context.Context, id uuid.UUID, unidad models.Unidad) (*models.Unidad, error) {
	if id == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de actualización de unidad con ID nulo")
		return nil, errors.New("el ID de la unidad es requerido")
	}
	if unidad.Nombre == "" {
		s.logger.WarnContext(ctx, "Intento de actualización de unidad con nombre vacío", slog.String("id", id.String()))
		return nil, errors.New("falta el nombre de la unidad")
	}
	return s.store.UpdateUnidad(ctx, id, &unidad)
}

func (s *ServiceUnidad) DeleteUnidad(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de eliminación de unidad con ID nulo")
		return errors.New("el ID de la unidad es requerido")
	}
	return s.store.DeleteUnidad(ctx, id)
}
