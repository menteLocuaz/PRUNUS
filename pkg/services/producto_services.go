package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/prunus/pkg/dto"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
)

type ServiceProducto struct {
	store  store.StoreProducto
	logger *slog.Logger
}

func NewServiceProducto(s store.StoreProducto, logger *slog.Logger) *ServiceProducto {
	return &ServiceProducto{
		store:  s,
		logger: logger,
	}
}

func (s *ServiceProducto) GetAllProductos(ctx context.Context, params dto.PaginationParams) ([]*models.Producto, error) {
	return s.store.GetAllProductos(ctx, params)
}

func (s *ServiceProducto) GetProductoByID(ctx context.Context, id uuid.UUID) (*models.Producto, error) {
	return s.store.GetProductoByID(ctx, id)
}

func (s *ServiceProducto) CreateProducto(ctx context.Context, producto models.Producto) (*models.Producto, error) {
	if err := s.validateProducto(&producto); err != nil {
		s.logger.WarnContext(ctx, "Fallo de validación al crear producto",
			slog.String("nombre", producto.Nombre),
			slog.Any("error", err),
		)
		return nil, err
	}

	res, err := s.store.CreateProducto(ctx, &producto)
	if err != nil {
		return nil, err
	}

	// Log preventivo: Producto creado con stock inicial bajo
	if res.Stock < 5 {
		s.logger.WarnContext(ctx, "Producto creado con stock bajo",
			slog.String("id_producto", res.IDProducto.String()),
			slog.String("nombre", res.Nombre),
			slog.Uint64("stock", uint64(res.Stock)),
		)
	}

	return res, nil
}

func (s *ServiceProducto) UpdateProducto(ctx context.Context, id uuid.UUID, producto models.Producto) (*models.Producto, error) {
	if producto.Nombre == "" {
		s.logger.WarnContext(ctx, "Intento de actualización con nombre vacío",
			slog.String("id_producto", id.String()),
		)
		return nil, errors.New("falta el nombre del producto")
	}

	res, err := s.store.UpdateProducto(ctx, id, &producto)
	if err != nil {
		return nil, err
	}

	// Log preventivo: Alerta de stock crítico tras actualización
	if res.Stock == 0 {
		s.logger.WarnContext(ctx, "Producto con stock agotado tras actualización",
			slog.String("id_producto", id.String()),
			slog.String("nombre", res.Nombre),
		)
	}

	return res, nil
}

func (s *ServiceProducto) DeleteProducto(ctx context.Context, id uuid.UUID) error {
	return s.store.DeleteProducto(ctx, id)
}

func (s *ServiceProducto) validateProducto(p *models.Producto) error {
	if p.Nombre == "" {
		return errors.New("falta el nombre del producto")
	}
	if p.IDSucursal == uuid.Nil {
		return errors.New("falta el id de la sucursal")
	}
	if p.IDCategoria == uuid.Nil {
		return errors.New("falta el id de la categoria")
	}
	if p.IDMoneda == uuid.Nil {
		return errors.New("falta el id de la moneda")
	}
	if p.IDUnidad == uuid.Nil {
		return errors.New("falta el id de la unidad")
	}
	if p.IDStatus == uuid.Nil {
		return errors.New("falta el id del estatus")
	}
	return nil
}
