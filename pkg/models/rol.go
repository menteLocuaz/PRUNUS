package models

import (
	"time"

	"github.com/google/uuid"
)

type Rol struct {
	IDRol      uuid.UUID `json:"id_rol"`
	RolNombre  string    `json:"nombre_rol"`
	IDSucursal uuid.UUID `json:"id_sucursal"`
	IDStatus   uuid.UUID `json:"id_status"`

	Sucursal *Sucursal `json:"sucursal,omitempty"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
