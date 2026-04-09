package services

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
)

type ServicePeriodo struct {
	store    store.StorePeriodo
	posStore store.StorePOS // Inyectado para validación cruzada
	logger   *slog.Logger
}

func NewServicePeriodo(s store.StorePeriodo, ps store.StorePOS, logger *slog.Logger) *ServicePeriodo {
	return &ServicePeriodo{
		store:    s,
		posStore: ps,
		logger:   logger,
	}
}

// CRUD Básico
func (s *ServicePeriodo) AbrirNuevoPeriodo(ctx context.Context, idUsuario uuid.UUID) (*models.Periodo, error) {
	// 1. Validar si ya hay un periodo activo
	activo, _ := s.store.GetActivePeriodo(ctx)
	if activo != nil {
		return nil, errors.New("ya existe un periodo abierto actualmente")
	}

	nuevo := &models.Periodo{
		PrdFechaApertura:   time.Now(),
		PrdUsuarioApertura: idUsuario,
		IDStatus:           models.EstatusActivo, // Usar estatus real del catálogo
	}

	return s.store.CreatePeriodo(ctx, nuevo)
}

// Servicio aparte del CRUD: Cierre Seguro de Periodo
func (s *ServicePeriodo) FinalizarPeriodo(ctx context.Context, idPeriodo uuid.UUID, idUsuarioCierre uuid.UUID) error {
	// 1. Validar que no existan estaciones (cajas) abiertas en este periodo
	// Usamos el posStore para consultar si hay "ControlEstacion" activos
	estacionesAbiertas, err := s.posStore.GetTotalActiveControls(ctx)
	if err != nil {
		return err
	}

	if estacionesAbiertas > 0 {
		return errors.New("no se puede cerrar el periodo: hay cajas o estaciones abiertas")
	}

	// 2. Ejecutar cierre en el store
	return s.store.CerrarPeriodo(ctx, idPeriodo, idUsuarioCierre)
}

func (s *ServicePeriodo) GetActivePeriodo(ctx context.Context) (*models.Periodo, error) {
	return s.store.GetActivePeriodo(ctx)
}
