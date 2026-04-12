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

type ServiceUnidad struct {
	store  store.StoreUnidad
	cache  *utils.CacheManager
	logger *slog.Logger
}

func NewServiceUnidad(s store.StoreUnidad, c *utils.CacheManager, logger *slog.Logger) *ServiceUnidad {
	return &ServiceUnidad{
		store:  s,
		cache:  c,
		logger: logger,
	}
}

const (
	cacheKeyUnidadesAll = "unidades:all"
	cacheKeyUnidadID    = "unidades:id:%s"
	cacheTTLUnidades    = 24 * time.Hour
)

func (s *ServiceUnidad) GetAllUnidades(ctx context.Context) ([]*models.Unidad, error) {
	return utils.GetOrSet(ctx, s.cache, cacheKeyUnidadesAll, cacheTTLUnidades, func() ([]*models.Unidad, error) {
		return s.store.GetAllUnidades(ctx)
	})
}

func (s *ServiceUnidad) GetUnidadByID(ctx context.Context, id uuid.UUID) (*models.Unidad, error) {
	if id == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de obtener unidad con ID nulo")
		return nil, errors.New("el ID de la unidad es requerido")
	}

	key := fmt.Sprintf(cacheKeyUnidadID, id.String())
	return utils.GetOrSet(ctx, s.cache, key, cacheTTLUnidades, func() (*models.Unidad, error) {
		return s.store.GetUnidadByID(ctx, id)
	})
}

func (s *ServiceUnidad) CreateUnidad(ctx context.Context, unidad models.Unidad) (*models.Unidad, error) {
	if unidad.Nombre == "" {
		s.logger.WarnContext(ctx, "Intento de creación de unidad con nombre vacío")
		return nil, errors.New("falta el nombre de la unidad")
	}
	if unidad.Abreviatura == "" {
		s.logger.WarnContext(ctx, "Intento de creación de unidad con abreviatura vacía")
		return nil, errors.New("falta la abreviatura de la unidad")
	}
	if unidad.IDSucursal == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de creación de unidad sin sucursal", slog.String("nombre", unidad.Nombre))
		return nil, errors.New("falta el id de la sucursal")
	}
	if unidad.IDStatus == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de creación de unidad sin estatus", slog.String("nombre", unidad.Nombre))
		return nil, errors.New("falta el id de estatus")
	}

	res, err := s.store.CreateUnidad(ctx, &unidad)
	if err != nil {
		return nil, err
	}

	// Invalidar caché
	s.cache.Invalidate(ctx, "unidades:")

	return res, nil
}

func (s *ServiceUnidad) UpdateUnidad(ctx context.Context, id uuid.UUID, unidad models.Unidad) (*models.Unidad, error) {
	if id == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de actualización de unidad con ID nulo")
		return nil, errors.New("el ID de la unidad es requerido")
	}
	if unidad.Nombre == "" {
		s.logger.WarnContext(ctx, "Intento de actualización de unidad con nombre vacío", slog.String("id", id.String()))
		return nil, errors.New("falta el nombre de la unidad")
	}

	res, err := s.store.UpdateUnidad(ctx, id, &unidad)
	if err != nil {
		return nil, err
	}

	// Invalidar caché
	s.cache.Invalidate(ctx, "unidades:")

	return res, nil
}

func (s *ServiceUnidad) DeleteUnidad(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de eliminación de unidad con ID nulo")
		return errors.New("el ID de la unidad es requerido")
	}

	if err := s.store.DeleteUnidad(ctx, id); err != nil {
		return err
	}

	// Invalidar caché
	s.cache.Invalidate(ctx, "unidades:")

	return nil
}
