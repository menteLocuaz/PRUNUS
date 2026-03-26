package models

import (
	"time"

	"github.com/google/uuid" // Nueva dependencia
)

type Estatus struct {
	IDStatus       uuid.UUID  `json:"id_status"`
	StdDescripcion string     `json:"std_descripcion"`
	StdTipoEstado  string     `json:"std_tipo_estado"`
	Factor         string     `json:"factor,omitempty"`
	Nivel          int        `json:"nivel,omitempty"`
	MdlID          int        `json:"mdl_id"`
	IsActive       bool       `json:"is_active"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}
