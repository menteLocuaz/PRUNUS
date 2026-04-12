package dto

import "github.com/google/uuid"

type MonedaCreateRequest struct {
	Nombre     string    `json:"nombre" validate:"required,min=1,max=50"`
	CodigoISO  string    `json:"codigo_iso" validate:"required,len=3"`
	Simbolo    string    `json:"simbolo" validate:"required,max=5"`
	IDSucursal uuid.UUID `json:"id_sucursal"`
	IDStatus   uuid.UUID `json:"id_status"`
}

type MonedaUpdateRequest struct {
	Nombre     string    `json:"nombre" validate:"required,min=1,max=50"`
	CodigoISO  string    `json:"codigo_iso" validate:"required,len=3"`
	Simbolo    string    `json:"simbolo" validate:"required,max=5"`
	IDSucursal uuid.UUID `json:"id_sucursal" validate:"required"`
	IDStatus   uuid.UUID `json:"id_status" validate:"required"`
}
