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
	IDSucursal  uint   `json:"id_sucursal"`
	IDRol       uint   `json:"id_rol"`
	UsuEmail    string `json:"email"`
	UsuNombre   string `json:"usu_nombre"`
	UsuDni      string `json:"usu_dni"`
	UsuTelefono string `json:"usu_telefono"`
	UsuPassword string `json:"password"`
	Estado      int    `json:"estado"`
}

// UsuarioUpdateRequest estructura de solicitud para actualizar un usuario
type UsuarioUpdateRequest struct {
	IDSucursal  uint   `json:"id_sucursal"`
	IDRol       uint   `json:"id_rol"`
	UsuEmail    string `json:"email"`
	UsuNombre   string `json:"usu_nombre"`
	UsuDni      string `json:"usu_dni"`
	UsuTelefono string `json:"usu_telefono"`
	UsuPassword string `json:"password"`
	Estado      int    `json:"estado"`
}
