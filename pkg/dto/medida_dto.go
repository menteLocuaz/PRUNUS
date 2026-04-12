package dto

import "github.com/google/uuid"

type MedidaCreateRequest struct {
	Nombre      string    `json:"nombre" validate:"required,min=1,max=50"`
	Abreviatura string    `json:"abreviatura" validate:"required,max=10"`
	IDStatus    uuid.UUID `json:"id_status" validate:"required"`
	IDSucursal  uuid.UUID `json:"id_sucursal" validate:"required"`
}

type MedidaUpdateRequest struct {
	Nombre      string    `json:"nombre" validate:"required,min=1,max=50"`
	Abreviatura string    `json:"abreviatura" validate:"required,max=10"`
	IDStatus    uuid.UUID `json:"id_status" validate:"required"`
	IDSucursal  uuid.UUID `json:"id_sucursal" validate:"required"`
}
