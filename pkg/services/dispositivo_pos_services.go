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

// ServiceDispositivoPos encapsula la lógica de negocio para la gestión de dispositivos POS.
type ServiceDispositivoPos struct {
	store  store.StoreDispositivoPos
	logger *zap.Logger
}

// NewServiceDispositivoPos crea una nueva instancia del servicio de dispositivos POS.
func NewServiceDispositivoPos(s store.StoreDispositivoPos, logger *zap.Logger) *ServiceDispositivoPos {
	return &ServiceDispositivoPos{
		store:  s,
		logger: logger,
	}
}

func (s *ServiceDispositivoPos) validate(d *models.DispositivoPos) error {
	if d.Nombre == "" {
		return errors.New("el nombre del dispositivo es requerido")
	}
	if d.TipoDispositivo == "" {
		return errors.New("el tipo del dispositivo es requerido")
	}
	if d.IDEstacion == uuid.Nil {
		return errors.New("el ID de la estación es requerido")
	}
	if d.IDStatus == uuid.Nil {
		return errors.New("el ID del estatus es requerido")
	}
	return nil
}

// GetAll obtiene todos los dispositivos POS registrados.
func (s *ServiceDispositivoPos) GetAll(ctx context.Context) ([]*models.DispositivoPos, error) {
	zaplogger.WithContext(ctx, s.logger).Info("Obteniendo todos los dispositivos POS")
	return s.store.GetAll(ctx)
}

// GetByID obtiene un dispositivo POS específico por su identificador único.
func (s *ServiceDispositivoPos) GetByID(ctx context.Context, id uuid.UUID) (*models.DispositivoPos, error) {
	if id == uuid.Nil {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de obtener dispositivo con ID nulo")
		return nil, errors.New("el ID del dispositivo es requerido")
	}

	dispositivo, err := s.store.GetByID(ctx, id)
	if err != nil {
		zaplogger.WithContext(ctx, s.logger).Error("Error al obtener dispositivo por ID",
			zap.Error(err),
			zap.String("id", id.String()),
		)
		return nil, err
	}
	return dispositivo, nil
}

// GetByEstacion obtiene todos los dispositivos asociados a una estación POS.
func (s *ServiceDispositivoPos) GetByEstacion(ctx context.Context, idEstacion uuid.UUID) ([]*models.DispositivoPos, error) {
	if idEstacion == uuid.Nil {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de obtener dispositivos con ID de estación nulo")
		return nil, errors.New("el ID de la estación es requerido")
	}

	return s.store.GetByEstacion(ctx, idEstacion)
}

// Create registra un nuevo dispositivo POS en el sistema.
func (s *ServiceDispositivoPos) Create(ctx context.Context, d models.DispositivoPos) (*models.DispositivoPos, error) {
	if err := s.validate(&d); err != nil {
		zaplogger.WithContext(ctx, s.logger).Warn("Validación fallida al crear dispositivo", zap.Error(err))
		return nil, err
	}

	res, err := s.store.Create(ctx, &d)
	if err != nil {
		zaplogger.WithContext(ctx, s.logger).Error("Error al crear dispositivo en la base de datos", zap.Error(err))
		return nil, fmt.Errorf("no se pudo crear el dispositivo: %w", err)
	}

	zaplogger.WithContext(ctx, s.logger).Info("Dispositivo POS creado exitosamente", zap.String("id", res.IDDispositivo.String()))
	return res, nil
}

// Update actualiza la información de un dispositivo POS existente.
func (s *ServiceDispositivoPos) Update(ctx context.Context, id uuid.UUID, d models.DispositivoPos) (*models.DispositivoPos, error) {
	if id == uuid.Nil {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de actualizar dispositivo con ID nulo")
		return nil, errors.New("el ID del dispositivo es requerido")
	}

	if err := s.validate(&d); err != nil {
		zaplogger.WithContext(ctx, s.logger).Warn("Validación fallida al actualizar dispositivo",
			zap.String("id", id.String()),
			zap.Error(err),
		)
		return nil, err
	}

	res, err := s.store.Update(ctx, id, &d)
	if err != nil {
		zaplogger.WithContext(ctx, s.logger).Error("Error al actualizar dispositivo",
			zap.Error(err),
			zap.String("id", id.String()),
		)
		return nil, fmt.Errorf("no se pudo actualizar el dispositivo: %w", err)
	}

	zaplogger.WithContext(ctx, s.logger).Info("Dispositivo POS actualizado exitosamente", zap.String("id", id.String()))
	return res, nil
}

// Delete realiza una eliminación lógica (soft delete) de un dispositivo POS.
func (s *ServiceDispositivoPos) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de eliminar dispositivo con ID nulo")
		return errors.New("el ID del dispositivo es requerido")
	}

	if err := s.store.Delete(ctx, id); err != nil {
		zaplogger.WithContext(ctx, s.logger).Error("Error al eliminar dispositivo",
			zap.Error(err),
			zap.String("id", id.String()),
		)
		return err
	}

	zaplogger.WithContext(ctx, s.logger).Info("Dispositivo POS eliminado exitosamente", zap.String("id", id.String()))
	return nil
}
