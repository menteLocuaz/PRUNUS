package services

import (
	"errors"

	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
)

type ServiceMoneda struct {
	store store.StoreMoneda
}

func NewServiceMoneda(s store.StoreMoneda) *ServiceMoneda {
	return &ServiceMoneda{store: s}
}

func (s *ServiceMoneda) GetAllMonedas() ([]*models.Moneda, error) {
	return s.store.GetAllMonedas()
}

func (s *ServiceMoneda) GetMonedaByID(id uint) (*models.Moneda, error) {
	return s.store.GetMonedaByID(id)
}

func (s *ServiceMoneda) CreateMoneda(moneda models.Moneda) (*models.Moneda, error) {
	if moneda.Nombre == "" {
		return nil, errors.New("falta el nombre de la moneda")
	}
	if moneda.IDSucursal == 0 {
		return nil, errors.New("falta el id de la sucursal")
	}
	return s.store.CreateMoneda(&moneda)
}

func (s *ServiceMoneda) UpdateMoneda(id uint, moneda models.Moneda) (*models.Moneda, error) {
	if moneda.Nombre == "" {
		return nil, errors.New("falta el nombre de la moneda")
	}
	return s.store.UpdateMoneda(id, &moneda)
}

func (s *ServiceMoneda) DeleteMoneda(id uint) error {
	return s.store.DeleteMoneda(id)
}
