package models

import (
	"time"

	"github.com/google/uuid"
)

type Inventario struct {
	IDInventario uuid.UUID  `json:"id_inventario"`
	IDProducto   uuid.UUID  `json:"id_producto"`
	IDSucursal   uuid.UUID  `json:"id_sucursal"`
	StockActual  float64    `json:"stock_actual"`
	StockMinimo  float64    `json:"stock_minimo"`
	StockMaximo  float64    `json:"stock_maximo"`
	PrecioCompra float64    `json:"precio_compra"`
	PrecioVenta  float64    `json:"precio_venta"`
	Ubicacion    string     `json:"ubicacion"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}
