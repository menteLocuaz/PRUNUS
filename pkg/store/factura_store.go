package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
)

type StoreFactura interface {
	CreateFactura(ctx context.Context, f *models.Factura, items []*models.DetalleFactura) (*models.Factura, error)
	GetFacturaByID(ctx context.Context, id uuid.UUID) (*models.Factura, []*models.DetalleFactura, error)
	GetAllFacturas(ctx context.Context) ([]*models.Factura, error)

	// Impuestos
	GetAllImpuestos(ctx context.Context) ([]*models.Impuesto, error)

	// Formas de Pago
	GetAllFormasPago(ctx context.Context) ([]*models.FormaPago, error)
}

type storeFactura struct {
	db *sql.DB
}

func NewFactura(db *sql.DB) StoreFactura {
	return &storeFactura{db: db}
}

func (s *storeFactura) CreateFactura(ctx context.Context, f *models.Factura, items []*models.DetalleFactura) (*models.Factura, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	metadataJSON, _ := json.Marshal(f.Metadata)

	queryFac := `INSERT INTO factura (
		fac_numero, cfac_subtotal, cfac_iva, cfac_total, cfac_observacion, 
		id_user_pos, id_estacion, id_orden_pedido, id_cliente, id_periodo, 
		id_control_estacion, id_status, fecha_operacion, base_impuesto, 
		impuesto, valor_impuesto, metadata
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17) 
	RETURNING id_factura, cfac_fecha_creacion, created_at, updated_at`

	err = tx.QueryRowContext(ctx, queryFac,
		f.FacNumero, f.CfacSubtotal, f.CfacIVA, f.CfacTotal, f.CfacObservacion,
		f.IDUserPos, f.IDEstacion, f.IDOrdenPedido, f.IDCliente, f.IDPeriodo,
		f.IDControlEstacion, f.IDStatus, f.FechaOperacion, f.BaseImpuesto,
		f.Impuesto, f.ValorImpuesto, metadataJSON,
	).Scan(&f.IDFactura, &f.CfacFechaCreacion, &f.CreatedAt, &f.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("error al insertar factura: %w", err)
	}

	queryDet := `INSERT INTO detalle_factura (id_factura, id_producto, cantidad, precio, subtotal, impuesto, total) 
	             VALUES ($1, $2, $3, $4, $5, $6, $7)`

	for _, item := range items {
		_, err = tx.ExecContext(ctx, queryDet, f.IDFactura, item.IDProducto, item.Cantidad, item.Precio, item.Subtotal, item.Impuesto, item.Total)
		if err != nil {
			return nil, fmt.Errorf("error al insertar detalle: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return f, nil
}

func (s *storeFactura) GetFacturaByID(ctx context.Context, id uuid.UUID) (*models.Factura, []*models.DetalleFactura, error) {
	queryFac := `SELECT id_factura, fac_numero, cfac_subtotal, cfac_iva, cfac_total, cfac_observacion, id_cliente, created_at FROM factura WHERE id_factura = $1`
	f := &models.Factura{}
	err := s.db.QueryRowContext(ctx, queryFac, id).Scan(&f.IDFactura, &f.FacNumero, &f.CfacSubtotal, &f.CfacIVA, &f.CfacTotal, &f.CfacObservacion, &f.IDCliente, &f.CreatedAt)
	if err != nil {
		return nil, nil, err
	}

	queryDet := `SELECT id_detalle_factura, id_producto, cantidad, precio, subtotal, impuesto, total FROM detalle_factura WHERE id_factura = $1`
	rows, err := s.db.QueryContext(ctx, queryDet, id)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var items []*models.DetalleFactura
	for rows.Next() {
		item := &models.DetalleFactura{}
		if err := rows.Scan(&item.IDDetalleFactura, &item.IDProducto, &item.Cantidad, &item.Precio, &item.Subtotal, &item.Impuesto, &item.Total); err != nil {
			return nil, nil, err
		}
		items = append(items, item)
	}

	return f, items, nil
}

func (s *storeFactura) GetAllFacturas(ctx context.Context) ([]*models.Factura, error) {
	query := `SELECT id_factura, fac_numero, cfac_total, id_cliente, created_at FROM factura WHERE deleted_at IS NULL ORDER BY created_at DESC`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var facturas []*models.Factura
	for rows.Next() {
		f := &models.Factura{}
		if err := rows.Scan(&f.IDFactura, &f.FacNumero, &f.CfacTotal, &f.IDCliente, &f.CreatedAt); err != nil {
			return nil, err
		}
		facturas = append(facturas, f)
	}
	return facturas, nil
}

func (s *storeFactura) GetAllImpuestos(ctx context.Context) ([]*models.Impuesto, error) {
	query := `SELECT id_impuesto, nombre, porcentaje, tipo FROM impuesto WHERE deleted_at IS NULL`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var impuestos []*models.Impuesto
	for rows.Next() {
		i := &models.Impuesto{}
		if err := rows.Scan(&i.IDImpuesto, &i.Nombre, &i.Porcentaje, &i.Tipo); err != nil {
			return nil, err
		}
		impuestos = append(impuestos, i)
	}
	return impuestos, nil
}

func (s *storeFactura) GetAllFormasPago(ctx context.Context) ([]*models.FormaPago, error) {
	query := `SELECT id_forma_pago, fmp_codigo, fmp_descripcion, id_status FROM forma_pago WHERE deleted_at IS NULL`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var formas []*models.FormaPago
	for rows.Next() {
		f := &models.FormaPago{}
		if err := rows.Scan(&f.IDFormaPago, &f.FmpCodigo, &f.FmpDescripcion, &f.IDStatus); err != nil {
			return nil, err
		}
		formas = append(formas, f)
	}
	return formas, nil
}
