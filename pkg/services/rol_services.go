package services

import (
	"errors"

	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
)

// ServiceRol servicio que encapsula la lógica de negocio para rol
type ServiceRol struct {
	store store.StoreRol
}

// NewServiceRol crea una nueva instancia del servicio de rol
func NewServiceRol(s store.StoreRol) *ServiceRol {
	return &ServiceRol{store: s}
}

// GetAllRoles obtiene todos los roles del sistema
func (s *ServiceRol) GetAllRoles() ([]*models.Rol, error) {
	return s.store.GetAllRoles()
}

// GetRolByID obtiene un rol por su ID
func (s *ServiceRol) GetRolByID(id uint) (*models.Rol, error) {
	if id == 0 {
		return nil, errors.New("el ID del rol es requerido")
	}
	return s.store.GetRolByID(id)
}

// CreateRol crea un nuevo rol con validaciones de negocio
func (s *ServiceRol) CreateRol(rol models.Rol) (*models.Rol, error) {
	// Validar campos obligatorios
	if rol.RolNombre == "" {
		return nil, errors.New("el nombre del rol es requerido")
	}
	if rol.IDSucursal == 0 {
		return nil, errors.New("el ID de la sucursal es requerido")
	}

	// Establecer estado por defecto si no está definido
	if rol.Estado == 0 {
		rol.Estado = 1
	}

	return s.store.CreateRol(&rol)
}

// UpdateRol actualiza un rol existente con validaciones
func (s *ServiceRol) UpdateRol(id uint, rol models.Rol) (*models.Rol, error) {
	if id == 0 {
		return nil, errors.New("el ID del rol es requerido")
	}
	if rol.RolNombre == "" {
		return nil, errors.New("el nombre del rol es requerido")
	}
	if rol.IDSucursal == 0 {
		return nil, errors.New("el ID de la sucursal es requerido")
	}

	return s.store.UpdateRol(id, &rol)
}

// DeleteRol elimina un rol (soft delete)
func (s *ServiceRol) DeleteRol(id uint) error {
	if id == 0 {
		return errors.New("el ID del rol es requerido")
	}
	return s.store.DeleteRol(id)
}
