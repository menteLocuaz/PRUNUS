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
	activo, err := s.store.GetActivePeriodo(ctx, idSucursal)
	if err != nil {
		s.logger.Error("Error al consultar periodo activo", zap.Error(err), zap.String("id_sucursal", idSucursal.String()))
		return nil, err
	}

	if activo != nil {
		return activo, nil
	}

	// 1.1 Resolver Estatus dinámicamente (Evita errores de FK si las constantes no coinciden)
	statusID, err := s.store.GetStatusIDByDesc(ctx, "Activo")
	if err != nil {
		s.logger.Error("Fallo al resolver estatus 'Activo' para periodo", zap.Error(err))
		return nil, errors.New("error interno: no se encontró el catálogo de estatus activo")
	}

	// 2. Preparar el nuevo periodo con auditoría
	ahora := time.Now()
	nombreDefecto := "PER-" + ahora.Format("2006-01-02-1504")

	nuevo := &models.Periodo{
		IDPeriodo:          uuid.New(),
		Nombre:             nombreDefecto,
		IDSucursal:         idSucursal,
		PrdFechaApertura:   ahora,
		PrdUsuarioApertura: idUsuario,
		PrdIPApertura:      ip,
		PrdMotivoApertura:  motivo,
		IDStatus:           statusID,
	}

	result, err := s.store.CreatePeriodo(ctx, nuevo)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// FinalizarPeriodo realiza el cierre contable con snapshot histórico
func (s *ServicePeriodo) FinalizarPeriodo(ctx context.Context, idPeriodo uuid.UUID, idUsuarioCierre uuid.UUID, ipCierre string) error {
	// 1. VALIDACIÓN OPERATIVA: No pueden haber estaciones abiertas
	estacionesAbiertas, err := s.posStore.GetTotalActiveControls(ctx)
	if err != nil {
		return err
	}

	if estacionesAbiertas > 0 {
		return errors.New("no se puede cerrar el periodo: hay estaciones de trabajo aún activas")
	}

	// 2. GENERAR SNAPSHOT DE AUDITORÍA
	snapshot, err := s.store.GenerarSnapshotPeriodo(ctx, idPeriodo)
	if err != nil {
		s.logger.Error("Fallo al generar snapshot de periodo", zap.Error(err), zap.String("id_periodo", idPeriodo.String()))
		// Continuamos con el cierre pero logueamos el error, o podrías abortar según requerimiento
	} else {
		snapshot.IDUsuarioCierre = idUsuarioCierre
		snapshot.DataJSON = map[string]interface{}{
			"ip_cierre": ipCierre,
			"timestamp": time.Now().Format(time.RFC3339),
			"version":   "1.0",
		}
		
		// Guardar snapshot persistente
		if err := s.store.GuardarSnapshot(ctx, snapshot); err != nil {
			s.logger.Error("No se pudo guardar el snapshot de cierre", zap.Error(err))
		}
	}

	// 3. EJECUTAR CIERRE DEFINITIVO
	return s.store.CerrarPeriodo(ctx, idPeriodo, idUsuarioCierre, ipCierre)
}

func (s *ServicePeriodo) GetActivePeriodo(ctx context.Context, idSucursal uuid.UUID) (*models.Periodo, error) {
	return s.store.GetActivePeriodo(ctx, idSucursal)
}
