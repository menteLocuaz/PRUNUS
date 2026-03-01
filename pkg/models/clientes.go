package models

import "time"

// Cliente representa la entidad de un cliente en el sistema.
type Cliente struct {
	IDCliente      uint   `json:"id_cliente"`
	EmpresaCliente string `json:"empresa_cliente"`
	Nombre         string `json:"nombre"`
	RUC            string `json:"ruc"`
	Direccion      string `json:"direccion"`
	Telefono       string `json:"telefono"`
	Email          string `json:"email"`
	Estado         int    `json:"estado"`

	// Campos para auditoría y eliminación lógica
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
