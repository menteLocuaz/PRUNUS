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

type StoreEstatus interface {
	GetAllEstatus(ctx context.Context) ([]*models.Estatus, error)
	GetEstatusByID(ctx context.Context, id uuid.UUID) (*models.Estatus, error)
	GetEstatusByTipo(ctx context.Context, tipo string) ([]*models.Estatus, error)
	GetEstatusByModulo(ctx context.Context, moduloID int) ([]*models.Estatus, error)
	CreateEstatus(ctx context.Context, estatus *models.Estatus) (*models.Estatus, error)
	UpdateEstatus(ctx context.Context, id uuid.UUID, estatus *models.Estatus) (*models.Estatus, error)
	DeleteEstatus(ctx context.Context, id uuid.UUID) error
}

type storeEstatus struct {
	db *sql.DB
}

func NewEstatus(db *sql.DB) StoreEstatus {
	return &storeEstatus{db: db}
}

func (s *storeEstatus) GetAllEstatus(ctx context.Context) ([]*models.Estatus, error) {
	defer performance.Trace(ctx, "store", "GetAllEstatus", performance.DbThreshold, time.Now())
	query := `
	SELECT 
		id_status, 
		std_descripcion, 
		stp_tipo_estado, 
		mdl_id, 
		created_at, 
		updated_at 
	FROM estatus 
	WHERE deleted_at IS NULL 
	ORDER BY created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener estatus: %w", err)
	}
	defer rows.Close()

	var estatusList []*models.Estatus
	for rows.Next() {
		e := &models.Estatus{}
		if err := rows.Scan(
			&e.IDStatus,
			&e.StdDescripcion,
			&e.StpTipoEstado,
			&e.MdlID,
			&e.CreatedAt,
			&e.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error al escanear estatus: %w", err)
		}
		estatusList = append(estatusList, e)
	}

	return estatusList, nil
}

func (s *storeEstatus) GetEstatusByID(ctx context.Context, id uuid.UUID) (*models.Estatus, error) {
	defer performance.Trace(ctx, "store", "GetEstatusByID", performance.DbThreshold, time.Now())
	query := `
	SELECT 
		id_status, 
		std_descripcion, 
		stp_tipo_estado, 
		mdl_id, 
		created_at, 
		updated_at 
	FROM estatus 
	WHERE id_status = $1 AND deleted_at IS NULL
	`

	e := &models.Estatus{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&e.IDStatus,
		&e.StdDescripcion,
		&e.StpTipoEstado,
		&e.MdlID,
		&e.CreatedAt,
		&e.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("estatus con ID %s no encontrado", id)
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener estatus: %w", err)
	}

	return e, nil
}

func (s *storeEstatus) GetEstatusByTipo(ctx context.Context, tipo string) ([]*models.Estatus, error) {
	defer performance.Trace(ctx, "store", "GetEstatusByTipo", performance.DbThreshold, time.Now())
	query := `
	SELECT 
		id_status, 
		std_descripcion, 
		stp_tipo_estado, 
		mdl_id, 
		created_at, 
		updated_at 
	FROM estatus 
	WHERE stp_tipo_estado = $1 AND deleted_at IS NULL
	ORDER BY std_descripcion ASC
	`

	rows, err := s.db.QueryContext(ctx, query, tipo)
	if err != nil {
		return nil, fmt.Errorf("error al obtener estatus por tipo: %w", err)
	}
	defer rows.Close()

	var estatusList []*models.Estatus
	for rows.Next() {
		e := &models.Estatus{}
		if err := rows.Scan(
			&e.IDStatus,
			&e.StdDescripcion,
			&e.StpTipoEstado,
			&e.MdlID,
			&e.CreatedAt,
			&e.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error al escanear estatus por tipo: %w", err)
		}
		estatusList = append(estatusList, e)
	}

	return estatusList, nil
}

func (s *storeEstatus) GetEstatusByModulo(ctx context.Context, moduloID int) ([]*models.Estatus, error) {
	defer performance.Trace(ctx, "store", "GetEstatusByModulo", performance.DbThreshold, time.Now())
	query := `
	SELECT 
		id_status, 
		std_descripcion, 
		stp_tipo_estado, 
		mdl_id, 
		created_at, 
		updated_at 
	FROM estatus 
	WHERE mdl_id = $1 AND deleted_at IS NULL
	ORDER BY std_descripcion ASC
	`

	rows, err := s.db.QueryContext(ctx, query, moduloID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener estatus por modulo: %w", err)
	}
	defer rows.Close()

	var estatusList []*models.Estatus
	for rows.Next() {
		e := &models.Estatus{}
		if err := rows.Scan(
			&e.IDStatus,
			&e.StdDescripcion,
			&e.StpTipoEstado,
			&e.MdlID,
			&e.CreatedAt,
			&e.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error al escanear estatus por modulo: %w", err)
		}
		estatusList = append(estatusList, e)
	}

	return estatusList, nil
}

func (s *storeEstatus) CreateEstatus(ctx context.Context, estatus *models.Estatus) (*models.Estatus, error) {
	defer performance.Trace(ctx, "store", "CreateEstatus", performance.DbThreshold, time.Now())
	query := `
	INSERT INTO estatus (std_descripcion, stp_tipo_estado, mdl_id)
	VALUES ($1, $2, $3)
	RETURNING id_status, created_at, updated_at
	`

	err := s.db.QueryRowContext(ctx, query, estatus.StdDescripcion, estatus.StpTipoEstado, estatus.MdlID).Scan(
		&estatus.IDStatus,
		&estatus.CreatedAt,
		&estatus.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error al crear estatus: %w", err)
	}

	return estatus, nil
}

func (s *storeEstatus) UpdateEstatus(ctx context.Context, id uuid.UUID, estatus *models.Estatus) (*models.Estatus, error) {
	defer performance.Trace(ctx, "store", "UpdateEstatus", performance.DbThreshold, time.Now())
	query := `
	UPDATE estatus
	SET 
		std_descripcion = $1, 
		stp_tipo_estado = $2, 
		mdl_id = $3, 
		updated_at = CURRENT_TIMESTAMP
	WHERE id_status = $4 AND deleted_at IS NULL
	RETURNING id_status, std_descripcion, stp_tipo_estado, mdl_id, created_at, updated_at
	`

	err := s.db.QueryRowContext(ctx, query, estatus.StdDescripcion, estatus.StpTipoEstado, estatus.MdlID, id).Scan(
		&estatus.IDStatus,
		&estatus.StdDescripcion,
		&estatus.StpTipoEstado,
		&estatus.MdlID,
		&estatus.CreatedAt,
		&estatus.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("estatus con ID %s no encontrado", id)
	}
	if err != nil {
		return nil, fmt.Errorf("error al actualizar estatus: %w", err)
	}

	return estatus, nil
}

func (s *storeEstatus) DeleteEstatus(ctx context.Context, id uuid.UUID) error {
	defer performance.Trace(ctx, "store", "DeleteEstatus", performance.DbThreshold, time.Now())
	query := `UPDATE estatus SET deleted_at = CURRENT_TIMESTAMP WHERE id_status = $1 AND deleted_at IS NULL`

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error al eliminar estatus: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
