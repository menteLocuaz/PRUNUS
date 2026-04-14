package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
	"github.com/prunus/pkg/utils"
	zaplogger "github.com/prunus/pkg/utils/logger"
	"go.uber.org/zap"
)

type ServiceMoneda struct {
	store  store.StoreMoneda
	cache  *utils.CacheManager
	logger *zap.Logger
}

func NewServiceMoneda(s store.StoreMoneda, c *utils.CacheManager, logger *zap.Logger) *ServiceMoneda {
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
	return utils.GetOrSet(ctx, s.cache, cacheKeyMonedasAll, cacheTTLMonedas, func() ([]*models.Moneda, error) {
		return s.store.GetAllMonedas(ctx)
	})
}

func (s *ServiceMoneda) GetMonedaByID(ctx context.Context, id uuid.UUID) (*models.Moneda, error) {
	if id == uuid.Nil {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de obtener moneda con ID nulo")
		return nil, errors.New("el ID de la moneda es requerido")
	}

	key := fmt.Sprintf(cacheKeyMonedaID, id.String())
	return utils.GetOrSet(ctx, s.cache, key, cacheTTLMonedas, func() (*models.Moneda, error) {
		return s.store.GetMonedaByID(ctx, id)
	})
}

func (s *ServiceMoneda) CreateMoneda(ctx context.Context, moneda models.Moneda) (*models.Moneda, error) {
	if moneda.Nombre == "" {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de creación de moneda con nombre vacío")
		return nil, errors.New("falta el nombre de la moneda")
	}
	if moneda.IDSucursal == uuid.Nil {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de creación de moneda sin sucursal", zap.String("nombre", moneda.Nombre))
		return nil, errors.New("falta el id de la sucursal")
	}

	if moneda.IDStatus == uuid.Nil {
		moneda.IDStatus = models.EstatusGlobalActivo
	}

	res, err := s.store.CreateMoneda(ctx, &moneda)
	if err != nil {
		return nil, err
	}

	s.cache.Invalidate(ctx, "monedas:")

	return res, nil
}

func (s *ServiceMoneda) UpdateMoneda(ctx context.Context, id uuid.UUID, moneda models.Moneda) (*models.Moneda, error) {
	if id == uuid.Nil {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de actualización de moneda con ID nulo")
		return nil, errors.New("el ID de la moneda es requerido")
	}
	if moneda.Nombre == "" {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de actualización de moneda con nombre vacío", zap.String("id", id.String()))
		return nil, errors.New("falta el nombre de la moneda")
	}

	res, err := s.store.UpdateMoneda(ctx, id, &moneda)
	if err != nil {
		return nil, err
	}

	s.cache.Invalidate(ctx, "monedas:")

	return res, nil
}

func (s *ServiceMoneda) DeleteMoneda(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de eliminación de moneda con ID nulo")
		return errors.New("el ID de la moneda es requerido")
	}

	if err := s.store.DeleteMoneda(ctx, id); err != nil {
		return err
	}

	s.cache.Invalidate(ctx, "monedas:")

	return nil
}
