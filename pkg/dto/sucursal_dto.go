package dto

import "github.com/google/uuid"

type SucursalResponse struct {
	ID             uuid.UUID `json:"id_sucursal"`
	IDEmpresa      uuid.UUID `json:"id_empresa"`
	NombreSucursal string    `json:"nombre_sucursal"`
	IDStatus       uuid.UUID `json:"id_status"`
}

type SucursalCreateRequest struct {
	IDEmpresa      uuid.UUID `json:"id_empresa" validate:"required"`
	NombreSucursal string    `json:"nombre_sucursal" validate:"required,min=3,max=150"`
	IDStatus       uuid.UUID `json:"id_status" validate:"required"`
}

type SucursalUpdateRequest struct {
	IDEmpresa      uuid.UUID `json:"id_empresa" validate:"required"`
	NombreSucursal string    `json:"nombre_sucursal" validate:"required,min=3,max=150"`
	IDStatus       uuid.UUID `json:"id_status" validate:"required"`
}
