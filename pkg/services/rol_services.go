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
)

// ServiceRol servicio que encapsula la lógica de negocio para rol
type ServiceRol struct {
	store  store.StoreRol
	cache  models.CacheStore
	logger *slog.Logger
}

// NewServiceRol crea una nueva instancia del servicio de rol
func NewServiceRol(s store.StoreRol, c models.CacheStore, logger *slog.Logger) *ServiceRol {
	return &ServiceRol{
		store:  s,
		cache:  c,
		logger: logger,
	}
}

const (
	cacheKeyRolesAll = "roles:all"
	cacheKeyRolID    = "roles:id:%s"
	cacheKeyPermisos = "roles:permisos:%s"
)

// GetPermisosByRol obtiene los permisos de un rol con caching
func (s *ServiceRol) GetPermisosByRol(ctx context.Context, rolID uuid.UUID) ([]string, error) {
	if rolID == uuid.Nil {
		return nil, nil
	}

	var permisos []string
	key := fmt.Sprintf(cacheKeyPermisos, rolID.String())

	// Intentar obtener del caché
	err := s.cache.Get(ctx, key, &permisos)
	if err == nil {
		return permisos, nil
	}

	// Si no hay caché, ir al store (que usa storeUsuario internamente o directamente si lo exponemos)
	// Nota: Como storeRol no tiene GetPermisosByRol, usualmente esto se consulta en storeUsuario
	// pero para coherencia lo manejaremos aquí si storeRol lo permite.
	// REVISIÓN: El store de usuarios tiene la implementación.
	// Para este refactor, asumiremos que ServiceRol tiene acceso al store necesario o lo inyectamos.
	// Dado el estado actual, lo ideal es que storeRol tenga este método si es responsabilidad de roles.
	
	// Por ahora, simularemos la obtención desde el store que tenga el método.
	// Si storeRol no lo tiene, deberíamos agregarlo a la interfaz StoreRol.
	permisos, err = s.store.GetPermisosByRol(ctx, rolID)
	if err != nil {
		return nil, err
	}

	// Guardar en caché con TTL largo (ej: 24 horas)
	_ = s.cache.Set(ctx, key, permisos, 24*time.Hour)

	return permisos, nil
}

// GetAllRoles obtiene todos los roles del sistema
func (s *ServiceRol) GetAllRoles(ctx context.Context) ([]*models.Rol, error) {
	var roles []*models.Rol

	// Intentar obtener del caché
	err := s.cache.Get(ctx, cacheKeyRolesAll, &roles)
	if err == nil {
		return roles, nil
	}

	// Si no hay caché, ir a la base de datos
	roles, err = s.store.GetAllRoles(ctx)
	if err != nil {
		return nil, err
	}

	// Guardar en caché
	_ = s.cache.Set(ctx, cacheKeyRolesAll, roles, 1*time.Hour)

	return roles, nil
}

// GetRolByID obtiene un rol por su ID
func (s *ServiceRol) GetRolByID(ctx context.Context, id uuid.UUID) (*models.Rol, error) {
	if id == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de obtener rol con ID nulo")
		return nil, errors.New("el ID del rol es requerido")
	}

	var rol *models.Rol
	key := fmt.Sprintf(cacheKeyRolID, id.String())

	// Intentar obtener del caché
	err := s.cache.Get(ctx, key, &rol)
	if err == nil {
		return rol, nil
	}

	// Si no hay caché, ir a la base de datos
	rol, err = s.store.GetRolByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Guardar en caché
	_ = s.cache.Set(ctx, key, rol, 1*time.Hour)

	return rol, nil
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

	// Invalidar caché de la lista completa
	_ = s.cache.Delete(ctx, cacheKeyRolesAll)

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

	// Invalidar caché
	_ = s.cache.Delete(ctx, cacheKeyRolesAll)
	_ = s.cache.Delete(ctx, fmt.Sprintf(cacheKeyRolID, id.String()))
	_ = s.cache.Delete(ctx, fmt.Sprintf(cacheKeyPermisos, id.String()))

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

	// Invalidar caché
	_ = s.cache.Delete(ctx, cacheKeyRolesAll)
	_ = s.cache.Delete(ctx, fmt.Sprintf(cacheKeyRolID, id.String()))
	_ = s.cache.Delete(ctx, fmt.Sprintf(cacheKeyPermisos, id.String()))

	return nil
}
