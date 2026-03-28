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
func (s *ServicePOS) DesmontarCajero(ctx context.Context, ctrlID uuid.UUID, usrID uuid.UUID, rstID string, motivoDescuadre string, accionInt int) error {
	// 1. Ejecutar actualizaciones en DB (Actualiza Control_Estacion y Retiros)
	err := s.store.DesmontarCajero(ctx, ctrlID, models.EstatusInactivo, models.EstatusRetiroTotal, models.EstatusDesmontado, motivoDescuadre)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error al desmontar cajero", slog.Any("error", err), slog.String("ctrl_id", ctrlID.String()))
		return err
	}

	// 2. Determinar descripción de Auditoría según accionInt (Migración de lógica de SP)
	var descAudit string
	switch accionInt {
	case 1: // Cierre normal con/sin motivo
		if motivoDescuadre != "" {
			descAudit = fmt.Sprintf("DESMONTADO DEL CAJERO CON MOTIVO DE DESCUADRE: %s", motivoDescuadre)
		} else {
			descAudit = fmt.Sprintf("DESMONTAR EN BACK OFFICE Cod Control Estacion: %s, Cod Usuario: %s, Cod Restaurante: %s",
				ctrlID.String(), usrID.String(), rstID)
		}
	case 2: // Por Administrador
		descAudit = fmt.Sprintf("DESMONTADO DEL CAJERO POR ADMINISTRADOR: %s", usrID.String())
	case 3: // Con Motivo (Explícito en SP original)
		descAudit = fmt.Sprintf("DESMONTADO DEL CAJERO CON MOTIVO DE DESCUADRE: %s", motivoDescuadre)
	default:
		descAudit = fmt.Sprintf("DESMONTADO DE CAJERO: %s", ctrlID.String())
	}

	audit := &models.AuditoriaCaja{
		IDControlEstacion: ctrlID,
		TipoMovimiento:    models.EstatusDesmontado,
		IDUsuario:         0, // El SP original pasaba 0 en el campo numérico
		Descripcion:       descAudit,
	}

	if err := s.logsStore.CreateAuditoriaCaja(ctx, audit); err != nil {
		s.logger.ErrorContext(ctx, "Error al registrar auditoría de desmontado", slog.Any("error", err))
	}

	return nil
}

// ActualizarValoresDeclarados actualiza el arqueo de caja para una forma de pago (Migración de SP)
func (s *ServicePOS) ActualizarValoresDeclarados(ctx context.Context, ctrlID, formaPagoID, userID uuid.UUID, valor float64, tpEnvID int) error {
	err := s.store.UpdateValoresDeclarados(ctx, ctrlID, formaPagoID, userID, valor, tpEnvID,
		models.EstatusRetiroEfectivo, models.EstatusRetiroTotal, models.EstatusDesmontado)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error al actualizar valores declarados", slog.Any("error", err), slog.String("ctrl_id", ctrlID.String()))
		return err
	}

	// Auditoría
	descAudit := fmt.Sprintf("ACTUALIZA VALORES DECLARADOS - Valor: %.2f, IdFormaPago: %s, Cod Usuario: %s, Control_Estación: %s",
		valor, formaPagoID.String(), userID.String(), ctrlID.String())

	audit := &models.AuditoriaCaja{
		IDControlEstacion: ctrlID,
		TipoMovimiento:    models.EstatusDesmontado,
		IDUsuario:         0,
		Descripcion:       descAudit,
	}

	if err := s.logsStore.CreateAuditoriaCaja(ctx, audit); err != nil {
		s.logger.ErrorContext(ctx, "Error al registrar auditoría de valores declarados", slog.Any("error", err))
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
