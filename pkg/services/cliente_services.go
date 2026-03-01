package services

import (
	"errors"

	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
)

type ServiceCliente struct {
	store store.StoreCliente
}

func NewServiceCliente(s store.StoreCliente) *ServiceCliente {
	return &ServiceCliente{store: s}
}

func (s *ServiceCliente) GetAllClientes() ([]*models.Cliente, error) {
	return s.store.GetAllClientes()
}

func (s *ServiceCliente) GetClienteByID(id uint) (*models.Cliente, error) {
	return s.store.GetClienteByID(id)
}

func (s *ServiceCliente) CreateCliente(cliente models.Cliente) (*models.Cliente, error) {
	if cliente.Nombre == "" {
		return nil, errors.New("falta el nombre del cliente")
	}
	if cliente.EmpresaCliente == "" {
		return nil, errors.New("falta el nombre de la empresa del cliente")
	}
	return s.store.CreateCliente(&cliente)
}

func (s *ServiceCliente) UpdateCliente(id uint, cliente models.Cliente) (*models.Cliente, error) {
	if cliente.Nombre == "" {
		return nil, errors.New("falta el nombre del cliente")
	}
	return s.store.UpdateCliente(id, &cliente)
}

func (s *ServiceCliente) DeleteCliente(id uint) error {
	return s.store.DeleteCliente(id)
}
