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

// StoreProducto define las operaciones de persistencia para el catálogo de productos.
type StoreProducto interface {
	GetAllProductos(ctx context.Context, params dto.PaginationParams) ([]*models.Producto, error)
	GetProductoByID(ctx context.Context, id uuid.UUID) (*models.Producto, error)
	GetProductoByCodigo(ctx context.Context, codigo string) (*models.Producto, error)
	CreateProducto(ctx context.Context, producto *models.Producto) (*models.Producto, error)
	UpdateProducto(ctx context.Context, id uuid.UUID, producto *models.Producto) (*models.Producto, error)
	DeleteProducto(ctx context.Context, id uuid.UUID) error
}

type storeProducto struct {
	db *sql.DB
}

// Campos base para SELECT de producto mapeados correctamente al esquema DB.
// Se eliminan id_moneda e id_unidad porque no existen en la tabla producto (están en inventario o lotes).
const productoSelectFields = `
	p.id_producto, p.pro_nombre, COALESCE(p.pro_descripcion, ''), COALESCE(p.pro_codigo, ''),
	COALESCE(p.sku, ''), p.id_status, p.id_categoria, p.created_at, p.updated_at,
	
	c.id_categoria, c.nombre
`

// scanRowProducto centraliza el escaneo de resultados para mantener consistencia.
func (s *storeProducto) scanRowProducto(scanner interface{ Scan(dest ...any) error }, p *models.Producto) error {
	if p.Categoria == nil {
		p.Categoria = &models.Categoria{}
	}

	return scanner.Scan(
		&p.IDProducto, &p.Nombre, &p.Descripcion, &p.CodigoBarras,
		&p.SKU, &p.IDStatus, &p.IDCategoria, &p.CreatedAt, &p.UpdatedAt,
		&p.Categoria.IDCategoria, &p.Categoria.Nombre,
	)
}

// NewProducto crea una nueva instancia del store de productos.
func NewProducto(db *sql.DB) StoreProducto {
	return &storeProducto{db: db}
}

// GetAllProductos obtiene la lista de productos con soporte para paginación.
func (s *storeProducto) GetAllProductos(ctx context.Context, params dto.PaginationParams) ([]*models.Producto, error) {
	defer performance.Trace(ctx, "store", "GetAllProductos", performance.DbThreshold, time.Now())

	if params.Limit <= 0 {
		params.Limit = dto.DefaultLimit
	}

	query := fmt.Sprintf(`
		SELECT %s
		FROM producto p
		LEFT JOIN categoria c ON c.id_categoria = p.id_categoria
		WHERE p.deleted_at IS NULL
	`, productoSelectFields)

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
		p := &models.Producto{}
		if err := s.scanRowProducto(rows, p); err != nil {
			return nil, fmt.Errorf("error al escanear producto: %w", err)
		}
		productos = append(productos, p)
	}

	return productos, nil
}

// GetProductoByID busca un producto específico por su identificador único.
func (s *storeProducto) GetProductoByID(ctx context.Context, id uuid.UUID) (*models.Producto, error) {
	defer performance.Trace(ctx, "store", "GetProductoByID", performance.DbThreshold, time.Now())
	query := fmt.Sprintf(`
		SELECT %s
		FROM producto p
		LEFT JOIN categoria c ON c.id_categoria = p.id_categoria
		WHERE p.id_producto = $1 AND p.deleted_at IS NULL
	`, productoSelectFields)

	p := &models.Producto{}
	err := s.scanRowProducto(s.db.QueryRowContext(ctx, query, id), p)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("producto con ID %s no encontrado", id)
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener producto: %w", err)
	}

	return p, nil
}

// GetProductoByCodigo busca un producto por pro_codigo o sku.
func (s *storeProducto) GetProductoByCodigo(ctx context.Context, codigo string) (*models.Producto, error) {
	defer performance.Trace(ctx, "store", "GetProductoByCodigo", performance.DbThreshold, time.Now())
	query := fmt.Sprintf(`
		SELECT %s
		FROM producto p
		LEFT JOIN categoria c ON c.id_categoria = p.id_categoria
		WHERE (p.pro_codigo = $1 OR p.sku = $1) AND p.deleted_at IS NULL
	`, productoSelectFields)

	p := &models.Producto{}
	err := s.scanRowProducto(s.db.QueryRowContext(ctx, query, codigo), p)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("producto con código %s no encontrado", codigo)
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener producto por código: %w", err)
	}

	return p, nil
}

// CreateProducto inserta un nuevo registro en la tabla producto utilizando ExecAudited.
func (s *storeProducto) CreateProducto(ctx context.Context, producto *models.Producto) (*models.Producto, error) {
	defer performance.Trace(ctx, "store", "CreateProducto", performance.DbThreshold, time.Now())

	err := ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		query := `
			INSERT INTO producto (pro_nombre, pro_descripcion, pro_codigo, sku, id_status, id_categoria)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id_producto, created_at, updated_at
		`
		return tx.QueryRowContext(ctx, query,
			producto.Nombre, producto.Descripcion, producto.CodigoBarras, producto.SKU,
			producto.IDStatus, producto.IDCategoria,
		).Scan(&producto.IDProducto, &producto.CreatedAt, &producto.UpdatedAt)
	})

	if err != nil {
		return nil, fmt.Errorf("error al crear producto: %w", err)
	}

	return producto, nil
}

// UpdateProducto actualiza la información de un producto existente.
func (s *storeProducto) UpdateProducto(ctx context.Context, id uuid.UUID, producto *models.Producto) (*models.Producto, error) {
	defer performance.Trace(ctx, "store", "UpdateProducto", performance.DbThreshold, time.Now())

	err := ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		query := `
			UPDATE producto
			SET
				pro_nombre = $1, pro_descripcion = $2, pro_codigo = $3, sku = $4,
				id_status = $5, id_categoria = $6, 
				updated_at = CURRENT_TIMESTAMP
			WHERE id_producto = $7 AND deleted_at IS NULL
			RETURNING created_at, updated_at
		`
		return tx.QueryRowContext(ctx, query,
			producto.Nombre, producto.Descripcion, producto.CodigoBarras, producto.SKU,
			producto.IDStatus, producto.IDCategoria, id,
		).Scan(&producto.CreatedAt, &producto.UpdatedAt)
	})

	if err != nil {
		return nil, fmt.Errorf("error al actualizar producto: %w", err)
	}

	producto.IDProducto = id
	return producto, nil
}

// DeleteProducto realiza un borrado lógico (soft delete) del producto.
func (s *storeProducto) DeleteProducto(ctx context.Context, id uuid.UUID) error {
	defer performance.Trace(ctx, "store", "DeleteProducto", performance.DbThreshold, time.Now())

	err := ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		query := `UPDATE producto SET deleted_at = CURRENT_TIMESTAMP WHERE id_producto = $1 AND deleted_at IS NULL`
		result, err := tx.ExecContext(ctx, query, id)
		if err != nil {
			return err
		}
		rows, _ := result.RowsAffected()
		if rows == 0 {
			return sql.ErrNoRows
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error al eliminar producto: %w", err)
	}

	return nil
}
