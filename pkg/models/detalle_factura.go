package models

import (
	"time"

	"github.com/google/uuid" // Nueva dependencia
)

type DetalleFactura struct {
	IDDetalleFactura uuid.UUID  `json:"id_detalle_factura"`
	IDFactura        uuid.UUID  `json:"id_factura"`
	IDProducto       uuid.UUID  `json:"id_producto"`
	Cantidad         float64    `json:"cantidad"`
	Precio           float64    `json:"precio"`
	Subtotal         float64    `json:"subtotal"`
	Impuesto         float64    `json:"impuesto"`
	Total            float64    `json:"total"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
}
