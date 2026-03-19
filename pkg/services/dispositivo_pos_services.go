package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
)

type ServiceDispositivoPos struct {
	store  store.StoreDispositivoPos
	logger *slog.Logger
}

func NewServiceDispositivoPos(s store.StoreDispositivoPos, logger *slog.Logger) *ServiceDispositivoPos {
	return &ServiceDispositivoPos{
		store:  s,
		logger: logger,
	}
}

func (s *ServiceDispositivoPos) GetAll(ctx context.Context) ([]*models.DispositivoPos, error) {
	return s.store.GetAll(ctx)
}

func (s *ServiceDispositivoPos) GetByID(ctx context.Context, id uuid.UUID) (*models.DispositivoPos, error) {
	if id == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de obtener dispositivo con ID nulo")
		return nil, errors.New("el ID del dispositivo es requerido")
	}
	return s.store.GetByID(ctx, id)
}

func (s *ServiceDispositivoPos) Create(ctx context.Context, d models.DispositivoPos) (*models.DispositivoPos, error) {
	if d.Nombre == "" {
		s.logger.WarnContext(ctx, "Intento de crear dispositivo con nombre vacío")
		return nil, errors.New("el nombre del dispositivo es requerido")
	}
	if d.Tipo == "" {
		s.logger.WarnContext(ctx, "Intento de crear dispositivo con tipo vacío")
		return nil, errors.New("el tipo del dispositivo es requerido")
	}
	if d.IDEstacion == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de crear dispositivo sin estación")
		return nil, errors.New("el ID de la estación es requerido")
	}
	return s.store.Create(ctx, &d)
}

func (s *ServiceDispositivoPos) Update(ctx context.Context, id uuid.UUID, d models.DispositivoPos) (*models.DispositivoPos, error) {
	if id == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de actualizar dispositivo con ID nulo")
		return nil, errors.New("el ID del dispositivo es requerido")
	}
	if d.Nombre == "" {
		s.logger.WarnContext(ctx, "Intento de actualizar dispositivo con nombre vacío")
		return nil, errors.New("el nombre del dispositivo es requerido")
	}
	return s.store.Update(ctx, id, &d)
}

func (s *ServiceDispositivoPos) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de eliminar dispositivo con ID nulo")
		return errors.New("el ID del dispositivo es requerido")
	}
	return s.store.Delete(ctx, id)
}
