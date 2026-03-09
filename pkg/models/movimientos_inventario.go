package models

import (
	"time"

	"github.com/google/uuid"
)

type MovimientoInventario struct {
	IDMovimiento   uuid.UUID  `json:"id_movimiento"`
	IDProducto     uuid.UUID  `json:"id_producto"`
	TipoMovimiento string     `json:"tipo_movimiento"` // VENTA, COMPRA, AJUSTE, DEVOLUCION
	Cantidad       float64    `json:"cantidad"`
	Fecha          time.Time  `json:"fecha"`
	IDUsuario      uuid.UUID  `json:"id_usuario"`
	Referencia     string     `json:"referencia"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}
