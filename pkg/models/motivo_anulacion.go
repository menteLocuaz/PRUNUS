package models

import (
	"time"

	"github.com/google/uuid"
)

type CabeceraMotivoAnulacion struct {
	IDCabeceraMotivo uuid.UUID  `json:"id_cabecera_motivo"`
	Descripcion      string     `json:"descripcion"`
	Estado           int        `json:"estado"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
}

type MotivoAnulacion struct {
	IDMotivoAnulacion uuid.UUID  `json:"id_motivo_anulacion"`
	IDCabeceraMotivo  uuid.UUID  `json:"id_cabecera_motivo"`
	Descripcion       string     `json:"descripcion"`
	IDStatus          uuid.UUID  `json:"id_status"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
}
