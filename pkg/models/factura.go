package models

import (
	"time"

	"github.com/google/uuid"
)

type Factura struct {
	IDFactura         uuid.UUID              `json:"id_factura"`
	FacNumero         string                 `json:"fac_numero"`
	Subtotal          float64                `json:"subtotal"`
	Impuesto          float64                `json:"impuesto"`
	Total             float64                `json:"total"`
	Observacion       string                 `json:"observacion"`
	IDUsuario         uuid.UUID              `json:"id_usuario"`
	IDEstacion        uuid.UUID              `json:"id_estacion"`
	IDOrdenPedido     uuid.UUID              `json:"id_orden_pedido"`
	IDCliente         uuid.UUID              `json:"id_cliente"`
	IDPeriodo         uuid.UUID              `json:"id_periodo"`
	IDControlEstacion uuid.UUID              `json:"id_control_estacion"`
	IDStatus          uuid.UUID              `json:"id_status"`
	IDSucursal        uuid.UUID              `json:"id_sucursal"`
	FechaOperacion    time.Time              `json:"fecha_operacion"`
	FechaVencimiento  *time.Time             `json:"fecha_vencimiento,omitempty"`
	BaseImpuesto      float64                `json:"base_impuesto"`
	ValorImpuesto     float64                `json:"valor_impuesto"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
	DeletedAt         *time.Time             `json:"deleted_at,omitempty"`
}
