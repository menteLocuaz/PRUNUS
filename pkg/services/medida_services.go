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

type ServiceUnidad struct {
	store  store.StoreUnidad
	cache  models.CacheStore
	logger *slog.Logger
}

func NewServiceUnidad(s store.StoreUnidad, c models.CacheStore, logger *slog.Logger) *ServiceUnidad {
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
	var unidades []*models.Unidad

	// Intentar caché
	if err := s.cache.Get(ctx, cacheKeyUnidadesAll, &unidades); err == nil {
		return unidades, nil
	}

	// BD
	unidades, err := s.store.GetAllUnidades(ctx)
	if err != nil {
		return nil, err
	}

	// Guardar en caché
	_ = s.cache.Set(ctx, cacheKeyUnidadesAll, unidades, cacheTTLUnidades)

	return unidades, nil
}

func (s *ServiceUnidad) GetUnidadByID(ctx context.Context, id uuid.UUID) (*models.Unidad, error) {
	if id == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de obtener unidad con ID nulo")
		return nil, errors.New("el ID de la unidad es requerido")
	}

	var unidad *models.Unidad
	key := fmt.Sprintf(cacheKeyUnidadID, id.String())

	// Intentar caché
	if err := s.cache.Get(ctx, key, &unidad); err == nil {
		return unidad, nil
	}

	// BD
	unidad, err := s.store.GetUnidadByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Guardar en caché
	_ = s.cache.Set(ctx, key, unidad, cacheTTLUnidades)

	return unidad, nil
}

func (s *ServiceUnidad) CreateUnidad(ctx context.Context, unidad models.Unidad) (*models.Unidad, error) {
	if unidad.Nombre == "" {
		s.logger.WarnContext(ctx, "Intento de creación de unidad con nombre vacío")
		return nil, errors.New("falta el nombre de la unidad")
	}
	if unidad.IDSucursal == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de creación de unidad sin sucursal", slog.String("nombre", unidad.Nombre))
		return nil, errors.New("falta el id de la sucursal")
	}
	
	res, err := s.store.CreateUnidad(ctx, &unidad)
	if err != nil {
		return nil, err
	}

	// Invalidar caché
	s.invalidateCache(ctx, res.IDUnidad)

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
	s.invalidateCache(ctx, id)

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
	s.invalidateCache(ctx, id)

	return nil
}

func (s *ServiceUnidad) invalidateCache(ctx context.Context, id uuid.UUID) {
	_ = s.cache.Delete(ctx, cacheKeyUnidadesAll)
	_ = s.cache.Delete(ctx, fmt.Sprintf(cacheKeyUnidadID, id.String()))
}
