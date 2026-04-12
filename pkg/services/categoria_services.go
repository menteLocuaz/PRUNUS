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

type ServiceCategoria struct {
	store  store.StoreCategoria
	cache  *utils.CacheManager
	logger *slog.Logger
}

func NewServiceCategoria(s store.StoreCategoria, c *utils.CacheManager, logger *slog.Logger) *ServiceCategoria {
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
	return utils.GetOrSet(ctx, s.cache, cacheKeyCategoriasAll, cacheExpiration, func() ([]*models.Categoria, error) {
		return s.store.GetAllCategorias(ctx)
	})
}

func (s *ServiceCategoria) GetCategoriaByID(ctx context.Context, id uuid.UUID) (*models.Categoria, error) {
	key := fmt.Sprintf(cacheKeyCategoriaID, id.String())
	return utils.GetOrSet(ctx, s.cache, key, cacheExpiration, func() (*models.Categoria, error) {
		return s.store.GetCategoriaByID(ctx, id)
	})
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
	if categoria.IDStatus == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de creación de categoría sin estatus", slog.String("nombre", categoria.Nombre))
		return nil, errors.New("falta el id de estatus")
	}

	result, err := s.store.CreateCategoria(ctx, &categoria)
	if err != nil {
		return nil, err
	}

	// Invalidar caché de la lista completa
	s.cache.Invalidate(ctx, "categorias:")

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
	s.cache.Invalidate(ctx, "categorias:")

	return result, nil
}

func (s *ServiceCategoria) DeleteCategoria(ctx context.Context, id uuid.UUID) error {
	err := s.store.DeleteCategoria(ctx, id)
	if err != nil {
		return err
	}

	// Invalidar caché
	s.cache.Invalidate(ctx, "categorias:")

	return nil
}
