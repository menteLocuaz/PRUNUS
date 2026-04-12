package dto

import "github.com/google/uuid"

type CategoriaCreateRequest struct {
	Nombre     string    `json:"nombre" validate:"required,min=3,max=100"`
	IDStatus   uuid.UUID `json:"id_status" validate:"required"`
	IDSucursal uuid.UUID `json:"id_sucursal" validate:"required"`
}

type CategoriaUpdateRequest struct {
	Nombre     string    `json:"nombre" validate:"required,min=3,max=100"`
	IDStatus   uuid.UUID `json:"id_status" validate:"required"`
	IDSucursal uuid.UUID `json:"id_sucursal" validate:"required"`
}
