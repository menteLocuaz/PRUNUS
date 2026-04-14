package models

import (
	"time"

	"github.com/google/uuid"
)

type DispositivoPos struct {
	IDDispositivo   uuid.UUID              `json:"id_dispositivo"`
	IDEstacion      uuid.UUID              `json:"id_estacion"`
	Nombre          string                 `json:"nombre"`
	TipoDispositivo string                 `json:"tipo_dispositivo"`
	Configuracion   map[string]interface{} `json:"configuracion,omitempty"`
	IDStatus        uuid.UUID              `json:"id_status"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	DeletedAt       *time.Time             `json:"deleted_at,omitempty"`
}
