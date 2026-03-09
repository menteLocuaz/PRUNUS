package dto

import "github.com/google/uuid"

type CategoriaCreateRequest struct {
	Nombre     string    `json:"nombre" validate:"required,min=3,max=100"`
	IDSucursal uuid.UUID `json:"id_sucursal" validate:"required"`
}

type CategoriaUpdateRequest struct {
	Nombre     string    `json:"nombre" validate:"required,min=3,max=100"`
	IDSucursal uuid.UUID `json:"id_sucursal" validate:"required"`
}
