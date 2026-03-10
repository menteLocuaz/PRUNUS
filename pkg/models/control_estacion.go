package models

import (
	"time"

	"github.com/google/uuid"
)

type ControlEstacion struct {
	IDControlEstacion  uuid.UUID  `json:"id_control_estacion"`
	IDEstacion         uuid.UUID  `json:"id_estacion"`
	FechaInicio        time.Time  `json:"fecha_inicio"`
	FechaSalida        *time.Time `json:"fecha_salida,omitempty"`
	FondoBase          float64    `json:"fondo_base"`
	UsuarioAsignado    uuid.UUID  `json:"usuario_asignado"`
	FechaFondoAceptado *time.Time `json:"fecha_fondo_aceptado,omitempty"`
	UsuarioRetiroFondo *uuid.UUID `json:"usuario_retiro_fondo,omitempty"`
	FondoRetirado      *float64   `json:"fondo_retirado,omitempty"`
	IDStatus           uuid.UUID  `json:"id_status"`
	IDUserPos          uuid.UUID  `json:"id_user_pos"`
	IDPeriodo          uuid.UUID  `json:"id_periodo"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty"`
}
