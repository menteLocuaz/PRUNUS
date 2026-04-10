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

type StoreInventario interface {
	GetAllInventario(ctx context.Context, params dto.PaginationParams) ([]*models.Inventario, error)
	GetInventarioByID(ctx context.Context, id uuid.UUID) (*models.Inventario, error)
	GetInventarioBySucursal(ctx context.Context, idSucursal uuid.UUID, params dto.PaginationParams) ([]*models.Inventario, error)
	GetInventarioByProductoYSucursal(ctx context.Context, idProducto, idSucursal uuid.UUID) (*models.Inventario, error)
	CreateInventario(ctx context.Context, inventario *models.Inventario) (*models.Inventario, error)
	UpdateInventario(ctx context.Context, id uuid.UUID, inventario *models.Inventario) (*models.Inventario, error)
	DeleteInventario(ctx context.Context, id uuid.UUID) error

	// Movimientos
	RegistrarMovimiento(ctx context.Context, m *models.MovimientoInventario) (*models.MovimientoInventario, error)
	RegistrarMovimientoMasivo(ctx context.Context, idSucursal, idUsuario uuid.UUID, tipoMov, referencia string, items []models.MovimientoItem) ([]*models.MovimientoInventario, error)
	GetMovimientosByProducto(ctx context.Context, productoID uuid.UUID, params dto.PaginationParams) ([]*models.MovimientoInventario, error)
	GetAlertasStock(ctx context.Context, sucursalID uuid.UUID) ([]*models.Inventario, error)
	GetValuacion(ctx context.Context, sucursalID uuid.UUID, metodo string) (float64, error)
	GetAnalisisRotacion(ctx context.Context, sucursalID uuid.UUID) (map[string][]uuid.UUID, error)

	// Lotes
	CreateLote(ctx context.Context, lote *models.Lote) (*models.Lote, error)
	GetLotesByProducto(ctx context.Context, idProducto, idSucursal uuid.UUID) ([]*models.Lote, error)
	UpdateLoteCantidad(ctx context.Context, idLote uuid.UUID, cantidad float64) error

	// Analítica operativa
	GetRotacionDetalle(ctx context.Context, sucursalID uuid.UUID, params dto.RotacionFiltroParams) ([]*dto.RotacionProductoResponse, error)
	GetComposicionCategoria(ctx context.Context, sucursalID uuid.UUID) ([]*dto.ComposicionCategoriaResponse, error)
	GetAlertasStockDetalle(ctx context.Context, sucursalID uuid.UUID) ([]*dto.AlertaStockResponse, error)

	// Estado financiero
	CapturarSnapshotInventario(ctx context.Context, sucursalID uuid.UUID) error
	GetValorHistorico(ctx context.Context, sucursalID uuid.UUID, params dto.RotacionFiltroParams) ([]*dto.ValorHistoricoResponse, error)
	GetPerdidas(ctx context.Context, sucursalID uuid.UUID, params dto.RotacionFiltroParams) ([]*dto.PerdidaResponse, error)
	GetMargenGanancia(ctx context.Context, sucursalID uuid.UUID, params dto.RotacionFiltroParams) ([]*dto.MargenProductoResponse, error)
}

type storeInventario struct {
	db *sql.DB
}

// Campos base para SELECT de inventario
const inventarioSelectFields = `
	id_inventario, id_producto, id_sucursal, stock_actual, stock_minimo, 
	stock_maximo, precio_compra, precio_venta, created_at, updated_at
`

// scanRowInventario helper para escanear inventario
func (s *storeInventario) scanRowInventario(scanner interface{ Scan(dest ...any) error }, i *models.Inventario) error {
	return scanner.Scan(
		&i.IDInventario, &i.IDProducto, &i.IDSucursal, &i.StockActual, &i.StockMinimo,
		&i.StockMaximo, &i.PrecioCompra, &i.PrecioVenta, &i.CreatedAt, &i.UpdatedAt,
	)
}

// Campos base para SELECT de movimientos
const movimientoSelectFields = `
	id_movimiento, id_producto, id_sucursal, tipo_movimiento, cantidad, 
	costo_unitario, precio_unitario, stock_anterior, stock_posterior, 
	fecha, id_usuario, referencia
`

