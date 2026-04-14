package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
	"go.uber.org/zap"
)

type ServicePeriodo struct {
	store    store.StorePeriodo
	posStore store.StorePOS
	logger   *zap.Logger
}

func NewServicePeriodo(s store.StorePeriodo, ps store.StorePOS, logger *zap.Logger) *ServicePeriodo {
	return &ServicePeriodo{
		store:    s,
		posStore: ps,
		logger:   logger,
	}
}

func (s *ServicePeriodo) AbrirNuevoPeriodo(ctx context.Context, idUsuario uuid.UUID) (*models.Periodo, error) {
	activo, _ := s.store.GetActivePeriodo(ctx)
	if activo != nil {
		return nil, errors.New("ya existe un periodo abierto actualmente")
	}

	nuevo := &models.Periodo{
		PrdFechaApertura:   time.Now(),
		PrdUsuarioApertura: idUsuario,
		IDStatus:           models.EstatusActivo,
	}

	return s.store.CreatePeriodo(ctx, nuevo)
}

func (s *ServicePeriodo) FinalizarPeriodo(ctx context.Context, idPeriodo uuid.UUID, idUsuarioCierre uuid.UUID) error {
	estacionesAbiertas, err := s.posStore.GetTotalActiveControls(ctx)
	if err != nil {
		return err
	}

	if estacionesAbiertas > 0 {
		return errors.New("no se puede cerrar el periodo: hay cajas o estaciones abiertas")
	}

	return s.store.CerrarPeriodo(ctx, idPeriodo, idUsuarioCierre)
}

func (s *ServicePeriodo) GetActivePeriodo(ctx context.Context) (*models.Periodo, error) {
	return s.store.GetActivePeriodo(ctx)
}
