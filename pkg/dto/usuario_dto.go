package dto

// UsuarioResponse estructura de respuesta para el usuario
type UsuarioResponse struct {
	IDUsuario   uint   `json:"id_usuario"`
	IDSucursal  uint   `json:"id_sucursal"`
	IDRol       uint   `json:"id_rol"`
	UsuEmail    string `json:"email"`
	UsuNombre   string `json:"usu_nombre"`
	UsuDni      string `json:"usu_dni"`
	UsuTelefono string `json:"usu_telefono"`
	Estado      int    `json:"estado"`
}

// UsuarioCreateRequest estructura de solicitud para crear un usuario
type UsuarioCreateRequest struct {
	IDSucursal  uint   `json:"id_sucursal" validate:"required"`
	IDRol       uint   `json:"id_rol" validate:"required"`
	UsuEmail    string `json:"email" validate:"required,email"`
	UsuNombre   string `json:"usu_nombre" validate:"required,min=3,max=100"`
	UsuDni      string `json:"usu_dni" validate:"required,min=8,max=15"`
	UsuTelefono string `json:"usu_telefono" validate:"omitempty,max=20"`
	UsuPassword string `json:"password" validate:"required,min=6"`
	Estado      int    `json:"estado" validate:"required,oneof=0 1"`
}

// UsuarioUpdateRequest estructura de solicitud para actualizar un usuario
type UsuarioUpdateRequest struct {
	IDSucursal  uint   `json:"id_sucursal" validate:"required"`
	IDRol       uint   `json:"id_rol" validate:"required"`
	UsuEmail    string `json:"email" validate:"required,email"`
	UsuNombre   string `json:"usu_nombre" validate:"required,min=3,max=100"`
	UsuDni      string `json:"usu_dni" validate:"required,min=8,max=15"`
	UsuTelefono string `json:"usu_telefono" validate:"omitempty,max=20"`
	UsuPassword string `json:"password" validate:"omitempty,min=6"`
	Estado      int    `json:"estado" validate:"required,oneof=0 1"`
}
