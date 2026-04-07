package models

import (
	"time"

	"github.com/google/uuid"
)

type GastoOperativo struct {
	IDGasto     uuid.UUID  `json:"id_gasto"`
	IDSucursal  uuid.UUID  `json:"id_sucursal"`
	Descripcion string     `json:"descripcion"`
	Monto       float64    `json:"monto"`
	Frecuencia  string     `json:"frecuencia"` // MENSUAL, ANUAL, UNICO
	FechaGasto  time.Time  `json:"fecha_gasto"`
	IDUsuario   uuid.UUID  `json:"id_usuario"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}
