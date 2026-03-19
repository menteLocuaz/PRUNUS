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

type StoreDispositivoPos interface {
	GetAll(ctx context.Context) ([]*models.DispositivoPos, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.DispositivoPos, error)
	Create(ctx context.Context, d *models.DispositivoPos) (*models.DispositivoPos, error)
	Update(ctx context.Context, id uuid.UUID, d *models.DispositivoPos) (*models.DispositivoPos, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type storeDispositivoPos struct {
	db *sql.DB
}

func NewDispositivoPosStore(db *sql.DB) StoreDispositivoPos {
	return &storeDispositivoPos{db: db}
}

func (s *storeDispositivoPos) GetAll(ctx context.Context) ([]*models.DispositivoPos, error) {
	defer performance.Trace(ctx, "store", "GetAllDispositivoPos", performance.DbThreshold, time.Now())
	query := `SELECT id_dispositivo, nombre, tipo, ip, id_estacion, created_at, updated_at FROM dispositivos_pos WHERE deleted_at IS NULL`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener dispositivos: %w", err)
	}
	defer rows.Close()

	var dispositivos []*models.DispositivoPos
	for rows.Next() {
		d := &models.DispositivoPos{}
		if err := rows.Scan(&d.IDDispositivo, &d.Nombre, &d.Tipo, &d.IP, &d.IDEstacion, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error al escanear dispositivo: %w", err)
		}
		dispositivos = append(dispositivos, d)
	}
	return dispositivos, nil
}

func (s *storeDispositivoPos) GetByID(ctx context.Context, id uuid.UUID) (*models.DispositivoPos, error) {
	defer performance.Trace(ctx, "store", "GetDispositivoPosByID", performance.DbThreshold, time.Now())
	query := `SELECT id_dispositivo, nombre, tipo, ip, id_estacion, created_at, updated_at FROM dispositivos_pos WHERE id_dispositivo = $1 AND deleted_at IS NULL`
	d := &models.DispositivoPos{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(&d.IDDispositivo, &d.Nombre, &d.Tipo, &d.IP, &d.IDEstacion, &d.CreatedAt, &d.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("dispositivo no encontrado")
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener dispositivo: %w", err)
	}
	return d, nil
}

func (s *storeDispositivoPos) Create(ctx context.Context, d *models.DispositivoPos) (*models.DispositivoPos, error) {
	defer performance.Trace(ctx, "store", "CreateDispositivoPos", performance.DbThreshold, time.Now())
	query := `INSERT INTO dispositivos_pos (nombre, tipo, ip, id_estacion) VALUES ($1, $2, $3, $4) RETURNING id_dispositivo, created_at, updated_at`
	err := s.db.QueryRowContext(ctx, query, d.Nombre, d.Tipo, d.IP, d.IDEstacion).Scan(&d.IDDispositivo, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("error al crear dispositivo: %w", err)
	}
	return d, nil
}

func (s *storeDispositivoPos) Update(ctx context.Context, id uuid.UUID, d *models.DispositivoPos) (*models.DispositivoPos, error) {
	defer performance.Trace(ctx, "store", "UpdateDispositivoPos", performance.DbThreshold, time.Now())
	query := `UPDATE dispositivos_pos SET nombre = $1, tipo = $2, ip = $3, id_estacion = $4, updated_at = CURRENT_TIMESTAMP WHERE id_dispositivo = $5 AND deleted_at IS NULL RETURNING id_dispositivo, created_at, updated_at`
	err := s.db.QueryRowContext(ctx, query, d.Nombre, d.Tipo, d.IP, d.IDEstacion, id).Scan(&d.IDDispositivo, &d.CreatedAt, &d.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("dispositivo no encontrado")
	}
	if err != nil {
		return nil, fmt.Errorf("error al actualizar dispositivo: %w", err)
	}
	return d, nil
}

func (s *storeDispositivoPos) Delete(ctx context.Context, id uuid.UUID) error {
	defer performance.Trace(ctx, "store", "DeleteDispositivoPos", performance.DbThreshold, time.Now())
	query := `UPDATE dispositivos_pos SET deleted_at = CURRENT_TIMESTAMP WHERE id_dispositivo = $1 AND deleted_at IS NULL`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error al eliminar dispositivo: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("dispositivo no encontrado")
	}
	return nil
}
