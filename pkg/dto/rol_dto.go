package dto

// RolResponse estructura de respuesta para el rol
type RolResponse struct {
	IDRol      uint   `json:"id_rol"`
	RolNombre  string `json:"nombre_rol"`
	IDSucursal uint   `json:"id_sucursal"`
	Estado     int    `json:"estado"`
}

// RolCreateRequest estructura de solicitud para crear un rol
type RolCreateRequest struct {
	RolNombre  string `json:"nombre_rol"`
	IDSucursal uint   `json:"id_sucursal"`
	Estado     int    `json:"estado"`
}

// RolUpdateRequest estructura de solicitud para actualizar un rol
type RolUpdateRequest struct {
	RolNombre  string `json:"nombre_rol"`
	IDSucursal uint   `json:"id_sucursal"`
	Estado     int    `json:"estado"`
}
