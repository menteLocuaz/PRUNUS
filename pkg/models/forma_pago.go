package models

import (
	"time"

	"github.com/google/uuid"
)

// FormaPago representa un método de pago aceptado en el sistema (Efectivo, Tarjeta, etc.)
type FormaPago struct {
	IDFormaPago uuid.UUID `json:"id_forma_pago"`
	Nombre      string    `json:"nombre"`       // DB: nombre
	RequiereRef bool      `json:"requiere_ref"` // DB: requiere_ref (Indica si necesita número de referencia)
	IDStatus    uuid.UUID `json:"id_status"`

	// Auditoría y Eliminación lógica
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
