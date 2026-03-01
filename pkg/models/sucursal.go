package models

import "time"

type Sucursal struct {
	IDSucursal     uint     `json:"id_sucursal"`
	IDEmpresa      uint     `json:"id_empresa"`
	NombreSucursal string   `json:"nombre_sucursal"`
	Estado         int      `json:"estado"`
	Empresa        *Empresa `json:"empresa,omitempty"`
	// Campos para auditoría y eliminación lógica
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
