package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
)

type ProductoCreateRequest struct {
	Nombre           string    `json:"nombre" validate:"required,min=3,max=150"`
	Descripcion      string    `json:"descripcion" validate:"omitempty,max=500"`
	CodigoBarras     string    `json:"codigo_barras" validate:"omitempty,max=50"`
	SKU              string    `json:"sku" validate:"omitempty,max=50"`
	PrecioCompra     float64   `json:"precio_compra" validate:"required,gte=0"`
	PrecioVenta      float64   `json:"precio_venta" validate:"required,gte=0"`
	Stock            float64   `json:"stock" validate:"required,gte=0"`
	FechaVencimiento time.Time `json:"fecha_vencimiento" validate:"omitempty"`
	Imagen           string    `json:"imagen" validate:"omitempty"`
	IDStatus         uuid.UUID `json:"id_status"`
	IDSucursal       uuid.UUID `json:"id_sucursal" validate:"required"`
	IDCategoria      uuid.UUID `json:"id_categoria" validate:"required"`
	IDMoneda         uuid.UUID `json:"id_moneda" validate:"required"`
	IDUnidad         uuid.UUID `json:"id_unidad" validate:"required"`
}

func (r *ProductoCreateRequest) ToModel() models.Producto {
	var fechaVencimiento *time.Time
	if !r.FechaVencimiento.IsZero() {
		fechaVencimiento = &r.FechaVencimiento
	}

	return models.Producto{
		Nombre:           r.Nombre,
		Descripcion:      r.Descripcion,
		CodigoBarras:     r.CodigoBarras,
		SKU:              r.SKU,
		FechaVencimiento: fechaVencimiento,
		Imagen:           r.Imagen,
		IDStatus:         r.IDStatus,
		IDCategoria:      r.IDCategoria,
		IDMoneda:         r.IDMoneda,
		IDUnidad:         r.IDUnidad,
	}
}

type ProductoUpdateRequest struct {
	Nombre           string    `json:"nombre" validate:"required,min=3,max=150"`
	Descripcion      string    `json:"descripcion" validate:"omitempty,max=500"`
	CodigoBarras     string    `json:"codigo_barras" validate:"omitempty,max=50"`
	SKU              string    `json:"sku" validate:"omitempty,max=50"`
	FechaVencimiento time.Time `json:"fecha_vencimiento" validate:"omitempty"`
	Imagen           string    `json:"imagen" validate:"omitempty"`
	IDStatus         uuid.UUID `json:"id_status" validate:"required"`
	IDCategoria      uuid.UUID `json:"id_categoria" validate:"required"`
	IDMoneda         uuid.UUID `json:"id_moneda" validate:"required"`
	IDUnidad         uuid.UUID `json:"id_unidad" validate:"required"`
}

func (r *ProductoUpdateRequest) ToModel() models.Producto {
	var fechaVencimiento *time.Time
	if !r.FechaVencimiento.IsZero() {
		fechaVencimiento = &r.FechaVencimiento
	}

	return models.Producto{
		Nombre:           r.Nombre,
		Descripcion:      r.Descripcion,
		CodigoBarras:     r.CodigoBarras,
		SKU:              r.SKU,
		FechaVencimiento: fechaVencimiento,
		Imagen:           r.Imagen,
		IDStatus:         r.IDStatus,
		IDCategoria:      r.IDCategoria,
		IDMoneda:         r.IDMoneda,
		IDUnidad:         r.IDUnidad,
	}
}