// scanRowMovimiento helper para escanear movimientos
func (s *storeInventario) scanRowMovimiento(scanner interface{ Scan(dest ...any) error }, m *models.MovimientoInventario) error {
	return scanner.Scan(
		&m.IDMovimiento, &m.IDProducto, &m.IDSucursal, &m.TipoMovimiento, &m.Cantidad,
		&m.CostoUnitario, &m.PrecioUnitario, &m.StockAnterior, &m.StockPosterior,
		&m.Fecha, &m.IDUsuario, &m.Referencia,
	)
}

// Campos base para SELECT de lotes
const loteSelectFields = `
	id_lote, id_producto, id_sucursal, codigo_lote, cantidad_inicial, 
	cantidad_actual, costo_compra, fecha_vencimiento, fecha_recepcion,
	created_at, updated_at
`

// scanRowLote helper para escanear lotes
func (s *storeInventario) scanRowLote(scanner interface{ Scan(dest ...any) error }, l *models.Lote) error {
	return scanner.Scan(
		&l.IDLote, &l.IDProducto, &l.IDSucursal, &l.CodigoLote, &l.CantidadInicial,
		&l.CantidadActual, &l.CostoCompra, &l.FechaVencimiento, &l.FechaRecepcion,
		&l.CreatedAt, &l.UpdatedAt,
	)
}

func NewInventario(db *sql.DB) StoreInventario {
	return &storeInventario{db: db}
}

func (s *storeInventario) GetAllInventario(ctx context.Context, params dto.PaginationParams) ([]*models.Inventario, error) {
	defer performance.Trace(ctx, "store", "GetAllInventario", performance.DbThreshold, time.Now())

	if params.Limit <= 0 {
		params.Limit = dto.DefaultLimit
	}

	query := fmt.Sprintf(`
		SELECT %s
		FROM inventario
		WHERE deleted_at IS NULL
	`, inventarioSelectFields)

	var args []interface{}
	if params.LastDate != nil {
		query += " AND created_at < $1"
		args = append(args, params.LastDate)
	}

	query += " ORDER BY created_at DESC LIMIT $" + fmt.Sprint(len(args)+1)
	args = append(args, params.Limit)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error al obtener inventario: %w", err)
	}
	defer rows.Close()

	var inventarios []*models.Inventario
	for rows.Next() {
		i := &models.Inventario{}
		if err := s.scanRowInventario(rows, i); err != nil {
			return nil, fmt.Errorf("error al escanear inventario: %w", err)
		}
		inventarios = append(inventarios, i)
	}

	return inventarios, nil
}

func (s *storeInventario) GetInventarioBySucursal(ctx context.Context, idSucursal uuid.UUID, params dto.PaginationParams) ([]*models.Inventario, error) {
	defer performance.Trace(ctx, "store", "GetInventarioBySucursal", performance.DbThreshold, time.Now())

	if params.Limit <= 0 {
		params.Limit = dto.DefaultLimit
	}

	query := fmt.Sprintf(`
		SELECT %s
		FROM inventario
		WHERE id_sucursal = $1 AND deleted_at IS NULL
	`, inventarioSelectFields)

	var args []interface{}
	args = append(args, idSucursal)

	if params.LastDate != nil {
		query += " AND created_at < $2"
		args = append(args, params.LastDate)
	}

	query += " ORDER BY created_at DESC LIMIT $" + fmt.Sprint(len(args)+1)
	args = append(args, params.Limit)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error al obtener inventario por sucursal: %w", err)
	}
	defer rows.Close()

	var inventarios []*models.Inventario
	for rows.Next() {
		i := &models.Inventario{}
		if err := s.scanRowInventario(rows, i); err != nil {
			return nil, fmt.Errorf("error al escanear inventario: %w", err)
		}
		inventarios = append(inventarios, i)
	}

	return inventarios, nil
}

func (s *storeInventario) GetInventarioByID(ctx context.Context, id uuid.UUID) (*models.Inventario, error) {
	defer performance.Trace(ctx, "store", "GetInventarioByID", performance.DbThreshold, time.Now())
	query := fmt.Sprintf(`
		SELECT %s
		FROM inventario
		WHERE id_inventario = $1 AND deleted_at IS NULL
	`, inventarioSelectFields)

	i := &models.Inventario{}
	err := s.scanRowInventario(s.db.QueryRowContext(ctx, query, id), i)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("inventario con ID %s no encontrado", id)
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener inventario: %w", err)
	}

	return i, nil
}

