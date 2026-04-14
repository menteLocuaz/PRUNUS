package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/prunus/pkg/dto"
	"github.com/prunus/pkg/models"
)

type StoreOrdenPedido interface {
	CreateOrden(ctx context.Context, o *models.OrdenPedido) (*models.OrdenPedido, error)
	GetOrdenByID(ctx context.Context, id uuid.UUID) (*models.OrdenPedido, error)
	GetAllOrdenes(ctx context.Context, params dto.PaginationParams) ([]*models.OrdenPedido, error)
	UpdateOrdenStatus(ctx context.Context, id uuid.UUID, statusID uuid.UUID) error
}

type storeOrdenPedido struct {
	db *sql.DB
}

func NewOrdenPedido(db *sql.DB) StoreOrdenPedido {
	return &storeOrdenPedido{db: db}
}

func (s *storeOrdenPedido) CreateOrden(ctx context.Context, o *models.OrdenPedido) (*models.OrdenPedido, error) {
	query := `INSERT INTO orden_pedido (odp_observacion, id_user_pos, id_periodo, id_estacion, id_status, direccion, canal, odp_total) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id_orden_pedido, odp_fecha_creacion, created_at, updated_at`
	err := s.db.QueryRowContext(ctx, query, o.OdpObservacion, o.IDUserPos, o.IDPeriodo, o.IDEstacion, o.IDStatus, o.Direccion, o.Canal, o.OdpTotal).
		Scan(&o.IDOrdenPedido, &o.OdpFechaCreacion, &o.CreatedAt, &o.UpdatedAt)
	return o, err
}

func (s *storeOrdenPedido) GetOrdenByID(ctx context.Context, id uuid.UUID) (*models.OrdenPedido, error) {
	query := `SELECT id_orden_pedido, odp_fecha_creacion, odp_observacion, id_user_pos, id_periodo, id_estacion, id_status, direccion, canal, odp_total FROM orden_pedido WHERE id_orden_pedido = $1 AND deleted_at IS NULL`
	o := &models.OrdenPedido{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(&o.IDOrdenPedido, &o.OdpFechaCreacion, &o.OdpObservacion, &o.IDUserPos, &o.IDPeriodo, &o.IDEstacion, &o.IDStatus, &o.Direccion, &o.Canal, &o.OdpTotal)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("orden no encontrada")
	}
	return o, err
}

func (s *storeOrdenPedido) GetAllOrdenes(ctx context.Context, params dto.PaginationParams) ([]*models.OrdenPedido, error) {
	if params.Limit <= 0 {
		params.Limit = dto.DefaultLimit
	}

	query := `SELECT id_orden_pedido, odp_fecha_creacion, odp_observacion, id_status, canal, odp_total 
	          FROM orden_pedido 
	          WHERE deleted_at IS NULL`

	var args []interface{}
	if params.LastDate != nil {
		query += " AND created_at < $1"
		args = append(args, params.LastDate)
	}

	query += " ORDER BY created_at DESC LIMIT $" + fmt.Sprint(len(args)+1)
	args = append(args, params.Limit)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ordenes []*models.OrdenPedido
	for rows.Next() {
		o := &models.OrdenPedido{}
		if err := rows.Scan(&o.IDOrdenPedido, &o.OdpFechaCreacion, &o.OdpObservacion, &o.IDStatus, &o.Canal, &o.OdpTotal); err != nil {
			return nil, err
		}
		ordenes = append(ordenes, o)
	}
	return ordenes, nil
}

func (s *storeOrdenPedido) UpdateOrdenStatus(ctx context.Context, id uuid.UUID, statusID uuid.UUID) error {
	query := `UPDATE orden_pedido SET id_status = $1, updated_at = CURRENT_TIMESTAMP WHERE id_orden_pedido = $2`
	_, err := s.db.ExecContext(ctx, query, statusID, id)
	return err
}
