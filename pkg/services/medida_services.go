package services

import (
	"errors"

	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
)

type ServiceUnidad struct {
	store store.StoreUnidad
}

func NewServiceUnidad(s store.StoreUnidad) *ServiceUnidad {
	return &ServiceUnidad{store: s}
}

func (s *ServiceUnidad) GetAllUnidades() ([]*models.Unidad, error) {
	return s.store.GetAllUnidades()
}

func (s *ServiceUnidad) GetUnidadByID(id uint) (*models.Unidad, error) {
	return s.store.GetUnidadByID(id)
}

func (s *ServiceUnidad) CreateUnidad(unidad models.Unidad) (*models.Unidad, error) {
	if unidad.Nombre == "" {
		return nil, errors.New("falta el nombre de la unidad")
	}
	if unidad.IDSucursal == 0 {
		return nil, errors.New("falta el id de la sucursal")
	}
	return s.store.CreateUnidad(&unidad)
}

func (s *ServiceUnidad) UpdateUnidad(id uint, unidad models.Unidad) (*models.Unidad, error) {
	if unidad.Nombre == "" {
		return nil, errors.New("falta el nombre de la unidad")
	}
	return s.store.UpdateUnidad(id, &unidad)
}

func (s *ServiceUnidad) DeleteUnidad(id uint) error {
	return s.store.DeleteUnidad(id)
}
