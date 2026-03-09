package models

import (
	"time"

	"github.com/google/uuid" // Nueva dependencia
)

type DispositivoPos struct {
	IDDispositivo uuid.UUID  `json:"id_dispositivo"`
	Nombre        string     `json:"nombre"`
	Tipo          string     `json:"tipo"`
	IP            string     `json:"ip"`
	IDEstacion    uint       `json:"id_estacion"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
}
