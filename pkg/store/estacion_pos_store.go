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

// StoreEstacionPos define la interfaz para el almacenamiento de estaciones POS.
type StoreEstacionPos interface {
	GetAll(ctx context.Context) ([]*models.EstacionPos, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.EstacionPos, error)
	GetBySucursal(ctx context.Context, idSucursal uuid.UUID) ([]*models.EstacionPos, error)
	Create(ctx context.Context, e *models.EstacionPos) (*models.EstacionPos, error)
	Update(ctx context.Context, id uuid.UUID, e *models.EstacionPos) (*models.EstacionPos, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type storeEstacionPos struct {
	db *sql.DB
}

// NewEstacionPosStore crea una nueva instancia del almacén de estaciones POS.
func NewEstacionPosStore(db *sql.DB) StoreEstacionPos {
	return &storeEstacionPos{db: db}
}

const estacionPosSelectFields = `
	id_estacion, codigo, nombre, ip, id_sucursal, id_status, created_at, updated_at
`

func (s *storeEstacionPos) scanRow(scanner interface{ Scan(dest ...any) error }, e *models.EstacionPos) error {
	return scanner.Scan(
		&e.IDEstacion, &e.Codigo, &e.Nombre, &e.IP, &e.IDSucursal, &e.IDStatus, &e.CreatedAt, &e.UpdatedAt,
	)
}

func (s *storeEstacionPos) GetAll(ctx context.Context) ([]*models.EstacionPos, error) {
	defer performance.Trace(ctx, "store", "GetAllEstacionPos", performance.DbThreshold, time.Now())
	query := fmt.Sprintf(`SELECT %s FROM estaciones_pos WHERE deleted_at IS NULL`, estacionPosSelectFields)
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener estaciones: %w", err)
	}
	defer rows.Close()

	var estaciones []*models.EstacionPos
	for rows.Next() {
		e := &models.EstacionPos{}
		if err := s.scanRow(rows, e); err != nil {
			return nil, fmt.Errorf("error al escanear estación: %w", err)
		}
		estaciones = append(estaciones, e)
	}
	return estaciones, nil
}

func (s *storeEstacionPos) GetByID(ctx context.Context, id uuid.UUID) (*models.EstacionPos, error) {
	defer performance.Trace(ctx, "store", "GetEstacionPosByID", performance.DbThreshold, time.Now())
	query := fmt.Sprintf(`SELECT %s FROM estaciones_pos WHERE id_estacion = $1 AND deleted_at IS NULL`, estacionPosSelectFields)
	e := &models.EstacionPos{}
	err := s.scanRow(s.db.QueryRowContext(ctx, query, id), e)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("estación no encontrada")
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener estación: %w", err)
	}
	return e, nil
}

func (s *storeEstacionPos) GetBySucursal(ctx context.Context, idSucursal uuid.UUID) ([]*models.EstacionPos, error) {
	defer performance.Trace(ctx, "store", "GetEstacionPosBySucursal", performance.DbThreshold, time.Now())
	query := fmt.Sprintf(`SELECT %s FROM estaciones_pos WHERE id_sucursal = $1 AND deleted_at IS NULL`, estacionPosSelectFields)
	rows, err := s.db.QueryContext(ctx, query, idSucursal)
	if err != nil {
		return nil, fmt.Errorf("error al obtener estaciones por sucursal: %w", err)
	}
	defer rows.Close()

	var estaciones []*models.EstacionPos
	for rows.Next() {
		e := &models.EstacionPos{}
		if err := s.scanRow(rows, e); err != nil {
			return nil, fmt.Errorf("error al escanear estación: %w", err)
		}
		estaciones = append(estaciones, e)
	}
	return estaciones, nil
}

func (s *storeEstacionPos) Create(ctx context.Context, e *models.EstacionPos) (*models.EstacionPos, error) {
	defer performance.Trace(ctx, "store", "CreateEstacionPos", performance.DbThreshold, time.Now())

	err := ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		query := `
			INSERT INTO estaciones_pos (codigo, nombre, ip, id_sucursal, id_status) 
			VALUES ($1, $2, $3, $4, $5) 
			RETURNING id_estacion, created_at, updated_at`

		return tx.QueryRowContext(ctx, query,
			e.Codigo, e.Nombre, e.IP, e.IDSucursal, e.IDStatus,
		).Scan(&e.IDEstacion, &e.CreatedAt, &e.UpdatedAt)
	})

	if err != nil {
		return nil, fmt.Errorf("error al crear estación: %w", err)
	}
	return e, nil
}

func (s *storeEstacionPos) Update(ctx context.Context, id uuid.UUID, e *models.EstacionPos) (*models.EstacionPos, error) {
	defer performance.Trace(ctx, "store", "UpdateEstacionPos", performance.DbThreshold, time.Now())

	err := ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		query := `
			UPDATE estaciones_pos 
			SET codigo = $1, nombre = $2, ip = $3, id_sucursal = $4, id_status = $5, updated_at = CURRENT_TIMESTAMP 
			WHERE id_estacion = $6 AND deleted_at IS NULL 
			RETURNING created_at, updated_at`

		return tx.QueryRowContext(ctx, query,
			e.Codigo, e.Nombre, e.IP, e.IDSucursal, e.IDStatus, id,
		).Scan(&e.CreatedAt, &e.UpdatedAt)
	})

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("estación no encontrada")
	}
	if err != nil {
		return nil, fmt.Errorf("error al actualizar estación: %w", err)
	}

	e.IDEstacion = id
	return e, nil
}

func (s *storeEstacionPos) Delete(ctx context.Context, id uuid.UUID) error {
	defer performance.Trace(ctx, "store", "DeleteEstacionPos", performance.DbThreshold, time.Now())

	return ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		query := `UPDATE estaciones_pos SET deleted_at = CURRENT_TIMESTAMP WHERE id_estacion = $1 AND deleted_at IS NULL`
		result, err := tx.ExecContext(ctx, query, id)
		if err != nil {
			return err
		}
		rows, _ := result.RowsAffected()
		if rows == 0 {
			return fmt.Errorf("estación no encontrada")
		}
		return nil
	})
}
