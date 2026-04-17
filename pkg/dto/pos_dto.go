package dto

import (
	"time"
	"github.com/google/uuid"
)

type AbrirCajaDTO struct {
	IDEstacion uuid.UUID `json:"id_estacion"`
	IDUserPos  uuid.UUID `json:"id_user_pos"`
	FondoBase  float64   `json:"fondo_base"`
}

type EstadoCajaDTO struct {
	NombreEstacion    string    `json:"nombre_estacion"`
	IDControlEstacion uuid.UUID `json:"id_control_estacion"`
	IDStatus          uuid.UUID `json:"id_status"`
	StatusDescripcion string    `json:"status_descripcion"`
	FechaInicio       time.Time `json:"fecha_inicio"`
	FondoBase         float64   `json:"fondo_base"`
}
