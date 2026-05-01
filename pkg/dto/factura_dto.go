package dto

import (
	"github.com/google/uuid"
)

type FacturaResponse struct {
	IDFactura uuid.UUID `json:"id_factura"`
	FacNumero string    `json:"fac_numero"`
	Total     float64   `json:"total"`
	StatusMsg string    `json:"status_msg"`
}

type FacturaCompletaRequest struct {
	Cabecera       FacturaCabeceraRequest  `json:"cabecera"   validate:"required"`
	Detalles       []FacturaDetalleRequest `json:"detalles"   validate:"required,min=1,dive"`
	Pagos          []FacturaPagoRequest    `json:"pagos"      validate:"required,min=1,dive"`
	IdempotencyKey *uuid.UUID              `json:"idempotency_key,omitempty"`
}

type FacturaCabeceraRequest struct {
	FacNumero         string         `json:"fac_numero"          validate:"omitempty,max=50"`
	IDSucursal        *uuid.UUID     `json:"id_sucursal,omitempty"`
	IDCliente         *uuid.UUID     `json:"id_cliente,omitempty"`
	IDEstacion        *uuid.UUID     `json:"id_estacion,omitempty"`
	IDPeriodo         *uuid.UUID     `json:"id_periodo,omitempty"`
	IDControlEstacion *uuid.UUID     `json:"id_control_estacion,omitempty"`
	Subtotal          float64        `json:"subtotal"            validate:"gte=0"`
	Impuesto          float64        `json:"impuesto"            validate:"gte=0"`
	Total             float64        `json:"total"               validate:"gt=0"`
	Metadata          map[string]any `json:"metadata"`
}

type FacturaDetalleRequest struct {
	IDProducto     *uuid.UUID `json:"id_producto,omitempty" validate:"required"`
	IDLote         *uuid.UUID `json:"id_lote,omitempty"`
	Cantidad       float64    `json:"cantidad"              validate:"required,gt=0"`
	PrecioUnitario float64    `json:"precio_unitario"       validate:"required,gte=0"`
	Subtotal       float64    `json:"subtotal"              validate:"gte=0"`
	Impuesto       float64    `json:"impuesto"              validate:"gte=0"`
	Total          float64    `json:"total"                 validate:"gt=0"`
}

type FacturaPagoRequest struct {
	MetodoPago string  `json:"metodo_pago" validate:"required,min=1,max=100"`
	Monto      float64 `json:"monto"       validate:"required,gt=0"`
	Referencia string  `json:"referencia"  validate:"omitempty,max=255"`
}
