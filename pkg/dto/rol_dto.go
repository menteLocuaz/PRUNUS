package dto

import "github.com/google/uuid"

// RolResponse estructura de respuesta para el rol
type RolResponse struct {
	IDRol      uuid.UUID `json:"id_rol"`
	RolNombre  string    `json:"nombre_rol"`
	IDSucursal uuid.UUID `json:"id_sucursal"`
	IDStatus   uuid.UUID `json:"id_status"`
}

// RolCreateRequest estructura de solicitud para crear un rol
type RolCreateRequest struct {
	RolNombre  string    `json:"nombre_rol" validate:"required,min=3,max=100"`
	IDSucursal uuid.UUID `json:"id_sucursal"`
	IDStatus   uuid.UUID `json:"id_status"`
}

// RolUpdateRequest estructura de solicitud para actualizar un rol
type RolUpdateRequest struct {
	RolNombre  string    `json:"nombre_rol" validate:"required,min=3,max=100"`
	IDSucursal uuid.UUID `json:"id_sucursal" validate:"required"`
	IDStatus   uuid.UUID `json:"id_status" validate:"required"`
}
