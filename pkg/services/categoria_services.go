package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
)

type ServiceCategoria struct {
	store  store.StoreCategoria
	logger *slog.Logger
}

func NewServiceCategoria(s store.StoreCategoria, logger *slog.Logger) *ServiceCategoria {
	return &ServiceCategoria{
		store:  s,
		logger: logger,
	}
}

func (s *ServiceCategoria) GetAllCategorias(ctx context.Context) ([]*models.Categoria, error) {
	return s.store.GetAllCategorias(ctx)
}

func (s *ServiceCategoria) GetCategoriaByID(ctx context.Context, id uuid.UUID) (*models.Categoria, error) {
	return s.store.GetCategoriaByID(ctx, id)
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

	return s.store.CreateCategoria(ctx, &categoria)
}

func (s *ServiceCategoria) UpdateCategoria(ctx context.Context, id uuid.UUID, categoria models.Categoria) (*models.Categoria, error) {
	if categoria.Nombre == "" {
		s.logger.WarnContext(ctx, "Intento de actualización de categoría con nombre vacío", slog.String("id", id.String()))
		return nil, errors.New("falta el nombre de la categoria")
	}

	return s.store.UpdateCategoria(ctx, id, &categoria)
}

func (s *ServiceCategoria) DeleteCategoria(ctx context.Context, id uuid.UUID) error {
	return s.store.DeleteCategoria(ctx, id)
}
