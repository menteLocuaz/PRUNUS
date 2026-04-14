package models

import (
	"time"

	"github.com/google/uuid"
)

type FormaPagoFactura struct {
	IDPagoFac   uuid.UUID  `json:"id_pago_fac"`
	IDFactura   uuid.UUID  `json:"id_factura"`
	MetodoPago  string     `json:"metodo_pago"`
	Monto       float64    `json:"monto"`
	Referencia  string     `json:"referencia"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}
