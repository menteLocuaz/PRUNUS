package models

import (
	"time"

	"github.com/google/uuid" // Nueva dependencia
)

type Empresa struct {
	IDEmpresa uuid.UUID  `json:"id"`
	Nombre    string     `json:"nombre"`
	RUT       string     `json:"rut"`
	IDStatus  uuid.UUID  `json:"id_status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
