package models

import (
	"time"

	"github.com/google/uuid" // Nueva dependencia
)

type Estatus struct {
	IDStatus       uuid.UUID  `json:"id_status"`
	StdDescripcion string     `json:"std_descripcion"`
	StpTipoEstado  string     `json:"stp_tipo_estado"`
	MdlID          int        `json:"mdl_id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}
