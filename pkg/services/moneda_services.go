package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
)

type ServiceMoneda struct {
	store  store.StoreMoneda
	logger *slog.Logger
}

func NewServiceMoneda(s store.StoreMoneda, logger *slog.Logger) *ServiceMoneda {
	return &ServiceMoneda{
		store:  s,
		logger: logger,
	}
}

func (s *ServiceMoneda) GetAllMonedas(ctx context.Context) ([]*models.Moneda, error) {
	return s.store.GetAllMonedas(ctx)
}

func (s *ServiceMoneda) GetMonedaByID(ctx context.Context, id uuid.UUID) (*models.Moneda, error) {
	if id == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de obtener moneda con ID nulo")
		return nil, errors.New("el ID de la moneda es requerido")
	}
	return s.store.GetMonedaByID(ctx, id)
}

func (s *ServiceMoneda) CreateMoneda(ctx context.Context, moneda models.Moneda) (*models.Moneda, error) {
	if moneda.Nombre == "" {
		s.logger.WarnContext(ctx, "Intento de creación de moneda con nombre vacío")
		return nil, errors.New("falta el nombre de la moneda")
	}
	if moneda.IDSucursal == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de creación de moneda sin sucursal", slog.String("nombre", moneda.Nombre))
		return nil, errors.New("falta el id de la sucursal")
	}
	if moneda.IDStatus == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de creación de moneda sin estatus", slog.String("nombre", moneda.Nombre))
		return nil, errors.New("falta el id del estatus")
	}
	return s.store.CreateMoneda(ctx, &moneda)
}

func (s *ServiceMoneda) UpdateMoneda(ctx context.Context, id uuid.UUID, moneda models.Moneda) (*models.Moneda, error) {
	if id == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de actualización de moneda con ID nulo")
		return nil, errors.New("el ID de la moneda es requerido")
	}
	if moneda.Nombre == "" {
		s.logger.WarnContext(ctx, "Intento de actualización de moneda con nombre vacío", slog.String("id", id.String()))
		return nil, errors.New("falta el nombre de la moneda")
	}
	return s.store.UpdateMoneda(ctx, id, &moneda)
}

func (s *ServiceMoneda) DeleteMoneda(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de eliminación de moneda con ID nulo")
		return errors.New("el ID de la moneda es requerido")
	}
	return s.store.DeleteMoneda(ctx, id)
}
