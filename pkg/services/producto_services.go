package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
)

type ServiceProducto struct {
	store store.StoreProducto
}

func NewServiceProducto(s store.StoreProducto) *ServiceProducto {
	return &ServiceProducto{store: s}
}

func (s *ServiceProducto) GetAllProductos() ([]*models.Producto, error) {
	return s.store.GetAllProductos()
}

func (s *ServiceProducto) GetProductoByID(id uuid.UUID) (*models.Producto, error) {
	return s.store.GetProductoByID(id)
}

func (s *ServiceProducto) CreateProducto(producto models.Producto) (*models.Producto, error) {
	if producto.Nombre == "" {
		return nil, errors.New("falta el nombre del producto")
	}
	if producto.IDSucursal == uuid.Nil {
		return nil, errors.New("falta el id de la sucursal")
	}
	if producto.IDCategoria == uuid.Nil {
		return nil, errors.New("falta el id de la categoria")
	}
	if producto.IDMoneda == uuid.Nil {
		return nil, errors.New("falta el id de la moneda")
	}
	if producto.IDUnidad == uuid.Nil {
		return nil, errors.New("falta el id de la unidad")
	}
	if producto.IDStatus == uuid.Nil {
		return nil, errors.New("falta el id del estatus")
	}
	return s.store.CreateProducto(&producto)
}

func (s *ServiceProducto) UpdateProducto(id uuid.UUID, producto models.Producto) (*models.Producto, error) {
	if producto.Nombre == "" {
		return nil, errors.New("falta el nombre del producto")
	}
	return s.store.UpdateProducto(id, &producto)
}

func (s *ServiceProducto) DeleteProducto(id uuid.UUID) error {
	return s.store.DeleteProducto(id)
}
