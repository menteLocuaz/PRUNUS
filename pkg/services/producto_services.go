package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/prunus/pkg/dto"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
	"github.com/prunus/pkg/utils"
	zaplogger "github.com/prunus/pkg/utils/logger"
	"go.uber.org/zap"
)

type ServiceProducto struct {
	store           store.StoreProducto
	inventarioStore store.StoreInventario
	cache           *utils.CacheManager
	logger          *zap.Logger
}

func NewServiceProducto(s store.StoreProducto, inv store.StoreInventario, c *utils.CacheManager, logger *zap.Logger) *ServiceProducto {
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

// CreateProducto es una operación coordinada entre Catálogo e Inventario
func (s *ServiceProducto) CreateProducto(ctx context.Context, req dto.ProductoCreateRequest) (*models.Producto, error) {
	if req.Nombre == "" {
		return nil, errors.New("falta el nombre del producto")
	}

	producto := req.ToModel()

	if producto.IDStatus == uuid.Nil {
		producto.IDStatus = uuid.MustParse("31f4e127-e7e1-414d-aaef-6e92e4c5d970")
	}

	res, err := s.store.CreateProducto(ctx, &producto)
	if err != nil {
		return nil, fmt.Errorf("error al crear catálogo de producto: %w", err)
	}

	inv := &models.Inventario{
		IDProducto:   res.IDProducto,
		IDSucursal:   req.IDSucursal,
		StockActual:  float64(req.Stock),
		PrecioCompra: req.PrecioCompra,
		PrecioVenta:  req.PrecioVenta,
	}

	_, err = s.inventarioStore.CreateInventario(ctx, inv)
	if err != nil {
		zaplogger.WithContext(ctx, s.logger).Error("Catálogo creado pero falló inicialización de inventario",
			zap.String("id_producto", res.IDProducto.String()),
			zap.Error(err),
		)
		return res, fmt.Errorf("producto creado pero sin stock inicial: %w", err)
	}

	zaplogger.WithContext(ctx, s.logger).Info("Producto e inventario creados exitosamente",
		zap.String("id_producto", res.IDProducto.String()),
		zap.String("id_sucursal", req.IDSucursal.String()),
	)

	s.cache.Invalidate(ctx, "productos:")

	return s.store.GetProductoByID(ctx, res.IDProducto)
}

func (s *ServiceProducto) UpdateProducto(ctx context.Context, id uuid.UUID, req dto.ProductoUpdateRequest) (*models.Producto, error) {
	if req.Nombre == "" {
		return nil, errors.New("falta el nombre del producto")
	}

	producto := req.ToModel()
	producto.IDProducto = id

	res, err := s.store.UpdateProducto(ctx, id, &producto)
	if err != nil {
		return nil, err
	}

	s.cache.Invalidate(ctx, "productos:")

	return s.store.GetProductoByID(ctx, res.IDProducto)
}

func (s *ServiceProducto) DeleteProducto(ctx context.Context, id uuid.UUID) error {
	err := s.store.DeleteProducto(ctx, id)
	if err != nil {
		return err
	}

	s.cache.Invalidate(ctx, "productos:")

	return nil
}
