package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/utils/performance"
)

type StoreInventario interface {
	GetAllInventario(ctx context.Context) ([]*models.Inventario, error)
	GetInventarioByID(ctx context.Context, id uuid.UUID) (*models.Inventario, error)
	GetInventarioByProductoYSucursal(ctx context.Context, idProducto, idSucursal uuid.UUID) (*models.Inventario, error)
	CreateInventario(ctx context.Context, inventario *models.Inventario) (*models.Inventario, error)
	UpdateInventario(ctx context.Context, id uuid.UUID, inventario *models.Inventario) (*models.Inventario, error)
	DeleteInventario(ctx context.Context, id uuid.UUID) error

	// Movimientos
	CreateMovimiento(ctx context.Context, m *models.MovimientoInventario) (*models.MovimientoInventario, error)
	GetMovimientosByProducto(ctx context.Context, productoID uuid.UUID) ([]*models.MovimientoInventario, error)
}

type storeInventario struct {
	db *sql.DB
}

func NewInventario(db *sql.DB) StoreInventario {
	return &storeInventario{db: db}
}

func (s *storeInventario) GetAllInventario(ctx context.Context) ([]*models.Inventario, error) {
	defer performance.Trace(ctx, "store", "GetAllInventario", performance.DbThreshold, time.Now())
	query := `
		SELECT 
			id_inventario, id_producto, id_sucursal, stock_actual, stock_minimo, 
			stock_maximo, precio_compra, precio_venta, created_at, updated_at
		FROM inventario
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query)
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

func (s *storeInventario) CreateMovimiento(ctx context.Context, m *models.MovimientoInventario) (*models.MovimientoInventario, error) {
	query := `INSERT INTO movimientos_inventario (id_producto, tipo_movimiento, cantidad, fecha, id_usuario, referencia) 
	          VALUES ($1, $2, $3, $4, $5, $6) RETURNING id_movimiento, created_at, updated_at`
	err := s.db.QueryRowContext(ctx, query, m.IDProducto, m.TipoMovimiento, m.Cantidad, time.Now(), m.IDUsuario, m.Referencia).
		Scan(&m.IDMovimiento, &m.CreatedAt, &m.UpdatedAt)
	m.Fecha = time.Now()
	return m, err
}

func (s *storeInventario) GetMovimientosByProducto(ctx context.Context, productoID uuid.UUID) ([]*models.MovimientoInventario, error) {
	query := `SELECT id_movimiento, id_producto, tipo_movimiento, cantidad, fecha, id_usuario, referencia FROM movimientos_inventario WHERE id_producto = $1 AND deleted_at IS NULL ORDER BY fecha DESC`
	rows, err := s.db.QueryContext(ctx, query, productoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movimientos []*models.MovimientoInventario
	for rows.Next() {
		m := &models.MovimientoInventario{}
		if err := rows.Scan(&m.IDMovimiento, &m.IDProducto, &m.TipoMovimiento, &m.Cantidad, &m.Fecha, &m.IDUsuario, &m.Referencia); err != nil {
			return nil, err
		}
		movimientos = append(movimientos, m)
	}
	return movimientos, nil
}
