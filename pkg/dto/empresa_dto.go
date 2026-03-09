package dto

import "github.com/google/uuid"

type EmpresaResponse struct {
	ID       uuid.UUID `json:"id"`
	Nombre   string    `json:"nombre"`
	RUT      string    `json:"rut"`
	IDStatus uuid.UUID `json:"id_status"`
}

type EmpresaCreateRequest struct {
	Nombre   string    `json:"nombre" validate:"required,min=3,max=150"`
	RUT      string    `json:"rut" validate:"required,max=20"`
	IDStatus uuid.UUID `json:"id_status" validate:"required"`
}

type EmpresaUpdateRequest struct {
	Nombre   string    `json:"nombre" validate:"required,min=3,max=150"`
	RUT      string    `json:"rut" validate:"required,max=20"`
	IDStatus uuid.UUID `json:"id_status" validate:"required"`
}
