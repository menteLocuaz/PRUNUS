package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/prunus/pkg/dto"
	"github.com/prunus/pkg/models"
)

type StoreFactura interface {
	CreateFactura(ctx context.Context, f *models.Factura, items []*models.DetalleFactura) (*models.Factura, error)
	RegistrarFacturaCompleta(ctx context.Context, req dto.FacturaCompletaRequest, idUsuario uuid.UUID) (*dto.FacturaResponse, error)
	GetFacturaByID(ctx context.Context, id uuid.UUID) (*models.Factura, []*models.DetalleFactura, error)
	GetAllFacturas(ctx context.Context, params dto.PaginationParams) ([]*models.Factura, error)

	// Impuestos
	GetAllImpuestos(ctx context.Context) ([]*models.Impuesto, error)

	// Formas de Pago
	GetAllFormasPago(ctx context.Context) ([]*models.FormaPago, error)
}

type storeFactura struct {
	db *sql.DB
}

// Campos base para SELECT de factura
const facturaSelectFields = `
	id_factura, fac_numero, cfac_subtotal, cfac_iva, cfac_total, cfac_observacion, 
	id_user_pos, id_estacion, id_orden_pedido, id_cliente, id_periodo, 
	id_control_estacion, id_status, fecha_operacion, base_impuesto, 
	impuesto, valor_impuesto, created_at, updated_at
`

// scanRowFactura helper para escanear facturas
func (s *storeFactura) scanRowFactura(scanner interface{ Scan(dest ...any) error }, f *models.Factura) error {
	return scanner.Scan(
		&f.IDFactura, &f.FacNumero, &f.CfacSubtotal, &f.CfacIVA, &f.CfacTotal, &f.CfacObservacion,
		&f.IDUserPos, &f.IDEstacion, &f.IDOrdenPedido, &f.IDCliente, &f.IDPeriodo,
		&f.IDControlEstacion, &f.IDStatus, &f.FechaOperacion, &f.BaseImpuesto,
		&f.Impuesto, &f.ValorImpuesto, &f.CreatedAt, &f.UpdatedAt,
	)
}

// Campos base para SELECT de detalle_factura
const detalleFacturaSelectFields = `
	id_detalle_factura, id_factura, id_producto, cantidad, precio, subtotal, impuesto, total
`

// scanRowDetalleFactura helper para escanear detalles de factura
func (s *storeFactura) scanRowDetalleFactura(scanner interface{ Scan(dest ...any) error }, d *models.DetalleFactura) error {
	return scanner.Scan(
		&d.IDDetalleFactura, &d.IDFactura, &d.IDProducto, &d.Cantidad, &d.Precio, &d.Subtotal, &d.Impuesto, &d.Total,
	)
}

func NewFactura(db *sql.DB) StoreFactura {
	return &storeFactura{db: db}
}

