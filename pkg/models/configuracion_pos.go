package models

import "time"

// CanalImpresion representa la tabla Canal_Impresion
type CanalImpresion struct {
	IDCanalImpresion string     `json:"id_canal_impresion"`
	Descripcion      string     `json:"descripcion"`
	CdnID            int        `json:"cdn_id"`
	IDStatus         string     `json:"id_status"` // Mantenemos IDStatus si es legacy, pero priorizamos deleted_at
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
}

// Impresora representa la tabla Impresora
type Impresora struct {
	IDImpresora string     `json:"id_impresora"`
	Nombre      string     `json:"nombre"`
	RstID       int        `json:"rst_id"`
	IDStatus    string     `json:"id_status"`
	CreatedAt   time.Time  `json:"created_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

// Puerto representa la tabla Puertos
type Puerto struct {
	IDPuertos   string     `json:"id_puertos"`
	Descripcion string     `json:"descripcion"`
	IDStatus    string     `json:"id_status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}
