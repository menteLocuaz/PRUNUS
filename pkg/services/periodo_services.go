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

func (s *ServicePeriodo) Logger() *zap.Logger {
	return s.logger
}

// AbrirNuevoPeriodo inicia un periodo con controles de concurrencia y auditoría
func (s *ServicePeriodo) AbrirNuevoPeriodo(ctx context.Context, idUsuario, idSucursal uuid.UUID, ip, motivo string) (*models.Periodo, error) {
	// 1. Verificar si ya existe un periodo abierto para la sucursal
	activo, err := s.store.GetActivePeriodo(ctx, idSucursal)
	if err != nil {
		s.logger.Error("Error al consultar periodo activo", zap.Error(err), zap.String("id_sucursal", idSucursal.String()))
		return nil, err
	}

	if activo != nil {
		s.logger.Info("Se intentó abrir un periodo, pero ya existe uno activo",
			zap.String("id_periodo", activo.IDPeriodo.String()),
			zap.String("id_sucursal", idSucursal.String()))
		return activo, nil
	}

	// 2. Preparar el nuevo periodo con auditoría
	nuevo := &models.Periodo{
		IDPeriodo:          uuid.New(),
		IDSucursal:         idSucursal,
		PrdFechaApertura:   time.Now(),
		PrdUsuarioApertura: idUsuario,
		PrdIPApertura:      ip,
		PrdMotivoApertura:  motivo,
		IDStatus:           models.EstatusActivo,
	}

	// 3. Persistir (El Store maneja la concurrencia real vía UNIQUE INDEX)
	result, err := s.store.CreatePeriodo(ctx, nuevo)
	if err != nil {
		s.logger.Error("Error al crear nuevo periodo (concurrencia)", zap.Error(err))
		return nil, err
	}

	s.logger.Info("Nuevo periodo contable abierto exitosamente",
		zap.String("id_periodo", result.IDPeriodo.String()),
		zap.String("usuario", idUsuario.String()),
		zap.String("sucursal", idSucursal.String()),
		zap.String("ip", ip))

	return result, nil
}

func (s *ServicePeriodo) FinalizarPeriodo(ctx context.Context, idPeriodo uuid.UUID, idUsuarioCierre uuid.UUID, ipCierre string) error {
	// 1. Validar que no haya estaciones abiertas (Control Operativo)
	estacionesAbiertas, err := s.posStore.GetTotalActiveControls(ctx)
	if err != nil {
		return err
	}

	if estacionesAbiertas > 0 {
		return errors.New("no se puede cerrar el periodo: hay cajas o estaciones abiertas en el sistema")
	}

	// 2. Ejecutar cierre con auditoría
	return s.store.CerrarPeriodo(ctx, idPeriodo, idUsuarioCierre, ipCierre)
}

func (s *ServicePeriodo) GetActivePeriodo(ctx context.Context, idSucursal uuid.UUID) (*models.Periodo, error) {
	return s.store.GetActivePeriodo(ctx, idSucursal)
}
