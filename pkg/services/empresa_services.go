// Package services contiene la lógica de negocio y reglas de validación
// para las entidades del sistema. En este caso, maneja las operaciones
// relacionadas con Empresa.
package services

import (
	"errors"

	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
)

// ServiceEmpresa define un servicio que encapsula la lógica de negocio
// para la entidad Empresa. Se apoya en un StoreEmpresa para interactuar
// con la capa de persistencia (base de datos).
type ServiceEmpresa struct {
	store store.StoreEmpresa
}

// NewServiceEmpresa crea una nueva instancia del servicio de Empresa,
// recibiendo como dependencia un StoreEmpresa. Esto facilita la inyección
// de dependencias y el testeo.
func NewServiceEmpresa(s store.StoreEmpresa) *ServiceEmpresa {
	return &ServiceEmpresa{store: s}
}

// GetAllEmpresa devuelve todas las empresas registradas en el sistema.
// Retorna un slice de punteros a Empresa y un posible error.
func (s *ServiceEmpresa) GetAllEmpresa() ([]*models.Empresa, error) {
	return s.store.GetAllEmpresa()
}

// GetByIDEmpresa busca una empresa por su ID único.
// Retorna la empresa encontrada o un error si no existe.
func (s *ServiceEmpresa) GetByIDEmpresa(id uint) (*models.Empresa, error) {
	return s.store.GetByIdEmpresa(id)
}

// CrearEmpresa valida y crea una nueva empresa en el sistema.
// Si el nombre está vacío, retorna un error de validación.
// En caso contrario, delega la creación al Store.
func (s *ServiceEmpresa) CrearEmpresa(empresa models.Empresa) (*models.Empresa, error) {
	if empresa.Nombre == "" {
		return nil, errors.New("Falta el nombre de la empresa")
	}
	return s.store.CreateEmpresa(&empresa)
}

// UpdateEmpresa valida y actualiza una empresa existente.
// Si el nombre está vacío, retorna un error de validación.
// En caso contrario, delega la actualización al Store.
func (s *ServiceEmpresa) UpdateEmpresa(id uint, empresa models.Empresa) (*models.Empresa, error) {
	if empresa.Nombre == "" {
		return nil, errors.New("Falta el nombre de la empresa")
	}
	return s.store.UpdateEmpresa(id, &empresa)
}

// ElimminarEmpresa elimina una empresa por su ID.
// Retorna un error si la operación falla.
func (s *ServiceEmpresa) ElimminarEmpresa(id uint) error {
	return s.store.DeleteEmpresa(id)
}
