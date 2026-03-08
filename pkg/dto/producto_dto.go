package dto

import "time"

type ProductoCreateRequest struct {
	Nombre           string    `json:"nombre" validate:"required,min=3,max=150"`
	Descripcion      string    `json:"descripcion" validate:"omitempty,max=500"`
	PrecioCompra     float64   `json:"precio_compra" validate:"required,gte=0"`
	PrecioVenta      float64   `json:"precio_venta" validate:"required,gte=0"`
	Stock            uint      `json:"stock" validate:"required,gte=0"`
	FechaVencimiento time.Time `json:"fecha_vencimiento" validate:"omitempty"`
	Imagen           string    `json:"imagen" validate:"omitempty"`
	Estado           int       `json:"estado" validate:"required,oneof=0 1"`
	IDSucursal       uint      `json:"id_sucursal" validate:"required"`
	IDCategoria      uint      `json:"id_categoria" validate:"required"`
	IDMoneda         uint      `json:"id_moneda" validate:"required"`
	IDUnidad         uint      `json:"id_unidad" validate:"required"`
}

type ProductoUpdateRequest struct {
	Nombre           string    `json:"nombre" validate:"required,min=3,max=150"`
	Descripcion      string    `json:"descripcion" validate:"omitempty,max=500"`
	PrecioCompra     float64   `json:"precio_compra" validate:"required,gte=0"`
	PrecioVenta      float64   `json:"precio_venta" validate:"required,gte=0"`
	Stock            uint      `json:"stock" validate:"required,gte=0"`
	FechaVencimiento time.Time `json:"fecha_vencimiento" validate:"omitempty"`
	Imagen           string    `json:"imagen" validate:"omitempty"`
	Estado           int       `json:"estado" validate:"required,oneof=0 1"`
	IDSucursal       uint      `json:"id_sucursal" validate:"required"`
	IDCategoria      uint      `json:"id_categoria" validate:"required"`
	IDMoneda         uint      `json:"id_moneda" validate:"required"`
	IDUnidad         uint      `json:"id_unidad" validate:"required"`
}
