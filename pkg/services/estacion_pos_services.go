package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
)

type ServiceEstacionPos struct {
	store  store.StoreEstacionPos
	logger *slog.Logger
}

func NewServiceEstacionPos(s store.StoreEstacionPos, logger *slog.Logger) *ServiceEstacionPos {
	return &ServiceEstacionPos{
		store:  s,
		logger: logger,
	}
}

func (s *ServiceEstacionPos) GetAll(ctx context.Context) ([]*models.EstacionPos, error) {
	return s.store.GetAll(ctx)
}

func (s *ServiceEstacionPos) GetByID(ctx context.Context, id uuid.UUID) (*models.EstacionPos, error) {
	if id == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de obtener estación con ID nulo")
		return nil, errors.New("el ID de la estación es requerido")
	}
	return s.store.GetByID(ctx, id)
}

func (s *ServiceEstacionPos) Create(ctx context.Context, e models.EstacionPos) (*models.EstacionPos, error) {
	if e.Codigo == "" {
		s.logger.WarnContext(ctx, "Intento de crear estación con código vacío")
		return nil, errors.New("el código de la estación es requerido")
	}
	if e.Nombre == "" {
		s.logger.WarnContext(ctx, "Intento de crear estación con nombre vacío")
		return nil, errors.New("el nombre de la estación es requerido")
	}
	if e.IDSucursal == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de crear estación sin sucursal")
		return nil, errors.New("el ID de la sucursal es requerido")
	}
	if e.IDStatus == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de crear estación sin estatus")
		return nil, errors.New("el ID del estatus es requerido")
	}
	return s.store.Create(ctx, &e)
}

func (s *ServiceEstacionPos) Update(ctx context.Context, id uuid.UUID, e models.EstacionPos) (*models.EstacionPos, error) {
	if id == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de actualizar estación con ID nulo")
		return nil, errors.New("el ID de la estación es requerido")
	}
	if e.Nombre == "" {
		s.logger.WarnContext(ctx, "Intento de actualizar estación con nombre vacío")
		return nil, errors.New("el nombre de la estación es requerido")
	}
	return s.store.Update(ctx, id, &e)
}

func (s *ServiceEstacionPos) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de eliminar estación con ID nulo")
		return errors.New("el ID de la estación es requerido")
	}
	return s.store.Delete(ctx, id)
}
