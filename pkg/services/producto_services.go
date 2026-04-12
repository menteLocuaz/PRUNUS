package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/prunus/pkg/dto"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
	"github.com/prunus/pkg/utils"
)

type ServiceProducto struct {
	store           store.StoreProducto
	inventarioStore store.StoreInventario
	cache           *utils.CacheManager
	logger          *slog.Logger
}

func NewServiceProducto(s store.StoreProducto, inv store.StoreInventario, c *utils.CacheManager, logger *slog.Logger) *ServiceProducto {
	return &ServiceProducto{
		store:           s,
		inventarioStore: inv,
		cache:           c,
		logger:          logger,
	}
}

const (
	cacheKeyProductosAll    = "productos:all:last:%s:limit:%d"
	cacheKeyProductoID      = "productos:id:%s"
	cacheKeyProductoByCode  = "productos:code:%s"
	cacheExpirationProducto = 1 * time.Hour
)

func (s *ServiceProducto) GetAllProductos(ctx context.Context, params dto.PaginationParams) ([]*models.Producto, error) {
	lastDateStr := "none"
	if params.LastDate != nil {
		lastDateStr = params.LastDate.Format(time.RFC3339)
	}

	key := fmt.Sprintf(cacheKeyProductosAll, lastDateStr, params.Limit)
	return utils.GetOrSet(ctx, s.cache, key, cacheExpirationProducto, func() ([]*models.Producto, error) {
		return s.store.GetAllProductos(ctx, params)
	})
}

func (s *ServiceProducto) GetProductoByID(ctx context.Context, id uuid.UUID) (*models.Producto, error) {
	key := fmt.Sprintf(cacheKeyProductoID, id.String())
	return utils.GetOrSet(ctx, s.cache, key, cacheExpirationProducto, func() (*models.Producto, error) {
		return s.store.GetProductoByID(ctx, id)
	})
}

func (s *ServiceProducto) GetProductoByCodigo(ctx context.Context, codigo string) (*models.Producto, error) {
	key := fmt.Sprintf(cacheKeyProductoByCode, codigo)
	return utils.GetOrSet(ctx, s.cache, key, cacheExpirationProducto, func() (*models.Producto, error) {
		return s.store.GetProductoByCodigo(ctx, codigo)
	})
}

// CreateProducto ahora es una operación coordinada entre Catálogo e Inventario
func (s *ServiceProducto) CreateProducto(ctx context.Context, req dto.ProductoCreateRequest) (*models.Producto, error) {
	if req.Nombre == "" {
		return nil, errors.New("falta el nombre del producto")
	}

	// 1. Convertir DTO a Modelo usando el helper centralizado
	producto := req.ToModel()

	// 2. Lógica de negocio adicional (estatus por defecto)
	if producto.IDStatus == uuid.Nil {
		// Estatus 'Disponible' para Módulo Producto (ID según catálogo)
		producto.IDStatus = uuid.MustParse("31f4e127-e7e1-414d-aaef-6e92e4c5d970")
	}

	// 3. Crear Producto en Catálogo Maestro
	res, err := s.store.CreateProducto(ctx, &producto)
	if err != nil {
		return nil, fmt.Errorf("error al crear catálogo de producto: %w", err)
	}

	// 4. Crear Inventario Inicial para la Sucursal enviada
	inv := &models.Inventario{
		IDProducto:   res.IDProducto,
		IDSucursal:   req.IDSucursal,
		StockActual:  float64(req.Stock),
		PrecioCompra: req.PrecioCompra,
		PrecioVenta:  req.PrecioVenta,
	}

	_, err = s.inventarioStore.CreateInventario(ctx, inv)
	if err != nil {
		s.logger.ErrorContext(ctx, "Catálogo creado pero falló inicialización de inventario",
			slog.String("id_producto", res.IDProducto.String()),
			slog.Any("error", err),
		)
		// Retornamos el producto pero informamos del error de inventario
		return res, fmt.Errorf("producto creado pero sin stock inicial: %w", err)
	}

	s.logger.InfoContext(ctx, "Producto e inventario creados exitosamente",
		slog.String("id_producto", res.IDProducto.String()),
		slog.String("id_sucursal", req.IDSucursal.String()),
	)

	// Invalidar caché
	s.cache.Invalidate(ctx, "productos:")

	// 5. Retornar el producto completamente poblado (con relaciones)
	return s.store.GetProductoByID(ctx, res.IDProducto)
}

func (s *ServiceProducto) UpdateProducto(ctx context.Context, id uuid.UUID, req dto.ProductoUpdateRequest) (*models.Producto, error) {
	if req.Nombre == "" {
		return nil, errors.New("falta el nombre del producto")
	}

	// Convertir DTO a Modelo
	producto := req.ToModel()
	producto.IDProducto = id

	res, err := s.store.UpdateProducto(ctx, id, &producto)
	if err != nil {
		return nil, err
	}

	// Invalidar caché
	s.cache.Invalidate(ctx, "productos:")

	// Retornar el producto completamente poblado (con relaciones)
	return s.store.GetProductoByID(ctx, res.IDProducto)
}

func (s *ServiceProducto) DeleteProducto(ctx context.Context, id uuid.UUID) error {
	err := s.store.DeleteProducto(ctx, id)
	if err != nil {
		return err
	}

	// Invalidar caché
	s.cache.Invalidate(ctx, "productos:")

	return nil
}
