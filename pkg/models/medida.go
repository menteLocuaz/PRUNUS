package models

import (
	"time"

	"github.com/google/uuid"
)

// Unidad representa una unidad operativa asociada a una sucursal.
type Unidad struct {
	IDUnidad    uuid.UUID `json:"id_unidad"`
	Nombre      string    `json:"nombre"`
	Abreviatura string    `json:"abreviatura"` // KG, UND, LTS, etc.
	IDStatus    uuid.UUID `json:"id_status"`
	IDSucursal  uuid.UUID `json:"id_sucursal"`

	// Relación de navegación (no se persiste directamente en la tabla)
	Sucursal *Sucursal `json:"sucursal,omitempty"`
	Status   *Estatus  `json:"status,omitempty"`

	// Eliminación lógica (soft delete)
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
