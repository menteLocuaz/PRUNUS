package dto

import "github.com/google/uuid"

type MonedaCreateRequest struct {
	Nombre     string    `json:"nombre" validate:"required,min=1,max=100"`
	IDSucursal uuid.UUID `json:"id_sucursal" validate:"required"`
	IDStatus   uuid.UUID `json:"id_status" validate:"required"`
}

type MonedaUpdateRequest struct {
	Nombre     string    `json:"nombre" validate:"required,min=1,max=100"`
	IDSucursal uuid.UUID `json:"id_sucursal" validate:"required"`
	IDStatus   uuid.UUID `json:"id_status" validate:"required"`
}
