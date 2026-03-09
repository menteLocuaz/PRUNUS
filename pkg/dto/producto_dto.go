package dto

import (
	"time"

	"github.com/google/uuid"
)

type ProductoCreateRequest struct {
	Nombre           string    `json:"nombre" validate:"required,min=3,max=150"`
	Descripcion      string    `json:"descripcion" validate:"omitempty,max=500"`
	PrecioCompra     float64   `json:"precio_compra" validate:"required,gte=0"`
	PrecioVenta      float64   `json:"precio_venta" validate:"required,gte=0"`
	Stock            uint      `json:"stock" validate:"required,gte=0"`
	FechaVencimiento time.Time `json:"fecha_vencimiento" validate:"omitempty"`
	Imagen           string    `json:"imagen" validate:"omitempty"`
	IDStatus         uuid.UUID `json:"id_status" validate:"required"`
	IDSucursal       uuid.UUID `json:"id_sucursal" validate:"required"`
	IDCategoria      uuid.UUID `json:"id_categoria" validate:"required"`
	IDMoneda         uuid.UUID `json:"id_moneda" validate:"required"`
	IDUnidad         uuid.UUID `json:"id_unidad" validate:"required"`
}

type ProductoUpdateRequest struct {
	Nombre           string    `json:"nombre" validate:"required,min=3,max=150"`
	Descripcion      string    `json:"descripcion" validate:"omitempty,max=500"`
	PrecioCompra     float64   `json:"precio_compra" validate:"required,gte=0"`
	PrecioVenta      float64   `json:"precio_venta" validate:"required,gte=0"`
	Stock            uint      `json:"stock" validate:"required,gte=0"`
	FechaVencimiento time.Time `json:"fecha_vencimiento" validate:"omitempty"`
	Imagen           string    `json:"imagen" validate:"omitempty"`
	IDStatus         uuid.UUID `json:"id_status" validate:"required"`
	IDSucursal       uuid.UUID `json:"id_sucursal" validate:"required"`
	IDCategoria      uuid.UUID `json:"id_categoria" validate:"required"`
	IDMoneda         uuid.UUID `json:"id_moneda" validate:"required"`
	IDUnidad         uuid.UUID `json:"id_unidad" validate:"required"`
}
