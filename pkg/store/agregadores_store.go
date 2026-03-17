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

type StoreAgregadores interface {
	GetAllAgregadores(ctx context.Context) ([]*models.Agregador, error)
	GetAgregadorByID(ctx context.Context, id uuid.UUID) (*models.Agregador, error)
	CreateAgregador(ctx context.Context, a *models.Agregador) (*models.Agregador, error)
	UpdateAgregador(ctx context.Context, id uuid.UUID, a *models.Agregador) (*models.Agregador, error)
	DeleteAgregador(ctx context.Context, id uuid.UUID) error

	CreateOrdenAgregador(ctx context.Context, o *models.OrdenAgregador) (*models.OrdenAgregador, error)
}

type storeAgregadores struct {
	db *sql.DB
}

func NewAgregadores(db *sql.DB) StoreAgregadores {
	return &storeAgregadores{db: db}
}

func (s *storeAgregadores) GetAllAgregadores(ctx context.Context) ([]*models.Agregador, error) {
	defer performance.Trace(ctx, "store", "GetAllAgregadores", performance.DbThreshold, time.Now())
	query := `SELECT id_agregador, nombre, descripcion, created_at, updated_at FROM agregadores WHERE deleted_at IS NULL`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agregadores []*models.Agregador
	for rows.Next() {
		a := &models.Agregador{}
		if err := rows.Scan(&a.IDAgregador, &a.Nombre, &a.Descripcion, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		agregadores = append(agregadores, a)
	}
	return agregadores, nil
}

func (s *storeAgregadores) GetAgregadorByID(ctx context.Context, id uuid.UUID) (*models.Agregador, error) {
	defer performance.Trace(ctx, "store", "GetAgregadorByID", performance.DbThreshold, time.Now())
	query := `SELECT id_agregador, nombre, descripcion, created_at, updated_at FROM agregadores WHERE id_agregador = $1 AND deleted_at IS NULL`
	a := &models.Agregador{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(&a.IDAgregador, &a.Nombre, &a.Descripcion, &a.CreatedAt, &a.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("agregador no encontrado")
	}
	return a, err
}

func (s *storeAgregadores) CreateAgregador(ctx context.Context, a *models.Agregador) (*models.Agregador, error) {
	query := `INSERT INTO agregadores (nombre, descripcion) VALUES ($1, $2) RETURNING id_agregador, created_at, updated_at`
	err := s.db.QueryRowContext(ctx, query, a.Nombre, a.Descripcion).Scan(&a.IDAgregador, &a.CreatedAt, &a.UpdatedAt)
	return a, err
}

func (s *storeAgregadores) UpdateAgregador(ctx context.Context, id uuid.UUID, a *models.Agregador) (*models.Agregador, error) {
	query := `UPDATE agregadores SET nombre = $1, descripcion = $2, updated_at = CURRENT_TIMESTAMP WHERE id_agregador = $3 AND deleted_at IS NULL RETURNING id_agregador, created_at, updated_at`
	err := s.db.QueryRowContext(ctx, query, a.Nombre, a.Descripcion, id).Scan(&a.IDAgregador, &a.CreatedAt, &a.UpdatedAt)
	return a, err
}

func (s *storeAgregadores) DeleteAgregador(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE agregadores SET deleted_at = CURRENT_TIMESTAMP WHERE id_agregador = $1 AND deleted_at IS NULL`
	_, err := s.db.ExecContext(ctx, query, id)
	return err
}

func (s *storeAgregadores) CreateOrdenAgregador(ctx context.Context, o *models.OrdenAgregador) (*models.OrdenAgregador, error) {
	datosJSON, err := json.Marshal(o.DatosAgregador)
	if err != nil {
		return nil, err
	}

	query := `INSERT INTO orden_agregador (id_orden_pedido, id_agregador, codigo_externo, datos_agregador, fecha) 
	          VALUES ($1, $2, $3, $4, $5) RETURNING id_orden_agregador, created_at, updated_at`
	err = s.db.QueryRowContext(ctx, query, o.IDOrdenPedido, o.IDAgregador, o.CodigoExterno, datosJSON, o.Fecha).Scan(&o.IDOrdenAgregador, &o.CreatedAt, &o.UpdatedAt)
	return o, err
}
