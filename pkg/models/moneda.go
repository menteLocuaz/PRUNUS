package models

import (
	"time"

	"github.com/google/uuid"
)

// Moneda representa una moneda asociada a una sucursal.
type Moneda struct {
	IDMoneda   uuid.UUID `json:"id_moneda"`
	Nombre     string    `json:"nombre"`
	CodigoISO  string    `json:"codigo_iso"` // USD, EUR, COP, etc.
	Simbolo    string    `json:"simbolo"`    // $, €, etc.
	IDSucursal uuid.UUID `json:"id_sucursal"`
	IDStatus   uuid.UUID `json:"id_status"`

	// Relación de navegación (no se persiste directamente en la tabla)
	Sucursal *Sucursal `json:"sucursal,omitempty"`

	// Eliminación lógica (soft delete)
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
