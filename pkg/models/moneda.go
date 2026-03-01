package models

import "time"

// Moneda representa una moneda asociada a una sucursal.
type Moneda struct {
	IDMoneda   uint   `json:"id_moneda"`
	Nombre     string `json:"nombre"`
	IDSucursal uint   `json:"id_sucursal"`
	Estado     int    `json:"estado"`

	// Relación de navegación (no se persiste directamente en la tabla)
	Sucursal *Sucursal `json:"sucursal,omitempty"`

	// Eliminación lógica (soft delete)
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
