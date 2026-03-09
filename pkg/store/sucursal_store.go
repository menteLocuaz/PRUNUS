package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
)

type StoreSucursal interface {
	GetAllSucursales() ([]*models.Sucursal, error)
	GetSucursalByID(id uuid.UUID) (*models.Sucursal, error)
	CreateSucursal(sucursal *models.Sucursal) (*models.Sucursal, error)
	UpdateSucursal(id uuid.UUID, sucursal *models.Sucursal) (*models.Sucursal, error)
	DeleteSucursal(id uuid.UUID) error
}

type storeSucursal struct {
	db *sql.DB
}

func NewSucursal(db *sql.DB) StoreSucursal {
	return &storeSucursal{db: db}
}

// OBTIENE TODAS LAS SUCURSALES
func (s *storeSucursal) GetAllSucursales() ([]*models.Sucursal, error) {
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

	rows, err := s.db.Query(query)
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
func (s *storeSucursal) GetSucursalByID(id uuid.UUID) (*models.Sucursal, error) {
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
	err := s.db.QueryRow(query, id).Scan(
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
func (s *storeSucursal) CreateSucursal(sucursal *models.Sucursal) (*models.Sucursal, error) {
	query := `INSERT INTO sucursal (id_empresa, nombre_sucursal, id_status) VALUES ($1, $2, $3) RETURNING id_sucursal`

	var id uuid.UUID
	err := s.db.QueryRow(query, sucursal.IDEmpresa, sucursal.NombreSucursal, sucursal.IDStatus).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("error al crear sucursal: %w", err)
	}

	sucursal.IDSucursal = id
	return sucursal, nil
}

// ACTUALIZAR LA SUCURSAL
func (s *storeSucursal) UpdateSucursal(id uuid.UUID, sucursal *models.Sucursal) (*models.Sucursal, error) {
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

	err := s.db.QueryRow(
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
func (s *storeSucursal) DeleteSucursal(id uuid.UUID) error {
	query := `UPDATE sucursal
	          SET deleted_at = $1
	          WHERE id_sucursal = $2 AND deleted_at IS NULL`

	result, err := s.db.Exec(query, time.Now(), id)
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
