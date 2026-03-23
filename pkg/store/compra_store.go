package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
)

type StoreCompra interface {
	CreateOrden(ctx context.Context, o *models.OrdenCompra) (*models.OrdenCompra, error)
	GetOrdenByID(ctx context.Context, id uuid.UUID) (*models.OrdenCompra, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, statusID uuid.UUID, fechaRecepcion *time.Time) error
	UpdateDetalleRecepcion(ctx context.Context, idDetalle uuid.UUID, cantidad float64) error
	GetAllOrdenes(ctx context.Context) ([]*models.OrdenCompra, error)
}

type storeCompra struct {
	db *sql.DB
}

func NewCompra(db *sql.DB) StoreCompra {
	return &storeCompra{db: db}
}

func (s *storeCompra) CreateOrden(ctx context.Context, o *models.OrdenCompra) (*models.OrdenCompra, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	queryCab := `INSERT INTO orden_compra (
		numero_orden, id_proveedor, id_sucursal, id_usuario, id_moneda, id_status, 
		subtotal, impuesto, total, observaciones
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
	RETURNING id_orden_compra, fecha_emision, created_at, updated_at`

	err = tx.QueryRowContext(ctx, queryCab,
		o.NumeroOrden, o.IDProveedor, o.IDSucursal, o.IDUsuario, o.IDMoneda, o.IDStatus,
		o.Subtotal, o.Impuesto, o.Total, o.Observaciones,
	).Scan(&o.IDOrdenCompra, &o.FechaEmision, &o.CreatedAt, &o.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("error al insertar cabecera de compra: %w", err)
	}

	queryDet := `INSERT INTO detalle_orden_compra (
		id_orden_compra, id_producto, cantidad_pedida, cantidad_recibida, 
		precio_unitario, impuesto, total
	) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id_detalle_compra`

	for _, d := range o.Detalles {
		err = tx.QueryRowContext(ctx, queryDet,
			o.IDOrdenCompra, d.IDProducto, d.CantidadPedida, 0,
			d.PrecioUnitario, d.Impuesto, d.Total,
		).Scan(&d.IDDetalleCompra)
		if err != nil {
			return nil, fmt.Errorf("error al insertar detalle de compra: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return o, nil
}

func (s *storeCompra) GetOrdenByID(ctx context.Context, id uuid.UUID) (*models.OrdenCompra, error) {
	query := `SELECT id_orden_compra, numero_orden, id_proveedor, id_sucursal, id_usuario, id_status, total FROM orden_compra WHERE id_orden_compra = $1`
	o := &models.OrdenCompra{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&o.IDOrdenCompra, &o.NumeroOrden, &o.IDProveedor, &o.IDSucursal, &o.IDUsuario, &o.IDStatus, &o.Total,
	)
	if err != nil {
		return nil, err
	}

	queryDet := `SELECT id_detalle_compra, id_producto, cantidad_pedida, cantidad_recibida, precio_unitario FROM detalle_orden_compra WHERE id_orden_compra = $1`
	rows, err := s.db.QueryContext(ctx, queryDet, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		d := &models.DetalleOrdenCompra{}
		if err := rows.Scan(&d.IDDetalleCompra, &d.IDProducto, &d.CantidadPedida, &d.CantidadRecibida, &d.PrecioUnitario); err != nil {
			return nil, err
		}
		o.Detalles = append(o.Detalles, d)
	}

	return o, nil
}

func (s *storeCompra) UpdateStatus(ctx context.Context, id uuid.UUID, statusID uuid.UUID, fechaRecepcion *time.Time) error {
	query := `UPDATE orden_compra SET id_status = $1, fecha_recepcion = $2, updated_at = CURRENT_TIMESTAMP WHERE id_orden_compra = $3`
	_, err := s.db.ExecContext(ctx, query, statusID, fechaRecepcion, id)
	return err
}

func (s *storeCompra) UpdateDetalleRecepcion(ctx context.Context, idDetalle uuid.UUID, cantidad float64) error {
	query := `UPDATE detalle_orden_compra SET cantidad_recibida = $1 WHERE id_detalle_compra = $2`
	_, err := s.db.ExecContext(ctx, query, cantidad, idDetalle)
	return err
}

func (s *storeCompra) GetAllOrdenes(ctx context.Context) ([]*models.OrdenCompra, error) {
	query := `SELECT id_orden_compra, numero_orden, id_proveedor, total, id_status, created_at FROM orden_compra WHERE deleted_at IS NULL ORDER BY created_at DESC`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ordenes []*models.OrdenCompra
	for rows.Next() {
		o := &models.OrdenCompra{}
		if err := rows.Scan(&o.IDOrdenCompra, &o.NumeroOrden, &o.IDProveedor, &o.Total, &o.IDStatus, &o.CreatedAt); err != nil {
			return nil, err
		}
		ordenes = append(ordenes, o)
	}
	return ordenes, nil
}
