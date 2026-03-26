package dto

import (
	"github.com/google/uuid"
)

// FacturaCompletaRequest representa la petición para registrar una factura con todo su detalle y pagos.
type FacturaCompletaRequest struct {
	Cabecera FacturaCabeceraJSON `json:"cabecera" validate:"required"`
	Detalles []FacturaDetalleJSON `json:"detalles" validate:"required,min=1"`
	Pagos    []FacturaPagoJSON    `json:"pagos" validate:"required,min=1"`
}

// FacturaCabeceraJSON mapea los campos que espera la función factura_registrar_completa para la cabecera.
type FacturaCabeceraJSON struct {
	FacNumero         string                 `json:"fac_numero"`
	Subtotal          float64                `json:"subtotal"`
	Iva               float64                `json:"iva"`
	Total             float64                `json:"total"`
	Observacion       string                 `json:"observacion"`
	IDEstacion        uuid.UUID              `json:"id_estacion"`
	IDOrdenPedido     uuid.UUID              `json:"id_orden_pedido"`
	IDCliente         uuid.UUID              `json:"id_cliente"`
	IDPeriodo         uuid.UUID              `json:"id_periodo"`
	IDControlEstacion uuid.UUID              `json:"id_control_estacion"`
	BaseImpuesto      float64                `json:"base_impuesto"`
	Impuesto          float64                `json:"impuesto"`
	ValorImpuesto     float64                `json:"valor_impuesto"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
}

// FacturaDetalleJSON mapea los items del detalle para la función de base de datos.
type FacturaDetalleJSON struct {
	IDProducto uuid.UUID `json:"id_producto"`
	Cantidad   float64   `json:"cantidad"`
	Precio     float64   `json:"precio"`
	Subtotal   float64   `json:"subtotal"`
	Impuesto   float64   `json:"impuesto"`
	Total      float64   `json:"total"`
}

// FacturaPagoJSON mapea los pagos realizados para la función de base de datos.
type FacturaPagoJSON struct {
	IDFormaPago  uuid.UUID `json:"id_forma_pago"`
	ValorBillete float64   `json:"valor_billete"`
	TotalPagar   float64   `json:"total_pagar"`
}

// FacturaResponse es la respuesta tras registrar una factura mediante la función almacenada.
type FacturaResponse struct {
	IDFactura  uuid.UUID `json:"id_factura"`
	FacNumero  string    `json:"fac_numero"`
	Total      float64   `json:"total"`
	StatusMsg  string    `json:"status_msg"`
}
