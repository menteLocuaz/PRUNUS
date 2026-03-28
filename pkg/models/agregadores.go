package models

import (
	"time"

	"github.com/google/uuid" // Nueva dependencia
)

type Agregador struct {
	IDAgregador uuid.UUID  `json:"id_agregador"`
	Nombre      string     `json:"nombre"`
	Descripcion string     `json:"descripcion"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type OrdenAgregador struct {
	IDOrdenAgregador  uuid.UUID              `json:"id_orden_agregador"`
	IDOrdenPedido     uuid.UUID              `json:"id_orden_pedido"`
	IDAgregador       uuid.UUID              `json:"id_agregador"`
	CodigoExterno     string                 `json:"codigo_externo"`
	DatosAgregador    map[string]interface{} `json:"datos_agregador"` // JSONB
	ComisionAgregador float64                `json:"comision_agregador"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
	DeletedAt         *time.Time             `json:"deleted_at,omitempty"`
}
