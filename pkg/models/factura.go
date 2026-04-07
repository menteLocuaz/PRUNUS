package models

import (
	"time"

	"github.com/google/uuid" // Nueva dependencia
)

type Factura struct {
	IDFactura         uuid.UUID              `json:"id_factura"`
	CfacFechaCreacion time.Time              `json:"cfac_fecha_creacion"`
	FacNumero         string                 `json:"fac_numero"`
	CfacSubtotal      float64                `json:"cfac_subtotal"`
	CfacIVA           float64                `json:"cfac_iva"`
	CfacTotal         float64                `json:"cfac_total"`
	CfacObservacion   string                 `json:"cfac_observacion"`
	IDUserPos         uuid.UUID              `json:"id_user_pos"`
	IDEstacion        uuid.UUID              `json:"id_estacion"`
	IDOrdenPedido     uuid.UUID              `json:"id_orden_pedido"`
	IDCliente         uuid.UUID              `json:"id_cliente"`
	IDMotivoAnulacion *uuid.UUID             `json:"id_motivo_anulacion,omitempty"`
	IDPeriodo         uuid.UUID              `json:"id_periodo"`
	IDControlEstacion uuid.UUID              `json:"id_control_estacion"`
	IDStatus          uuid.UUID              `json:"id_status"`
	FechaOperacion    time.Time              `json:"fecha_operacion"`
	FechaVencimiento  *time.Time             `json:"fecha_vencimiento,omitempty"`
	BaseImpuesto      float64                `json:"base_impuesto"`
	Impuesto          float64                `json:"impuesto"`
	ValorImpuesto     float64                `json:"valor_impuesto"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
	DeletedAt         *time.Time             `json:"deleted_at,omitempty"`
}
