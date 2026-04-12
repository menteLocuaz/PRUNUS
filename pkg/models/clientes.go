package models

import (
	"time"

	"github.com/google/uuid"
)

// Cliente representa la entidad de un cliente en el catálogo maestro.
type Cliente struct {
	IDCliente      uuid.UUID              `json:"id_cliente"`
	NombreCompleto string                 `json:"nombre_completo"` // DB: nombre_completo
	TipoDocumento  string                 `json:"tipo_documento"`  // DB: tipo_documento (CEDULA, RUC, NIT)
	Documento      string                 `json:"documento"`       // DB: documento
	Email          string                 `json:"email"`
	Telefono       string                 `json:"telefono"`
	Direccion      string                 `json:"direccion"`
	IDStatus       uuid.UUID              `json:"id_status"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`

	// Campos para auditoría y eliminación lógica
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
