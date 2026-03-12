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

type ServiceCategoria struct {
	store  store.StoreCategoria
	cache  models.CacheStore
	logger *slog.Logger
}

func NewServiceCategoria(s store.StoreCategoria, c models.CacheStore, logger *slog.Logger) *ServiceCategoria {
	return &ServiceCategoria{
		store:  s,
		cache:  c,
		logger: logger,
	}
}

const (
	cacheKeyCategoriasAll = "categorias:all"
	cacheKeyCategoriaID   = "categorias:id:%s"
	cacheExpiration       = 1 * time.Hour
)

func (s *ServiceCategoria) GetAllCategorias(ctx context.Context) ([]*models.Categoria, error) {
	var categorias []*models.Categoria

	// Intentar obtener del caché
	err := s.cache.Get(ctx, cacheKeyCategoriasAll, &categorias)
	if err == nil {
		return categorias, nil
	}

	// Si no hay caché, ir a la base de datos
	categorias, err = s.store.GetAllCategorias(ctx)
	if err != nil {
		return nil, err
	}

	// Guardar en caché
	_ = s.cache.Set(ctx, cacheKeyCategoriasAll, categorias, cacheExpiration)

	return categorias, nil
}

func (s *ServiceCategoria) GetCategoriaByID(ctx context.Context, id uuid.UUID) (*models.Categoria, error) {
	var categoria *models.Categoria
	key := fmt.Sprintf(cacheKeyCategoriaID, id.String())

	// Intentar obtener del caché
	err := s.cache.Get(ctx, key, &categoria)
	if err == nil {
		return categoria, nil
	}

	// Si no hay caché, ir a la base de datos
	categoria, err = s.store.GetCategoriaByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Guardar en caché
	_ = s.cache.Set(ctx, key, categoria, cacheExpiration)

	return categoria, nil
}

func (s *ServiceCategoria) CreateCategoria(ctx context.Context, categoria models.Categoria) (*models.Categoria, error) {
	if categoria.Nombre == "" {
		s.logger.WarnContext(ctx, "Intento de creación de categoría con nombre vacío")
		return nil, errors.New("falta el nombre de la categoria")
	}
	if categoria.IDSucursal == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de creación de categoría sin sucursal", slog.String("nombre", categoria.Nombre))
		return nil, errors.New("falta el id de la sucursal")
	}

	result, err := s.store.CreateCategoria(ctx, &categoria)
	if err != nil {
		return nil, err
	}

	// Invalidar caché de la lista completa
	_ = s.cache.Delete(ctx, cacheKeyCategoriasAll)

	return result, nil
}

func (s *ServiceCategoria) UpdateCategoria(ctx context.Context, id uuid.UUID, categoria models.Categoria) (*models.Categoria, error) {
	if categoria.Nombre == "" {
		s.logger.WarnContext(ctx, "Intento de actualización de categoría con nombre vacío", slog.String("id", id.String()))
		return nil, errors.New("falta el nombre de la categoria")
	}

	result, err := s.store.UpdateCategoria(ctx, id, &categoria)
	if err != nil {
		return nil, err
	}

	// Invalidar caché
	_ = s.cache.Delete(ctx, cacheKeyCategoriasAll)
	_ = s.cache.Delete(ctx, fmt.Sprintf(cacheKeyCategoriaID, id.String()))

	return result, nil
}

func (s *ServiceCategoria) DeleteCategoria(ctx context.Context, id uuid.UUID) error {
	err := s.store.DeleteCategoria(ctx, id)
	if err != nil {
		return err
	}

	// Invalidar caché
	_ = s.cache.Delete(ctx, cacheKeyCategoriasAll)
	_ = s.cache.Delete(ctx, fmt.Sprintf(cacheKeyCategoriaID, id.String()))

	return nil
}
