package models

import "time"

type Rol struct {
	IDRol      uint   `json:"id_rol"`
	RolNombre  string `json:"nombre_rol"`
	IDSucursal uint   `json:"id_sucursal"`
	Estado     int    `json:"estado"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
