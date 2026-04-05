package store

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
)

const (
	cacheKeyCategoriasAll = "categorias:all"
	cacheKeyCategoriaID   = "categorias:id:%s"
	cacheExpiration       = 1 * time.Hour
)

type CategoriaCacheDecorator struct {
	base  StoreCategoria
	cache models.CacheStore
}

func NewCategoriaCacheDecorator(base StoreCategoria, cache models.CacheStore) StoreCategoria {
	return &CategoriaCacheDecorator{
		base:  base,
		cache: cache,
	}
}

func (d *CategoriaCacheDecorator) GetAllCategorias(ctx context.Context) ([]*models.Categoria, error) {
	if d.cache == nil {
		return d.base.GetAllCategorias(ctx)
	}

	var categorias []*models.Categoria
	if err := d.cache.Get(ctx, cacheKeyCategoriasAll, &categorias); err == nil {
		return categorias, nil
	}

	categorias, err := d.base.GetAllCategorias(ctx)
	if err != nil {
		return nil, err
	}

	_ = d.cache.Set(ctx, cacheKeyCategoriasAll, categorias, cacheExpiration)
	return categorias, nil
}

func (d *CategoriaCacheDecorator) GetCategoriaByID(ctx context.Context, id uuid.UUID) (*models.Categoria, error) {
	if d.cache == nil {
		return d.base.GetCategoriaByID(ctx, id)
	}

	var categoria *models.Categoria
	key := fmt.Sprintf(cacheKeyCategoriaID, id.String())
	if err := d.cache.Get(ctx, key, &categoria); err == nil {
		return categoria, nil
	}

	categoria, err := d.base.GetCategoriaByID(ctx, id)
	if err != nil {
		return nil, err
	}

	_ = d.cache.Set(ctx, key, categoria, cacheExpiration)
	return categoria, nil
}

func (d *CategoriaCacheDecorator) CreateCategoria(ctx context.Context, categoria *models.Categoria) (*models.Categoria, error) {
	res, err := d.base.CreateCategoria(ctx, categoria)
	if err == nil && d.cache != nil {
		_ = d.cache.Delete(ctx, cacheKeyCategoriasAll)
	}
	return res, err
}

func (d *CategoriaCacheDecorator) UpdateCategoria(ctx context.Context, id uuid.UUID, categoria *models.Categoria) (*models.Categoria, error) {
	res, err := d.base.UpdateCategoria(ctx, id, categoria)
	if err == nil && d.cache != nil {
		_ = d.cache.Delete(ctx, cacheKeyCategoriasAll)
		_ = d.cache.Delete(ctx, fmt.Sprintf(cacheKeyCategoriaID, id.String()))
	}
	return res, err
}

func (d *CategoriaCacheDecorator) DeleteCategoria(ctx context.Context, id uuid.UUID) error {
	err := d.base.DeleteCategoria(ctx, id)
	if err == nil && d.cache != nil {
		_ = d.cache.Delete(ctx, cacheKeyCategoriasAll)
		_ = d.cache.Delete(ctx, fmt.Sprintf(cacheKeyCategoriaID, id.String()))
	}
	return err
}
