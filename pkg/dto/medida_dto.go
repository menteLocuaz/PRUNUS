package dto

import "github.com/google/uuid"

type MedidaCreateRequest struct {
	Nombre     string    `json:"nombre" validate:"required,min=1,max=100"`
	IDSucursal uuid.UUID `json:"id_sucursal" validate:"required"`
}

type MedidaUpdateRequest struct {
	Nombre     string    `json:"nombre" validate:"required,min=1,max=100"`
	IDSucursal uuid.UUID `json:"id_sucursal" validate:"required"`
}
