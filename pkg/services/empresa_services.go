package services

import (
	"errors"

	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
)

type ServiceEmpresa struct {
	store store.StoreEmpresa
}

func NewServiceEmpresa(s store.StoreEmpresa) *ServiceEmpresa {
	return &ServiceEmpresa{store: s}
}

func (s *ServiceEmpresa) GetAllEmpresa() ([]*models.Empresa, error) {
	return s.store.GetAllEmpresa()
}

func (s *ServiceEmpresa) GetByIDEmpresa(id uint) (*models.Empresa, error) {
	return s.store.GetByIdEmpresa(id)
}

func (s *ServiceEmpresa) CrearEmpresa(empresa models.Empresa) (*models.Empresa, error) {
	if empresa.Nombre == "" {
		return nil, errors.New("Falta el nombre del empresa")
	}
	return s.store.CreateEmpresa(&empresa)
}

func (s *ServiceEmpresa) UpdateEmpresa(id uint, empresa models.Empresa) (*models.Empresa, error) {
	if empresa.Nombre == "" {
		return nil, errors.New("Falta el nombre del empresa")
	}
	return s.store.UpdateEmpresa(id, &empresa)
}

func (s *ServiceEmpresa) ElimminarEmpresa(id uint) error {
	return s.store.DeleteEmpresa(id)
}
