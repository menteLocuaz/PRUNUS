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

// StoreEstatus define la interfaz para el acceso a datos de la tabla estatus.
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

// NewEstatus crea una nueva instancia de StoreEstatus.
func NewEstatus(db *sql.DB) StoreEstatus {
	return &storeEstatus{db: db}
}

// estatusSelectFields incluye todos los campos necesarios para poblar el modelo models.Estatus.
// Se utiliza COALESCE para manejar posibles valores nulos en columnas agregadas recientemente.
const estatusSelectFields = `
	id_status, 
	std_descripcion, 
	COALESCE(std_tipo_estado, '') as std_tipo_estado, 
	COALESCE(factor, '') as factor,
	COALESCE(nivel, 0) as nivel,
	mdl_id, 
	is_active,
	created_at, 
	updated_at,
	deleted_at
`

// scanRowEstatus centraliza el mapeo de las columnas de la base de datos al struct del modelo.
func (s *storeEstatus) scanRowEstatus(scanner interface{ Scan(dest ...any) error }, e *models.Estatus) error {
	return scanner.Scan(
		&e.IDStatus,
		&e.StdDescripcion,
		&e.StdTipoEstado,
		&e.Factor,
		&e.Nivel,
		&e.MdlID,
		&e.IsActive,
		&e.CreatedAt,
		&e.UpdatedAt,
		&e.DeletedAt,
	)
}

func (s *storeEstatus) GetAllEstatus(ctx context.Context) ([]*models.Estatus, error) {
	defer performance.Trace(ctx, "store", "GetAllEstatus", performance.DbThreshold, time.Now())
	query := fmt.Sprintf(`
		SELECT %s 
		FROM estatus 
		WHERE deleted_at IS NULL 
		ORDER BY mdl_id, nivel, std_descripcion
	`, estatusSelectFields)

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener estatus: %w", err)
	}
	defer rows.Close()

	var estatusList []*models.Estatus
	for rows.Next() {
		e := &models.Estatus{}
		if err := s.scanRowEstatus(rows, e); err != nil {
			return nil, fmt.Errorf("error al escanear estatus: %w", err)
		}
		estatusList = append(estatusList, e)
	}

	return estatusList, nil
}

func (s *storeEstatus) GetEstatusByID(ctx context.Context, id uuid.UUID) (*models.Estatus, error) {
	defer performance.Trace(ctx, "store", "GetEstatusByID", performance.DbThreshold, time.Now())
	query := fmt.Sprintf(`
		SELECT %s 
		FROM estatus 
		WHERE id_status = $1 AND deleted_at IS NULL
	`, estatusSelectFields)

	e := &models.Estatus{}
	err := s.scanRowEstatus(s.db.QueryRowContext(ctx, query, id), e)

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
	query := fmt.Sprintf(`
		SELECT %s 
		FROM estatus 
		WHERE std_tipo_estado = $1 AND deleted_at IS NULL
		ORDER BY std_descripcion ASC
	`, estatusSelectFields)

	rows, err := s.db.QueryContext(ctx, query, tipo)
	if err != nil {
		return nil, fmt.Errorf("error al obtener estatus por tipo: %w", err)
	}
	defer rows.Close()

	var estatusList []*models.Estatus
	for rows.Next() {
		e := &models.Estatus{}
		if err := s.scanRowEstatus(rows, e); err != nil {
			return nil, fmt.Errorf("error al escanear estatus por tipo: %w", err)
		}
		estatusList = append(estatusList, e)
	}

	return estatusList, nil
}

func (s *storeEstatus) GetEstatusByModulo(ctx context.Context, moduloID int) ([]*models.Estatus, error) {
	defer performance.Trace(ctx, "store", "GetEstatusByModulo", performance.DbThreshold, time.Now())
	query := fmt.Sprintf(`
		SELECT %s 
		FROM estatus 
		WHERE mdl_id = $1 AND deleted_at IS NULL
		ORDER BY nivel, std_descripcion ASC
	`, estatusSelectFields)

	rows, err := s.db.QueryContext(ctx, query, moduloID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener estatus por modulo: %w", err)
	}
	defer rows.Close()

	var estatusList []*models.Estatus
	for rows.Next() {
		e := &models.Estatus{}
		if err := s.scanRowEstatus(rows, e); err != nil {
			return nil, fmt.Errorf("error al escanear estatus por modulo: %w", err)
		}
		estatusList = append(estatusList, e)
	}

	return estatusList, nil
}

func (s *storeEstatus) CreateEstatus(ctx context.Context, estatus *models.Estatus) (*models.Estatus, error) {
	defer performance.Trace(ctx, "store", "CreateEstatus", performance.DbThreshold, time.Now())
	
	err := ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		query := `
			INSERT INTO estatus (std_descripcion, std_tipo_estado, factor, nivel, mdl_id, is_active)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id_status, created_at, updated_at
		`
		return tx.QueryRowContext(ctx, query, 
			estatus.StdDescripcion, 
			estatus.StdTipoEstado, 
			estatus.Factor,
			estatus.Nivel,
			estatus.MdlID,
			estatus.IsActive,
		).Scan(
			&estatus.IDStatus,
			&estatus.CreatedAt,
			&estatus.UpdatedAt,
		)
	})

	if err != nil {
		return nil, fmt.Errorf("error al crear estatus: %w", err)
	}

	return estatus, nil
}

func (s *storeEstatus) UpdateEstatus(ctx context.Context, id uuid.UUID, estatus *models.Estatus) (*models.Estatus, error) {
	defer performance.Trace(ctx, "store", "UpdateEstatus", performance.DbThreshold, time.Now())
	
	err := ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		query := `
			UPDATE estatus
			SET 
				std_descripcion = $1, 
				std_tipo_estado = $2, 
				factor = $3,
				nivel = $4,
				mdl_id = $5, 
				is_active = $6,
				updated_at = CURRENT_TIMESTAMP
			WHERE id_status = $7 AND deleted_at IS NULL
			RETURNING id_status, created_at, updated_at
		`
		return tx.QueryRowContext(ctx, query, 
			estatus.StdDescripcion, 
			estatus.StdTipoEstado, 
			estatus.Factor,
			estatus.Nivel,
			estatus.MdlID, 
			estatus.IsActive,
			id,
		).Scan(
			&estatus.IDStatus,
			&estatus.CreatedAt,
			&estatus.UpdatedAt,
		)
	})

	if err != nil {
		return nil, fmt.Errorf("error al actualizar estatus: %w", err)
	}

	return estatus, nil
}

func (s *storeEstatus) DeleteEstatus(ctx context.Context, id uuid.UUID) error {
	defer performance.Trace(ctx, "store", "DeleteEstatus", performance.DbThreshold, time.Now())
	
	return ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		query := `UPDATE estatus SET deleted_at = CURRENT_TIMESTAMP WHERE id_status = $1 AND deleted_at IS NULL`
		result, err := tx.ExecContext(ctx, query, id)
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
	})
}
