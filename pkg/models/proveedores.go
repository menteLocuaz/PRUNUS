package models

import (
	"time"

	"github.com/google/uuid"
)

// Proveedor representa un proveedor asociado a una sucursal y empresa.
type Proveedor struct {
	IDProveedor uuid.UUID `json:"id_proveedor"`
	Nombre      string    `json:"nombre"`
	RUC         string    `json:"ruc"`
	Telefono    string    `json:"telefono"`
	Direccion   string    `json:"direccion"`
	Email       string    `json:"email"`
	IDStatus    uuid.UUID `json:"id_status"`

	// Claves foráneas
	IDSucursal uuid.UUID `json:"id_sucursal"`
	IDEmpresa  uuid.UUID `json:"id_empresa"`

	// Relaciones de navegación (no se persisten directamente en la tabla)
	Sucursal *Sucursal `json:"sucursal,omitempty"`
	Empresa  *Empresa  `json:"empresa,omitempty"`

	// Campos para auditoría y eliminación lógica
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
