package models

import (
	"time"

	"github.com/google/uuid"
)

type OrdenCompra struct {
	IDOrdenCompra  uuid.UUID  `json:"id_orden_compra"`
	NumeroOrden    string     `json:"numero_orden"`
	IDProveedor    uuid.UUID  `json:"id_proveedor"`
	IDSucursal     uuid.UUID  `json:"id_sucursal"`
	IDUsuario      uuid.UUID  `json:"id_usuario"`
	IDMoneda       uuid.UUID  `json:"id_moneda"`
	IDStatus       uuid.UUID  `json:"id_status"`
	FechaEmision   time.Time  `json:"fecha_emision"`
	FechaRecepcion *time.Time `json:"fecha_recepcion,omitempty"`
	Subtotal       float64    `json:"subtotal"`
	Impuesto       float64    `json:"impuesto"`
	Total          float64    `json:"total"`
	Observaciones  string     `json:"observaciones"`

	// Relaciones
	Proveedor *Proveedor            `json:"proveedor,omitempty"`
	Sucursal  *Sucursal             `json:"sucursal,omitempty"`
	Usuario   *Usuario              `json:"usuario,omitempty"`
	Detalles  []*DetalleOrdenCompra `json:"detalles,omitempty"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type DetalleOrdenCompra struct {
	IDDetalleCompra  uuid.UUID `json:"id_detalle_compra"`
	IDOrdenCompra    uuid.UUID `json:"id_orden_compra"`
	IDProducto       uuid.UUID `json:"id_producto"`
	CantidadPedida   float64   `json:"cantidad_pedida"`
	CantidadRecibida float64   `json:"cantidad_recibida"`
	PrecioUnitario   float64   `json:"precio_unitario"`
	Impuesto         float64   `json:"impuesto"`
	Total            float64   `json:"total"`

	// Relación
	Producto *Producto `json:"producto,omitempty"`
}
