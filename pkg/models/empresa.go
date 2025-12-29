package models

import "time"

type Empresa struct {
	IDEmpresa uint   `json:"id"`
	Nombre    string `json:"nombre"`
	RUT       string `json:"rut"`
	Estado    int    `json:"estado"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
