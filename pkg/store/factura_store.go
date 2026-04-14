package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/prunus/pkg/dto"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/utils/performance"
)

type StoreFactura interface {
	CreateFactura(ctx context.Context, f *models.Factura, items []*models.DetalleFactura) (*models.Factura, error)
	RegistrarFacturaCompleta(ctx context.Context, req dto.FacturaCompletaRequest, idUsuario uuid.UUID) (*dto.FacturaResponse, error)
	GetFacturaByID(ctx context.Context, id uuid.UUID) (*models.Factura, []*models.DetalleFactura, error)
	GetAllFacturas(ctx context.Context, params dto.PaginationParams) ([]*models.Factura, error)

	// Impuestos
	GetAllImpuestos(ctx context.Context) ([]*models.Impuesto, error)
	GetImpuestoByID(ctx context.Context, id uuid.UUID) (*models.Impuesto, error)
	CreateImpuesto(ctx context.Context, i *models.Impuesto) (*models.Impuesto, error)
	UpdateImpuesto(ctx context.Context, id uuid.UUID, i *models.Impuesto) (*models.Impuesto, error)
	DeleteImpuesto(ctx context.Context, id uuid.UUID) error

	// Formas de Pago
	GetAllFormasPago(ctx context.Context) ([]*models.FormaPago, error)
	GetFormaPagoByID(ctx context.Context, id uuid.UUID) (*models.FormaPago, error)
	CreateFormaPago(ctx context.Context, f *models.FormaPago) (*models.FormaPago, error)
	UpdateFormaPago(ctx context.Context, id uuid.UUID, f *models.FormaPago) (*models.FormaPago, error)
	DeleteFormaPago(ctx context.Context, id uuid.UUID) error
}

type storeFactura struct {
	db *sql.DB
}

// Campos base para SELECT de factura actualizados al esquema moderno con blindaje contra NULLs
const facturaSelectFields = `
	id_factura, 
	COALESCE(fac_numero, ''), 
	COALESCE(subtotal, 0), 
	COALESCE(impuesto, 0), 
	COALESCE(total, 0), 
	COALESCE(observacion, ''), 
	COALESCE(id_usuario, '00000000-0000-0000-0000-000000000000'), 
	COALESCE(id_estacion, '00000000-0000-0000-0000-000000000000'), 
	COALESCE(id_orden_pedido, '00000000-0000-0000-0000-000000000000'), 
	COALESCE(id_cliente, '00000000-0000-0000-0000-000000000000'), 
	COALESCE(id_periodo, '00000000-0000-0000-0000-000000000000'), 
	COALESCE(id_control_estacion, '00000000-0000-0000-0000-000000000000'), 
	COALESCE(id_status, '00000000-0000-0000-0000-000000000000'), 
	COALESCE(id_sucursal, '00000000-0000-0000-0000-000000000000'), 
	COALESCE(fecha_operacion, created_at), 
	fecha_vencimiento,
	COALESCE(base_impuesto, 0), 
	COALESCE(valor_impuesto, 0), 
	metadata, 
	created_at, 
	updated_at
`

// scanRowFactura helper para escanear facturas con soporte para metadata JSONB
func (s *storeFactura) scanRowFactura(scanner interface{ Scan(dest ...any) error }, f *models.Factura) error {
	var metadataJSON []byte
	err := scanner.Scan(
		&f.IDFactura, &f.FacNumero, &f.Subtotal, &f.Impuesto, &f.Total, &f.Observacion,
		&f.IDUsuario, &f.IDEstacion, &f.IDOrdenPedido, &f.IDCliente, &f.IDPeriodo,
		&f.IDControlEstacion, &f.IDStatus, &f.IDSucursal, &f.FechaOperacion, &f.FechaVencimiento,
		&f.BaseImpuesto, &f.ValorImpuesto, &metadataJSON, &f.CreatedAt, &f.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("error escanear factura: %w", err)
	}
	if len(metadataJSON) > 0 {
		json.Unmarshal(metadataJSON, &f.Metadata)
	}
	return nil
}

