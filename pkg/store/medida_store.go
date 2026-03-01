package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/prunus/pkg/models"
)

type StoreUnidad interface {
	GetAllUnidades() ([]*models.Unidad, error)
	GetUnidadByID(id uint) (*models.Unidad, error)
	CreateUnidad(unidad *models.Unidad) (*models.Unidad, error)
	UpdateUnidad(id uint, unidad *models.Unidad) (*models.Unidad, error)
	DeleteUnidad(id uint) error
}

type storeUnidad struct {
	db *sql.DB
}

func NewUnidad(db *sql.DB) StoreUnidad {
	return &storeUnidad{db: db}
}

func (s *storeUnidad) GetAllUnidades() ([]*models.Unidad, error) {
	query := `
	SELECT
		u.id_unidad,
		u.nombre,
		u.id_sucursal,
		u.created_at,
		u.updated_at,

		su.id_sucursal,
		su.nombre_sucursal,
		su.estado
	FROM unidad u
	LEFT JOIN sucursal su ON su.id_sucursal = u.id_sucursal
	WHERE u.deleted_at IS NULL
	  AND su.deleted_at IS NULL
	ORDER BY u.id_unidad
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener unidades: %w", err)
	}
	defer rows.Close()

	var unidades []*models.Unidad

	for rows.Next() {
		u := &models.Unidad{
			Sucursal: &models.Sucursal{},
		}

		if err := rows.Scan(
			&u.IDUnidad,
			&u.Nombre,
			&u.IDSucursal,
			&u.CreatedAt,
			&u.UpdatedAt,

			&u.Sucursal.IDSucursal,
			&u.Sucursal.NombreSucursal,
			&u.Sucursal.Estado,
		); err != nil {
			return nil, fmt.Errorf("error al escanear unidad: %w", err)
		}

		unidades = append(unidades, u)
	}

	return unidades, nil
}

func (s *storeUnidad) GetUnidadByID(id uint) (*models.Unidad, error) {
	query := `
	SELECT
		u.id_unidad,
		u.nombre,
		u.id_sucursal,
		u.created_at,
		u.updated_at,

		su.id_sucursal,
		su.nombre_sucursal,
		su.estado
	FROM unidad u
	LEFT JOIN sucursal su ON su.id_sucursal = u.id_sucursal
	WHERE u.id_unidad = $1
	  AND u.deleted_at IS NULL
	  AND su.deleted_at IS NULL
	`

	u := &models.Unidad{
		Sucursal: &models.Sucursal{},
	}

	err := s.db.QueryRow(query, id).Scan(
		&u.IDUnidad,
		&u.Nombre,
		&u.IDSucursal,
		&u.CreatedAt,
		&u.UpdatedAt,

		&u.Sucursal.IDSucursal,
		&u.Sucursal.NombreSucursal,
		&u.Sucursal.Estado,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("unidad con ID %d no encontrada", id)
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener unidad: %w", err)
	}

	return u, nil
}

func (s *storeUnidad) CreateUnidad(unidad *models.Unidad) (*models.Unidad, error) {
	query := `INSERT INTO unidad (nombre, id_sucursal) VALUES ($1, $2) RETURNING id_unidad`

	var id uint
	err := s.db.QueryRow(query, unidad.Nombre, unidad.IDSucursal).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("error al crear unidad: %w", err)
	}

	unidad.IDUnidad = id
	return unidad, nil
}

func (s *storeUnidad) UpdateUnidad(id uint, unidad *models.Unidad) (*models.Unidad, error) {
	query := `
		UPDATE unidad
		SET
			nombre = $1,
			id_sucursal = $2,
			updated_at = CURRENT_TIMESTAMP
		WHERE id_unidad = $3
		  AND deleted_at IS NULL
		RETURNING
			id_unidad,
			nombre,
			id_sucursal,
			created_at,
			updated_at
	`

	err := s.db.QueryRow(query, unidad.Nombre, unidad.IDSucursal, id).Scan(
		&unidad.IDUnidad,
		&unidad.Nombre,
		&unidad.IDSucursal,
		&unidad.CreatedAt,
		&unidad.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("unidad con ID %d no encontrada", id)
	}
	if err != nil {
		return nil, fmt.Errorf("error al actualizar unidad: %w", err)
	}

	return unidad, nil
}

func (s *storeUnidad) DeleteUnidad(id uint) error {
	query := `UPDATE unidad SET deleted_at = $1 WHERE id_unidad = $2 AND deleted_at IS NULL`

	result, err := s.db.Exec(query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("error al eliminar unidad: %w", err)
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
