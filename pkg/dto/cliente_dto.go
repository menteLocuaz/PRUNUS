package dto

import "github.com/google/uuid"

// ClienteCreateRequest define los datos necesarios para crear un nuevo cliente.
type ClienteCreateRequest struct {
	NombreCompleto string                 `json:"nombre_completo" validate:"required,max=255"`
	TipoDocumento  string                 `json:"tipo_documento" validate:"max=20"`
	Documento      string                 `json:"documento" validate:"max=50"`
	Email          string                 `json:"email" validate:"omitempty,email,max=255"`
	Telefono       string                 `json:"telefono" validate:"max=50"`
	Direccion      string                 `json:"direccion"`
	IDStatus       uuid.UUID              `json:"id_status" validate:"required"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// ClienteUpdateRequest define los datos que pueden ser actualizados en un cliente.
type ClienteUpdateRequest struct {
	NombreCompleto string                 `json:"nombre_completo" validate:"required,max=255"`
	TipoDocumento  string                 `json:"tipo_documento" validate:"max=20"`
	Documento      string                 `json:"documento" validate:"max=50"`
	Email          string                 `json:"email" validate:"omitempty,email,max=255"`
	Telefono       string                 `json:"telefono" validate:"max=50"`
	Direccion      string                 `json:"direccion"`
	IDStatus       uuid.UUID              `json:"id_status" validate:"required"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}
