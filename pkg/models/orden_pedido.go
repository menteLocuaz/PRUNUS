package models

import (
	"time"

	"github.com/google/uuid"
)

type OrdenPedido struct {
	IDOrdenPedido    uuid.UUID  `json:"id_orden_pedido"`
	OdpFechaCreacion time.Time  `json:"odp_fecha_creacion"`
	OdpObservacion   string     `json:"odp_observacion"`
	IDUserPos        uuid.UUID  `json:"id_user_pos"`
	IDPeriodo        uuid.UUID  `json:"id_periodo"`
	IDEstacion       uuid.UUID  `json:"id_estacion"`
	IDStatus         uuid.UUID  `json:"id_status"`
	Direccion        string     `json:"direccion"`
	Canal            string     `json:"canal"`
	OdpTotal         float64    `json:"odp_total"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
}
