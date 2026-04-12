package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/prunus/pkg/dto"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
	"github.com/prunus/pkg/utils"
)

type ServiceFactura struct {
	store  store.StoreFactura
	cache  *utils.CacheManager
	logger *slog.Logger
}

func NewServiceFactura(s store.StoreFactura, c *utils.CacheManager, logger *slog.Logger) *ServiceFactura {
	return &ServiceFactura{
		store:  s,
		cache:  c,
		logger: logger,
	}
}

const (
	cacheKeyImpuestos  = "facturas:impuestos"
	cacheKeyFormasPago = "facturas:formas_pago"
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
	return utils.GetOrSet(ctx, s.cache, cacheKeyImpuestos, cacheTTLStatic, func() ([]*models.Impuesto, error) {
		return s.store.GetAllImpuestos(ctx)
	})
}

func (s *ServiceFactura) GetFormasPago(ctx context.Context) ([]*models.FormaPago, error) {
	return utils.GetOrSet(ctx, s.cache, cacheKeyFormasPago, cacheTTLStatic, func() ([]*models.FormaPago, error) {
		return s.store.GetAllFormasPago(ctx)
	})
}

func (s *ServiceFactura) CreateFormaPago(ctx context.Context, f models.FormaPago) (*models.FormaPago, error) {
	res, err := s.store.CreateFormaPago(ctx, &f)
	if err == nil {
		s.cache.Invalidate(ctx, "facturas:")
	}
	return res, err
}

func (s *ServiceFactura) UpdateFormaPago(ctx context.Context, id uuid.UUID, f models.FormaPago) (*models.FormaPago, error) {
	res, err := s.store.UpdateFormaPago(ctx, id, &f)
	if err == nil {
		s.cache.Invalidate(ctx, "facturas:")
	}
	return res, err
}

func (s *ServiceFactura) DeleteFormaPago(ctx context.Context, id uuid.UUID) error {
	err := s.store.DeleteFormaPago(ctx, id)
	if err == nil {
		s.cache.Invalidate(ctx, "facturas:")
	}
	return err
}

func (s *ServiceFactura) GetFormaPagoByID(ctx context.Context, id uuid.UUID) (*models.FormaPago, error) {
	return s.store.GetFormaPagoByID(ctx, id)
}
