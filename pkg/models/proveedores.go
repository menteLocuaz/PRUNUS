package models

import (
	"time"

	"github.com/google/uuid"
)

// Proveedor representa un proveedor en el catálogo maestro.
type Proveedor struct {
	IDProveedor    uuid.UUID              `json:"id_proveedor"`
	RazonSocial    string                 `json:"razon_social"`    // DB: razon_social
	NitRut         string                 `json:"nit_rut"`         // DB: nit_rut
	ContactoNombre string                 `json:"contacto_nombre"` // DB: contacto_nombre
	Telefono       string                 `json:"telefono"`
	Direccion      string                 `json:"direccion"`
	Email          string                 `json:"email"`
	IDStatus       uuid.UUID              `json:"id_status"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`

	// Relaciones de navegación (Opcional)
	Status *Estatus `json:"status,omitempty"`

	// Campos para auditoría y eliminación lógica
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
