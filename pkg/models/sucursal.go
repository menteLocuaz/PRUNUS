package models

import (
	"time"

	"github.com/google/uuid"
)

type Sucursal struct {
	IDSucursal     uuid.UUID `json:"id_sucursal"`
	IDEmpresa      uuid.UUID `json:"id_empresa"`
	NombreSucursal string    `json:"nombre_sucursal"`
	IDStatus       uuid.UUID `json:"id_status"`

	// Relación de navegación (no se persiste directamente en la tabla)
	Empresa *Empresa `json:"empresa,omitempty"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
