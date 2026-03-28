package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/prunus/pkg/dto"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
)

type ServicePOS struct {
	store     store.StorePOS
	logsStore store.StoreLogs
	logger    *slog.Logger
}

func NewServicePOS(s store.StorePOS, l store.StoreLogs, logger *slog.Logger) *ServicePOS {
	return &ServicePOS{
		store:     s,
		logsStore: l,
		logger:    logger,
	}
}

// DesmontarCajero cierra la sesión de una estación y actualiza retiros (Migración de SP)
func (s *ServicePOS) DesmontarCajero(ctx context.Context, ctrlID uuid.UUID, usrID uuid.UUID, rstID string) error {
	// 1. Ejecutar actualizaciones en DB
	err := s.store.DesmontarCajero(ctx, ctrlID, models.EstatusInactivo, models.EstatusRetiroTotal, models.EstatusDesmontado)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error al desmontar cajero", slog.Any("error", err), slog.String("ctrl_id", ctrlID.String()))
		return err
	}

	// 2. Registrar Auditoría
	descAudit := fmt.Sprintf("DESMONTAR EN BACK OFFICE Cod Control Estacion: %s, Cod Usuario: %s, Cod Restaurante: %s",
		ctrlID.String(), usrID.String(), rstID)

	audit := &models.AuditoriaCaja{
		IDControlEstacion: ctrlID,
		TipoMovimiento:    models.EstatusDesmontado,
		IDUsuario:         0, // El SP original pasaba 0 en el campo numérico
		Descripcion:       descAudit,
	}

	if err := s.logsStore.CreateAuditoriaCaja(ctx, audit); err != nil {
		s.logger.ErrorContext(ctx, "Error al registrar auditoría de desmontado", slog.Any("error", err))
		// No retornamos error aquí para no revertir la operación principal si falla solo el log
	}

	return nil
}

// AbrirCaja realiza la apertura de una estación de POS
func (s *ServicePOS) AbrirCaja(ctx context.Context, input dto.AbrirCajaDTO, idUsuario uuid.UUID) (*models.ControlEstacion, error) {
	// 1. Validar que la estación existe
	_, err := s.store.GetEstacionByID(ctx, input.IDEstacion)
	if err != nil {
		s.logger.WarnContext(ctx, "Intento de abrir caja en estación inexistente", slog.String("id_estacion", input.IDEstacion.String()))
		return nil, fmt.Errorf("estación no encontrada: %w", err)
	}

	// 2. Validar que no haya una sesión activa en esta estación
	activo, err := s.store.GetActiveControlByEstacion(ctx, input.IDEstacion)
	if err != nil {
		return nil, err
	}
	if activo != nil {
		s.logger.WarnContext(ctx, "Intento de abrir caja en estación con sesión activa", slog.String("id_estacion", input.IDEstacion.String()))
		return nil, errors.New("la estación ya tiene una sesión abierta")
	}

	// 3. Validar que exista un periodo activo
	periodo, err := s.store.GetActivePeriodo(ctx)
	if err != nil || periodo == nil {
		s.logger.WarnContext(ctx, "Intento de abrir caja sin periodo activo")
		return nil, errors.New("operación no permitida: no hay un periodo contable abierto por administración")
	}

	// 4. Crear el registro de control de estación
	control := &models.ControlEstacion{
		IDEstacion:      input.IDEstacion,
		FondoBase:       input.FondoBase,
		UsuarioAsignado: idUsuario,
		IDStatus:        models.EstatusFondoAsignado,
		IDUserPos:       input.IDUserPos,
		IDPeriodo:       periodo.IDPeriodo,
	}

	result, err := s.store.CreateControlEstacion(ctx, control)
	if err != nil {
		return nil, err
	}

	// 5. Actualizar estatus de la estación
	err = s.store.UpdateEstacionStatus(ctx, input.IDEstacion, models.EstatusFondoAsignado)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error al actualizar estatus de estación", slog.Any("error", err))
	}

	return result, nil
}

// GetEstadoCaja obtiene el estado actual de una estación
func (s *ServicePOS) GetEstadoCaja(ctx context.Context, idEstacion uuid.UUID) (*dto.EstadoCajaDTO, error) {
	estacion, err := s.store.GetEstacionByID(ctx, idEstacion)
	if err != nil {
		return nil, err
	}

	control, err := s.store.GetActiveControlByEstacion(ctx, idEstacion)
	if err != nil {
		return nil, err
	}

	if control == nil {
		return &dto.EstadoCajaDTO{
			NombreEstacion:    estacion.Nombre,
			StatusDescripcion: "Cerrada",
		}, nil
	}

	return &dto.EstadoCajaDTO{
		IDControlEstacion: control.IDControlEstacion,
		NombreEstacion:    estacion.Nombre,
		FondoBase:         control.FondoBase,
		IDStatus:          control.IDStatus,
		FechaInicio:       control.FechaInicio,
	}, nil
}
