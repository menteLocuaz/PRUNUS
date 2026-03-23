package models

import (
	"time"

	"github.com/google/uuid"
)

type MovimientoInventario struct {
	IDMovimiento   uuid.UUID  `json:"id_movimiento"`
	IDProducto     uuid.UUID  `json:"id_producto"`
	IDSucursal     uuid.UUID  `json:"id_sucursal"`
	TipoMovimiento string     `json:"tipo_movimiento"` // VENTA, COMPRA, AJUSTE, DEVOLUCION, TRASLADO
	Cantidad       float64    `json:"cantidad"`
	CostoUnitario  float64    `json:"costo_unitario"`
	PrecioUnitario float64    `json:"precio_unitario"`
	StockAnterior  float64    `json:"stock_anterior"`
	StockPosterior float64    `json:"stock_posterior"`
	Fecha          time.Time  `json:"fecha"`
	IDUsuario      uuid.UUID  `json:"id_usuario"`
	Referencia     string     `json:"referencia"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}