func (s *storeFactura) CreateFactura(ctx context.Context, f *models.Factura, items []*models.DetalleFactura) (*models.Factura, error) {
	err := ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		// Si el número de factura viene vacío, la base de datos (o la lógica del store) lo manejará
		if f.FacNumero == "" {
			f.FacNumero = "AUTO"
		}

		// Si f.FacNumero es "AUTO", intentamos generar el secuencial antes de insertar
		if f.FacNumero == "AUTO" {
			querySec := `SELECT fn_get_next_secuencial($1, 'FACTURA')`
			err := tx.QueryRowContext(ctx, querySec, f.IDEstacion).Scan(&f.FacNumero)
			if err != nil {
				return fmt.Errorf("error al generar secuencial automático: %w", err)
			}
		}

		metadataJSON, _ := json.Marshal(f.Metadata)

		queryFac := `INSERT INTO factura (
			fac_numero, cfac_subtotal, cfac_iva, cfac_total, cfac_observacion, 
			id_user_pos, id_estacion, id_orden_pedido, id_cliente, id_periodo, 
			id_control_estacion, id_status, fecha_operacion, base_impuesto, 
			impuesto, valor_impuesto, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17) 
		RETURNING id_factura, created_at, updated_at`

		err := tx.QueryRowContext(ctx, queryFac,
			f.FacNumero, f.CfacSubtotal, f.CfacIVA, f.CfacTotal, f.CfacObservacion,
			f.IDUserPos, f.IDEstacion, f.IDOrdenPedido, f.IDCliente, f.IDPeriodo,
			f.IDControlEstacion, f.IDStatus, f.FechaOperacion, f.BaseImpuesto,
			f.Impuesto, f.ValorImpuesto, metadataJSON,
		).Scan(&f.IDFactura, &f.CreatedAt, &f.UpdatedAt)

		if err != nil {
			return fmt.Errorf("error al insertar factura: %w", err)
		}

		queryDet := `INSERT INTO detalle_factura (id_factura, id_producto, cantidad, precio, subtotal, impuesto, total) 
		             VALUES ($1, $2, $3, $4, $5, $6, $7)`

		for _, item := range items {
			_, err = tx.ExecContext(ctx, queryDet, f.IDFactura, item.IDProducto, item.Cantidad, item.Precio, item.Subtotal, item.Impuesto, item.Total)
			if err != nil {
				return fmt.Errorf("error al insertar detalle: %w", err)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return f, nil
}

func (s *storeFactura) RegistrarFacturaCompleta(ctx context.Context, req dto.FacturaCompletaRequest, idUsuario uuid.UUID) (*dto.FacturaResponse, error) {
	var res dto.FacturaResponse
	err := ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		if req.Cabecera.FacNumero == "" {
			req.Cabecera.FacNumero = "AUTO"
		}

		cabeceraJSON, err := json.Marshal(req.Cabecera)
		if err != nil {
			return fmt.Errorf("error al serializar cabecera: %w", err)
		}

		detallesJSON, err := json.Marshal(req.Detalles)
		if err != nil {
			return fmt.Errorf("error al serializar detalles: %w", err)
		}

		pagosJSON, err := json.Marshal(req.Pagos)
		if err != nil {
			return fmt.Errorf("error al serializar pagos: %w", err)
		}

		query := `SELECT id_factura, fac_numero, total, status_msg 
		          FROM factura_registrar_completa($1, $2, $3, $4)`

		return tx.QueryRowContext(ctx, query, cabeceraJSON, detallesJSON, pagosJSON, idUsuario).Scan(
			&res.IDFactura, &res.FacNumero, &res.Total, &res.StatusMsg,
		)
	})

	if err != nil {
		return nil, fmt.Errorf("error al registrar factura completa: %w", err)
	}

	return &res, nil
}

func (s *storeFactura) GetFacturaByID(ctx context.Context, id uuid.UUID) (*models.Factura, []*models.DetalleFactura, error) {
	queryFac := fmt.Sprintf("SELECT %s FROM factura WHERE id_factura = $1", facturaSelectFields)
	f := &models.Factura{}
	err := s.scanRowFactura(s.db.QueryRowContext(ctx, queryFac, id), f)
	if err != nil {
		return nil, nil, err
	}

	queryDet := fmt.Sprintf("SELECT %s FROM detalle_factura WHERE id_factura = $1", detalleFacturaSelectFields)
	rows, err := s.db.QueryContext(ctx, queryDet, id)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var items []*models.DetalleFactura
	for rows.Next() {
		item := &models.DetalleFactura{}
		if err := s.scanRowDetalleFactura(rows, item); err != nil {
			return nil, nil, err
		}
		items = append(items, item)
	}

	return f, items, nil
}

func (s *storeFactura) GetAllFacturas(ctx context.Context, params dto.PaginationParams) ([]*models.Factura, error) {
	if params.Limit <= 0 {
		params.Limit = dto.DefaultLimit
	}

	query := fmt.Sprintf(`
		SELECT %s 
		FROM factura 
		WHERE deleted_at IS NULL
	`, facturaSelectFields)

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

	var facturas []*models.Factura
	for rows.Next() {
		f := &models.Factura{}
		if err := s.scanRowFactura(rows, f); err != nil {
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
