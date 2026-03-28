package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/prunus/pkg/dto"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
)

type ServiceFactura struct {
	store  store.StoreFactura
	cache  models.CacheStore
	logger *slog.Logger
}

func NewServiceFactura(s store.StoreFactura, c models.CacheStore, logger *slog.Logger) *ServiceFactura {
	return &ServiceFactura{
		store:  s,
		cache:  c,
		logger: logger,
	}
}

const (
	cacheKeyImpuestos  = "catalog:impuestos"
	cacheKeyFormasPago = "catalog:formas_pago"
	cacheTTLStatic     = 24 * time.Hour
)

func (s *ServiceFactura) CreateFactura(ctx context.Context, f models.Factura, items []*models.DetalleFactura) (*models.Factura, error) {
	return s.store.CreateFactura(ctx, &f, items)
}

func (s *ServiceFactura) RegistrarFacturaCompleta(ctx context.Context, req dto.FacturaCompletaRequest, idUsuario uuid.UUID) (*dto.FacturaResponse, error) {
	return s.store.RegistrarFacturaCompleta(ctx, req, idUsuario)
}

func (s *ServiceFactura) GetFactura(ctx context.Context, id uuid.UUID) (*models.Factura, []*models.DetalleFactura, error) {
	return s.store.GetFacturaByID(ctx, id)
}

func (s *ServiceFactura) GetAllFacturas(ctx context.Context, params dto.PaginationParams) ([]*models.Factura, error) {
	return s.store.GetAllFacturas(ctx, params)
}

func (s *ServiceFactura) GetImpuestos(ctx context.Context) ([]*models.Impuesto, error) {
	var impuestos []*models.Impuesto

	// Intentar caché
	if err := s.cache.Get(ctx, cacheKeyImpuestos, &impuestos); err == nil {
		return impuestos, nil
	}

	// BD
	impuestos, err := s.store.GetAllImpuestos(ctx)
	if err != nil {
		return nil, err
	}

	// Guardar en caché
	_ = s.cache.Set(ctx, cacheKeyImpuestos, impuestos, cacheTTLStatic)

	return impuestos, nil
}

func (s *ServiceFactura) GetFormasPago(ctx context.Context) ([]*models.FormaPago, error) {
	var formas []*models.FormaPago

	// Intentar caché
	if err := s.cache.Get(ctx, cacheKeyFormasPago, &formas); err == nil {
		return formas, nil
	}

	// BD
	formas, err := s.store.GetAllFormasPago(ctx)
	if err != nil {
		return nil, err
	}

	// Guardar en caché
	_ = s.cache.Set(ctx, cacheKeyFormasPago, formas, cacheTTLStatic)

	return formas, nil
}
