package dto

import (
	"time"

	"github.com/google/uuid"
)

// AbrirCajaDTO entrada para iniciar una sesión de caja (Control Estación)
type AbrirCajaDTO struct {
	IDEstacion uuid.UUID `json:"id_estacion" validate:"required"`
	FondoBase  float64   `json:"fondo_base" validate:"required,min=0"`
	IDUserPos  uuid.UUID `json:"id_user_pos" validate:"required"`
}

// MovimientoCajaDTO entrada para registrar un retiro o ingreso
type MovimientoCajaDTO struct {
	Monto  float64 `json:"monto" validate:"required,gt=0"`
	Motivo string  `json:"motivo" validate:"required"`
}

// CierreCajaDTO entrada para finalizar el turno
type CierreCajaDTO struct {
	FondoRetirado float64 `json:"fondo_retirado" validate:"required,min=0"`
}

// EstadoCajaDTO respuesta con el resumen actual de la caja
type EstadoCajaDTO struct {
	IDControlEstacion uuid.UUID `json:"id_control_estacion"`
	NombreEstacion    string    `json:"nombre_estacion"`
	FondoBase         float64   `json:"fondo_base"`
	IDStatus          uuid.UUID `json:"id_status"`
	StatusDescripcion string    `json:"status_descripcion"`
	FechaInicio       time.Time `json:"fecha_inicio"`
}