func (s *storeInventario) GetInventarioByProductoYSucursal(ctx context.Context, idProducto, idSucursal uuid.UUID) (*models.Inventario, error) {
	defer performance.Trace(ctx, "store", "GetInventarioByProductoYSucursal", performance.DbThreshold, time.Now())
	query := fmt.Sprintf(`
		SELECT %s
		FROM inventario
		WHERE id_producto = $1 AND id_sucursal = $2 AND deleted_at IS NULL
	`, inventarioSelectFields)

	i := &models.Inventario{}
	err := s.scanRowInventario(s.db.QueryRowContext(ctx, query, idProducto, idSucursal), i)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener inventario por producto y sucursal: %w", err)
	}

	return i, nil
}

func (s *storeInventario) CreateInventario(ctx context.Context, inventario *models.Inventario) (*models.Inventario, error) {
	defer performance.Trace(ctx, "store", "CreateInventario", performance.DbThreshold, time.Now())

	err := ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		query := `
			INSERT INTO inventario (
				id_producto, id_sucursal, stock_actual, stock_minimo, 
				stock_maximo, precio_compra, precio_venta
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id_inventario, created_at, updated_at
		`
		return tx.QueryRowContext(ctx, query,
			inventario.IDProducto, inventario.IDSucursal, inventario.StockActual, inventario.StockMinimo,
			inventario.StockMaximo, inventario.PrecioCompra, inventario.PrecioVenta,
		).Scan(&inventario.IDInventario, &inventario.CreatedAt, &inventario.UpdatedAt)
	})

	if err != nil {
		return nil, fmt.Errorf("error al crear inventario: %w", err)
	}

	return inventario, nil
}

func (s *storeInventario) UpdateInventario(ctx context.Context, id uuid.UUID, inventario *models.Inventario) (*models.Inventario, error) {
	defer performance.Trace(ctx, "store", "UpdateInventario", performance.DbThreshold, time.Now())

	err := ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		query := `
			UPDATE inventario
			SET 
				stock_actual = $1, stock_minimo = $2, stock_maximo = $3, 
				precio_compra = $4, precio_venta = $5, updated_at = CURRENT_TIMESTAMP
			WHERE id_inventario = $6 AND deleted_at IS NULL
			RETURNING id_inventario, id_producto, id_sucursal, created_at, updated_at
		`
		return tx.QueryRowContext(ctx, query,
			inventario.StockActual, inventario.StockMinimo, inventario.StockMaximo,
			inventario.PrecioCompra, inventario.PrecioVenta, id,
		).Scan(&inventario.IDInventario, &inventario.IDProducto, &inventario.IDSucursal, &inventario.CreatedAt, &inventario.UpdatedAt)
	})

	if err != nil {
		return nil, fmt.Errorf("error al actualizar inventario: %w", err)
	}

	return inventario, nil
}

func (s *storeInventario) DeleteInventario(ctx context.Context, id uuid.UUID) error {
	defer performance.Trace(ctx, "store", "DeleteInventario", performance.DbThreshold, time.Now())

	return ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		query := `UPDATE inventario SET deleted_at = CURRENT_TIMESTAMP WHERE id_inventario = $1 AND deleted_at IS NULL`
		result, err := tx.ExecContext(ctx, query, id)
		if err != nil {
			return err
		}
		rows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rows == 0 {
			return sql.ErrNoRows
		}
		return nil
	})
}

func (s *storeInventario) RegistrarMovimiento(ctx context.Context, m *models.MovimientoInventario) (*models.MovimientoInventario, error) {
	defer performance.Trace(ctx, "store", "RegistrarMovimiento", performance.DbThreshold, time.Now())

	err := ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		// Creamos el JSON de items para la función almacenada.
		itemsJSON := fmt.Sprintf(`[{"id_producto": "%s", "cantidad": %f}]`, m.IDProducto, m.Cantidad)

		query := `SELECT id_movimiento, id_producto, stock_anterior, cantidad, stock_posterior 
		          FROM inventario_ia_movimiento($1, $2, $3, $4, $5)`

		return tx.QueryRowContext(ctx, query,
			m.IDSucursal, m.IDUsuario, m.TipoMovimiento, m.Referencia, itemsJSON,
		).Scan(&m.IDMovimiento, &m.IDProducto, &m.StockAnterior, &m.Cantidad, &m.StockPosterior)
	})

	if err != nil {
		return nil, fmt.Errorf("error al registrar movimiento mediante función: %w", err)
	}

	m.Fecha = time.Now()
	m.CreatedAt = m.Fecha
	m.UpdatedAt = m.Fecha

	return m, nil
}

