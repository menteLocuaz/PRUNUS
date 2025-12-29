package store

import (
	"database/sql"
	"time"

	"github.com/prunus/pkg/models"
)

type StoreEmpresa interface {
	GetAllEmpresa() ([]*models.Empresa, error)
	GetByIdEmpresa(id uint) (*models.Empresa, error)
	CreateEmpresa(empresa *models.Empresa) (*models.Empresa, error)
	UpdateEmpresa(id uint, empresa *models.Empresa) (*models.Empresa, error)
	DeleteEmpresa(id uint) error
}

type store struct {
	db *sql.DB
}

func NewEmpresa(db *sql.DB) StoreEmpresa {
	return &store{db: db}
}

// OBTIENE TODAS LAS EMMPRESA
func (s *store) GetAllEmpresa() ([]*models.Empresa, error) {
	query := `SELECT id_empresa, nombre, rut, estado FROM empresa WHERE deleted_at IS NULL`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var empresas []*models.Empresa

	for rows.Next() {
		e := &models.Empresa{}
		if err := rows.Scan(&e.IDEmpresa, &e.Nombre, &e.RUT, &e.Estado); err != nil {
			return nil, err
		}
		empresas = append(empresas, e)
	}

	return empresas, nil
}

// ONTIEN UNA SOLA EMPRESA
func (s *store) GetByIdEmpresa(id uint) (*models.Empresa, error) {
	query := `SELECT id_empresa, nombre, rut, estado FROM empresa
	          WHERE id_empresa = $1 AND deleted_at IS NULL`

	e := &models.Empresa{}
	err := s.db.QueryRow(query, id).
		Scan(&e.IDEmpresa, &e.Nombre, &e.RUT, &e.Estado)

	if err != nil {
		return nil, err
	}

	return e, nil
}

// CREAR EMMPRESA
func (s *store) CreateEmpresa(empresa *models.Empresa) (*models.Empresa, error) {
	query := `INSERT INTO empresa (nombre, rut, estado) VALUES ($1, $2, $3) RETURNING id_empresa`

	var id uint
	err := s.db.QueryRow(query, empresa.Nombre, empresa.RUT, empresa.Estado).Scan(&id)
	if err != nil {
		return nil, err
	}

	empresa.IDEmpresa = id
	return empresa, nil
}

// ACTULIZAR LA EMPRESA

func (s *store) UpdateEmpresa(id uint, empresa *models.Empresa) (*models.Empresa, error) {
	query := `UPDATE empresa
	          SET nombre = $1, rut = $2, estado = $3
	          WHERE id_empresa = $4 AND deleted_at IS NULL`

	_, err := s.db.Exec(query, empresa.Nombre, empresa.RUT, empresa.Estado, id)
	if err != nil {
		return nil, err
	}

	empresa.IDEmpresa = id
	return empresa, nil
}

// ELIMMINAR
func (s *store) DeleteEmpresa(id uint) error {
	query := `UPDATE empresa
	          SET deleted_at = $1
	          WHERE id_empresa = $2 AND deleted_at IS NULL`

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
