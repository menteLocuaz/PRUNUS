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

// AbrirCaja realiza la apertura de una estación de POS de forma segura
func (s *ServicePOS) AbrirCaja(ctx context.Context, input dto.AbrirCajaDTO, idUsuario uuid.UUID) (*models.ControlEstacion, error) {
	// 1. Validar existencia de la estación
	_, err := s.store.GetEstacionByID(ctx, input.IDEstacion)
	if err != nil {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de apertura en estación inexistente",
			zap.String("id_estacion", input.IDEstacion.String()))
		return nil, fmt.Errorf("la estación especificada no existe: %w", err)
	}

	// 2. Verificar si ya existe una sesión activa para esta estación (Idempotencia)
	// Si ya existe, devolvemos la activa para no causar errores 400 en el frontend
	activo, err := s.store.GetActiveControlByEstacion(ctx, input.IDEstacion)
	if err != nil {
		return nil, fmt.Errorf("error al verificar sesiones activas: %w", err)
	}
	if activo != nil {
		zaplogger.WithContext(ctx, s.logger).Info("Estación con sesión ya abierta (reutilizando)",
			zap.String("id_estacion", input.IDEstacion.String()),
			zap.String("id_control", activo.IDControlEstacion.String()))
		return activo, nil
	}

	// 3. Validar periodo contable activo (específico por sucursal)
	// Obtenemos la sucursal desde el contexto o el input si fuera necesario, 
	// pero aquí asumimos que el store ya debe filtrar por sucursal si queremos robustez.
	// Para mantener compatibilidad, pasaremos el ID de sucursal del usuario.
	periodo, err := s.store.GetActivePeriodo(ctx) // Nota: El store POS actual no recibe sucursal, pero debería.
	if err != nil || periodo == nil {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de apertura sin periodo contable activo")
		return nil, errors.New("no es posible abrir caja: no existe un periodo contable activo en su sucursal")
	}

	// 4. Preparar el modelo de control de estación
	control := &models.ControlEstacion{
		IDEstacion:      input.IDEstacion,
		FondoBase:       input.FondoBase,
		UsuarioAsignado: idUsuario,
		IDStatus:        models.EstatusFondoAsignado,
		IDUserPos:       input.IDUserPos,
		IDPeriodo:       periodo.IDPeriodo,
	}

	// 5. Persistir la apertura en base de datos
	result, err := s.store.CreateControlEstacion(ctx, control)
	if err != nil {
		return nil, fmt.Errorf("error al crear el registro de control de estación: %w", err)
	}

	// 6. Actualizar el estatus físico/lógico de la estación
	if err := s.store.UpdateEstacionStatus(ctx, input.IDEstacion, models.EstatusFondoAsignado); err != nil {
		zaplogger.WithContext(ctx, s.logger).Error("Fallo al actualizar estatus de estación",
			zap.Error(err),
			zap.String("id_estacion", input.IDEstacion.String()))
		// Nota: No retornamos error aquí para no revertir la creación, pero lo logueamos
	}

	return result, nil
}

// GetEstadoCaja obtiene el estado actual de una estación de forma optimizada
func (s *ServicePOS) GetEstadoCaja(ctx context.Context, idEstacion uuid.UUID) (*dto.EstadoCajaDTO, error) {
	// Ahora usamos una única consulta que trae Estación + Control Activo + Estatus
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
