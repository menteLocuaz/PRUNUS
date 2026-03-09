package models

import (
	"time"

	"github.com/google/uuid" // Nueva dependencia
)

type EstacionPos struct {
	IDEstacion uuid.UUID  `json:"id_estacion"`
	Codigo     string     `json:"codigo"`
	Nombre     string     `json:"nombre"`
	IP         string     `json:"ip"`
	IDSucursal uuid.UUID  `json:"id_sucursal"`
	IDStatus   uuid.UUID  `json:"id_status"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}
