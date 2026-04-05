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
}

type storeInventario struct {
	db *sql.DB
}

func NewInventario(db *sql.DB) StoreInventario {
	return &storeInventario{db: db}
}

func (s *storeInventario) GetAllInventario(ctx context.Context, params dto.PaginationParams) ([]*models.Inventario, error) {
	defer performance.Trace(ctx, "store", "GetAllInventario", performance.DbThreshold, time.Now())

	if params.Limit <= 0 {
		params.Limit = dto.DefaultLimit
	}

	query := `
		SELECT 
			id_inventario, id_producto, id_sucursal, stock_actual, stock_minimo, 
			stock_maximo, precio_compra, precio_venta, created_at, updated_at
		FROM inventario
		WHERE deleted_at IS NULL
	`

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
		if err := rows.Scan(
			&i.IDInventario, &i.IDProducto, &i.IDSucursal, &i.StockActual, &i.StockMinimo,
			&i.StockMaximo, &i.PrecioCompra, &i.PrecioVenta, &i.CreatedAt, &i.UpdatedAt,
		); err != nil {
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

	query := `
		SELECT 
			id_inventario, id_producto, id_sucursal, stock_actual, stock_minimo, 
			stock_maximo, precio_compra, precio_venta, created_at, updated_at
		FROM inventario
		WHERE id_sucursal = $1 AND deleted_at IS NULL
	`

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
		if err := rows.Scan(
			&i.IDInventario, &i.IDProducto, &i.IDSucursal, &i.StockActual, &i.StockMinimo,
			&i.StockMaximo, &i.PrecioCompra, &i.PrecioVenta, &i.CreatedAt, &i.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error al escanear inventario: %w", err)
		}
		inventarios = append(inventarios, i)
	}

	return inventarios, nil
}

func (s *storeInventario) GetInventarioByID(ctx context.Context, id uuid.UUID) (*models.Inventario, error) {
	defer performance.Trace(ctx, "store", "GetInventarioByID", performance.DbThreshold, time.Now())
	query := `
		SELECT 
			id_inventario, id_producto, id_sucursal, stock_actual, stock_minimo, 
			stock_maximo, precio_compra, precio_venta, created_at, updated_at
		FROM inventario
		WHERE id_inventario = $1 AND deleted_at IS NULL
	`

	i := &models.Inventario{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&i.IDInventario, &i.IDProducto, &i.IDSucursal, &i.StockActual, &i.StockMinimo,
		&i.StockMaximo, &i.PrecioCompra, &i.PrecioVenta, &i.CreatedAt, &i.UpdatedAt,
	)

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
	query := `
		SELECT 
			id_inventario, id_producto, id_sucursal, stock_actual, stock_minimo, 
			stock_maximo, precio_compra, precio_venta, created_at, updated_at
		FROM inventario
		WHERE id_producto = $1 AND id_sucursal = $2 AND deleted_at IS NULL
	`

	i := &models.Inventario{}
	err := s.db.QueryRowContext(ctx, query, idProducto, idSucursal).Scan(
		&i.IDInventario, &i.IDProducto, &i.IDSucursal, &i.StockActual, &i.StockMinimo,
		&i.StockMaximo, &i.PrecioCompra, &i.PrecioVenta, &i.CreatedAt, &i.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // No es un error, simplemente no existe
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener inventario por producto y sucursal: %w", err)
	}

	return i, nil
}

func (s *storeInventario) CreateInventario(ctx context.Context, inventario *models.Inventario) (*models.Inventario, error) {
	defer performance.Trace(ctx, "store", "CreateInventario", performance.DbThreshold, time.Now())
	query := `
		INSERT INTO inventario (
			id_producto, id_sucursal, stock_actual, stock_minimo, 
			stock_maximo, precio_compra, precio_venta
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id_inventario, created_at, updated_at
	`

	err := s.db.QueryRowContext(ctx, query,
		inventario.IDProducto, inventario.IDSucursal, inventario.StockActual, inventario.StockMinimo,
		inventario.StockMaximo, inventario.PrecioCompra, inventario.PrecioVenta,
	).Scan(&inventario.IDInventario, &inventario.CreatedAt, &inventario.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("error al crear inventario: %w", err)
	}

	return inventario, nil
}

func (s *storeInventario) UpdateInventario(ctx context.Context, id uuid.UUID, inventario *models.Inventario) (*models.Inventario, error) {
	defer performance.Trace(ctx, "store", "UpdateInventario", performance.DbThreshold, time.Now())
	query := `
		UPDATE inventario
		SET 
			stock_actual = $1, stock_minimo = $2, stock_maximo = $3, 
			precio_compra = $4, precio_venta = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id_inventario = $6 AND deleted_at IS NULL
		RETURNING id_inventario, id_producto, id_sucursal, created_at, updated_at
	`

	err := s.db.QueryRowContext(ctx, query,
		inventario.StockActual, inventario.StockMinimo, inventario.StockMaximo,
		inventario.PrecioCompra, inventario.PrecioVenta, id,
	).Scan(&inventario.IDInventario, &inventario.IDProducto, &inventario.IDSucursal, &inventario.CreatedAt, &inventario.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("inventario con ID %s no encontrado", id)
	}
	if err != nil {
		return nil, fmt.Errorf("error al actualizar inventario: %w", err)
	}

	return inventario, nil
}

