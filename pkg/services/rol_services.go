package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
)

// ServiceRol servicio que encapsula la lógica de negocio para rol
type ServiceRol struct {
	store store.StoreRol
	cache models.CacheStore
}

// NewServiceRol crea una nueva instancia del servicio de rol
func NewServiceRol(s store.StoreRol, c models.CacheStore) *ServiceRol {
	return &ServiceRol{store: s, cache: c}
}

const (
	cacheKeyRolesAll = "roles:all"
	cacheKeyRolID    = "roles:id:%d"
)

// GetAllRoles obtiene todos los roles del sistema
func (s *ServiceRol) GetAllRoles() ([]*models.Rol, error) {
	ctx := context.Background()
	var roles []*models.Rol

	// Intentar obtener del caché
	err := s.cache.Get(ctx, cacheKeyRolesAll, &roles)
	if err == nil {
		return roles, nil
	}

	// Si no hay caché, ir a la base de datos
	roles, err = s.store.GetAllRoles()
	if err != nil {
		return nil, err
	}

	// Guardar en caché
	_ = s.cache.Set(ctx, cacheKeyRolesAll, roles, 1*time.Hour)

	return roles, nil
}

// GetRolByID obtiene un rol por su ID
func (s *ServiceRol) GetRolByID(id uint) (*models.Rol, error) {
	if id == 0 {
		return nil, errors.New("el ID del rol es requerido")
	}

	ctx := context.Background()
	var rol *models.Rol
	key := fmt.Sprintf(cacheKeyRolID, id)

	// Intentar obtener del caché
	err := s.cache.Get(ctx, key, &rol)
	if err == nil {
		return rol, nil
	}

	// Si no hay caché, ir a la base de datos
	rol, err = s.store.GetRolByID(id)
	if err != nil {
		return nil, err
	}

	// Guardar en caché
	_ = s.cache.Set(ctx, key, rol, 1*time.Hour)

	return rol, nil
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

	result, err := s.store.CreateRol(&rol)
	if err != nil {
		return nil, err
	}

	// Invalidar caché de la lista completa
	_ = s.cache.Delete(context.Background(), cacheKeyRolesAll)

	return result, nil
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

	result, err := s.store.UpdateRol(id, &rol)
	if err != nil {
		return nil, err
	}

	// Invalidar caché
	ctx := context.Background()
	_ = s.cache.Delete(ctx, cacheKeyRolesAll)
	_ = s.cache.Delete(ctx, fmt.Sprintf(cacheKeyRolID, id))

	return result, nil
}

// DeleteRol elimina un rol (soft delete)
func (s *ServiceRol) DeleteRol(id uint) error {
	if id == 0 {
		return errors.New("el ID del rol es requerido")
	}

	err := s.store.DeleteRol(id)
	if err != nil {
		return err
	}

	// Invalidar caché
	ctx := context.Background()
	_ = s.cache.Delete(ctx, cacheKeyRolesAll)
	_ = s.cache.Delete(ctx, fmt.Sprintf(cacheKeyRolID, id))

	return nil
}
