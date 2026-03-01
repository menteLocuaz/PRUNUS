package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/prunus/pkg/models"
)

type StoreMoneda interface {
	GetAllMonedas() ([]*models.Moneda, error)
	GetMonedaByID(id uint) (*models.Moneda, error)
	CreateMoneda(moneda *models.Moneda) (*models.Moneda, error)
	UpdateMoneda(id uint, moneda *models.Moneda) (*models.Moneda, error)
	DeleteMoneda(id uint) error
}

type storeMoneda struct {
	db *sql.DB
}

func NewMoneda(db *sql.DB) StoreMoneda {
	return &storeMoneda{db: db}
}

func (s *storeMoneda) GetAllMonedas() ([]*models.Moneda, error) {
	query := `
	SELECT
		m.id_moneda,
		m.nombre,
		m.id_sucursal,
		m.estado,
		m.created_at,
		m.updated_at,

		su.id_sucursal,
		su.nombre_sucursal,
		su.estado
	FROM moneda m
	LEFT JOIN sucursal su ON su.id_sucursal = m.id_sucursal
	WHERE m.deleted_at IS NULL
	  AND su.deleted_at IS NULL
	ORDER BY m.id_moneda
	`

	rows, err := s.db.Query(query)
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
			&m.Estado,
			&m.CreatedAt,
			&m.UpdatedAt,

			&m.Sucursal.IDSucursal,
			&m.Sucursal.NombreSucursal,
			&m.Sucursal.Estado,
		); err != nil {
			return nil, fmt.Errorf("error al escanear moneda: %w", err)
		}

		monedas = append(monedas, m)
	}

	return monedas, nil
}

func (s *storeMoneda) GetMonedaByID(id uint) (*models.Moneda, error) {
	query := `
	SELECT
		m.id_moneda,
		m.nombre,
		m.id_sucursal,
		m.estado,
		m.created_at,
		m.updated_at,

		su.id_sucursal,
		su.nombre_sucursal,
		su.estado
	FROM moneda m
	LEFT JOIN sucursal su ON su.id_sucursal = m.id_sucursal
	WHERE m.id_moneda = $1
	  AND m.deleted_at IS NULL
	  AND su.deleted_at IS NULL
	`

	m := &models.Moneda{
		Sucursal: &models.Sucursal{},
	}

	err := s.db.QueryRow(query, id).Scan(
		&m.IDMoneda,
		&m.Nombre,
		&m.IDSucursal,
		&m.Estado,
		&m.CreatedAt,
		&m.UpdatedAt,

		&m.Sucursal.IDSucursal,
		&m.Sucursal.NombreSucursal,
		&m.Sucursal.Estado,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("moneda con ID %d no encontrada", id)
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener moneda: %w", err)
	}

	return m, nil
}

func (s *storeMoneda) CreateMoneda(moneda *models.Moneda) (*models.Moneda, error) {
	query := `INSERT INTO moneda (nombre, id_sucursal, estado) VALUES ($1, $2, $3) RETURNING id_moneda`

	var id uint
	err := s.db.QueryRow(query, moneda.Nombre, moneda.IDSucursal, moneda.Estado).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("error al crear moneda: %w", err)
	}

	moneda.IDMoneda = id
	return moneda, nil
}

func (s *storeMoneda) UpdateMoneda(id uint, moneda *models.Moneda) (*models.Moneda, error) {
	query := `
		UPDATE moneda
		SET
			nombre = $1,
			id_sucursal = $2,
			estado = $3,
			updated_at = CURRENT_TIMESTAMP
		WHERE id_moneda = $4
		  AND deleted_at IS NULL
		RETURNING
			id_moneda,
			nombre,
			id_sucursal,
			estado,
			created_at,
			updated_at
	`

	err := s.db.QueryRow(query, moneda.Nombre, moneda.IDSucursal, moneda.Estado, id).Scan(
		&moneda.IDMoneda,
		&moneda.Nombre,
		&moneda.IDSucursal,
		&moneda.Estado,
		&moneda.CreatedAt,
		&moneda.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("moneda con ID %d no encontrada", id)
	}
	if err != nil {
		return nil, fmt.Errorf("error al actualizar moneda: %w", err)
	}

	return moneda, nil
}

func (s *storeMoneda) DeleteMoneda(id uint) error {
	query := `UPDATE moneda SET deleted_at = $1 WHERE id_moneda = $2 AND deleted_at IS NULL`

	result, err := s.db.Exec(query, time.Now(), id)
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
