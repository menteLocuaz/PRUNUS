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

type StoreEmpresa interface {
	GetAllEmpresa(ctx context.Context) ([]*models.Empresa, error)
	GetByIdEmpresa(ctx context.Context, id uuid.UUID) (*models.Empresa, error)
	CreateEmpresa(ctx context.Context, empresa *models.Empresa) (*models.Empresa, error)
	UpdateEmpresa(ctx context.Context, id uuid.UUID, empresa *models.Empresa) (*models.Empresa, error)
	DeleteEmpresa(ctx context.Context, id uuid.UUID) error
}

type storeEmpresa struct {
	db *sql.DB
}

func NewEmpresa(db *sql.DB) StoreEmpresa {
	return &storeEmpresa{db: db}
}

// OBTIENE TODAS LAS EMMPRESA
func (s *storeEmpresa) GetAllEmpresa(ctx context.Context) ([]*models.Empresa, error) {
	defer performance.Trace(ctx, "store", "GetAllEmpresa", performance.DbThreshold, time.Now())
	query := `SELECT id_empresa, nombre, rut, id_status FROM empresa WHERE deleted_at IS NULL`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener empresas: %w", err)
	}
	defer rows.Close()

	var empresas []*models.Empresa

	for rows.Next() {
		e := &models.Empresa{}
		if err := rows.Scan(&e.IDEmpresa, &e.Nombre, &e.RUT, &e.IDStatus); err != nil {
			return nil, fmt.Errorf("error al escanear empresa: %w", err)
		}
		empresas = append(empresas, e)
	}

	return empresas, nil
}

// ONTIEN UNA SOLA EMPRESA
func (s *storeEmpresa) GetByIdEmpresa(ctx context.Context, id uuid.UUID) (*models.Empresa, error) {
	defer performance.Trace(ctx, "store", "GetByIdEmpresa", performance.DbThreshold, time.Now())
	query := `SELECT id_empresa, nombre, rut, id_status FROM empresa
	          WHERE id_empresa = $1 AND deleted_at IS NULL`

	e := &models.Empresa{}
	err := s.db.QueryRowContext(ctx, query, id).
		Scan(&e.IDEmpresa, &e.Nombre, &e.RUT, &e.IDStatus)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("empresa con ID %s no encontrada", id)
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener empresa: %w", err)
	}

	return e, nil
}

// CREAR EMMPRESA
func (s *storeEmpresa) CreateEmpresa(ctx context.Context, empresa *models.Empresa) (*models.Empresa, error) {
	defer performance.Trace(ctx, "store", "CreateEmpresa", performance.DbThreshold, time.Now())
	query := `INSERT INTO empresa (nombre, rut, id_status) VALUES ($1, $2, $3) RETURNING id_empresa`

	var id uuid.UUID
	err := s.db.QueryRowContext(ctx, query, empresa.Nombre, empresa.RUT, empresa.IDStatus).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("error al crear empresa: %w", err)
	}

	empresa.IDEmpresa = id
	return empresa, nil
}

// ACTULIZAR LA EMPRESA
func (s *storeEmpresa) UpdateEmpresa(ctx context.Context, id uuid.UUID, empresa *models.Empresa) (*models.Empresa, error) {
	defer performance.Trace(ctx, "store", "UpdateEmpresa", performance.DbThreshold, time.Now())
	query := `UPDATE empresa
	          SET nombre = $1, rut = $2, id_status = $3, updated_at = CURRENT_TIMESTAMP
	          WHERE id_empresa = $4 AND deleted_at IS NULL`

	_, err := s.db.ExecContext(ctx, query, empresa.Nombre, empresa.RUT, empresa.IDStatus, id)
	if err != nil {
		return nil, fmt.Errorf("error al actualizar empresa: %w", err)
	}

	empresa.IDEmpresa = id
	return empresa, nil
}

// ELIMMINAR
func (s *storeEmpresa) DeleteEmpresa(ctx context.Context, id uuid.UUID) error {
	defer performance.Trace(ctx, "store", "DeleteEmpresa", performance.DbThreshold, time.Now())
	query := `UPDATE empresa
	          SET deleted_at = $1
	          WHERE id_empresa = $2 AND deleted_at IS NULL`

	result, err := s.db.ExecContext(ctx, query, time.Now(), id)

	if err != nil {
		return fmt.Errorf("error al eliminar empresa: %w", err)
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
