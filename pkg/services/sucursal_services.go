package services

import (
	"errors"

	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
)

type ServiceSucursal struct {
	store store.StoreSucursal
}

func NewServiceSucursal(s store.StoreSucursal) *ServiceSucursal {
	return &ServiceSucursal{store: s}
}

// obtine todas las sucursales
func (s *ServiceSucursal) GetAllSucursales() ([]*models.Sucursal, error) {
	return s.store.GetAllSucursales()
}

// obtien solo una sucursla
func (s *ServiceSucursal) GetSucursalByID(id uint) (*models.Sucursal, error) {
	return s.store.GetSucursalByID(id)
}

// crea sucursal
func (s *ServiceSucursal) CreateSucursal(sucursal models.Sucursal) (*models.Sucursal, error) {
	if sucursal.NombreSucursal == "" {
		return nil, errors.New("Falta el nombre de la sucursal")
	}

	return s.store.CreateSucursal(&sucursal)
}

// actualizar empresa
func (s *ServiceSucursal) UpdateSucursal(id uint, sucursal models.Sucursal) (*models.Sucursal, error) {
	if sucursal.NombreSucursal == "" {
		return nil, errors.New("Falta el nombre de la sucursal")
	}
	return s.store.UpdateSucursal(id, &sucursal)
}

// eliminar empresa
func (s *ServiceSucursal) DeleteSucursal(id uint) error {
	return s.store.DeleteSucursal(id)
}
