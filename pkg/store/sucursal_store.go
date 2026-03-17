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

type StoreSucursal interface {
	GetAllSucursales(ctx context.Context) ([]*models.Sucursal, error)
	GetSucursalByID(ctx context.Context, id uuid.UUID) (*models.Sucursal, error)
	CreateSucursal(ctx context.Context, sucursal *models.Sucursal) (*models.Sucursal, error)
	UpdateSucursal(ctx context.Context, id uuid.UUID, sucursal *models.Sucursal) (*models.Sucursal, error)
	DeleteSucursal(ctx context.Context, id uuid.UUID) error
}

type storeSucursal struct {
	db *sql.DB
}

func NewSucursal(db *sql.DB) StoreSucursal {
	return &storeSucursal{db: db}
}

// OBTIENE TODAS LAS SUCURSALES
func (s *storeSucursal) GetAllSucursales(ctx context.Context) ([]*models.Sucursal, error) {
	defer performance.Trace(ctx, "store", "GetAllSucursales", performance.DbThreshold, time.Now())
	query := `
	SELECT 
		s.id_sucursal,
		s.id_empresa,
		s.nombre_sucursal,
		s.id_status,

		e.id_empresa,
		e.nombre,
		e.rut,
		e.id_status
	FROM sucursal s
	LEFT JOIN empresa e ON e.id_empresa = s.id_empresa
	WHERE s.deleted_at IS NULL
	  AND e.deleted_at IS NULL
	ORDER BY s.id_sucursal
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener sucursales: %w", err)
	}
	defer rows.Close()

	var sucursales []*models.Sucursal

	for rows.Next() {
		s := &models.Sucursal{
			Empresa: &models.Empresa{},
		}

		if err := rows.Scan(
			&s.IDSucursal,
			&s.IDEmpresa,
			&s.NombreSucursal,
			&s.IDStatus,

			&s.Empresa.IDEmpresa,
			&s.Empresa.Nombre,
			&s.Empresa.RUT,
			&s.Empresa.IDStatus,
		); err != nil {
			return nil, fmt.Errorf("error al escanear sucursal: %w", err)
		}

		sucursales = append(sucursales, s)
	}

	return sucursales, nil
}

// OBTIENE UNA SOLA SUCURSAL
func (s *storeSucursal) GetSucursalByID(ctx context.Context, id uuid.UUID) (*models.Sucursal, error) {
	defer performance.Trace(ctx, "store", "GetSucursalByID", performance.DbThreshold, time.Now())
	query := `
	SELECT 
		s.id_sucursal,
		s.id_empresa,
		s.nombre_sucursal,
		s.id_status,

		e.id_empresa,
		e.nombre,
		e.rut,
		e.id_status
	FROM sucursal s
	LEFT JOIN empresa e ON e.id_empresa = s.id_empresa
	WHERE s.id_sucursal = $1
	  AND s.deleted_at IS NULL
	  AND e.deleted_at IS NULL
	`
	sucursal := &models.Sucursal{
		Empresa: &models.Empresa{},
	}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&sucursal.IDSucursal,
		&sucursal.IDEmpresa,
		&sucursal.NombreSucursal,
		&sucursal.IDStatus,

		&sucursal.Empresa.IDEmpresa,
		&sucursal.Empresa.Nombre,
		&sucursal.Empresa.RUT,
		&sucursal.Empresa.IDStatus,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("sucursal con ID %s no encontrada", id)
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener sucursal: %w", err)
	}

	return sucursal, nil
}

// CREAR SUCURSAL
func (s *storeSucursal) CreateSucursal(ctx context.Context, sucursal *models.Sucursal) (*models.Sucursal, error) {
	defer performance.Trace(ctx, "store", "CreateSucursal", performance.DbThreshold, time.Now())
	query := `INSERT INTO sucursal (id_empresa, nombre_sucursal, id_status) VALUES ($1, $2, $3) RETURNING id_sucursal`

	var id uuid.UUID
	err := s.db.QueryRowContext(ctx, query, sucursal.IDEmpresa, sucursal.NombreSucursal, sucursal.IDStatus).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("error al crear sucursal: %w", err)
	}

	sucursal.IDSucursal = id
	return sucursal, nil
}

// ACTUALIZAR LA SUCURSAL
func (s *storeSucursal) UpdateSucursal(ctx context.Context, id uuid.UUID, sucursal *models.Sucursal) (*models.Sucursal, error) {
	defer performance.Trace(ctx, "store", "UpdateSucursal", performance.DbThreshold, time.Now())
	query := `
		UPDATE sucursal
		SET
			id_empresa = $1,
			nombre_sucursal = $2,
			id_status = $3,
			updated_at = CURRENT_TIMESTAMP
		WHERE id_sucursal = $4
		  AND deleted_at IS NULL
		RETURNING
			id_sucursal,
			id_empresa,
			nombre_sucursal,
			id_status,
			created_at,
			updated_at
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		sucursal.IDEmpresa,
		sucursal.NombreSucursal,
		sucursal.IDStatus,
		id,
	).Scan(
		&sucursal.IDSucursal,
		&sucursal.IDEmpresa,
		&sucursal.NombreSucursal,
		&sucursal.IDStatus,
		&sucursal.CreatedAt,
		&sucursal.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("sucursal con ID %s no encontrada", id)
	}
	if err != nil {
		return nil, fmt.Errorf("error al actualizar sucursal: %w", err)
	}

	return sucursal, nil
}

// ELIMINAR
func (s *storeSucursal) DeleteSucursal(ctx context.Context, id uuid.UUID) error {
	defer performance.Trace(ctx, "store", "DeleteSucursal", performance.DbThreshold, time.Now())
	query := `UPDATE sucursal
	          SET deleted_at = $1
	          WHERE id_sucursal = $2 AND deleted_at IS NULL`

	result, err := s.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("error al eliminar sucursal: %w", err)
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
