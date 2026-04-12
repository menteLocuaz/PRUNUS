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

// StoreUnidad define las operaciones de persistencia para el catálogo de unidades.
type StoreUnidad interface {
	GetAllUnidades(ctx context.Context) ([]*models.Unidad, error)
	GetUnidadByID(ctx context.Context, id uuid.UUID) (*models.Unidad, error)
	GetUnidadesBySucursal(ctx context.Context, sucursalID uuid.UUID) ([]*models.Unidad, error)
	CreateUnidad(ctx context.Context, unidad *models.Unidad) (*models.Unidad, error)
	UpdateUnidad(ctx context.Context, id uuid.UUID, unidad *models.Unidad) (*models.Unidad, error)
	DeleteUnidad(ctx context.Context, id uuid.UUID) error
}

type storeUnidad struct {
	db *sql.DB
}

// NewUnidad crea una nueva instancia del store de unidades.
func NewUnidad(db *sql.DB) StoreUnidad {
	return &storeUnidad{db: db}
}

// Campos base para SELECT de unidad mapeados al esquema DB.
const unidadSelectFields = `
	u.id_unidad, u.nombre, u.abreviatura, u.id_status, u.id_sucursal, u.created_at, u.updated_at,
	s.id_sucursal, s.nombre_sucursal, s.id_status
`

// scanRowUnidad centraliza el escaneo de resultados para mantener consistencia.
func (s *storeUnidad) scanRowUnidad(scanner interface{ Scan(dest ...any) error }, u *models.Unidad) error {
	if u.Sucursal == nil {
		u.Sucursal = &models.Sucursal{}
	}

	return scanner.Scan(
		&u.IDUnidad, &u.Nombre, &u.Abreviatura, &u.IDStatus, &u.IDSucursal, &u.CreatedAt, &u.UpdatedAt,
		&u.Sucursal.IDSucursal, &u.Sucursal.NombreSucursal, &u.Sucursal.IDStatus,
	)
}

func (s *storeUnidad) GetAllUnidades(ctx context.Context) ([]*models.Unidad, error) {
	defer performance.Trace(ctx, "store", "GetAllUnidades", performance.DbThreshold, time.Now())
	query := fmt.Sprintf(`
		SELECT %s
		FROM unidad u
		LEFT JOIN sucursal s ON s.id_sucursal = u.id_sucursal
		WHERE u.deleted_at IS NULL
		ORDER BY u.nombre ASC
	`, unidadSelectFields)

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener unidades: %w", err)
	}
	defer rows.Close()

	var unidades []*models.Unidad
	for rows.Next() {
		u := &models.Unidad{}
		if err := s.scanRowUnidad(rows, u); err != nil {
			return nil, fmt.Errorf("error al escanear unidad: %w", err)
		}
		unidades = append(unidades, u)
	}

	return unidades, nil
}

func (s *storeUnidad) GetUnidadByID(ctx context.Context, id uuid.UUID) (*models.Unidad, error) {
	defer performance.Trace(ctx, "store", "GetUnidadByID", performance.DbThreshold, time.Now())
	query := fmt.Sprintf(`
		SELECT %s
		FROM unidad u
		LEFT JOIN sucursal s ON s.id_sucursal = u.id_sucursal
		WHERE u.id_unidad = $1 AND u.deleted_at IS NULL
	`, unidadSelectFields)

	u := &models.Unidad{}
	err := s.scanRowUnidad(s.db.QueryRowContext(ctx, query, id), u)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("unidad con ID %s no encontrada", id)
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener unidad: %w", err)
	}

	return u, nil
}

func (s *storeUnidad) GetUnidadesBySucursal(ctx context.Context, sucursalID uuid.UUID) ([]*models.Unidad, error) {
	defer performance.Trace(ctx, "store", "GetUnidadesBySucursal", performance.DbThreshold, time.Now())
	query := fmt.Sprintf(`
		SELECT %s
		FROM unidad u
		LEFT JOIN sucursal s ON s.id_sucursal = u.id_sucursal
		WHERE u.id_sucursal = $1 AND u.deleted_at IS NULL
		ORDER BY u.nombre ASC
	`, unidadSelectFields)

	rows, err := s.db.QueryContext(ctx, query, sucursalID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener unidades por sucursal: %w", err)
	}
	defer rows.Close()

	var unidades []*models.Unidad
	for rows.Next() {
		u := &models.Unidad{}
		if err := s.scanRowUnidad(rows, u); err != nil {
			return nil, fmt.Errorf("error al escanear unidad: %w", err)
		}
		unidades = append(unidades, u)
	}

	return unidades, nil
}

func (s *storeUnidad) CreateUnidad(ctx context.Context, unidad *models.Unidad) (*models.Unidad, error) {
	defer performance.Trace(ctx, "store", "CreateUnidad", performance.DbThreshold, time.Now())

	err := ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		query := `
			INSERT INTO unidad (nombre, abreviatura, id_status, id_sucursal)
			VALUES ($1, $2, $3, $4)
			RETURNING id_unidad, created_at, updated_at
		`
		return tx.QueryRowContext(ctx, query,
			unidad.Nombre, unidad.Abreviatura, unidad.IDStatus, unidad.IDSucursal,
		).Scan(&unidad.IDUnidad, &unidad.CreatedAt, &unidad.UpdatedAt)
	})

	if err != nil {
		return nil, fmt.Errorf("error al crear unidad: %w", err)
	}

	return unidad, nil
}

func (s *storeUnidad) UpdateUnidad(ctx context.Context, id uuid.UUID, unidad *models.Unidad) (*models.Unidad, error) {
	defer performance.Trace(ctx, "store", "UpdateUnidad", performance.DbThreshold, time.Now())

	err := ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		query := `
			UPDATE unidad
			SET nombre = $1, abreviatura = $2, id_status = $3, id_sucursal = $4, updated_at = CURRENT_TIMESTAMP
			WHERE id_unidad = $5 AND deleted_at IS NULL
			RETURNING created_at, updated_at
		`
		return tx.QueryRowContext(ctx, query,
			unidad.Nombre, unidad.Abreviatura, unidad.IDStatus, unidad.IDSucursal, id,
		).Scan(&unidad.CreatedAt, &unidad.UpdatedAt)
	})

	if err != nil {
		return nil, fmt.Errorf("error al actualizar unidad: %w", err)
	}

	unidad.IDUnidad = id
	return unidad, nil
}

func (s *storeUnidad) DeleteUnidad(ctx context.Context, id uuid.UUID) error {
	defer performance.Trace(ctx, "store", "DeleteUnidad", performance.DbThreshold, time.Now())

	err := ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		query := `UPDATE unidad SET deleted_at = CURRENT_TIMESTAMP WHERE id_unidad = $1 AND deleted_at IS NULL`
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
		return fmt.Errorf("error al eliminar unidad: %w", err)
	}

	return nil
}
