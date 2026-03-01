package services

import (
	"errors"

	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
)

type ServiceProveedor struct {
	store store.StoreProveedor
}

func NewServiceProveedor(s store.StoreProveedor) *ServiceProveedor {
	return &ServiceProveedor{store: s}
}

func (s *ServiceProveedor) GetAllProveedores() ([]*models.Proveedor, error) {
	return s.store.GetAllProveedores()
}

func (s *ServiceProveedor) GetProveedorByID(id uint) (*models.Proveedor, error) {
	return s.store.GetProveedorByID(id)
}

func (s *ServiceProveedor) CreateProveedor(proveedor models.Proveedor) (*models.Proveedor, error) {
	if proveedor.Nombre == "" {
		return nil, errors.New("falta el nombre del proveedor")
	}
	if proveedor.IDSucursal == 0 {
		return nil, errors.New("falta el id de la sucursal")
	}
	if proveedor.IDEmpresa == 0 {
		return nil, errors.New("falta el id de la empresa")
	}
	return s.store.CreateProveedor(&proveedor)
}

func (s *ServiceProveedor) UpdateProveedor(id uint, proveedor models.Proveedor) (*models.Proveedor, error) {
	if proveedor.Nombre == "" {
		return nil, errors.New("falta el nombre del proveedor")
	}
	return s.store.UpdateProveedor(id, &proveedor)
}

func (s *ServiceProveedor) DeleteProveedor(id uint) error {
	return s.store.DeleteProveedor(id)
}
