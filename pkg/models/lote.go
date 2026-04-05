package models

import (
	"time"

	"github.com/google/uuid"
)

type Lote struct {
	IDLote           uuid.UUID  `json:"id_lote"`
	IDProducto       uuid.UUID  `json:"id_producto"`
	IDSucursal       uuid.UUID  `json:"id_sucursal"`
	CodigoLote       string     `json:"codigo_lote"`
	CantidadInicial  float64    `json:"cantidad_inicial"`
	CantidadActual   float64    `json:"cantidad_actual"`
	CostoCompra      float64    `json:"costo_compra"`
	FechaVencimiento *time.Time `json:"fecha_vencimiento,omitempty"`
	FechaRecepcion   time.Time  `json:"fecha_recepcion"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
}
