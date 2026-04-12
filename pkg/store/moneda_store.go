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

// StoreMoneda define las operaciones de persistencia para el catálogo de monedas.
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

// NewMoneda crea una nueva instancia del store de monedas.
func NewMoneda(db *sql.DB) StoreMoneda {
	return &storeMoneda{db: db}
}

// Campos base para SELECT de moneda mapeados al esquema DB.
const monedaSelectFields = `
	m.id_moneda, m.nombre, m.codigo_iso, m.simbolo, m.id_sucursal, m.id_status, m.created_at, m.updated_at,
	s.id_sucursal, s.nombre_sucursal, s.id_status
`

// scanRowMoneda centraliza el escaneo de resultados para mantener consistencia.
func (s *storeMoneda) scanRowMoneda(scanner interface{ Scan(dest ...any) error }, m *models.Moneda) error {
	if m.Sucursal == nil {
		m.Sucursal = &models.Sucursal{}
	}

	return scanner.Scan(
		&m.IDMoneda, &m.Nombre, &m.CodigoISO, &m.Simbolo, &m.IDSucursal, &m.IDStatus, &m.CreatedAt, &m.UpdatedAt,
		&m.Sucursal.IDSucursal, &m.Sucursal.NombreSucursal, &m.Sucursal.IDStatus,
	)
}

func (s *storeMoneda) GetAllMonedas(ctx context.Context) ([]*models.Moneda, error) {
	defer performance.Trace(ctx, "store", "GetAllMonedas", performance.DbThreshold, time.Now())
	query := fmt.Sprintf(`
		SELECT %s
		FROM moneda m
		LEFT JOIN sucursal s ON s.id_sucursal = m.id_sucursal
		WHERE m.deleted_at IS NULL
		ORDER BY m.nombre ASC
	`, monedaSelectFields)

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener monedas: %w", err)
	}
	defer rows.Close()

	var monedas []*models.Moneda
	for rows.Next() {
		m := &models.Moneda{}
		if err := s.scanRowMoneda(rows, m); err != nil {
			return nil, fmt.Errorf("error al escanear moneda: %w", err)
		}
		monedas = append(monedas, m)
	}

	return monedas, nil
}

func (s *storeMoneda) GetMonedaByID(ctx context.Context, id uuid.UUID) (*models.Moneda, error) {
	defer performance.Trace(ctx, "store", "GetMonedaByID", performance.DbThreshold, time.Now())
	query := fmt.Sprintf(`
		SELECT %s
		FROM moneda m
		LEFT JOIN sucursal s ON s.id_sucursal = m.id_sucursal
		WHERE m.id_moneda = $1 AND m.deleted_at IS NULL
	`, monedaSelectFields)

	m := &models.Moneda{}
	err := s.scanRowMoneda(s.db.QueryRowContext(ctx, query, id), m)

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

	err := ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		query := `
			INSERT INTO moneda (nombre, codigo_iso, simbolo, id_sucursal, id_status)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id_moneda, created_at, updated_at
		`
		return tx.QueryRowContext(ctx, query,
			moneda.Nombre, moneda.CodigoISO, moneda.Simbolo, moneda.IDSucursal, moneda.IDStatus,
		).Scan(&moneda.IDMoneda, &moneda.CreatedAt, &moneda.UpdatedAt)
	})

	if err != nil {
		return nil, fmt.Errorf("error al crear moneda: %w", err)
	}

	return moneda, nil
}

func (s *storeMoneda) UpdateMoneda(ctx context.Context, id uuid.UUID, moneda *models.Moneda) (*models.Moneda, error) {
	defer performance.Trace(ctx, "store", "UpdateMoneda", performance.DbThreshold, time.Now())

	err := ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		query := `
			UPDATE moneda
			SET nombre = $1, codigo_iso = $2, simbolo = $3, id_sucursal = $4, id_status = $5, updated_at = CURRENT_TIMESTAMP
			WHERE id_moneda = $6 AND deleted_at IS NULL
			RETURNING created_at, updated_at
		`
		return tx.QueryRowContext(ctx, query,
			moneda.Nombre, moneda.CodigoISO, moneda.Simbolo, moneda.IDSucursal, moneda.IDStatus, id,
		).Scan(&moneda.CreatedAt, &moneda.UpdatedAt)
	})

	if err != nil {
		return nil, fmt.Errorf("error al actualizar moneda: %w", err)
	}

	moneda.IDMoneda = id
	return moneda, nil
}

func (s *storeMoneda) DeleteMoneda(ctx context.Context, id uuid.UUID) error {
	defer performance.Trace(ctx, "store", "DeleteMoneda", performance.DbThreshold, time.Now())

	err := ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		query := `UPDATE moneda SET deleted_at = CURRENT_TIMESTAMP WHERE id_moneda = $1 AND deleted_at IS NULL`
		result, err := tx.ExecContext(ctx, query, id)
		if err != nil {
			return err
		}
		rows, _ := result.RowsAffected()
		if rows == 0 {
			return sql.ErrNoRows
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error al eliminar moneda: %w", err)
	}

	return nil
}
