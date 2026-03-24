package dto

import (
	"github.com/google/uuid"
)

type OrdenCompraCreateRequest struct {
	NumeroOrden   string                `json:"numero_orden" validate:"required"`
	IDProveedor   uuid.UUID             `json:"id_proveedor" validate:"required"`
	IDSucursal    uuid.UUID             `json:"id_sucursal" validate:"required"`
	IDMoneda      uuid.UUID             `json:"id_moneda" validate:"required"`
	IDStatus      uuid.UUID             `json:"id_status" validate:"required"`
	Observaciones string                `json:"observaciones"`
	Detalles      []DetalleCompraRequest `json:"detalles" validate:"required,min=1"`
}

type DetalleCompraRequest struct {
	IDProducto     uuid.UUID `json:"id_producto" validate:"required"`
	CantidadPedida float64   `json:"cantidad_pedida" validate:"required,gt=0"`
	PrecioUnitario float64   `json:"precio_unitario" validate:"required,gt=0"`
	Impuesto       float64   `json:"impuesto"`
}

type RecepcionCompraRequest struct {
	IDOrdenCompra uuid.UUID              `json:"id_orden_compra" validate:"required"`
	IDStatus      uuid.UUID              `json:"id_status" validate:"required"` // Estado "RECIBIDO"
	Items         []DetalleRecepcionRequest `json:"items" validate:"required,min=1"`
}

type DetalleRecepcionRequest struct {
	IDDetalleCompra  uuid.UUID `json:"id_detalle_compra" validate:"required"`
	IDProducto       uuid.UUID `json:"id_producto" validate:"required"`
	CantidadRecibida float64   `json:"cantidad_recibida" validate:"required,gte=0"`
}
