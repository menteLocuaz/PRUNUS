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

type StoreMoneda interface {
	GetAllMonedas(ctx context.Context) ([]*models.Moneda, error)
	GetMonedaByID(ctx context.Context, id uuid.UUID) (*models.Moneda, error)
	CreateMoneda(ctx context.Context, moneda *models.Moneda) (*models.Moneda, error)
	UpdateMoneda(ctx context.Context, id uuid.UUID, moneda *models.Moneda) (*models.Moneda, error)
	DeleteMoneda(ctx context.Context, id uuid.UUID) error
}

type storeMoneda struct {
	db *sql.DB
}

func NewMoneda(db *sql.DB) StoreMoneda {
	return &storeMoneda{db: db}
}

func (s *storeMoneda) GetAllMonedas(ctx context.Context) ([]*models.Moneda, error) {
	defer performance.Trace(ctx, "store", "GetAllMonedas", performance.DbThreshold, time.Now())
	query := `
	SELECT
		m.id_moneda,
		m.nombre,
		m.id_sucursal,
		m.id_status,
		m.created_at,
		m.updated_at,

		su.id_sucursal,
		su.nombre_sucursal,
		su.id_status
	FROM moneda m
	LEFT JOIN sucursal su ON su.id_sucursal = m.id_sucursal
	WHERE m.deleted_at IS NULL
	  AND su.deleted_at IS NULL
	ORDER BY m.id_moneda
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener monedas: %w", err)
	}
	defer rows.Close()

	var monedas []*models.Moneda

	for rows.Next() {
		m := &models.Moneda{
			Sucursal: &models.Sucursal{},
		}

		if err := rows.Scan(
			&m.IDMoneda,
			&m.Nombre,
			&m.IDSucursal,
			&m.IDStatus,
			&m.CreatedAt,
			&m.UpdatedAt,

			&m.Sucursal.IDSucursal,
			&m.Sucursal.NombreSucursal,
			&m.Sucursal.IDStatus,
		); err != nil {
			return nil, fmt.Errorf("error al escanear moneda: %w", err)
		}

		monedas = append(monedas, m)
	}

	return monedas, nil
}

func (s *storeMoneda) GetMonedaByID(ctx context.Context, id uuid.UUID) (*models.Moneda, error) {
	defer performance.Trace(ctx, "store", "GetMonedaByID", performance.DbThreshold, time.Now())
	query := `
	SELECT
		m.id_moneda,
		m.nombre,
		m.id_sucursal,
		m.id_status,
		m.created_at,
		m.updated_at,

		su.id_sucursal,
		su.nombre_sucursal,
		su.id_status
	FROM moneda m
	LEFT JOIN sucursal su ON su.id_sucursal = m.id_sucursal
	WHERE m.id_moneda = $1
	  AND m.deleted_at IS NULL
	  AND su.deleted_at IS NULL
	`

	m := &models.Moneda{
		Sucursal: &models.Sucursal{},
	}

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&m.IDMoneda,
		&m.Nombre,
		&m.IDSucursal,
		&m.IDStatus,
		&m.CreatedAt,
		&m.UpdatedAt,

		&m.Sucursal.IDSucursal,
		&m.Sucursal.NombreSucursal,
		&m.Sucursal.IDStatus,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("moneda con ID %s no encontrada", id)
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener moneda: %w", err)
	}

	return m, nil
}

func (s *storeMoneda) CreateMoneda(ctx context.Context, moneda *models.Moneda) (*models.Moneda, error) {
	defer performance.Trace(ctx, "store", "CreateMoneda", performance.DbThreshold, time.Now())
	query := `INSERT INTO moneda (nombre, id_sucursal, id_status) VALUES ($1, $2, $3) RETURNING id_moneda`

	var id uuid.UUID
	err := s.db.QueryRowContext(ctx, query, moneda.Nombre, moneda.IDSucursal, moneda.IDStatus).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("error al crear moneda: %w", err)
	}

	moneda.IDMoneda = id
	return moneda, nil
}

func (s *storeMoneda) UpdateMoneda(ctx context.Context, id uuid.UUID, moneda *models.Moneda) (*models.Moneda, error) {
	defer performance.Trace(ctx, "store", "UpdateMoneda", performance.DbThreshold, time.Now())
	query := `
		UPDATE moneda
		SET
			nombre = $1,
			id_sucursal = $2,
			id_status = $3,
			updated_at = CURRENT_TIMESTAMP
		WHERE id_moneda = $4
		  AND deleted_at IS NULL
		RETURNING
			id_moneda,
			nombre,
			id_sucursal,
			id_status,
			created_at,
			updated_at
	`

	err := s.db.QueryRowContext(ctx, query, moneda.Nombre, moneda.IDSucursal, moneda.IDStatus, id).Scan(
		&moneda.IDMoneda,
		&moneda.Nombre,
		&moneda.IDSucursal,
		&moneda.IDStatus,
		&moneda.CreatedAt,
		&moneda.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("moneda con ID %s no encontrada", id)
	}
	if err != nil {
		return nil, fmt.Errorf("error al actualizar moneda: %w", err)
	}

	return moneda, nil
}

func (s *storeMoneda) DeleteMoneda(ctx context.Context, id uuid.UUID) error {
	defer performance.Trace(ctx, "store", "DeleteMoneda", performance.DbThreshold, time.Now())
	query := `UPDATE moneda SET deleted_at = $1 WHERE id_moneda = $2 AND deleted_at IS NULL`

	result, err := s.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("error al eliminar moneda: %w", err)
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

