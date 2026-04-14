package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/prunus/pkg/dto"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/utils/performance"
)

type StoreCompra interface {
	CreateOrden(ctx context.Context, o *models.OrdenCompra) (*models.OrdenCompra, error)
	GetOrdenByID(ctx context.Context, id uuid.UUID) (*models.OrdenCompra, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, statusID uuid.UUID, fechaRecepcion *time.Time) error
	UpdateDetalleRecepcion(ctx context.Context, idDetalle uuid.UUID, cantidad float64) error
	GetAllOrdenes(ctx context.Context, params dto.PaginationParams) ([]*models.OrdenCompra, error)
}

type storeCompra struct {
	db *sql.DB
}

// Campos base para SELECT de orden_compra con blindaje contra NULLs
const ordenCompraSelectFields = `
	id_orden_compra, 
	COALESCE(numero_orden, ''), 
	id_proveedor, 
	id_sucursal, 
	id_usuario, 
	id_moneda, 
	id_status, 
	COALESCE(subtotal, 0), 
	COALESCE(impuesto, 0), 
	COALESCE(total, 0), 
	COALESCE(observaciones, ''), 
	fecha_emision, 
	fecha_recepcion, 
	fecha_vencimiento, 
	created_at, 
	updated_at
`

// scanRowOrden helper para escanear órdenes de compra
func (s *storeCompra) scanRowOrden(scanner interface{ Scan(dest ...any) error }, o *models.OrdenCompra) error {
	return scanner.Scan(
		&o.IDOrdenCompra, &o.NumeroOrden, &o.IDProveedor, &o.IDSucursal, &o.IDUsuario,
		&o.IDMoneda, &o.IDStatus, &o.Subtotal, &o.Impuesto, &o.Total, &o.Observaciones,
		&o.FechaEmision, &o.FechaRecepcion, &o.FechaVencimiento, &o.CreatedAt, &o.UpdatedAt,
	)
}

// Campos base para SELECT de detalle_orden_compra
const detalleCompraSelectFields = `
	id_detalle_compra, id_orden_compra, id_producto, cantidad_pedida, 
	cantidad_recibida, precio_unitario, impuesto, total
`

// scanRowDetalle helper para escanear detalles de orden de compra
func (s *storeCompra) scanRowDetalle(scanner interface{ Scan(dest ...any) error }, d *models.DetalleOrdenCompra) error {
	return scanner.Scan(
		&d.IDDetalleCompra, &d.IDOrdenCompra, &d.IDProducto, &d.CantidadPedida,
		&d.CantidadRecibida, &d.PrecioUnitario, &d.Impuesto, &d.Total,
	)
}

func NewCompra(db *sql.DB) StoreCompra {
	return &storeCompra{db: db}
}

func (s *storeCompra) CreateOrden(ctx context.Context, o *models.OrdenCompra) (*models.OrdenCompra, error) {
	defer performance.Trace(ctx, "store", "CreateOrden", performance.DbThreshold, time.Now())

	err := ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		queryCab := `INSERT INTO orden_compra (
			numero_orden, id_proveedor, id_sucursal, id_usuario, id_moneda, id_status, 
			subtotal, impuesto, total, observaciones, fecha_vencimiento
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) 
		RETURNING id_orden_compra, fecha_emision, created_at, updated_at`

		err := tx.QueryRowContext(ctx, queryCab,
			o.NumeroOrden, o.IDProveedor, o.IDSucursal, o.IDUsuario, o.IDMoneda, o.IDStatus,
			o.Subtotal, o.Impuesto, o.Total, o.Observaciones, o.FechaVencimiento,
		).Scan(&o.IDOrdenCompra, &o.FechaEmision, &o.CreatedAt, &o.UpdatedAt)

		if err != nil {
			return fmt.Errorf("error al insertar cabecera de compra: %w", err)
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
				return fmt.Errorf("error al insertar detalle de compra: %w", err)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return o, nil
}

func (s *storeCompra) GetOrdenByID(ctx context.Context, id uuid.UUID) (*models.OrdenCompra, error) {
	defer performance.Trace(ctx, "store", "GetOrdenByID", performance.DbThreshold, time.Now())
	
	query := fmt.Sprintf("SELECT %s FROM orden_compra WHERE id_orden_compra = $1 AND deleted_at IS NULL", ordenCompraSelectFields)
	o := &models.OrdenCompra{}
	err := s.scanRowOrden(s.db.QueryRowContext(ctx, query, id), o)
	if err != nil {
		return nil, err
	}

	queryDet := fmt.Sprintf("SELECT %s FROM detalle_orden_compra WHERE id_orden_compra = $1", detalleCompraSelectFields)
	rows, err := s.db.QueryContext(ctx, queryDet, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		d := &models.DetalleOrdenCompra{}
		if err := s.scanRowDetalle(rows, d); err != nil {
			return nil, err
		}
		o.Detalles = append(o.Detalles, d)
	}

	return o, nil
}

func (s *storeCompra) UpdateStatus(ctx context.Context, id uuid.UUID, statusID uuid.UUID, fechaRecepcion *time.Time) error {
	defer performance.Trace(ctx, "store", "UpdateStatus", performance.DbThreshold, time.Now())
	
	query := `UPDATE orden_compra SET id_status = $1, fecha_recepcion = $2, updated_at = CURRENT_TIMESTAMP WHERE id_orden_compra = $3 AND deleted_at IS NULL`
	_, err := s.db.ExecContext(ctx, query, statusID, fechaRecepcion, id)
	return err
}

func (s *storeCompra) UpdateDetalleRecepcion(ctx context.Context, idDetalle uuid.UUID, cantidad float64) error {
	defer performance.Trace(ctx, "store", "UpdateDetalleRecepcion", performance.DbThreshold, time.Now())
	
	query := `UPDATE detalle_orden_compra SET cantidad_recibida = $1 WHERE id_detalle_compra = $2`
	_, err := s.db.ExecContext(ctx, query, cantidad, idDetalle)
	return err
}

func (s *storeCompra) GetAllOrdenes(ctx context.Context, params dto.PaginationParams) ([]*models.OrdenCompra, error) {
	defer performance.Trace(ctx, "store", "GetAllOrdenes", performance.DbThreshold, time.Now())
	
	if params.Limit <= 0 {
		params.Limit = dto.DefaultLimit
	}

	query := fmt.Sprintf(`
		SELECT %s 
		FROM orden_compra 
		WHERE deleted_at IS NULL 
	`, ordenCompraSelectFields)

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

	var ordenes []*models.OrdenCompra
	for rows.Next() {
		o := &models.OrdenCompra{}
		if err := s.scanRowOrden(rows, o); err != nil {
			return nil, err
		}
		ordenes = append(ordenes, o)
	}
	return ordenes, nil
}
