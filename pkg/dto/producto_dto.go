package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
)

type ProductoCreateRequest struct {
	Nombre           string          `json:"pro_nombre" validate:"required,min=3,max=150"`
	Descripcion      string          `json:"pro_descripcion" validate:"omitempty,max=500"`
	CodigoBarras     string          `json:"pro_codigo" validate:"omitempty,max=50"`
	SKU              string          `json:"sku" validate:"omitempty,max=50"`
	PrecioCompra     float64         `json:"precio_compra" validate:"required,gte=0"`
	PrecioVenta      float64         `json:"precio_venta" validate:"required,gte=0"`
	Stock            float64         `json:"stock" validate:"required,gte=0"`
	FechaVencimiento models.JSONDate `json:"fecha_vencimiento" validate:"omitempty"`
	Imagen           string          `json:"imagen" validate:"omitempty"`
	IDStatus         uuid.UUID       `json:"id_status"`
	IDSucursal       uuid.UUID       `json:"id_sucursal" validate:"required"`
	IDCategoria      uuid.UUID       `json:"id_categoria" validate:"required"`
	IDMoneda         uuid.UUID       `json:"id_moneda" validate:"required"`
	IDUnidad         uuid.UUID       `json:"id_unidad" validate:"required"`
}

func (r *ProductoCreateRequest) ToModel() models.Producto {
	fechaVencimiento := r.FechaVencimiento.ToTime()
	var fechaPtr *time.Time
	if !fechaVencimiento.IsZero() {
		fechaPtr = &fechaVencimiento
	}

	return models.Producto{
		Nombre:           r.Nombre,
		Descripcion:      r.Descripcion,
		CodigoBarras:     r.CodigoBarras,
		SKU:              r.SKU,
		FechaVencimiento: fechaPtr,
		Imagen:           r.Imagen,
		IDStatus:         r.IDStatus,
		IDCategoria:      r.IDCategoria,
		IDMoneda:         r.IDMoneda,
		IDUnidad:         r.IDUnidad,
	}
}

type ProductoUpdateRequest struct {
	Nombre           string          `json:"pro_nombre" validate:"required,min=3,max=150"`
	Descripcion      string          `json:"pro_descripcion" validate:"omitempty,max=500"`
	CodigoBarras     string          `json:"pro_codigo" validate:"omitempty,max=50"`
	SKU              string          `json:"sku" validate:"omitempty,max=50"`
	FechaVencimiento models.JSONDate `json:"fecha_vencimiento" validate:"omitempty"`
	Imagen           string          `json:"imagen" validate:"omitempty"`
	IDStatus         uuid.UUID       `json:"id_status" validate:"required"`
	IDCategoria      uuid.UUID       `json:"id_categoria" validate:"required"`
	IDMoneda         uuid.UUID       `json:"id_moneda" validate:"required"`
	IDUnidad         uuid.UUID       `json:"id_unidad" validate:"required"`
}

func (r *ProductoUpdateRequest) ToModel() models.Producto {
	fechaVencimiento := r.FechaVencimiento.ToTime()
	var fechaPtr *time.Time
	if !fechaVencimiento.IsZero() {
		fechaPtr = &fechaVencimiento
	}

	return models.Producto{
		Nombre:           r.Nombre,
		Descripcion:      r.Descripcion,
		CodigoBarras:     r.CodigoBarras,
		SKU:              r.SKU,
		FechaVencimiento: fechaPtr,
		Imagen:           r.Imagen,
		IDStatus:         r.IDStatus,
		IDCategoria:      r.IDCategoria,
		IDMoneda:         r.IDMoneda,
		IDUnidad:         r.IDUnidad,
	}
}
