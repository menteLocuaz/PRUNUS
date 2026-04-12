package dto

import "github.com/google/uuid"

type ProveedorCreateRequest struct {
	RazonSocial    string    `json:"razon_social" validate:"required,min=3,max=255"`
	NitRut         string    `json:"nit_rut" validate:"required,max=50"`
	ContactoNombre string    `json:"contacto_nombre" validate:"omitempty,max=255"`
	Telefono       string    `json:"telefono" validate:"omitempty,max=50"`
	Direccion      string    `json:"direccion" validate:"omitempty"`
	Email          string    `json:"email" validate:"omitempty,email,max=255"`
	IDStatus       uuid.UUID `json:"id_status" validate:"required"`
}

type ProveedorUpdateRequest struct {
	RazonSocial    string    `json:"razon_social" validate:"required,min=3,max=255"`
	NitRut         string    `json:"nit_rut" validate:"required,max=50"`
	ContactoNombre string    `json:"contacto_nombre" validate:"omitempty,max=255"`
	Telefono       string    `json:"telefono" validate:"omitempty,max=50"`
	Direccion      string    `json:"direccion" validate:"omitempty"`
	Email          string    `json:"email" validate:"omitempty,email,max=255"`
	IDStatus       uuid.UUID `json:"id_status" validate:"required"`
}
