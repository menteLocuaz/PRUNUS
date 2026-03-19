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

type StoreEstacionPos interface {
	GetAll(ctx context.Context) ([]*models.EstacionPos, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.EstacionPos, error)
	Create(ctx context.Context, e *models.EstacionPos) (*models.EstacionPos, error)
	Update(ctx context.Context, id uuid.UUID, e *models.EstacionPos) (*models.EstacionPos, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type storeEstacionPos struct {
	db *sql.DB
}

func NewEstacionPosStore(db *sql.DB) StoreEstacionPos {
	return &storeEstacionPos{db: db}
}

func (s *storeEstacionPos) GetAll(ctx context.Context) ([]*models.EstacionPos, error) {
	defer performance.Trace(ctx, "store", "GetAllEstacionPos", performance.DbThreshold, time.Now())
	query := `SELECT id_estacion, codigo, nombre, ip, id_sucursal, id_status, created_at, updated_at FROM estaciones_pos WHERE deleted_at IS NULL`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener estaciones: %w", err)
	}
	defer rows.Close()

	var estaciones []*models.EstacionPos
	for rows.Next() {
		e := &models.EstacionPos{}
		if err := rows.Scan(&e.IDEstacion, &e.Codigo, &e.Nombre, &e.IP, &e.IDSucursal, &e.IDStatus, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error al escanear estación: %w", err)
		}
		estaciones = append(estaciones, e)
	}
	return estaciones, nil
}

func (s *storeEstacionPos) GetByID(ctx context.Context, id uuid.UUID) (*models.EstacionPos, error) {
	defer performance.Trace(ctx, "store", "GetEstacionPosByID", performance.DbThreshold, time.Now())
	query := `SELECT id_estacion, codigo, nombre, ip, id_sucursal, id_status, created_at, updated_at FROM estaciones_pos WHERE id_estacion = $1 AND deleted_at IS NULL`
	e := &models.EstacionPos{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(&e.IDEstacion, &e.Codigo, &e.Nombre, &e.IP, &e.IDSucursal, &e.IDStatus, &e.CreatedAt, &e.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("estación no encontrada")
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener estación: %w", err)
	}
	return e, nil
}

func (s *storeEstacionPos) Create(ctx context.Context, e *models.EstacionPos) (*models.EstacionPos, error) {
	defer performance.Trace(ctx, "store", "CreateEstacionPos", performance.DbThreshold, time.Now())
	query := `INSERT INTO estaciones_pos (codigo, nombre, ip, id_sucursal, id_status) VALUES ($1, $2, $3, $4, $5) RETURNING id_estacion, created_at, updated_at`
	err := s.db.QueryRowContext(ctx, query, e.Codigo, e.Nombre, e.IP, e.IDSucursal, e.IDStatus).Scan(&e.IDEstacion, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("error al crear estación: %w", err)
	}
	return e, nil
}

func (s *storeEstacionPos) Update(ctx context.Context, id uuid.UUID, e *models.EstacionPos) (*models.EstacionPos, error) {
	defer performance.Trace(ctx, "store", "UpdateEstacionPos", performance.DbThreshold, time.Now())
	query := `UPDATE estaciones_pos SET codigo = $1, nombre = $2, ip = $3, id_sucursal = $4, id_status = $5, updated_at = CURRENT_TIMESTAMP WHERE id_estacion = $6 AND deleted_at IS NULL RETURNING id_estacion, created_at, updated_at`
	err := s.db.QueryRowContext(ctx, query, e.Codigo, e.Nombre, e.IP, e.IDSucursal, e.IDStatus, id).Scan(&e.IDEstacion, &e.CreatedAt, &e.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("estación no encontrada")
	}
	if err != nil {
		return nil, fmt.Errorf("error al actualizar estación: %w", err)
	}
	return e, nil
}

func (s *storeEstacionPos) Delete(ctx context.Context, id uuid.UUID) error {
	defer performance.Trace(ctx, "store", "DeleteEstacionPos", performance.DbThreshold, time.Now())
	query := `UPDATE estaciones_pos SET deleted_at = CURRENT_TIMESTAMP WHERE id_estacion = $1 AND deleted_at IS NULL`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error al eliminar estación: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("estación no encontrada")
	}
	return nil
}