func (s *storeInventario) RegistrarMovimientoMasivo(ctx context.Context, idSucursal, idUsuario uuid.UUID, tipoMov, referencia string, items []models.MovimientoItem) ([]*models.MovimientoInventario, error) {
	defer performance.Trace(ctx, "store", "RegistrarMovimientoMasivo", performance.DbThreshold, time.Now())

	itemsJSON, err := json.Marshal(items)
	if err != nil {
		return nil, fmt.Errorf("error al serializar items: %w", err)
	}

	var resultados []*models.MovimientoInventario
	err = ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		query := `SELECT id_movimiento, id_producto, stock_anterior, cantidad, stock_posterior 
		          FROM inventario_ia_movimiento($1, $2, $3, $4, $5)`

		rows, err := tx.QueryContext(ctx, query, idSucursal, idUsuario, tipoMov, referencia, itemsJSON)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			m := &models.MovimientoInventario{
				IDSucursal:     idSucursal,
				IDUsuario:      idUsuario,
				TipoMovimiento: tipoMov,
				Referencia:     referencia,
				Fecha:          time.Now(),
			}
			if err := rows.Scan(&m.IDMovimiento, &m.IDProducto, &m.StockAnterior, &m.Cantidad, &m.StockPosterior); err != nil {
				return err
			}
			m.CreatedAt = m.Fecha
			m.UpdatedAt = m.Fecha
			resultados = append(resultados, m)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error al ejecutar movimiento masivo: %w", err)
	}

	return resultados, nil
}

func (s *storeInventario) GetMovimientosByProducto(ctx context.Context, productoID uuid.UUID, params dto.PaginationParams) ([]*models.MovimientoInventario, error) {
	defer performance.Trace(ctx, "store", "GetMovimientosByProducto", performance.DbThreshold, time.Now())

	if params.Limit <= 0 {
		params.Limit = dto.DefaultLimit
	}

	query := fmt.Sprintf(`
		SELECT %s 
		FROM movimientos_inventario 
		WHERE id_producto = $1 AND deleted_at IS NULL
	`, movimientoSelectFields)

	var args []interface{}
	args = append(args, productoID)

	if params.LastDate != nil {
		query += " AND fecha < $2"
		args = append(args, params.LastDate)
	}

	query += " ORDER BY fecha DESC LIMIT $" + fmt.Sprint(len(args)+1)
	args = append(args, params.Limit)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movimientos []*models.MovimientoInventario
	for rows.Next() {
		m := &models.MovimientoInventario{}
		if err := s.scanRowMovimiento(rows, m); err != nil {
			return nil, err
		}
		movimientos = append(movimientos, m)
	}
	return movimientos, nil
}

