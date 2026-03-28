package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/utils/performance"
)

type StorePOS interface {
	// Control Estación
	GetActiveControlByEstacion(ctx context.Context, idEstacion uuid.UUID) (*models.ControlEstacion, error)
	CreateControlEstacion(ctx context.Context, control *models.ControlEstacion) (*models.ControlEstacion, error)
	UpdateControlEstacion(ctx context.Context, control *models.ControlEstacion) error

	// Estaciones
	GetEstacionByID(ctx context.Context, id uuid.UUID) (*models.EstacionPos, error)
	UpdateEstacionStatus(ctx context.Context, id uuid.UUID, idStatus uuid.UUID) error

	// Periodos
	GetActivePeriodo(ctx context.Context) (*models.Periodo, error)
	GetTotalActiveControls(ctx context.Context) (int, error)

	// Desmontar (Migración SP)
	DesmontarCajero(ctx context.Context, ctrlID uuid.UUID, idStatusInactivo uuid.UUID, idStatusRetiroTotal uuid.UUID, idStatusDesmontado uuid.UUID, motivoDescuadre string) error
}

type storePOS struct {
	db *sql.DB
}

func NewPOSStore(db *sql.DB) StorePOS {
	return &storePOS{db: db}
}

func (s *storePOS) GetTotalActiveControls(ctx context.Context) (int, error) {
	defer performance.Trace(ctx, "store", "GetTotalActiveControls", performance.DbThreshold, time.Now())
	query := `SELECT COUNT(*) FROM control_estacion WHERE fecha_salida IS NULL AND deleted_at IS NULL`
	var count int
	err := s.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

func (s *storePOS) GetActiveControlByEstacion(ctx context.Context, idEstacion uuid.UUID) (*models.ControlEstacion, error) {
	defer performance.Trace(ctx, "store", "GetActiveControlByEstacion", performance.DbThreshold, time.Now())
	query := `
		SELECT id_control_estacion, id_estacion, fecha_inicio, fecha_salida, fondo_base, 
		       usuario_asignado, id_status, id_user_pos, id_periodo, created_at, updated_at
		FROM control_estacion
		WHERE id_estacion = $1 AND fecha_salida IS NULL AND deleted_at IS NULL
		LIMIT 1
	`
	c := &models.ControlEstacion{}
	err := s.db.QueryRowContext(ctx, query, idEstacion).Scan(
		&c.IDControlEstacion, &c.IDEstacion, &c.FechaInicio, &c.FechaSalida, &c.FondoBase,
		&c.UsuarioAsignado, &c.IDStatus, &c.IDUserPos, &c.IDPeriodo, &c.CreatedAt, &c.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error al buscar control activo: %w", err)
	}
	return c, nil
}

func (s *storePOS) CreateControlEstacion(ctx context.Context, control *models.ControlEstacion) (*models.ControlEstacion, error) {
	defer performance.Trace(ctx, "store", "CreateControlEstacion", performance.DbThreshold, time.Now())
	query := `
		INSERT INTO control_estacion (id_estacion, fondo_base, usuario_asignado, id_status, id_user_pos, id_periodo)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id_control_estacion, fecha_inicio, created_at, updated_at
	`
	err := s.db.QueryRowContext(ctx, query,
		control.IDEstacion, control.FondoBase, control.UsuarioAsignado,
		control.IDStatus, control.IDUserPos, control.IDPeriodo,
	).Scan(&control.IDControlEstacion, &control.FechaInicio, &control.CreatedAt, &control.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("error al crear control estacion: %w", err)
	}
	return control, nil
}

func (s *storePOS) UpdateControlEstacion(ctx context.Context, control *models.ControlEstacion) error {
	defer performance.Trace(ctx, "store", "UpdateControlEstacion", performance.DbThreshold, time.Now())
	query := `
		UPDATE control_estacion
		SET fecha_salida = $1, fondo_retirado = $2, usuario_retiro_fondo = $3, 
		    id_status = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id_control_estacion = $5 AND deleted_at IS NULL
	`
	_, err := s.db.ExecContext(ctx, query,
		control.FechaSalida, control.FondoRetirado, control.UsuarioRetiroFondo,
		control.IDStatus, control.IDControlEstacion,
	)
	return err
}

func (s *storePOS) GetEstacionByID(ctx context.Context, id uuid.UUID) (*models.EstacionPos, error) {
	defer performance.Trace(ctx, "store", "GetEstacionByID", performance.DbThreshold, time.Now())
	query := `SELECT id_estacion, codigo, nombre, ip, id_sucursal, id_status FROM estaciones_pos WHERE id_estacion = $1 AND deleted_at IS NULL`
	e := &models.EstacionPos{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(&e.IDEstacion, &e.Codigo, &e.Nombre, &e.IP, &e.IDSucursal, &e.IDStatus)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (s *storePOS) UpdateEstacionStatus(ctx context.Context, id uuid.UUID, idStatus uuid.UUID) error {
	defer performance.Trace(ctx, "store", "UpdateEstacionStatus", performance.DbThreshold, time.Now())
	query := `UPDATE estaciones_pos SET id_status = $1, updated_at = CURRENT_TIMESTAMP WHERE id_estacion = $2`
	_, err := s.db.ExecContext(ctx, query, idStatus, id)
	return err
}

func (s *storePOS) GetActivePeriodo(ctx context.Context) (*models.Periodo, error) {
	defer performance.Trace(ctx, "store", "GetActivePeriodo", performance.DbThreshold, time.Now())
	query := `SELECT id_periodo, prd_fecha_apertura, id_status FROM periodo WHERE prd_fecha_cierre IS NULL AND deleted_at IS NULL LIMIT 1`
	p := &models.Periodo{}
	err := s.db.QueryRowContext(ctx, query).Scan(&p.IDPeriodo, &p.PrdFechaApertura, &p.IDStatus)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no hay un periodo activo abierto")
	}
	return p, err
}

func (s *storePOS) DesmontarCajero(ctx context.Context, ctrlID uuid.UUID, idStatusInactivo uuid.UUID, idStatusRetiroTotal uuid.UUID, idStatusDesmontado uuid.UUID, motivoDescuadre string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Cerrar sesión en Control_Estacion y grabar motivo de descuadre si existe
	queryCtrl := `UPDATE control_estacion 
	              SET id_status = $1, fecha_salida = CURRENT_TIMESTAMP, ctrc_motivo_descuadre = $2, updated_at = CURRENT_TIMESTAMP 
	              WHERE id_control_estacion = $3 AND deleted_at IS NULL`
	if _, err := tx.ExecContext(ctx, queryCtrl, idStatusInactivo, motivoDescuadre, ctrlID); err != nil {
		return fmt.Errorf("error al cerrar sesión en control_estacion: %w", err)
	}

	// 2. Actualizar retiros a "Retiro Total"
	queryRet := `UPDATE retiros SET id_status = $1, updated_at = CURRENT_TIMESTAMP WHERE id_control_estacion = $2 AND id_status = $3 AND deleted_at IS NULL`
	if _, err := tx.ExecContext(ctx, queryRet, idStatusRetiroTotal, ctrlID, idStatusDesmontado); err != nil {
		return fmt.Errorf("error al actualizar retiros: %w", err)
	}

	return tx.Commit()
}
