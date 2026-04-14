package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
	zaplogger "github.com/prunus/pkg/utils/logger"
	"go.uber.org/zap"
)

// ServiceEstacionPos encapsula la lógica de negocio para la gestión de estaciones POS.
type ServiceEstacionPos struct {
	store  store.StoreEstacionPos
	logger *zap.Logger
}

// NewServiceEstacionPos crea una nueva instancia del servicio de estaciones POS.
func NewServiceEstacionPos(s store.StoreEstacionPos, logger *zap.Logger) *ServiceEstacionPos {
	return &ServiceEstacionPos{
		store:  s,
		logger: logger,
	}
}

func (s *ServiceEstacionPos) validate(e *models.EstacionPos) error {
	if e.Codigo == "" {
		return errors.New("el código de la estación es requerido")
	}
	if e.Nombre == "" {
		return errors.New("el nombre de la estación es requerido")
	}
	if e.IDSucursal == uuid.Nil {
		return errors.New("el ID de la sucursal es requerido")
	}
	if e.IDStatus == uuid.Nil {
		return errors.New("el ID del estatus es requerido")
	}
	return nil
}

// GetAll obtiene todas las estaciones POS registradas.
func (s *ServiceEstacionPos) GetAll(ctx context.Context) ([]*models.EstacionPos, error) {
	zaplogger.WithContext(ctx, s.logger).Info("Obteniendo todas las estaciones POS")
	return s.store.GetAll(ctx)
}

// GetByID obtiene una estación POS específica por su identificador único.
func (s *ServiceEstacionPos) GetByID(ctx context.Context, id uuid.UUID) (*models.EstacionPos, error) {
	if id == uuid.Nil {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de obtener estación con ID nulo")
		return nil, errors.New("el ID de la estación es requerido")
	}

	estacion, err := s.store.GetByID(ctx, id)
	if err != nil {
		zaplogger.WithContext(ctx, s.logger).Error("Error al obtener estación por ID",
			zap.Error(err),
			zap.String("id", id.String()),
		)
		return nil, err
	}
	return estacion, nil
}

// GetBySucursal obtiene todas las estaciones asociadas a una sucursal específica.
func (s *ServiceEstacionPos) GetBySucursal(ctx context.Context, idSucursal uuid.UUID) ([]*models.EstacionPos, error) {
	if idSucursal == uuid.Nil {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de obtener estaciones con ID de sucursal nulo")
		return nil, errors.New("el ID de la sucursal es requerido")
	}

	return s.store.GetBySucursal(ctx, idSucursal)
}

// Create registra una nueva estación POS en el sistema.
func (s *ServiceEstacionPos) Create(ctx context.Context, e models.EstacionPos) (*models.EstacionPos, error) {
	if err := s.validate(&e); err != nil {
		zaplogger.WithContext(ctx, s.logger).Warn("Validación fallida al crear estación", zap.Error(err))
		return nil, err
	}

	res, err := s.store.Create(ctx, &e)
	if err != nil {
		zaplogger.WithContext(ctx, s.logger).Error("Error al crear estación en la base de datos", zap.Error(err))
		return nil, fmt.Errorf("no se pudo crear la estación: %w", err)
	}

	zaplogger.WithContext(ctx, s.logger).Info("Estación POS creada exitosamente", zap.String("id", res.IDEstacion.String()))
	return res, nil
}

// Update actualiza la información de una estación POS existente.
func (s *ServiceEstacionPos) Update(ctx context.Context, id uuid.UUID, e models.EstacionPos) (*models.EstacionPos, error) {
	if id == uuid.Nil {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de actualizar estación con ID nulo")
		return nil, errors.New("el ID de la estación es requerido")
	}

	if err := s.validate(&e); err != nil {
		zaplogger.WithContext(ctx, s.logger).Warn("Validación fallida al actualizar estación",
			zap.String("id", id.String()),
			zap.Error(err),
		)
		return nil, err
	}

	res, err := s.store.Update(ctx, id, &e)
	if err != nil {
		zaplogger.WithContext(ctx, s.logger).Error("Error al actualizar estación",
			zap.Error(err),
			zap.String("id", id.String()),
		)
		return nil, fmt.Errorf("no se pudo actualizar la estación: %w", err)
	}

	zaplogger.WithContext(ctx, s.logger).Info("Estación POS actualizada exitosamente", zap.String("id", id.String()))
	return res, nil
}

// Delete realiza una eliminación lógica (soft delete) de una estación POS.
func (s *ServiceEstacionPos) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de eliminar estación con ID nulo")
		return errors.New("el ID de la estación es requerido")
	}

	if err := s.store.Delete(ctx, id); err != nil {
		zaplogger.WithContext(ctx, s.logger).Error("Error al eliminar estación",
			zap.Error(err),
			zap.String("id", id.String()),
		)
		return err
	}

	zaplogger.WithContext(ctx, s.logger).Info("Estación POS eliminada exitosamente", zap.String("id", id.String()))
	return nil
}
