package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/utils/performance"
)

type StoreDispositivoPos interface {
	GetAll(ctx context.Context) ([]*models.DispositivoPos, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.DispositivoPos, error)
	GetByEstacion(ctx context.Context, idEstacion uuid.UUID) ([]*models.DispositivoPos, error)
	Create(ctx context.Context, d *models.DispositivoPos) (*models.DispositivoPos, error)
	Update(ctx context.Context, id uuid.UUID, d *models.DispositivoPos) (*models.DispositivoPos, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type storeDispositivoPos struct {
	db *sql.DB
}

const dispositivoPosSelectFields = `
	id_dispositivo, id_estacion, nombre, tipo_dispositivo, configuracion, 
	id_status, created_at, updated_at
`

func (s *storeDispositivoPos) scanRow(scanner interface{ Scan(dest ...any) error }, d *models.DispositivoPos) error {
	var configJSON []byte
	err := scanner.Scan(
		&d.IDDispositivo, &d.IDEstacion, &d.Nombre, &d.TipoDispositivo, &configJSON,
		&d.IDStatus, &d.CreatedAt, &d.UpdatedAt,
	)
	if err != nil {
		return err
	}
	if len(configJSON) > 0 {
		json.Unmarshal(configJSON, &d.Configuracion)
	}
	return nil
}

func NewDispositivoPosStore(db *sql.DB) StoreDispositivoPos {
	return &storeDispositivoPos{db: db}
}

func (s *storeDispositivoPos) GetAll(ctx context.Context) ([]*models.DispositivoPos, error) {
	defer performance.Trace(ctx, "store", "GetAllDispositivoPos", performance.DbThreshold, time.Now())
	query := fmt.Sprintf(`SELECT %s FROM dispositivos_pos WHERE deleted_at IS NULL`, dispositivoPosSelectFields)
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener dispositivos: %w", err)
	}
	defer rows.Close()

	var dispositivos []*models.DispositivoPos
	for rows.Next() {
		d := &models.DispositivoPos{}
		if err := s.scanRow(rows, d); err != nil {
			return nil, fmt.Errorf("error al escanear dispositivo: %w", err)
		}
		dispositivos = append(dispositivos, d)
	}
	return dispositivos, nil
}

func (s *storeDispositivoPos) GetByID(ctx context.Context, id uuid.UUID) (*models.DispositivoPos, error) {
	defer performance.Trace(ctx, "store", "GetDispositivoPosByID", performance.DbThreshold, time.Now())
	query := fmt.Sprintf(`SELECT %s FROM dispositivos_pos WHERE id_dispositivo = $1 AND deleted_at IS NULL`, dispositivoPosSelectFields)
	d := &models.DispositivoPos{}
	err := s.scanRow(s.db.QueryRowContext(ctx, query, id), d)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("dispositivo no encontrado")
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener dispositivo: %w", err)
	}
	return d, nil
}

func (s *storeDispositivoPos) GetByEstacion(ctx context.Context, idEstacion uuid.UUID) ([]*models.DispositivoPos, error) {
	defer performance.Trace(ctx, "store", "GetDispositivoPosByEstacion", performance.DbThreshold, time.Now())
	query := fmt.Sprintf(`SELECT %s FROM dispositivos_pos WHERE id_estacion = $1 AND deleted_at IS NULL`, dispositivoPosSelectFields)
	rows, err := s.db.QueryContext(ctx, query, idEstacion)
	if err != nil {
		return nil, fmt.Errorf("error al obtener dispositivos por estación: %w", err)
	}
	defer rows.Close()

	var dispositivos []*models.DispositivoPos
	for rows.Next() {
		d := &models.DispositivoPos{}
		if err := s.scanRow(rows, d); err != nil {
			return nil, fmt.Errorf("error al escanear dispositivo: %w", err)
		}
		dispositivos = append(dispositivos, d)
	}
	return dispositivos, nil
}

func (s *storeDispositivoPos) Create(ctx context.Context, d *models.DispositivoPos) (*models.DispositivoPos, error) {
	defer performance.Trace(ctx, "store", "CreateDispositivoPos", performance.DbThreshold, time.Now())
	
	err := ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		configJSON, _ := json.Marshal(d.Configuracion)
		query := `
			INSERT INTO dispositivos_pos (id_estacion, nombre, tipo_dispositivo, configuracion, id_status) 
			VALUES ($1, $2, $3, $4, $5) 
			RETURNING id_dispositivo, created_at, updated_at`
		
		return tx.QueryRowContext(ctx, query, 
			d.IDEstacion, d.Nombre, d.TipoDispositivo, configJSON, d.IDStatus,
		).Scan(&d.IDDispositivo, &d.CreatedAt, &d.UpdatedAt)
	})

	if err != nil {
		return nil, fmt.Errorf("error al crear dispositivo: %w", err)
	}
	return d, nil
}

func (s *storeDispositivoPos) Update(ctx context.Context, id uuid.UUID, d *models.DispositivoPos) (*models.DispositivoPos, error) {
	defer performance.Trace(ctx, "store", "UpdateDispositivoPos", performance.DbThreshold, time.Now())
	
	err := ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		configJSON, _ := json.Marshal(d.Configuracion)
		query := `
			UPDATE dispositivos_pos 
			SET id_estacion = $1, nombre = $2, tipo_dispositivo = $3, configuracion = $4, id_status = $5, updated_at = CURRENT_TIMESTAMP 
			WHERE id_dispositivo = $6 AND deleted_at IS NULL 
			RETURNING created_at, updated_at`
		
		return tx.QueryRowContext(ctx, query, 
			d.IDEstacion, d.Nombre, d.TipoDispositivo, configJSON, d.IDStatus, id,
		).Scan(&d.CreatedAt, &d.UpdatedAt)
	})

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("dispositivo no encontrado")
	}
	if err != nil {
		return nil, fmt.Errorf("error al actualizar dispositivo: %w", err)
	}
	
	d.IDDispositivo = id
	return d, nil
}

func (s *storeDispositivoPos) Delete(ctx context.Context, id uuid.UUID) error {
	defer performance.Trace(ctx, "store", "DeleteDispositivoPos", performance.DbThreshold, time.Now())
	
	return ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		query := `UPDATE dispositivos_pos SET deleted_at = CURRENT_TIMESTAMP WHERE id_dispositivo = $1 AND deleted_at IS NULL`
		result, err := tx.ExecContext(ctx, query, id)
		if err != nil {
			return err
		}
		rows, _ := result.RowsAffected()
		if rows == 0 {
			return fmt.Errorf("dispositivo no encontrado")
		}
		return nil
	})
}
