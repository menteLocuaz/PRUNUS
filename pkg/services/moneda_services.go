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

type ServiceMoneda struct {
	store  store.StoreMoneda
	cache  models.CacheStore
	logger *slog.Logger
}

func NewServiceMoneda(s store.StoreMoneda, c models.CacheStore, logger *slog.Logger) *ServiceMoneda {
	return &ServiceMoneda{
		store:  s,
		cache:  c,
		logger: logger,
	}
}

const (
	cacheKeyMonedasAll = "monedas:all"
	cacheKeyMonedaID   = "monedas:id:%s"
	cacheTTLMonedas    = 24 * time.Hour
)

func (s *ServiceMoneda) GetAllMonedas(ctx context.Context) ([]*models.Moneda, error) {
	var monedas []*models.Moneda

	// Intentar caché
	if err := s.cache.Get(ctx, cacheKeyMonedasAll, &monedas); err == nil {
		return monedas, nil
	}

	// BD
	monedas, err := s.store.GetAllMonedas(ctx)
	if err != nil {
		return nil, err
	}

	// Guardar en caché
	_ = s.cache.Set(ctx, cacheKeyMonedasAll, monedas, cacheTTLMonedas)

	return monedas, nil
}

func (s *ServiceMoneda) GetMonedaByID(ctx context.Context, id uuid.UUID) (*models.Moneda, error) {
	if id == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de obtener moneda con ID nulo")
		return nil, errors.New("el ID de la moneda es requerido")
	}

	var moneda *models.Moneda
	key := fmt.Sprintf(cacheKeyMonedaID, id.String())

	// Intentar caché
	if err := s.cache.Get(ctx, key, &moneda); err == nil {
		return moneda, nil
	}

	// BD
	moneda, err := s.store.GetMonedaByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Guardar en caché
	_ = s.cache.Set(ctx, key, moneda, cacheTTLMonedas)

	return moneda, nil
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

	res, err := s.store.CreateMoneda(ctx, &moneda)
	if err != nil {
		return nil, err
	}

	// Invalidar caché
	s.invalidateCache(ctx, res.IDMoneda)

	return res, nil
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

	res, err := s.store.UpdateMoneda(ctx, id, &moneda)
	if err != nil {
		return nil, err
	}

	// Invalidar caché
	s.invalidateCache(ctx, id)

	return res, nil
}

func (s *ServiceMoneda) DeleteMoneda(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de eliminación de moneda con ID nulo")
		return errors.New("el ID de la moneda es requerido")
	}

	if err := s.store.DeleteMoneda(ctx, id); err != nil {
		return err
	}

	// Invalidar caché
	s.invalidateCache(ctx, id)

	return nil
}

func (s *ServiceMoneda) invalidateCache(ctx context.Context, id uuid.UUID) {
	_ = s.cache.Delete(ctx, cacheKeyMonedasAll)
	_ = s.cache.Delete(ctx, fmt.Sprintf(cacheKeyMonedaID, id.String()))
}
