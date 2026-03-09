package services

import (
	"errors"

	"github.com/google/uuid"
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

func (s *ServiceUnidad) GetUnidadByID(id uuid.UUID) (*models.Unidad, error) {
	if id == uuid.Nil {
		return nil, errors.New("el ID de la unidad es requerido")
	}
	return s.store.GetUnidadByID(id)
}

func (s *ServiceUnidad) CreateUnidad(unidad models.Unidad) (*models.Unidad, error) {
	if unidad.Nombre == "" {
		return nil, errors.New("falta el nombre de la unidad")
	}
	if unidad.IDSucursal == uuid.Nil {
		return nil, errors.New("falta el id de la sucursal")
	}
	return s.store.CreateUnidad(&unidad)
}

func (s *ServiceUnidad) UpdateUnidad(id uuid.UUID, unidad models.Unidad) (*models.Unidad, error) {
	if id == uuid.Nil {
		return nil, errors.New("el ID de la unidad es requerido")
	}
	if unidad.Nombre == "" {
		return nil, errors.New("falta el nombre de la unidad")
	}
	return s.store.UpdateUnidad(id, &unidad)
}

func (s *ServiceUnidad) DeleteUnidad(id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("el ID de la unidad es requerido")
	}
	return s.store.DeleteUnidad(id)
}
