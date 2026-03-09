package models

import (
	"time"

	"github.com/google/uuid" // Nueva dependencia
)

type FormaPagoFactura struct {
	IDPagoFactura uuid.UUID  `json:"id_pago_factura"`
	IDFactura     uuid.UUID  `json:"id_factura"`
	IDFormaPago   uuid.UUID  `json:"id_forma_pago"`
	ValorBillete  float64    `json:"valor_billete"`
	TotalPagar    float64    `json:"total_pagar"`
	Fecha         time.Time  `json:"fecha"`
	IDUsuario     uuid.UUID  `json:"id_usuario"`
	IDStatus      uuid.UUID  `json:"id_status"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
}
