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

type StoreProducto interface {
	GetAllProductos(ctx context.Context, params dto.PaginationParams) ([]*models.Producto, error)
	GetProductoByID(ctx context.Context, id uuid.UUID) (*models.Producto, error)
	CreateProducto(ctx context.Context, producto *models.Producto) (*models.Producto, error)
	UpdateProducto(ctx context.Context, id uuid.UUID, producto *models.Producto) (*models.Producto, error)
	DeleteProducto(ctx context.Context, id uuid.UUID) error
}

type storeProducto struct {
	db *sql.DB
}

func NewProducto(db *sql.DB) StoreProducto {
	return &storeProducto{db: db}
}

func (s *storeProducto) GetAllProductos(ctx context.Context, params dto.PaginationParams) ([]*models.Producto, error) {
	defer performance.Trace(ctx, "store", "GetAllProductos", performance.DbThreshold, time.Now())
	
	if params.Limit <= 0 {
		params.Limit = dto.DefaultLimit
	}

	query := `
	SELECT
		p.id_producto,
		p.nombre,
		p.descripcion,
		p.precio_compra,
		p.precio_venta,
		p.stock,
		p.fecha_vencimiento,
		p.imagen,
		p.id_status,
		p.id_sucursal,
		p.id_categoria,
		p.id_moneda,
		p.id_unidad,
		p.created_at,
		p.updated_at,

		su.id_sucursal,
		su.nombre_sucursal,
		su.id_status,

		c.id_categoria,
		c.nombre,

		m.id_moneda,
		m.nombre,
		m.id_status,

		u.id_unidad,
		u.nombre
	FROM producto p
	LEFT JOIN sucursal su ON su.id_sucursal = p.id_sucursal
	LEFT JOIN categoria c ON c.id_categoria = p.id_categoria
	LEFT JOIN moneda m ON m.id_moneda = p.id_moneda
	LEFT JOIN unidad u ON u.id_unidad = p.id_unidad
	WHERE p.deleted_at IS NULL
	`
	
	var args []interface{}
	
	if params.LastDate != nil {
		query += " AND p.created_at < $1"
		args = append(args, params.LastDate)
	}

	query += " ORDER BY p.created_at DESC LIMIT $" + fmt.Sprint(len(args)+1)
	args = append(args, params.Limit)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error al obtener productos: %w", err)
	}
	defer rows.Close()

	var productos []*models.Producto

	for rows.Next() {
		p := &models.Producto{
			Sucursal:  &models.Sucursal{},
			Categoria: &models.Categoria{},
			Moneda:    &models.Moneda{},
			Unidad:    &models.Unidad{},
		}

		if err := rows.Scan(
			&p.IDProducto,
			&p.Nombre,
			&p.Descripcion,
			&p.PrecioCompra,
			&p.PrecioVenta,
			&p.Stock,
			&p.FechaVencimiento,
			&p.Imagen,
			&p.IDStatus,
			&p.IDSucursal,
			&p.IDCategoria,
			&p.IDMoneda,
			&p.IDUnidad,
			&p.CreatedAt,
			&p.UpdatedAt,

			&p.Sucursal.IDSucursal,
			&p.Sucursal.NombreSucursal,
			&p.Sucursal.IDStatus,

			&p.Categoria.IDCategoria,
			&p.Categoria.Nombre,

			&p.Moneda.IDMoneda,
			&p.Moneda.Nombre,
			&p.Moneda.IDStatus,

			&p.Unidad.IDUnidad,
			&p.Unidad.Nombre,
		); err != nil {
			return nil, fmt.Errorf("error al escanear producto: %w", err)
		}

		productos = append(productos, p)
	}

	return productos, nil
}

