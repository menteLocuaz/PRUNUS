package models

import (
	"time"

	"github.com/google/uuid"
)

type Periodo struct {
	IDPeriodo          uuid.UUID  `json:"id_periodo"`
	PrdFechaApertura   time.Time  `json:"prd_fecha_apertura"`
	PrdFechaCierre     *time.Time `json:"prd_fecha_cierre,omitempty"`
	PrdUsuarioApertura uint       `json:"prd_usuario_apertura"`
	PrdUsuarioCierre   *uint      `json:"prd_usuario_cierre,omitempty"`
	IDStatus           uuid.UUID  `json:"id_status"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty"`
}
