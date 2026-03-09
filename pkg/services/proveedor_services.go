package services

import (
	"errors"

	"github.com/google/uuid"
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

func (s *ServiceProveedor) GetProveedorByID(id uuid.UUID) (*models.Proveedor, error) {
	if id == uuid.Nil {
		return nil, errors.New("el ID del proveedor es requerido")
	}
	return s.store.GetProveedorByID(id)
}

func (s *ServiceProveedor) CreateProveedor(proveedor models.Proveedor) (*models.Proveedor, error) {
	if proveedor.Nombre == "" {
		return nil, errors.New("falta el nombre del proveedor")
	}
	if proveedor.IDSucursal == uuid.Nil {
		return nil, errors.New("falta el id de la sucursal")
	}
	if proveedor.IDEmpresa == uuid.Nil {
		return nil, errors.New("falta el id de la empresa")
	}
	if proveedor.IDStatus == uuid.Nil {
		return nil, errors.New("falta el id de estatus")
	}
	return s.store.CreateProveedor(&proveedor)
}

func (s *ServiceProveedor) UpdateProveedor(id uuid.UUID, proveedor models.Proveedor) (*models.Proveedor, error) {
	if id == uuid.Nil {
		return nil, errors.New("el ID del proveedor es requerido")
	}
	if proveedor.Nombre == "" {
		return nil, errors.New("falta el nombre del proveedor")
	}
	return s.store.UpdateProveedor(id, &proveedor)
}

func (s *ServiceProveedor) DeleteProveedor(id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("el ID del proveedor es requerido")
	}
	return s.store.DeleteProveedor(id)
}
