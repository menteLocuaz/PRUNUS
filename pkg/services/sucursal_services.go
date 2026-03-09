package services

import (
	"errors"

	"github.com/google/uuid"
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
func (s *ServiceSucursal) GetSucursalByID(id uuid.UUID) (*models.Sucursal, error) {
	return s.store.GetSucursalByID(id)
}

// crea sucursal
func (s *ServiceSucursal) CreateSucursal(sucursal models.Sucursal) (*models.Sucursal, error) {
	if sucursal.NombreSucursal == "" {
		return nil, errors.New("Falta el nombre de la sucursal")
	}
	if sucursal.IDEmpresa == uuid.Nil {
		return nil, errors.New("Falta el id de la empresa")
	}
	if sucursal.IDStatus == uuid.Nil {
		return nil, errors.New("Falta el id del estatus")
	}

	return s.store.CreateSucursal(&sucursal)
}

// actualizar empresa
func (s *ServiceSucursal) UpdateSucursal(id uuid.UUID, sucursal models.Sucursal) (*models.Sucursal, error) {
	if sucursal.NombreSucursal == "" {
		return nil, errors.New("Falta el nombre de la sucursal")
	}
	return s.store.UpdateSucursal(id, &sucursal)
}

// eliminar empresa
func (s *ServiceSucursal) DeleteSucursal(id uuid.UUID) error {
	return s.store.DeleteSucursal(id)
}
