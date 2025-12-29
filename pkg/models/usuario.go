package models

import "time"

type Usuario struct {
	IDUsuario  uint `json:"id_usuario"`
	IDSucursal uint `json:"id_sucursal"`

	UsuEmail    string `json:"email"`
	UsuNombre   string `json:"usu_nombre"`
	UsuDni      string `json:"usu_dni"`
	UsuTelefono string `json:"usu_telefono"`
	UsuPassword string `json:"password"`
	Estado      int    `json:"estado"`

	Rol *Rol `json:"rol,omitempty"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
