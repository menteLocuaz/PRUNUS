package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
)

type ServiceCategoria struct {
	store store.StoreCategoria
	cache models.CacheStore
}

func NewServiceCategoria(s store.StoreCategoria, c models.CacheStore) *ServiceCategoria {
	return &ServiceCategoria{store: s, cache: c}
}

const (
	cacheKeyCategoriasAll = "categorias:all"
	cacheKeyCategoriaID   = "categorias:id:%s"
	cacheExpiration       = 1 * time.Hour
)

func (s *ServiceCategoria) GetAllCategorias() ([]*models.Categoria, error) {
	ctx := context.Background()
	var categorias []*models.Categoria

	// Intentar obtener del caché
	err := s.cache.Get(ctx, cacheKeyCategoriasAll, &categorias)
	if err == nil {
		return categorias, nil
	}

	// Si no hay caché, ir a la base de datos
	categorias, err = s.store.GetAllCategorias()
	if err != nil {
		return nil, err
	}

	// Guardar en caché
	_ = s.cache.Set(ctx, cacheKeyCategoriasAll, categorias, cacheExpiration)

	return categorias, nil
}

func (s *ServiceCategoria) GetCategoriaByID(id uuid.UUID) (*models.Categoria, error) {
	ctx := context.Background()
	var categoria *models.Categoria
	key := fmt.Sprintf(cacheKeyCategoriaID, id.String())

	// Intentar obtener del caché
	err := s.cache.Get(ctx, key, &categoria)
	if err == nil {
		return categoria, nil
	}

	// Si no hay caché, ir a la base de datos
	categoria, err = s.store.GetCategoriaByID(id)
	if err != nil {
		return nil, err
	}

	// Guardar en caché
	_ = s.cache.Set(ctx, key, categoria, cacheExpiration)

	return categoria, nil
}

func (s *ServiceCategoria) CreateCategoria(categoria models.Categoria) (*models.Categoria, error) {
	if categoria.Nombre == "" {
		return nil, errors.New("falta el nombre de la categoria")
	}
	if categoria.IDSucursal == uuid.Nil {
		return nil, errors.New("falta el id de la sucursal")
	}

	result, err := s.store.CreateCategoria(&categoria)
	if err != nil {
		return nil, err
	}

	// Invalidar caché de la lista completa
	_ = s.cache.Delete(context.Background(), cacheKeyCategoriasAll)

	return result, nil
}

func (s *ServiceCategoria) UpdateCategoria(id uuid.UUID, categoria models.Categoria) (*models.Categoria, error) {
	if categoria.Nombre == "" {
		return nil, errors.New("falta el nombre de la categoria")
	}

	result, err := s.store.UpdateCategoria(id, &categoria)
	if err != nil {
		return nil, err
	}

	// Invalidar caché
	ctx := context.Background()
	_ = s.cache.Delete(ctx, cacheKeyCategoriasAll)
	_ = s.cache.Delete(ctx, fmt.Sprintf(cacheKeyCategoriaID, id.String()))

	return result, nil
}

func (s *ServiceCategoria) DeleteCategoria(id uuid.UUID) error {
	err := s.store.DeleteCategoria(id)
	if err != nil {
		return err
	}

	// Invalidar caché
	ctx := context.Background()
	_ = s.cache.Delete(ctx, cacheKeyCategoriasAll)
	_ = s.cache.Delete(ctx, fmt.Sprintf(cacheKeyCategoriaID, id.String()))

	return nil
}
