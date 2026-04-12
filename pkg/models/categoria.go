package models

import (
	"time"

	"github.com/google/uuid"
)

// Categoria representa una categoría de productos o servicios asociada a una sucursal.
type Categoria struct {
	IDCategoria uuid.UUID `json:"id_categoria"`
	Nombre      string    `json:"nombre"`
	IDStatus    uuid.UUID `json:"id_status"`
	IDSucursal  uuid.UUID `json:"id_sucursal"`

	// Relación de navegación (no se persiste directamente en la tabla)
	Sucursal *Sucursal `json:"sucursal,omitempty"`
	Status   *Estatus  `json:"status,omitempty"`

	// Campos para auditoría y eliminación lógica
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
