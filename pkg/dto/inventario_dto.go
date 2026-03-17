package dto

import (
	"github.com/google/uuid"
)

type InventarioCreateRequest struct {
	IDProducto   uuid.UUID `json:"id_producto" validate:"required"`
	IDSucursal   uuid.UUID `json:"id_sucursal" validate:"required"`
	StockActual  float64   `json:"stock_actual" validate:"required,gte=0"`
	StockMinimo  float64   `json:"stock_minimo" validate:"required,gte=0"`
	StockMaximo  float64   `json:"stock_maximo" validate:"required,gte=0"`
	PrecioCompra float64   `json:"precio_compra" validate:"required,gte=0"`
	PrecioVenta  float64   `json:"precio_venta" validate:"required,gte=0"`
}

type InventarioUpdateRequest struct {
	StockActual  float64 `json:"stock_actual" validate:"required,gte=0"`
	StockMinimo  float64 `json:"stock_minimo" validate:"required,gte=0"`
	StockMaximo  float64 `json:"stock_maximo" validate:"required,gte=0"`
	PrecioCompra float64 `json:"precio_compra" validate:"required,gte=0"`
	PrecioVenta  float64 `json:"precio_venta" validate:"required,gte=0"`
}
