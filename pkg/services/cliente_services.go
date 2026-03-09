package services

import (
	"errors"

	"github.com/google/uuid"
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

func (s *ServiceCliente) GetClienteByID(id uuid.UUID) (*models.Cliente, error) {
	if id == uuid.Nil {
		return nil, errors.New("el ID del cliente es requerido")
	}
	return s.store.GetClienteByID(id)
}

func (s *ServiceCliente) CreateCliente(cliente models.Cliente) (*models.Cliente, error) {
	if cliente.Nombre == "" {
		return nil, errors.New("falta el nombre del cliente")
	}
	if cliente.EmpresaCliente == "" {
		return nil, errors.New("falta el nombre de la empresa del cliente")
	}
	if cliente.IDStatus == uuid.Nil {
		return nil, errors.New("falta el ID de estatus")
	}
	return s.store.CreateCliente(&cliente)
}

func (s *ServiceCliente) UpdateCliente(id uuid.UUID, cliente models.Cliente) (*models.Cliente, error) {
	if id == uuid.Nil {
		return nil, errors.New("el ID del cliente es requerido")
	}
	if cliente.Nombre == "" {
		return nil, errors.New("falta el nombre del cliente")
	}
	return s.store.UpdateCliente(id, &cliente)
}

func (s *ServiceCliente) DeleteCliente(id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("el ID del cliente es requerido")
	}
	return s.store.DeleteCliente(id)
}
