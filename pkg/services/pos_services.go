package services

import (
	"context"
	"database/sql"
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

func (s *ServicePOS) Logger() *zap.Logger {
	return s.logger
}

// AbrirCaja realiza la apertura de una estación de POS con validaciones de consistencia
func (s *ServicePOS) AbrirCaja(ctx context.Context, input dto.AbrirCajaDTO, idUsuario uuid.UUID) (*models.ControlEstacion, error) {
	// 1. Validar fondo base (Control Operativo)
	if input.FondoBase < 0 {
		return nil, errors.New("el fondo base no puede ser negativo")
	}

	// 2. Validar existencia y estado de la estación (Máquina de Estados)
	estacion, err := s.store.GetEstacionByID(ctx, input.IDEstacion)
	if err != nil {
		return nil, fmt.Errorf("la estación especificada no existe")
	}

	if estacion.IDStatus == models.EstatusInactivo {
		return nil, errors.New("la estación está deshabilitada (Inactiva)")
	}
	if estacion.IDStatus == models.EstatusBloqueada {
		return nil, errors.New("la estación está BLOQUEADA por administración")
	}

	// 3. Validar consistencia del Usuario (Un cajero -> Una caja)
	usuarioActivo, err := s.store.GetActiveControlByUser(ctx, idUsuario)
	if err != nil {
		return nil, fmt.Errorf("error al validar sesión del usuario: %w", err)
	}
	if usuarioActivo != nil {
		// Si el usuario ya tiene esta misma estación abierta, es idempotente
		if usuarioActivo.IDEstacion == input.IDEstacion {
			return usuarioActivo, nil
		}
		return nil, fmt.Errorf("el usuario ya tiene una sesión activa en la estación: %s", input.IDEstacion.String())
	}

	// 4. Validar concurrencia de la Estación (Idempotencia / Doble apertura)
	activo, err := s.store.GetActiveControlByEstacion(ctx, input.IDEstacion)
	if err != nil {
		return nil, fmt.Errorf("error al verificar sesiones activas: %w", err)
	}
	if activo != nil {
		// Si está ocupada por OTRO usuario
		if activo.UsuarioAsignado != idUsuario {
			return nil, errors.New("la estación ya está ocupada por otro usuario")
		}
		return activo, nil
	}

	// 5. Validar periodo contable activo
	periodo, err := s.store.GetActivePeriodo(ctx)
	if err != nil || periodo == nil {
		return nil, errors.New("no es posible abrir caja: no existe un periodo contable activo")
	}

	// 6. Preparar y Persistir
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
		return nil, fmt.Errorf("error al crear el registro de control: %w", err)
	}

	// 7. Actualizar estatus de estación a ocupada (Fondo Asignado)
	_ = s.store.UpdateEstacionStatus(ctx, input.IDEstacion, models.EstatusFondoAsignado)

	return result, nil
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

	// Al desmontar, la estación vuelve a estar Activa (Disponible)
	// Pero primero obtenemos el control para saber qué estación era
	// (Nota: En un flujo real, DesmontarCajero debería recibir el ID de la estación o el Store debería manejarlo)
	
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

// GetEstadoCaja obtiene el estado actual de una estación de forma optimizada
func (s *ServicePOS) GetEstadoCaja(ctx context.Context, idEstacion uuid.UUID) (*dto.EstadoCajaDTO, error) {
	result, err := s.store.GetEstadoCompletoEstacion(ctx, idEstacion)
	if err != nil {
		if err != sql.ErrNoRows {
			zaplogger.WithContext(ctx, s.logger).Error("Error al obtener estado completo de estación",
				zap.Error(err),
				zap.String("id_estacion", idEstacion.String()),
			)
		}
		return nil, err
	}

	return result, nil
}
