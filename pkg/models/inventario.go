package models

import (
	"time"

	"github.com/google/uuid"
)

// Inventario representa el stock de un producto en una sucursal específica.
// Sincronizado con la tabla 'inventario' (Migración 000003, 000017, 000021).
type Inventario struct {
	IDInventario uuid.UUID  `json:"id_inventario"`
	IDProducto   uuid.UUID  `json:"id_producto"`
	IDSucursal   uuid.UUID  `json:"id_sucursal"`
	StockActual  float64    `json:"stock_actual"` // DB: NUMERIC(12,2)
	StockMinimo  float64    `json:"stock_minimo"` // DB: NUMERIC(12,2)
	PrecioCompra float64    `json:"precio_compra"` // DB: NUMERIC(18,2)
	PrecioVenta  float64    `json:"precio_venta"`  // DB: NUMERIC(18,2)
	Ubicacion    string     `json:"ubicacion"`     // DB: VARCHAR(100)
	
	// Campos para auditoría y eliminación lógica
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}
