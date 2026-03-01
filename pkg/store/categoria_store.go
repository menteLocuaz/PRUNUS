package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/prunus/pkg/models"
)

type StoreCategoria interface {
	GetAllCategorias() ([]*models.Categoria, error)
	GetCategoriaByID(id uint) (*models.Categoria, error)
	CreateCategoria(categoria *models.Categoria) (*models.Categoria, error)
	UpdateCategoria(id uint, categoria *models.Categoria) (*models.Categoria, error)
	DeleteCategoria(id uint) error
}

type storeCategoria struct {
	db *sql.DB
}

func NewCategoria(db *sql.DB) StoreCategoria {
	return &storeCategoria{db: db}
}

func (s *storeCategoria) GetAllCategorias() ([]*models.Categoria, error) {
	query := `
	SELECT
		c.id_categoria,
		c.nombre,
		c.id_sucursal,
		c.created_at,
		c.updated_at,

		su.id_sucursal,
		su.nombre_sucursal,
		su.estado
	FROM categoria c
	LEFT JOIN sucursal su ON su.id_sucursal = c.id_sucursal
	WHERE c.deleted_at IS NULL
	  AND su.deleted_at IS NULL
	ORDER BY c.id_categoria
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener categorias: %w", err)
	}
	defer rows.Close()

	var categorias []*models.Categoria

	for rows.Next() {
		c := &models.Categoria{
			Sucursal: &models.Sucursal{},
		}

		if err := rows.Scan(
			&c.IDCategoria,
			&c.Nombre,
			&c.IDSucursal,
			&c.CreatedAt,
			&c.UpdatedAt,

			&c.Sucursal.IDSucursal,
			&c.Sucursal.NombreSucursal,
			&c.Sucursal.Estado,
		); err != nil {
			return nil, fmt.Errorf("error al escanear categoria: %w", err)
		}

		categorias = append(categorias, c)
	}

	return categorias, nil
}

func (s *storeCategoria) GetCategoriaByID(id uint) (*models.Categoria, error) {
	query := `
	SELECT
		c.id_categoria,
		c.nombre,
		c.id_sucursal,
		c.created_at,
		c.updated_at,

		su.id_sucursal,
		su.nombre_sucursal,
		su.estado
	FROM categoria c
	LEFT JOIN sucursal su ON su.id_sucursal = c.id_sucursal
	WHERE c.id_categoria = $1
	  AND c.deleted_at IS NULL
	  AND su.deleted_at IS NULL
	`

	c := &models.Categoria{
		Sucursal: &models.Sucursal{},
	}

	err := s.db.QueryRow(query, id).Scan(
		&c.IDCategoria,
		&c.Nombre,
		&c.IDSucursal,
		&c.CreatedAt,
		&c.UpdatedAt,

		&c.Sucursal.IDSucursal,
		&c.Sucursal.NombreSucursal,
		&c.Sucursal.Estado,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("categoria con ID %d no encontrada", id)
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener categoria: %w", err)
	}

	return c, nil
}

func (s *storeCategoria) CreateCategoria(categoria *models.Categoria) (*models.Categoria, error) {
	query := `INSERT INTO categoria (nombre, id_sucursal) VALUES ($1, $2) RETURNING id_categoria`

	var id uint
	err := s.db.QueryRow(query, categoria.Nombre, categoria.IDSucursal).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("error al crear categoria: %w", err)
	}

	categoria.IDCategoria = id
	return categoria, nil
}

func (s *storeCategoria) UpdateCategoria(id uint, categoria *models.Categoria) (*models.Categoria, error) {
	query := `
		UPDATE categoria
		SET
			nombre = $1,
			id_sucursal = $2,
			updated_at = CURRENT_TIMESTAMP
		WHERE id_categoria = $3
		  AND deleted_at IS NULL
		RETURNING
			id_categoria,
			nombre,
			id_sucursal,
			created_at,
			updated_at
	`

	err := s.db.QueryRow(query, categoria.Nombre, categoria.IDSucursal, id).Scan(
		&categoria.IDCategoria,
		&categoria.Nombre,
		&categoria.IDSucursal,
		&categoria.CreatedAt,
		&categoria.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("categoria con ID %d no encontrada", id)
	}
	if err != nil {
		return nil, fmt.Errorf("error al actualizar categoria: %w", err)
	}

	return categoria, nil
}

func (s *storeCategoria) DeleteCategoria(id uint) error {
	query := `UPDATE categoria SET deleted_at = $1 WHERE id_categoria = $2 AND deleted_at IS NULL`

	result, err := s.db.Exec(query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("error al eliminar categoria: %w", err)
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
