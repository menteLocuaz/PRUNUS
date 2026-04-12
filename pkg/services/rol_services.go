package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
	"github.com/prunus/pkg/utils"
)

// ServiceRol servicio que encapsula la lógica de negocio para rol
type ServiceRol struct {
	store  store.StoreRol
	cache  *utils.CacheManager
	logger *slog.Logger
}

// NewServiceRol crea una nueva instancia del servicio de rol
func NewServiceRol(s store.StoreRol, c *utils.CacheManager, logger *slog.Logger) *ServiceRol {
	return &ServiceRol{
		store:  s,
		cache:  c,
		logger: logger,
	}
}

const (
	cacheTagRoles    = "roles:"
	cacheKeyRolesAll = "roles:all"
	cacheKeyRolID    = "roles:id:%s"
	cacheKeyPermisos = "roles:permisos:%s"
)

// GetPermisosByRol obtiene los permisos de un rol con caching
func (s *ServiceRol) GetPermisosByRol(ctx context.Context, rolID uuid.UUID) ([]string, error) {
	if rolID == uuid.Nil {
		return nil, nil
	}

	key := fmt.Sprintf(cacheKeyPermisos, rolID.String())

	return utils.GetOrSet(ctx, s.cache, key, 24*time.Hour, func() ([]string, error) {
		return s.store.GetPermisosByRol(ctx, rolID)
	})
}

// GetAllRoles obtiene todos los roles del sistema
func (s *ServiceRol) GetAllRoles(ctx context.Context) ([]*models.Rol, error) {
	return utils.GetOrSet(ctx, s.cache, cacheKeyRolesAll, 1*time.Hour, func() ([]*models.Rol, error) {
		return s.store.GetAllRoles(ctx)
	})
}

// GetRolByID obtiene un rol por su ID
func (s *ServiceRol) GetRolByID(ctx context.Context, id uuid.UUID) (*models.Rol, error) {
	if id == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de obtener rol con ID nulo")
		return nil, errors.New("el ID del rol es requerido")
	}

	key := fmt.Sprintf(cacheKeyRolID, id.String())

	return utils.GetOrSet(ctx, s.cache, key, 1*time.Hour, func() (*models.Rol, error) {
		return s.store.GetRolByID(ctx, id)
	})
}

// CreateRol crea un nuevo rol con validaciones de negocio
func (s *ServiceRol) CreateRol(ctx context.Context, rol models.Rol) (*models.Rol, error) {
	// Validar campos obligatorios
	if rol.RolNombre == "" {
		s.logger.WarnContext(ctx, "Intento de creación de rol con nombre vacío")
		return nil, errors.New("el nombre del rol es requerido")
	}
	if rol.IDSucursal == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de creación de rol sin sucursal", slog.String("nombre", rol.RolNombre))
		return nil, errors.New("el ID de la sucursal es requerido")
	}

	// Asignar estatus automático si no se proporciona
	if rol.IDStatus == uuid.Nil {
		rol.IDStatus = models.EstatusActivo
	}

	result, err := s.store.CreateRol(ctx, &rol)
	if err != nil {
		return nil, err
	}

	// Invalidar caché del grupo roles
	s.cache.Invalidate(ctx, cacheTagRoles)

	return result, nil
}

// UpdateRol actualiza un rol existente con validaciones
func (s *ServiceRol) UpdateRol(ctx context.Context, id uuid.UUID, rol models.Rol) (*models.Rol, error) {
	if id == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de actualización de rol con ID nulo")
		return nil, errors.New("el ID del rol es requerido")
	}
	if rol.RolNombre == "" {
		s.logger.WarnContext(ctx, "Intento de actualización de rol con nombre vacío", slog.String("id", id.String()))
		return nil, errors.New("el nombre del rol es requerido")
	}
	if rol.IDSucursal == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de actualización de rol sin sucursal", slog.String("id", id.String()))
		return nil, errors.New("el ID de la sucursal es requerido")
	}

	result, err := s.store.UpdateRol(ctx, id, &rol)
	if err != nil {
		return nil, err
	}

	// Invalidar caché del grupo roles
	s.cache.Invalidate(ctx, cacheTagRoles)

	return result, nil
}

// DeleteRol elimina un rol (soft delete)
func (s *ServiceRol) DeleteRol(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de eliminación de rol con ID nulo")
		return errors.New("el ID del rol es requerido")
	}

	err := s.store.DeleteRol(ctx, id)
	if err != nil {
		return err
	}

	// Invalidar caché del grupo roles
	s.cache.Invalidate(ctx, cacheTagRoles)

	return nil
}
