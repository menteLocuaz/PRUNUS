package dto

import "github.com/google/uuid"

// UsuarioResponse estructura de respuesta para el usuario
type UsuarioResponse struct {
	IDUsuario   uuid.UUID `json:"id_usuario"`
	IDSucursal  uuid.UUID `json:"id_sucursal"`
	IDRol       uuid.UUID `json:"id_rol"`
	UsuEmail    string    `json:"email"`
	UsuNombre   string    `json:"usu_nombre"`
	UsuDni      string    `json:"usu_dni"`
	UsuTelefono string    `json:"usu_telefono"`
	IDStatus    uuid.UUID `json:"id_status"`
}

// UsuarioCreateRequest estructura de solicitud para crear un usuario
type UsuarioCreateRequest struct {
	IDSucursal  uuid.UUID `json:"id_sucursal" validate:"required"`
	IDRol       uuid.UUID `json:"id_rol" validate:"required"`
	UsuEmail    string    `json:"email" validate:"required,email"`
	UsuNombre   string    `json:"usu_nombre" validate:"required,min=3,max=100"`
	UsuDni      string    `json:"usu_dni" validate:"required,min=8,max=15"`
	UsuTelefono string    `json:"usu_telefono" validate:"omitempty,max=20"`
	UsuPassword string    `json:"password" validate:"required,min=6"`
	IDStatus    uuid.UUID `json:"id_status" validate:"required"`
}

// UsuarioUpdateRequest estructura de solicitud para actualizar un usuario
type UsuarioUpdateRequest struct {
	IDSucursal  uuid.UUID `json:"id_sucursal" validate:"required"`
	IDRol       uuid.UUID `json:"id_rol" validate:"required"`
	UsuEmail    string    `json:"email" validate:"required,email"`
	UsuNombre   string    `json:"usu_nombre" validate:"required,min=3,max=100"`
	UsuDni      string    `json:"usu_dni" validate:"required,min=8,max=15"`
	UsuTelefono string    `json:"usu_telefono" validate:"omitempty,max=20"`
	UsuPassword string    `json:"password" validate:"omitempty,min=6"`
	IDStatus    uuid.UUID `json:"id_status" validate:"required"`
}
