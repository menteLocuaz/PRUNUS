package services

import (
	"errors"

	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
)

type ServiceCategoria struct {
	store store.StoreCategoria
}

func NewServiceCategoria(s store.StoreCategoria) *ServiceCategoria {
	return &ServiceCategoria{store: s}
}

func (s *ServiceCategoria) GetAllCategorias() ([]*models.Categoria, error) {
	return s.store.GetAllCategorias()
}

func (s *ServiceCategoria) GetCategoriaByID(id uint) (*models.Categoria, error) {
	return s.store.GetCategoriaByID(id)
}

func (s *ServiceCategoria) CreateCategoria(categoria models.Categoria) (*models.Categoria, error) {
	if categoria.Nombre == "" {
		return nil, errors.New("falta el nombre de la categoria")
	}
	if categoria.IDSucursal == 0 {
		return nil, errors.New("falta el id de la sucursal")
	}
	return s.store.CreateCategoria(&categoria)
}

func (s *ServiceCategoria) UpdateCategoria(id uint, categoria models.Categoria) (*models.Categoria, error) {
	if categoria.Nombre == "" {
		return nil, errors.New("falta el nombre de la categoria")
	}
	return s.store.UpdateCategoria(id, &categoria)
}

func (s *ServiceCategoria) DeleteCategoria(id uint) error {
	return s.store.DeleteCategoria(id)
}
