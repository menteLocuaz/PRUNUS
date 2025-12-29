package store

import (
	"database/sql"
	"time"

	"github.com/prunus/pkg/models"
)

type StoreSucursal interface {
	GetAllSucursales() ([]*models.Sucursal, error)
	GetSucursalByID(id uint) (*models.Sucursal, error)
	CreateSucursal(sucursal *models.Sucursal) (*models.Sucursal, error)
	UpdateSucursal(id uint, sucursal *models.Sucursal) (*models.Sucursal, error)
	DeleteSucursal(id uint) error
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
		s.estado,

		e.id_empresa,
		e.nombre,
		e.rut,
		e.estado
	FROM sucursal s
	JOIN empresa e ON e.id_empresa = s.id_empresa
	WHERE s.deleted_at IS NULL
	  AND e.deleted_at IS NULL
	ORDER BY s.id_sucursal
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
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
			&s.Estado,

			&s.Empresa.IDEmpresa,
			&s.Empresa.Nombre,
			&s.Empresa.RUT,
			&s.Empresa.Estado,
		); err != nil {
			return nil, err
		}

		sucursales = append(sucursales, s)
	}

	return sucursales, nil
}

// OBTIENE UNA SOLA SUCURSAL
func (s *storeSucursal) GetSucursalByID(id uint) (*models.Sucursal, error) {
	query := `
	SELECT 
		s.id_sucursal,
		s.id_empresa,
		s.nombre_sucursal,
		s.estado,

		e.id_empresa,
		e.nombre,
		e.rut,
		e.estado
	FROM sucursal s
	JOIN empresa e ON e.id_empresa = s.id_empresa
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
		&sucursal.Estado,

		&sucursal.Empresa.IDEmpresa,
		&sucursal.Empresa.Nombre,
		&sucursal.Empresa.RUT,
		&sucursal.Empresa.Estado,
	)

	if err != nil {
		return nil, err
	}

	return sucursal, nil
}

// CREAR SUCURSAL
func (s *storeSucursal) CreateSucursal(sucursal *models.Sucursal) (*models.Sucursal, error) {
	query := `INSERT INTO sucursal (id_empresa, nombre_sucursal, estado) VALUES ($1, $2, $3) RETURNING id_sucursal`

	var id uint
	err := s.db.QueryRow(query, sucursal.IDEmpresa, sucursal.NombreSucursal, sucursal.Estado).Scan(&id)
	if err != nil {
		return nil, err
	}

	sucursal.IDSucursal = id
	return sucursal, nil
}

// ACTUALIZAR LA SUCURSAL

func (s *storeSucursal) UpdateSucursal(id uint, sucursal *models.Sucursal) (*models.Sucursal, error) {
	query := `UPDATE sucursal
	          SET id_empresa = $1, nombre_sucursal = $2, estado = $3
	          WHERE id_sucursal = $4 AND deleted_at IS NULL`

	_, err := s.db.Exec(query, sucursal.IDEmpresa, sucursal.NombreSucursal, sucursal.Estado, id)
	if err != nil {
		return nil, err
	}

	sucursal.IDSucursal = id
	return sucursal, nil
}

// ELIMINAR
func (s *storeSucursal) DeleteSucursal(id uint) error {
	query := `UPDATE sucursal
	          SET deleted_at = $1
	          WHERE id_sucursal = $2 AND deleted_at IS NULL`

	result, err := s.db.Exec(query, time.Now(), id)
	if err != nil {
		return err
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
