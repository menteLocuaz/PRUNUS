package store

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/prunus/pkg/dto"
	"github.com/prunus/pkg/models"
)

const (
	cacheKeyImpuestos  = "catalog:impuestos"
	cacheKeyFormasPago = "catalog:formas_pago"
	cacheTTLStatic     = 24 * time.Hour
)

type FacturaCacheDecorator struct {
	base  StoreFactura
	cache models.CacheStore
}

func NewFacturaCacheDecorator(base StoreFactura, cache models.CacheStore) StoreFactura {
	return &FacturaCacheDecorator{
		base:  base,
		cache: cache,
	}
}

func (d *FacturaCacheDecorator) CreateFactura(ctx context.Context, f *models.Factura, items []*models.DetalleFactura) (*models.Factura, error) {
	return d.base.CreateFactura(ctx, f, items)
}

func (d *FacturaCacheDecorator) RegistrarFacturaCompleta(ctx context.Context, req dto.FacturaCompletaRequest, idUsuario uuid.UUID) (*dto.FacturaResponse, error) {
	return d.base.RegistrarFacturaCompleta(ctx, req, idUsuario)
}

func (d *FacturaCacheDecorator) GetFacturaByID(ctx context.Context, id uuid.UUID) (*models.Factura, []*models.DetalleFactura, error) {
	return d.base.GetFacturaByID(ctx, id)
}

func (d *FacturaCacheDecorator) GetAllFacturas(ctx context.Context, params dto.PaginationParams) ([]*models.Factura, error) {
	return d.base.GetAllFacturas(ctx, params)
}

func (d *FacturaCacheDecorator) GetAllImpuestos(ctx context.Context) ([]*models.Impuesto, error) {
	if d.cache == nil {
		return d.base.GetAllImpuestos(ctx)
	}

	var impuestos []*models.Impuesto
	if err := d.cache.Get(ctx, cacheKeyImpuestos, &impuestos); err == nil {
		return impuestos, nil
	}

	impuestos, err := d.base.GetAllImpuestos(ctx)
	if err != nil {
		return nil, err
	}

	_ = d.cache.Set(ctx, cacheKeyImpuestos, impuestos, cacheTTLStatic)
	return impuestos, nil
}

func (d *FacturaCacheDecorator) GetAllFormasPago(ctx context.Context) ([]*models.FormaPago, error) {
	if d.cache == nil {
		return d.base.GetAllFormasPago(ctx)
	}

	var formas []*models.FormaPago
	if err := d.cache.Get(ctx, cacheKeyFormasPago, &formas); err == nil {
		return formas, nil
	}

	formas, err := d.base.GetAllFormasPago(ctx)
	if err != nil {
		return nil, err
	}

	_ = d.cache.Set(ctx, cacheKeyFormasPago, formas, cacheTTLStatic)
	return formas, nil
}