// Campos base para SELECT de detalle_factura (Corregido precio -> precio_unitario)
const detalleFacturaSelectFields = `
	id_detalle, id_factura, id_producto, cantidad, precio_unitario, subtotal, impuesto, total
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
	defer performance.Trace(ctx, "store", "CreateFactura", performance.DbThreshold, time.Now())

	err := ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		if f.FacNumero == "" {
			f.FacNumero = "AUTO"
		}

		if f.FacNumero == "AUTO" {
			querySec := `SELECT fn_get_next_secuencial($1, 'FACTURA')`
			err := tx.QueryRowContext(ctx, querySec, f.IDEstacion).Scan(&f.FacNumero)
			if err != nil {
				return fmt.Errorf("error al generar secuencial automático: %w", err)
			}
		}

		metadataJSON, _ := json.Marshal(f.Metadata)

		queryFac := `INSERT INTO factura (
			fac_numero, subtotal, impuesto, total, observacion, 
			id_usuario, id_estacion, id_orden_pedido, id_cliente, id_periodo, 
			id_control_estacion, id_status, id_sucursal, fecha_operacion, fecha_vencimiento,
			base_impuesto, valor_impuesto, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18) 
		RETURNING id_factura, created_at, updated_at`

		err := tx.QueryRowContext(ctx, queryFac,
			f.FacNumero, f.Subtotal, f.Impuesto, f.Total, f.Observacion,
			f.IDUsuario, f.IDEstacion, f.IDOrdenPedido, f.IDCliente, f.IDPeriodo,
			f.IDControlEstacion, f.IDStatus, f.IDSucursal, f.FechaOperacion, f.FechaVencimiento,
			f.BaseImpuesto, f.ValorImpuesto, metadataJSON,
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
	defer performance.Trace(ctx, "store", "RegistrarFacturaCompleta", performance.DbThreshold, time.Now())
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
	defer performance.Trace(ctx, "store", "GetFacturaByID", performance.DbThreshold, time.Now())
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
	defer performance.Trace(ctx, "store", "GetAllFacturas", performance.DbThreshold, time.Now())
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
	defer performance.Trace(ctx, "store", "GetAllImpuestos", performance.DbThreshold, time.Now())
	query := `SELECT id_impuesto, nombre, porcentaje, id_status, created_at, updated_at FROM impuesto WHERE deleted_at IS NULL`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var impuestos []*models.Impuesto
	for rows.Next() {
		i := &models.Impuesto{}
		if err := rows.Scan(&i.IDImpuesto, &i.Nombre, &i.Porcentaje, &i.IDStatus, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, err
		}
		impuestos = append(impuestos, i)
	}
	return impuestos, nil
}

func (s *storeFactura) GetImpuestoByID(ctx context.Context, id uuid.UUID) (*models.Impuesto, error) {
	defer performance.Trace(ctx, "store", "GetImpuestoByID", performance.DbThreshold, time.Now())
	query := `SELECT id_impuesto, nombre, porcentaje, id_status, created_at, updated_at FROM impuesto WHERE id_impuesto = $1 AND deleted_at IS NULL`
	i := &models.Impuesto{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(&i.IDImpuesto, &i.Nombre, &i.Porcentaje, &i.IDStatus, &i.CreatedAt, &i.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return i, nil
}

func (s *storeFactura) CreateImpuesto(ctx context.Context, i *models.Impuesto) (*models.Impuesto, error) {
	defer performance.Trace(ctx, "store", "CreateImpuesto", performance.DbThreshold, time.Now())
	err := ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		query := `INSERT INTO impuesto (nombre, porcentaje, id_status) 
		          VALUES ($1, $2, $3) 
		          RETURNING id_impuesto, created_at, updated_at`
		
		return tx.QueryRowContext(ctx, query, i.Nombre, i.Porcentaje, i.IDStatus).Scan(
			&i.IDImpuesto, &i.CreatedAt, &i.UpdatedAt,
		)
	})
	if err != nil {
		return nil, err
	}
	return i, nil
}

func (s *storeFactura) UpdateImpuesto(ctx context.Context, id uuid.UUID, i *models.Impuesto) (*models.Impuesto, error) {
	defer performance.Trace(ctx, "store", "UpdateImpuesto", performance.DbThreshold, time.Now())
	err := ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		query := `UPDATE impuesto SET nombre = $1, porcentaje = $2, id_status = $3, updated_at = CURRENT_TIMESTAMP 
		          WHERE id_impuesto = $4 AND deleted_at IS NULL
		          RETURNING created_at, updated_at`
		
		return tx.QueryRowContext(ctx, query, i.Nombre, i.Porcentaje, i.IDStatus, id).Scan(&i.CreatedAt, &i.UpdatedAt)
	})
	if err != nil {
		return nil, err
	}
	i.IDImpuesto = id
	return i, nil
}

func (s *storeFactura) DeleteImpuesto(ctx context.Context, id uuid.UUID) error {
	defer performance.Trace(ctx, "store", "DeleteImpuesto", performance.DbThreshold, time.Now())
	return ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		query := `UPDATE impuesto SET deleted_at = CURRENT_TIMESTAMP WHERE id_impuesto = $1 AND deleted_at IS NULL`
		_, err := tx.ExecContext(ctx, query, id)
		return err
	})
}

func (s *storeFactura) GetAllFormasPago(ctx context.Context) ([]*models.FormaPago, error) {
	defer performance.Trace(ctx, "store", "GetAllFormasPago", performance.DbThreshold, time.Now())
	query := `SELECT id_forma_pago, nombre, requiere_ref, id_status, created_at, updated_at FROM forma_pago WHERE deleted_at IS NULL ORDER BY created_at DESC`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var formas []*models.FormaPago
	for rows.Next() {
		f := &models.FormaPago{}
		if err := rows.Scan(&f.IDFormaPago, &f.Nombre, &f.RequiereRef, &f.IDStatus, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, err
		}
		formas = append(formas, f)
	}
	return formas, nil
}

func (s *storeFactura) GetFormaPagoByID(ctx context.Context, id uuid.UUID) (*models.FormaPago, error) {
	defer performance.Trace(ctx, "store", "GetFormaPagoByID", performance.DbThreshold, time.Now())
	query := `SELECT id_forma_pago, nombre, requiere_ref, id_status, created_at, updated_at FROM forma_pago WHERE id_forma_pago = $1 AND deleted_at IS NULL`
	f := &models.FormaPago{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(&f.IDFormaPago, &f.Nombre, &f.RequiereRef, &f.IDStatus, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (s *storeFactura) CreateFormaPago(ctx context.Context, f *models.FormaPago) (*models.FormaPago, error) {
	defer performance.Trace(ctx, "store", "CreateFormaPago", performance.DbThreshold, time.Now())
	err := ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		query := `INSERT INTO forma_pago (nombre, requiere_ref, id_status) 
		          VALUES ($1, $2, $3) 
		          RETURNING id_forma_pago, created_at, updated_at`
		
		return tx.QueryRowContext(ctx, query, f.Nombre, f.RequiereRef, f.IDStatus).Scan(
			&f.IDFormaPago, &f.CreatedAt, &f.UpdatedAt,
		)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (s *storeFactura) UpdateFormaPago(ctx context.Context, id uuid.UUID, f *models.FormaPago) (*models.FormaPago, error) {
	defer performance.Trace(ctx, "store", "UpdateFormaPago", performance.DbThreshold, time.Now())
	err := ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		query := `UPDATE forma_pago SET nombre = $1, requiere_ref = $2, id_status = $3, updated_at = NOW() 
		          WHERE id_forma_pago = $4 AND deleted_at IS NULL
		          RETURNING created_at, updated_at`
		
		return tx.QueryRowContext(ctx, query, f.Nombre, f.RequiereRef, f.IDStatus, id).Scan(&f.CreatedAt, &f.UpdatedAt)
	})
	if err != nil {
		return nil, err
	}
	f.IDFormaPago = id
	return f, nil
}

func (s *storeFactura) DeleteFormaPago(ctx context.Context, id uuid.UUID) error {
	defer performance.Trace(ctx, "store", "DeleteFormaPago", performance.DbThreshold, time.Now())
	return ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		query := `UPDATE forma_pago SET deleted_at = NOW() WHERE id_forma_pago = $1 AND deleted_at IS NULL`
		_, err := tx.ExecContext(ctx, query, id)
		return err
	})
}
