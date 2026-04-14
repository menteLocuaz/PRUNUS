package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/prunus/pkg/dto"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
	zaplogger "github.com/prunus/pkg/utils/logger"
	"go.uber.org/zap"
)

type ServicePOS struct {
	store     store.StorePOS
	logsStore store.StoreLogs
	logger    *zap.Logger
}

func NewServicePOS(s store.StorePOS, l store.StoreLogs, logger *zap.Logger) *ServicePOS {
	return &ServicePOS{
		store:     s,
		logsStore: l,
		logger:    logger,
	}
}

// DesmontarCajero cierra la sesión de una estación y actualiza retiros
func (s *ServicePOS) DesmontarCajero(ctx context.Context, ctrlID uuid.UUID, usrID uuid.UUID, rstID string, motivoDescuadre string, accionInt int) error {
	err := s.store.DesmontarCajero(ctx, ctrlID, models.EstatusInactivo, models.EstatusRetiroTotal, models.EstatusDesmontado, motivoDescuadre)
	if err != nil {
		zaplogger.WithContext(ctx, s.logger).Error("Error al desmontar cajero",
			zap.Error(err),
			zap.String("ctrl_id", ctrlID.String()),
		)
		return err
	}

	var descAudit string
	switch accionInt {
	case 1:
		if motivoDescuadre != "" {
			descAudit = fmt.Sprintf("DESMONTADO DEL CAJERO CON MOTIVO DE DESCUADRE: %s", motivoDescuadre)
		} else {
			descAudit = fmt.Sprintf("DESMONTAR EN BACK OFFICE Cod Control Estacion: %s, Cod Usuario: %s, Cod Restaurante: %s",
				ctrlID.String(), usrID.String(), rstID)
		}
	case 2:
		descAudit = fmt.Sprintf("DESMONTADO DEL CAJERO POR ADMINISTRADOR: %s", usrID.String())
	case 3:
		descAudit = fmt.Sprintf("DESMONTADO DEL CAJERO CON MOTIVO DE DESCUADRE: %s", motivoDescuadre)
	default:
		descAudit = fmt.Sprintf("DESMONTADO DE CAJERO: %s", ctrlID.String())
	}

	audit := &models.AuditoriaCaja{
		IDControlEstacion: ctrlID,
		TipoMovimiento:    models.EstatusDesmontado,
		IDUsuario:         usrID,
		Descripcion:       descAudit,
	}

	if err := s.logsStore.CreateAuditoriaCaja(ctx, audit); err != nil {
		zaplogger.WithContext(ctx, s.logger).Error("Error al registrar auditoría de desmontado", zap.Error(err))
	}

	return nil
}

// ActualizarValoresDeclarados actualiza el arqueo de caja para una forma de pago
func (s *ServicePOS) ActualizarValoresDeclarados(ctx context.Context, ctrlID, formaPagoID, userID uuid.UUID, valor float64, tpEnvID int) error {
	err := s.store.UpdateValoresDeclarados(ctx, ctrlID, formaPagoID, userID, valor, tpEnvID,
		models.EstatusRetiroEfectivo, models.EstatusRetiroTotal, models.EstatusDesmontado)
	if err != nil {
		zaplogger.WithContext(ctx, s.logger).Error("Error al actualizar valores declarados",
			zap.Error(err),
			zap.String("ctrl_id", ctrlID.String()),
		)
		return err
	}

	descAudit := fmt.Sprintf("ACTUALIZA VALORES DECLARADOS - Valor: %.2f, IdFormaPago: %s, Cod Usuario: %s, Control_Estación: %s",
		valor, formaPagoID.String(), userID.String(), ctrlID.String())

	audit := &models.AuditoriaCaja{
		IDControlEstacion: ctrlID,
		TipoMovimiento:    models.EstatusDesmontado,
		IDUsuario:         userID,
		Descripcion:       descAudit,
	}

	if err := s.logsStore.CreateAuditoriaCaja(ctx, audit); err != nil {
		zaplogger.WithContext(ctx, s.logger).Error("Error al registrar auditoría de valores declarados", zap.Error(err))
	}

	return nil
}

// AbrirCaja realiza la apertura de una estación de POS
func (s *ServicePOS) AbrirCaja(ctx context.Context, input dto.AbrirCajaDTO, idUsuario uuid.UUID) (*models.ControlEstacion, error) {
	_, err := s.store.GetEstacionByID(ctx, input.IDEstacion)
	if err != nil {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de abrir caja en estación inexistente", zap.String("id_estacion", input.IDEstacion.String()))
		return nil, fmt.Errorf("estación no encontrada: %w", err)
	}

	activo, err := s.store.GetActiveControlByEstacion(ctx, input.IDEstacion)
	if err != nil {
		return nil, err
	}
	if activo != nil {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de abrir caja en estación con sesión activa", zap.String("id_estacion", input.IDEstacion.String()))
		return nil, errors.New("la estación ya tiene una sesión abierta")
	}

	periodo, err := s.store.GetActivePeriodo(ctx)
	if err != nil || periodo == nil {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de abrir caja sin periodo activo")
		return nil, errors.New("operación no permitida: no hay un periodo contable abierto por administración")
	}

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

	err = s.store.UpdateEstacionStatus(ctx, input.IDEstacion, models.EstatusFondoAsignado)
	if err != nil {
		zaplogger.WithContext(ctx, s.logger).Error("Error al actualizar estatus de estación", zap.Error(err))
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
