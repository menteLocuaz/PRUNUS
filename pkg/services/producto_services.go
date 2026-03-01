package services

import (
	"errors"

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

func (s *ServiceProducto) GetProductoByID(id uint) (*models.Producto, error) {
	return s.store.GetProductoByID(id)
}

func (s *ServiceProducto) CreateProducto(producto models.Producto) (*models.Producto, error) {
	if producto.Nombre == "" {
		return nil, errors.New("falta el nombre del producto")
	}
	if producto.IDSucursal == 0 {
		return nil, errors.New("falta el id de la sucursal")
	}
	if producto.IDCategoria == 0 {
		return nil, errors.New("falta el id de la categoria")
	}
	if producto.IDMoneda == 0 {
		return nil, errors.New("falta el id de la moneda")
	}
	if producto.IDUnidad == 0 {
		return nil, errors.New("falta el id de la unidad")
	}
	return s.store.CreateProducto(&producto)
}

func (s *ServiceProducto) UpdateProducto(id uint, producto models.Producto) (*models.Producto, error) {
	if producto.Nombre == "" {
		return nil, errors.New("falta el nombre del producto")
	}
	return s.store.UpdateProducto(id, &producto)
}

func (s *ServiceProducto) DeleteProducto(id uint) error {
	return s.store.DeleteProducto(id)
}