func (s *storeInventario) DeleteInventario(ctx context.Context, id uuid.UUID) error {
	defer performance.Trace(ctx, "store", "DeleteInventario", performance.DbThreshold, time.Now())
	query := `UPDATE inventario SET deleted_at = CURRENT_TIMESTAMP WHERE id_inventario = $1 AND deleted_at IS NULL`

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error al eliminar inventario: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s *storeInventario) RegistrarMovimiento(ctx context.Context, m *models.MovimientoInventario) (*models.MovimientoInventario, error) {
	defer performance.Trace(ctx, "store", "RegistrarMovimiento", performance.DbThreshold, time.Now())

	// Creamos el JSON de items para la función almacenada.
	// Aunque el modelo actual solo tiene un producto, la función soporta múltiples.
	itemsJSON := fmt.Sprintf(`[{"id_producto": "%s", "cantidad": %f}]`, m.IDProducto, m.Cantidad)

	query := `SELECT id_movimiento, id_producto, stock_anterior, cantidad, stock_posterior 
	          FROM inventario_ia_movimiento($1, $2, $3, $4, $5)`

	err := s.db.QueryRowContext(ctx, query,
		m.IDSucursal, m.IDUsuario, m.TipoMovimiento, m.Referencia, itemsJSON,
	).Scan(&m.IDMovimiento, &m.IDProducto, &m.StockAnterior, &m.Cantidad, &m.StockPosterior)

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

	query := `SELECT id_movimiento, id_producto, stock_anterior, cantidad, stock_posterior 
	          FROM inventario_ia_movimiento($1, $2, $3, $4, $5)`

	rows, err := s.db.QueryContext(ctx, query, idSucursal, idUsuario, tipoMov, referencia, itemsJSON)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar movimiento masivo: %w", err)
	}
	defer rows.Close()

	var resultados []*models.MovimientoInventario
	for rows.Next() {
		m := &models.MovimientoInventario{
			IDSucursal:     idSucursal,
			IDUsuario:      idUsuario,
			TipoMovimiento: tipoMov,
			Referencia:     referencia,
			Fecha:          time.Now(),
		}
		if err := rows.Scan(&m.IDMovimiento, &m.IDProducto, &m.StockAnterior, &m.Cantidad, &m.StockPosterior); err != nil {
			return nil, fmt.Errorf("error al escanear resultado de movimiento: %w", err)
		}
		m.CreatedAt = m.Fecha
		m.UpdatedAt = m.Fecha
		resultados = append(resultados, m)
	}

	return resultados, nil
}

func (s *storeInventario) GetMovimientosByProducto(ctx context.Context, productoID uuid.UUID, params dto.PaginationParams) ([]*models.MovimientoInventario, error) {
	defer performance.Trace(ctx, "store", "GetMovimientosByProducto", performance.DbThreshold, time.Now())

	if params.Limit <= 0 {
		params.Limit = dto.DefaultLimit
	}

	query := `SELECT 
		id_movimiento, id_producto, id_sucursal, tipo_movimiento, cantidad, 
		costo_unitario, precio_unitario, stock_anterior, stock_posterior, 
		fecha, id_usuario, referencia 
	FROM movimientos_inventario 
	WHERE id_producto = $1 AND deleted_at IS NULL`

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
		if err := rows.Scan(
			&m.IDMovimiento, &m.IDProducto, &m.IDSucursal, &m.TipoMovimiento, &m.Cantidad,
			&m.CostoUnitario, &m.PrecioUnitario, &m.StockAnterior, &m.StockPosterior,
			&m.Fecha, &m.IDUsuario, &m.Referencia,
		); err != nil {
			return nil, err
		}
		movimientos = append(movimientos, m)
	}
	return movimientos, nil
}

func (s *storeInventario) GetAlertasStock(ctx context.Context, sucursalID uuid.UUID) ([]*models.Inventario, error) {
	defer performance.Trace(ctx, "store", "GetAlertasStock", performance.DbThreshold, time.Now())
	query := `
		SELECT 
			id_inventario, id_producto, id_sucursal, stock_actual, stock_minimo, 
			stock_maximo, precio_compra, precio_venta, created_at, updated_at
		FROM inventario
		WHERE id_sucursal = $1 AND stock_actual <= stock_minimo AND deleted_at IS NULL
	`

	rows, err := s.db.QueryContext(ctx, query, sucursalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alertas []*models.Inventario
	for rows.Next() {
		i := &models.Inventario{}
		if err := rows.Scan(
			&i.IDInventario, &i.IDProducto, &i.IDSucursal, &i.StockActual, &i.StockMinimo,
			&i.StockMaximo, &i.PrecioCompra, &i.PrecioVenta, &i.CreatedAt, &i.UpdatedAt,
		); err != nil {
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
		// Para PEPS y UEPS usaremos la tabla de lotes si está disponible,
		// o una aproximación basada en el historial de movimientos.
		// Por ahora, usaremos los lotes.
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
	query := `
		SELECT 
			id_lote, id_producto, id_sucursal, codigo_lote, cantidad_inicial, 
			cantidad_actual, costo_compra, fecha_vencimiento, fecha_recepcion,
			created_at, updated_at
		FROM lotes
		WHERE id_producto = $1 AND id_sucursal = $2 AND deleted_at IS NULL AND cantidad_actual > 0
		ORDER BY fecha_recepcion ASC
	`

	rows, err := s.db.QueryContext(ctx, query, idProducto, idSucursal)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lotes []*models.Lote
	for rows.Next() {
		l := &models.Lote{}
		if err := rows.Scan(
			&l.IDLote, &l.IDProducto, &l.IDSucursal, &l.CodigoLote, &l.CantidadInicial,
			&l.CantidadActual, &l.CostoCompra, &l.FechaVencimiento, &l.FechaRecepcion,
			&l.CreatedAt, &l.UpdatedAt,
		); err != nil {
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
