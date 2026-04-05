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
)

type ServiceProducto struct {
	store           store.StoreProducto
	inventarioStore store.StoreInventario
	logger          *slog.Logger
}

func NewServiceProducto(s store.StoreProducto, inv store.StoreInventario, logger *slog.Logger) *ServiceProducto {
	return &ServiceProducto{
		store:           s,
		inventarioStore: inv,
		logger:          logger,
	}
}

func (s *ServiceProducto) GetAllProductos(ctx context.Context, params dto.PaginationParams) ([]*models.Producto, error) {
	return s.store.GetAllProductos(ctx, params)
}

func (s *ServiceProducto) GetProductoByID(ctx context.Context, id uuid.UUID) (*models.Producto, error) {
	return s.store.GetProductoByID(ctx, id)
}

func (s *ServiceProducto) GetProductoByCodigo(ctx context.Context, codigo string) (*models.Producto, error) {
	return s.store.GetProductoByCodigo(ctx, codigo)
}

// CreateProducto ahora es una operación coordinada entre Catálogo e Inventario
func (s *ServiceProducto) CreateProducto(ctx context.Context, req dto.ProductoCreateRequest) (*models.Producto, error) {
	if req.Nombre == "" {
		return nil, errors.New("falta el nombre del producto")
	}

	// 1. Mapear DTO a Modelo de Catálogo
	// Asignar estatus automático si no se envía
	idStatus := req.IDStatus
	if idStatus == uuid.Nil {
		// Estatus 'Disponible' para Módulo Producto (ID según catálogo)
		idStatus = uuid.MustParse("31f4e127-e7e1-414d-aaef-6e92e4c5d970")
	}

	// Manejar fecha de vencimiento opcional (si es zero value en Go, se guarda como nil en BD)
	var fechaVencimiento *time.Time
	if !req.FechaVencimiento.IsZero() {
		fechaVencimiento = &req.FechaVencimiento
	}

	producto := &models.Producto{
		Nombre:           req.Nombre,
		Descripcion:      req.Descripcion,
		CodigoBarras:     req.CodigoBarras,
		SKU:              req.SKU,
		FechaVencimiento: fechaVencimiento,
		Imagen:           req.Imagen,
		IDStatus:         idStatus,
		IDCategoria:      req.IDCategoria,
		IDMoneda:         req.IDMoneda,
		IDUnidad:         req.IDUnidad,
	}

	// 2. Crear Producto en Catálogo Maestro
	res, err := s.store.CreateProducto(ctx, producto)
	if err != nil {
		return nil, fmt.Errorf("error al crear catálogo de producto: %w", err)
	}

	// 3. Crear Inventario Inicial para la Sucursal enviada
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

	// 4. Retornar el producto completamente poblado (con relaciones) para consistencia con el frontend
	return s.store.GetProductoByID(ctx, res.IDProducto)
}

func (s *ServiceProducto) UpdateProducto(ctx context.Context, id uuid.UUID, req dto.ProductoUpdateRequest) (*models.Producto, error) {
	if req.Nombre == "" {
		return nil, errors.New("falta el nombre del producto")
	}

	// Manejar fecha de vencimiento opcional
	var fechaVencimiento *time.Time
	if !req.FechaVencimiento.IsZero() {
		fechaVencimiento = &req.FechaVencimiento
	}

	// Actualizar datos maestros
	producto := &models.Producto{
		Nombre:           req.Nombre,
		Descripcion:      req.Descripcion,
		CodigoBarras:     req.CodigoBarras,
		SKU:              req.SKU,
		FechaVencimiento: fechaVencimiento,
		Imagen:           req.Imagen,
		IDStatus:         req.IDStatus,
		IDCategoria:      req.IDCategoria,
		IDMoneda:         req.IDMoneda,
		IDUnidad:         req.IDUnidad,
	}

	res, err := s.store.UpdateProducto(ctx, id, producto)
	if err != nil {
		return nil, err
	}

	// Retornar el producto completamente poblado (con relaciones)
	return s.store.GetProductoByID(ctx, res.IDProducto)
}

func (s *ServiceProducto) DeleteProducto(ctx context.Context, id uuid.UUID) error {
	return s.store.DeleteProducto(ctx, id)
}
