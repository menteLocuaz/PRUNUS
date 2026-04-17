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
	Cabecera       FacturaCabeceraRequest  `json:"cabecera"`
	Detalles       []FacturaDetalleRequest `json:"detalles"`
	Pagos          []FacturaPagoRequest    `json:"pagos"`
	IdempotencyKey *uuid.UUID              `json:"idempotency_key,omitempty"`
}

type FacturaCabeceraRequest struct {
	FacNumero         string                 `json:"fac_numero"`
	IDSucursal        uuid.UUID              `json:"id_sucursal"`
	IDCliente         uuid.UUID              `json:"id_cliente"`
	IDEstacion        uuid.UUID              `json:"id_estacion"`
	IDPeriodo         uuid.UUID              `json:"id_periodo"`
	IDControlEstacion uuid.UUID              `json:"id_control_estacion"`
	Subtotal          float64                `json:"subtotal"`
	Impuesto          float64                `json:"impuesto"`
	Total             float64                `json:"total"`
	Metadata          map[string]interface{} `json:"metadata"`
}

type FacturaDetalleRequest struct {
	IDProducto     uuid.UUID  `json:"id_producto"`
	IDLote         *uuid.UUID `json:"id_lote,omitempty"`
	Cantidad       float64    `json:"cantidad"`
	PrecioUnitario float64    `json:"precio_unitario"`
	Subtotal       float64    `json:"subtotal"`
	Impuesto       float64    `json:"impuesto"`
	Total          float64    `json:"total"`
}

type FacturaPagoRequest struct {
	MetodoPago string  `json:"metodo_pago"`
	Monto      float64 `json:"monto"`
	Referencia string  `json:"referencia"`
}