func (s *storeInventario) GetAlertasStock(ctx context.Context, sucursalID uuid.UUID) ([]*models.Inventario, error) {
	defer performance.Trace(ctx, "store", "GetAlertasStock", performance.DbThreshold, time.Now())
	query := fmt.Sprintf(`
		SELECT %s
		FROM inventario
		WHERE id_sucursal = $1 AND stock_actual <= stock_minimo AND deleted_at IS NULL
	`, inventarioSelectFields)

	rows, err := s.db.QueryContext(ctx, query, sucursalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alertas []*models.Inventario
	for rows.Next() {
		i := &models.Inventario{}
		if err := s.scanRowInventario(rows, i); err != nil {
			return nil, err
		}
		alertas = append(alertas, i)
	}
	return alertas, nil
}

func (s *storeInventario) GetValuacion(ctx context.Context, sucursalID uuid.UUID, metodo string) (float64, error) {
	defer performance.Trace(ctx, "store", "GetValuacion", performance.DbThreshold, time.Now())

	var query string
	switch metodo {
	case "promedio":
		query = `SELECT COALESCE(SUM(stock_actual * precio_compra), 0) FROM inventario WHERE id_sucursal = $1 AND deleted_at IS NULL`
	case "peps", "ueps":
		query = `SELECT COALESCE(SUM(cantidad_actual * costo_compra), 0) FROM lotes WHERE id_sucursal = $1 AND deleted_at IS NULL AND cantidad_actual > 0`
	default:
		query = `SELECT COALESCE(SUM(stock_actual * precio_compra), 0) FROM inventario WHERE id_sucursal = $1 AND deleted_at IS NULL`
	}

	var total float64
	err := s.db.QueryRowContext(ctx, query, sucursalID).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (s *storeInventario) GetAnalisisRotacion(ctx context.Context, sucursalID uuid.UUID) (map[string][]uuid.UUID, error) {
	defer performance.Trace(ctx, "store", "GetAnalisisRotacion", performance.DbThreshold, time.Now())

	query := `
		WITH ranking_productos AS (
			SELECT 
				id_producto,
				(stock_actual * precio_compra) as valor_total,
				SUM(stock_actual * precio_compra) OVER (ORDER BY (stock_actual * precio_compra) DESC) as acumulado,
				SUM(stock_actual * precio_compra) OVER () as total
			FROM inventario
			WHERE id_sucursal = $1 AND deleted_at IS NULL AND stock_actual > 0
		)
		SELECT 
			id_producto,
			CASE 
				WHEN (acumulado / NULLIF(total, 0)) <= 0.80 THEN 'A'
				WHEN (acumulado / NULLIF(total, 0)) <= 0.95 THEN 'B'
				ELSE 'C'
			END as clase
		FROM ranking_productos
		ORDER BY valor_total DESC;
	`

	rows, err := s.db.QueryContext(ctx, query, sucursalID)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar análisis ABC: %w", err)
	}
	defer rows.Close()

	abc := make(map[string][]uuid.UUID)
	for rows.Next() {
		var id uuid.UUID
		var clase string
		if err := rows.Scan(&id, &clase); err != nil {
			return nil, err
		}
		abc[clase] = append(abc[clase], id)
	}

	return abc, nil
}

func (s *storeInventario) CreateLote(ctx context.Context, lote *models.Lote) (*models.Lote, error) {
	defer performance.Trace(ctx, "store", "CreateLote", performance.DbThreshold, time.Now())
	query := `
		INSERT INTO lotes (
			id_producto, id_sucursal, codigo_lote, cantidad_inicial, 
			cantidad_actual, costo_compra, fecha_vencimiento, fecha_recepcion
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id_lote, created_at, updated_at
	`

	err := s.db.QueryRowContext(ctx, query,
		lote.IDProducto, lote.IDSucursal, lote.CodigoLote, lote.CantidadInicial,
		lote.CantidadActual, lote.CostoCompra, lote.FechaVencimiento, lote.FechaRecepcion,
	).Scan(&lote.IDLote, &lote.CreatedAt, &lote.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("error al crear lote: %w", err)
	}

	return lote, nil
}

func (s *storeInventario) GetLotesByProducto(ctx context.Context, idProducto, idSucursal uuid.UUID) ([]*models.Lote, error) {
	defer performance.Trace(ctx, "store", "GetLotesByProducto", performance.DbThreshold, time.Now())
	query := fmt.Sprintf(`
		SELECT %s
		FROM lotes
		WHERE id_producto = $1 AND id_sucursal = $2 AND deleted_at IS NULL AND cantidad_actual > 0
		ORDER BY fecha_recepcion ASC
	`, loteSelectFields)

	rows, err := s.db.QueryContext(ctx, query, idProducto, idSucursal)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lotes []*models.Lote
	for rows.Next() {
		l := &models.Lote{}
		if err := s.scanRowLote(rows, l); err != nil {
			return nil, err
		}
		lotes = append(lotes, l)
	}
	return lotes, nil
}

func (s *storeInventario) UpdateLoteCantidad(ctx context.Context, idLote uuid.UUID, cantidad float64) error {
	defer performance.Trace(ctx, "store", "UpdateLoteCantidad", performance.DbThreshold, time.Now())
	query := `UPDATE lotes SET cantidad_actual = $1, updated_at = CURRENT_TIMESTAMP WHERE id_lote = $2`
	_, err := s.db.ExecContext(ctx, query, cantidad, idLote)
	return err
}

func (s *storeInventario) GetRotacionDetalle(ctx context.Context, sucursalID uuid.UUID, params dto.RotacionFiltroParams) ([]*dto.RotacionProductoResponse, error) {
	defer performance.Trace(ctx, "store", "GetRotacionDetalle", performance.DbThreshold, time.Now())

	query := `
		WITH ventas AS (
			SELECT
				id_producto,
				SUM(cantidad * COALESCE(costo_unitario, 0)) AS cogs,
				SUM(cantidad)                               AS unidades_vendidas
			FROM movimientos_inventario
			WHERE id_sucursal     = $1
			  AND fecha           BETWEEN $2 AND $3
			  AND tipo_movimiento IN ('SALIDA', 'VENTA')
			  AND deleted_at      IS NULL
			GROUP BY id_producto
		),
		primer_mov AS (
			SELECT DISTINCT ON (id_producto)
				id_producto,
				stock_anterior AS stock_inicio
			FROM movimientos_inventario
			WHERE id_sucursal = $1
			  AND fecha       BETWEEN $2 AND $3
			  AND deleted_at  IS NULL
			ORDER BY id_producto, fecha ASC
		),
		ultimo_mov AS (
			SELECT DISTINCT ON (id_producto)
				id_producto,
				stock_posterior AS stock_fin
			FROM movimientos_inventario
			WHERE id_sucursal = $1
			  AND fecha       BETWEEN $2 AND $3
			  AND deleted_at  IS NULL
			ORDER BY id_producto, fecha DESC
		)
		SELECT
			v.id_producto,
			v.cogs,
			v.unidades_vendidas,
			(COALESCE(p.stock_inicio, 0) + COALESCE(u.stock_fin, 0)) / 2.0 AS inventario_promedio,
			CASE
				WHEN (COALESCE(p.stock_inicio, 0) + COALESCE(u.stock_fin, 0)) = 0 THEN 0
				ELSE v.unidades_vendidas / ((COALESCE(p.stock_inicio, 0) + COALESCE(u.stock_fin, 0)) / 2.0)
			END AS indice_rotacion
		FROM ventas v
		LEFT JOIN primer_mov p ON v.id_producto = p.id_producto
		LEFT JOIN ultimo_mov u ON v.id_producto = u.id_producto
		ORDER BY indice_rotacion DESC
	`

	rows, err := s.db.QueryContext(ctx, query, sucursalID, params.FechaInicio, params.FechaFin)
	if err != nil {
		return nil, fmt.Errorf("error al calcular rotación de inventario: %w", err)
	}
	defer rows.Close()

	var resultado []*dto.RotacionProductoResponse
	for rows.Next() {
		r := &dto.RotacionProductoResponse{}
		if err := rows.Scan(&r.IDProducto, &r.COGS, &r.UnidadesVendidas, &r.InventarioPromedio, &r.IndiceRotacion); err != nil {
			return nil, fmt.Errorf("error al escanear rotación: %w", err)
		}
		resultado = append(resultado, r)
	}
	return resultado, nil
}

func (s *storeInventario) GetComposicionCategoria(ctx context.Context, sucursalID uuid.UUID) ([]*dto.ComposicionCategoriaResponse, error) {
	defer performance.Trace(ctx, "store", "GetComposicionCategoria", performance.DbThreshold, time.Now())

	query := `
		WITH valores AS (
			SELECT
				c.id_categoria,
				c.nombre                                          AS nombre_categoria,
				COUNT(DISTINCT p.id_producto)                    AS num_productos,
				COALESCE(SUM(i.stock_actual), 0)                 AS cantidad_total,
				COALESCE(SUM(i.stock_actual * i.precio_compra), 0) AS valor_total
			FROM categoria c
			LEFT JOIN producto p
				ON p.id_categoria = c.id_categoria
				AND p.deleted_at IS NULL
			LEFT JOIN inventario i
				ON i.id_producto = p.id_producto
				AND i.id_sucursal = $1
				AND i.deleted_at IS NULL
			WHERE c.id_sucursal = $1
			  AND c.deleted_at IS NULL
			GROUP BY c.id_categoria, c.nombre
		),
		total AS (
			SELECT NULLIF(SUM(valor_total), 0) AS gran_total FROM valores
		)
		SELECT
			v.id_categoria,
			v.nombre_categoria,
			v.num_productos,
			v.cantidad_total,
			v.valor_total,
			COALESCE(ROUND((v.valor_total * 100.0 / t.gran_total)::numeric, 2), 0) AS porcentaje_valor
		FROM valores v, total t
		ORDER BY v.valor_total DESC
	`

	rows, err := s.db.QueryContext(ctx, query, sucursalID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener composición por categoría: %w", err)
	}
	defer rows.Close()

	var resultado []*dto.ComposicionCategoriaResponse
	for rows.Next() {
		r := &dto.ComposicionCategoriaResponse{}
		if err := rows.Scan(
			&r.IDCategoria, &r.NombreCategoria, &r.NumProductos,
			&r.CantidadTotal, &r.ValorTotal, &r.PorcentajeValor,
		); err != nil {
			return nil, fmt.Errorf("error al escanear composición: %w", err)
		}
		resultado = append(resultado, r)
	}
	return resultado, nil
}

func (s *storeInventario) CapturarSnapshotInventario(ctx context.Context, sucursalID uuid.UUID) error {
	defer performance.Trace(ctx, "store", "CapturarSnapshotInventario", performance.DbThreshold, time.Now())
	_, err := s.db.ExecContext(ctx, `SELECT fn_snapshot_inventario($1)`, sucursalID)
	if err != nil {
		return fmt.Errorf("error al capturar snapshot de inventario: %w", err)
	}
	return nil
}

func (s *storeInventario) GetValorHistorico(ctx context.Context, sucursalID uuid.UUID, params dto.RotacionFiltroParams) ([]*dto.ValorHistoricoResponse, error) {
	defer performance.Trace(ctx, "store", "GetValorHistorico", performance.DbThreshold, time.Now())

	query := `
		SELECT fecha_snapshot, valor_total, cantidad_total, num_productos
		FROM inventario_historico
		WHERE id_sucursal    = $1
		  AND fecha_snapshot BETWEEN $2 AND $3
		ORDER BY fecha_snapshot ASC
	`

	rows, err := s.db.QueryContext(ctx, query, sucursalID, params.FechaInicio, params.FechaFin)
	if err != nil {
		return nil, fmt.Errorf("error al obtener historial de valor: %w", err)
	}
	defer rows.Close()

	var resultado []*dto.ValorHistoricoResponse
	for rows.Next() {
		r := &dto.ValorHistoricoResponse{}
		if err := rows.Scan(&r.FechaSnapshot, &r.ValorTotal, &r.CantidadTotal, &r.NumProductos); err != nil {
			return nil, fmt.Errorf("error al escanear historial: %w", err)
		}
		resultado = append(resultado, r)
	}
	return resultado, nil
}

func (s *storeInventario) GetPerdidas(ctx context.Context, sucursalID uuid.UUID, params dto.RotacionFiltroParams) ([]*dto.PerdidaResponse, error) {
	defer performance.Trace(ctx, "store", "GetPerdidas", performance.DbThreshold, time.Now())

	query := `
		SELECT
			mi.id_producto,
			p.nombre                             AS nombre_producto,
			mi.tipo_movimiento,
			SUM(mi.cantidad)                     AS unidades_perdidas,
			SUM(mi.cantidad * mi.costo_unitario) AS valor_perdido
		FROM movimientos_inventario mi
		INNER JOIN producto p
			ON mi.id_producto = p.id_producto
			AND p.deleted_at IS NULL
		WHERE mi.id_sucursal     = $1
		  AND mi.fecha           BETWEEN $2 AND $3
		  AND mi.tipo_movimiento IN ('MERMA', 'CADUCADO')
		  AND mi.deleted_at      IS NULL
		GROUP BY mi.id_producto, p.nombre, mi.tipo_movimiento
		ORDER BY valor_perdido DESC
	`

	rows, err := s.db.QueryContext(ctx, query, sucursalID, params.FechaInicio, params.FechaFin)
	if err != nil {
		return nil, fmt.Errorf("error al obtener pérdidas: %w", err)
	}
	defer rows.Close()

	var resultado []*dto.PerdidaResponse
	for rows.Next() {
		r := &dto.PerdidaResponse{}
		if err := rows.Scan(
			&r.IDProducto, &r.NombreProducto, &r.TipoMovimiento,
			&r.UnidadesPerdidas, &r.ValorPerdido,
		); err != nil {
			return nil, fmt.Errorf("error al escanear pérdida: %w", err)
		}
		resultado = append(resultado, r)
	}
	return resultado, nil
}

func (s *storeInventario) GetMargenGanancia(ctx context.Context, sucursalID uuid.UUID, params dto.RotacionFiltroParams) ([]*dto.MargenProductoResponse, error) {
	defer performance.Trace(ctx, "store", "GetMargenGanancia", performance.DbThreshold, time.Now())

	query := `
		SELECT
			mi.id_producto,
			p.nombre                                                          AS nombre_producto,
			ROUND(AVG(mi.costo_unitario)::numeric, 4)                        AS costo_prom,
			ROUND(AVG(mi.precio_unitario)::numeric, 4)                       AS precio_venta_prom,
			ROUND((AVG(mi.precio_unitario) - AVG(mi.costo_unitario))::numeric, 4) AS margen_bruto,
			CASE
				WHEN AVG(mi.precio_unitario) = 0 THEN 0
				ELSE ROUND(
					((AVG(mi.precio_unitario) - AVG(mi.costo_unitario)) / AVG(mi.precio_unitario) * 100)::numeric,
					2
				)
			END AS margen_porcentaje,
			SUM(mi.cantidad) AS unidades_vendidas
		FROM movimientos_inventario mi
		INNER JOIN producto p
			ON mi.id_producto = p.id_producto
			AND p.deleted_at IS NULL
		WHERE mi.id_sucursal     = $1
		  AND mi.fecha           BETWEEN $2 AND $3
		  AND mi.tipo_movimiento IN ('VENTA', 'SALIDA')
		  AND mi.precio_unitario > 0
		  AND mi.deleted_at      IS NULL
		GROUP BY mi.id_producto, p.nombre
		ORDER BY margen_porcentaje DESC
	`

	rows, err := s.db.QueryContext(ctx, query, sucursalID, params.FechaInicio, params.FechaFin)
	if err != nil {
		return nil, fmt.Errorf("error al calcular margen de ganancia: %w", err)
	}
	defer rows.Close()

	var resultado []*dto.MargenProductoResponse
	for rows.Next() {
		r := &dto.MargenProductoResponse{}
		if err := rows.Scan(
			&r.IDProducto, &r.NombreProducto,
			&r.CostoProm, &r.PrecioVentaProm, &r.MargenBruto, &r.MargenPorcentaje,
			&r.UnidadesVendidas,
		); err != nil {
			return nil, fmt.Errorf("error al escanear margen: %w", err)
		}
		resultado = append(resultado, r)
	}
	return resultado, nil
}

func (s *storeInventario) GetAlertasStockDetalle(ctx context.Context, sucursalID uuid.UUID) ([]*dto.AlertaStockResponse, error) {
	defer performance.Trace(ctx, "store", "GetAlertasStockDetalle", performance.DbThreshold, time.Now())

	query := `
		SELECT
			i.id_inventario,
			i.id_producto,
			p.nombre       AS nombre_producto,
			p.sku,
			i.id_sucursal,
			i.stock_actual,
			i.stock_minimo,
			(i.stock_minimo - i.stock_actual) AS deficit,
			i.precio_compra
		FROM inventario i
		INNER JOIN producto p
			ON i.id_producto = p.id_producto
			AND p.deleted_at IS NULL
		WHERE i.id_sucursal  = $1
		  AND i.stock_actual <= i.stock_minimo
		  AND i.deleted_at   IS NULL
		ORDER BY deficit DESC, i.stock_actual ASC
	`

	rows, err := s.db.QueryContext(ctx, query, sucursalID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener alertas de stock: %w", err)
	}
	defer rows.Close()

	var alertas []*dto.AlertaStockResponse
	for rows.Next() {
		a := &dto.AlertaStockResponse{}
		if err := rows.Scan(
			&a.IDInventario, &a.IDProducto, &a.NombreProducto, &a.SKU,
			&a.IDSucursal, &a.StockActual, &a.StockMinimo, &a.Deficit, &a.PrecioCompra,
		); err != nil {
			return nil, fmt.Errorf("error al escanear alerta: %w", err)
		}
		alertas = append(alertas, a)
	}
	return alertas, nil
}
