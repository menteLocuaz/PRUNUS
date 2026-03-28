package dto

import (
	"github.com/google/uuid"
)

type AgregadorCreateRequest struct {
	Nombre      string `json:"nombre" validate:"required,min=2,max=100"`
	Descripcion string `json:"descripcion" validate:"omitempty,max=255"`
}

type AgregadorUpdateRequest struct {
	Nombre      string `json:"nombre" validate:"required,min=2,max=100"`
	Descripcion string `json:"descripcion" validate:"omitempty,max=255"`
}

type OrdenAgregadorCreateRequest struct {
	IDOrdenPedido     uuid.UUID              `json:"id_orden_pedido" validate:"required"`
	IDAgregador       uuid.UUID              `json:"id_agregador" validate:"required"`
	CodigoExterno     string                 `json:"codigo_externo" validate:"required"`
	DatosAgregador    map[string]interface{} `json:"datos_agregador" validate:"omitempty"`
	ComisionAgregador float64                `json:"comision_agregador" validate:"omitempty"`
}
