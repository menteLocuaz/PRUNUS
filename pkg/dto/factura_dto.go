package dto

import (
	"github.com/google/uuid"
)

// FacturaCompletaRequest representa la petición para registrar una factura con todo su detalle y pagos.
type FacturaCompletaRequest struct {
	Cabecera FacturaCabeceraJSON  `json:"cabecera" validate:"required"`
	Detalles []FacturaDetalleJSON `json:"detalles" validate:"required,min=1"`
	Pagos    []FacturaPagoJSON    `json:"pagos" validate:"required,min=1"`
}

// FacturaCabeceraJSON mapea los campos que espera la función factura_registrar_completa para la cabecera.
type FacturaCabeceraJSON struct {
	FacNumero         string                 `json:"fac_numero"`
	Subtotal          float64                `json:"subtotal" validate:"required,gte=0"`
	Iva               float64                `json:"iva" validate:"gte=0"`
	Total             float64                `json:"total" validate:"required,gt=0"`
	Observacion       string                 `json:"observacion"`
	IDEstacion        uuid.UUID              `json:"id_estacion" validate:"required,uuid"`
	IDOrdenPedido     uuid.UUID              `json:"id_orden_pedido" validate:"omitempty,uuid"`
	IDCliente         uuid.UUID              `json:"id_cliente" validate:"required,uuid"`
	IDPeriodo         uuid.UUID              `json:"id_periodo" validate:"required,uuid"`
	IDControlEstacion uuid.UUID              `json:"id_control_estacion" validate:"required,uuid"`
	BaseImpuesto      float64                `json:"base_impuesto" validate:"gte=0"`
	Impuesto          float64                `json:"impuesto" validate:"gte=0"`
	ValorImpuesto     float64                `json:"valor_impuesto" validate:"gte=0"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
}

// FacturaDetalleJSON mapea los items del detalle para la función de base de datos.
type FacturaDetalleJSON struct {
	IDProducto uuid.UUID `json:"id_producto" validate:"required,uuid"`
	Cantidad   float64   `json:"cantidad" validate:"required,gt=0"`
	Precio     float64   `json:"precio" validate:"required,gte=0"`
	Subtotal   float64   `json:"subtotal" validate:"required,gte=0"`
	Impuesto   float64   `json:"impuesto" validate:"gte=0"`
	Total      float64   `json:"total" validate:"required,gte=0"`
}

// FacturaPagoJSON mapea los pagos realizados para la función de base de datos.
type FacturaPagoJSON struct {
	IDFormaPago  uuid.UUID `json:"id_forma_pago" validate:"required,uuid"`
	ValorBillete float64   `json:"valor_billete" validate:"required,gte=0"`
	TotalPagar   float64   `json:"total_pagar" validate:"required,gt=0"`
}

// FacturaResponse es la respuesta tras registrar una factura mediante la función almacenada.
type FacturaResponse struct {
	IDFactura uuid.UUID `json:"id_factura"`
	FacNumero string    `json:"fac_numero"`
	Total     float64   `json:"total"`
	StatusMsg string    `json:"status_msg"`
}
