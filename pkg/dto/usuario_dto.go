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
	UsuTarjetaNFC string  `json:"usu_tarjeta_nfc,omitempty"`
	NombreTicket  string  `json:"nombre_ticket,omitempty"`
	IDStatus    uuid.UUID `json:"id_status"`
	SucursalesAcceso []uuid.UUID `json:"sucursales_acceso,omitempty"`
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
	UsuTarjetaNFC string  `json:"usu_tarjeta_nfc" validate:"omitempty,max=100"`
	UsuPinPOS     string  `json:"usu_pin_pos" validate:"omitempty,max=10"`
	NombreTicket  string  `json:"nombre_ticket" validate:"omitempty,max=50"`
	IDStatus    uuid.UUID `json:"id_status" validate:"required"`
	SucursalesAcceso []uuid.UUID `json:"sucursales_acceso" validate:"omitempty"`
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
	UsuTarjetaNFC string  `json:"usu_tarjeta_nfc" validate:"omitempty,max=100"`
	UsuPinPOS     string  `json:"usu_pin_pos" validate:"omitempty,max=10"`
	NombreTicket  string  `json:"nombre_ticket" validate:"omitempty,max=50"`
	IDStatus    uuid.UUID `json:"id_status" validate:"required"`
	SucursalesAcceso []uuid.UUID `json:"sucursales_acceso" validate:"omitempty"`
}