func (s *storeProducto) GetProductoByID(ctx context.Context, id uuid.UUID) (*models.Producto, error) {
	defer performance.Trace(ctx, "store", "GetProductoByID", performance.DbThreshold, time.Now())
	query := `
	SELECT
		p.id_producto,
		p.nombre,
		p.descripcion,
		p.precio_compra,
		p.precio_venta,
		p.stock,
		p.fecha_vencimiento,
		p.imagen,
		p.id_status,
		p.id_sucursal,
		p.id_categoria,
		p.id_moneda,
		p.id_unidad,
		p.created_at,
		p.updated_at,

		su.id_sucursal,
		su.nombre_sucursal,
		su.id_status,

		c.id_categoria,
		c.nombre,

		m.id_moneda,
		m.nombre,
		m.id_status,

		u.id_unidad,
		u.nombre
	FROM producto p
	LEFT JOIN sucursal su ON su.id_sucursal = p.id_sucursal
	LEFT JOIN categoria c ON c.id_categoria = p.id_categoria
	LEFT JOIN moneda m ON m.id_moneda = p.id_moneda
	LEFT JOIN unidad u ON u.id_unidad = p.id_unidad
	WHERE p.id_producto = $1
	  AND p.deleted_at IS NULL
	`

	p := &models.Producto{
		Sucursal:  &models.Sucursal{},
		Categoria: &models.Categoria{},
		Moneda:    &models.Moneda{},
		Unidad:    &models.Unidad{},
	}

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&p.IDProducto,
		&p.Nombre,
		&p.Descripcion,
		&p.PrecioCompra,
		&p.PrecioVenta,
		&p.Stock,
		&p.FechaVencimiento,
		&p.Imagen,
		&p.IDStatus,
		&p.IDSucursal,
		&p.IDCategoria,
		&p.IDMoneda,
		&p.IDUnidad,
		&p.CreatedAt,
		&p.UpdatedAt,

		&p.Sucursal.IDSucursal,
		&p.Sucursal.NombreSucursal,
		&p.Sucursal.IDStatus,

		&p.Categoria.IDCategoria,
		&p.Categoria.Nombre,

		&p.Moneda.IDMoneda,
		&p.Moneda.Nombre,
		&p.Moneda.IDStatus,

		&p.Unidad.IDUnidad,
		&p.Unidad.Nombre,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("producto con ID %s no encontrado", id)
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener producto: %w", err)
	}

	return p, nil
}

func (s *storeProducto) CreateProducto(ctx context.Context, producto *models.Producto) (*models.Producto, error) {
	defer performance.Trace(ctx, "store", "CreateProducto", performance.DbThreshold, time.Now())
	query := `
		INSERT INTO producto (nombre, descripcion, precio_compra, precio_venta, stock, fecha_vencimiento, imagen, id_status, id_sucursal, id_categoria, id_moneda, id_unidad)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id_producto
	`

	var id uuid.UUID
	err := s.db.QueryRowContext(ctx, query,
		producto.Nombre,
		producto.Descripcion,
		producto.PrecioCompra,
		producto.PrecioVenta,
		producto.Stock,
		producto.FechaVencimiento,
		producto.Imagen,
		producto.IDStatus,
		producto.IDSucursal,
		producto.IDCategoria,
		producto.IDMoneda,
		producto.IDUnidad,
	).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("error al crear producto: %w", err)
	}

	producto.IDProducto = id
	return producto, nil
}

func (s *storeProducto) UpdateProducto(ctx context.Context, id uuid.UUID, producto *models.Producto) (*models.Producto, error) {
	defer performance.Trace(ctx, "store", "UpdateProducto", performance.DbThreshold, time.Now())
	query := `
		UPDATE producto
		SET
			nombre = $1,
			descripcion = $2,
			precio_compra = $3,
			precio_venta = $4,
			stock = $5,
			fecha_vencimiento = $6,
			imagen = $7,
			id_status = $8,
			id_sucursal = $9,
			id_categoria = $10,
			id_moneda = $11,
			id_unidad = $12,
			updated_at = CURRENT_TIMESTAMP
		WHERE id_producto = $13
		  AND deleted_at IS NULL
		RETURNING
			id_producto,
			nombre,
			descripcion,
			precio_compra,
			precio_venta,
			stock,
			fecha_vencimiento,
			imagen,
			id_status,
			id_sucursal,
			id_categoria,
			id_moneda,
			id_unidad,
			created_at,
			updated_at
	`

	err := s.db.QueryRowContext(ctx, query,
		producto.Nombre,
		producto.Descripcion,
		producto.PrecioCompra,
		producto.PrecioVenta,
		producto.Stock,
		producto.FechaVencimiento,
		producto.Imagen,
		producto.IDStatus,
		producto.IDSucursal,
		producto.IDCategoria,
		producto.IDMoneda,
		producto.IDUnidad,
		id,
	).Scan(
		&producto.IDProducto,
		&producto.Nombre,
		&producto.Descripcion,
		&producto.PrecioCompra,
		&producto.PrecioVenta,
		&producto.Stock,
		&producto.FechaVencimiento,
		&producto.Imagen,
		&producto.IDStatus,
		&producto.IDSucursal,
		&producto.IDCategoria,
		&producto.IDMoneda,
		&producto.IDUnidad,
		&producto.CreatedAt,
		&producto.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("producto con ID %s no encontrado", id)
	}
	if err != nil {
		return nil, fmt.Errorf("error al actualizar producto: %w", err)
	}

	return producto, nil
}

func (s *storeProducto) DeleteProducto(ctx context.Context, id uuid.UUID) error {
	defer performance.Trace(ctx, "store", "DeleteProducto", performance.DbThreshold, time.Now())
	query := `UPDATE producto SET deleted_at = $1 WHERE id_producto = $2 AND deleted_at IS NULL`

	result, err := s.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("error al eliminar producto: %w", err)
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
