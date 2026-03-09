package dto

import "time"

// AbrirCajaDTO entrada para iniciar una sesión de caja
type AbrirCajaDTO struct {
	IDCaja        uint    `json:"id_caja" validate:"required"`
	MontoApertura float64 `json:"monto_apertura" validate:"required,min=0"`
}

// MovimientoCajaDTO entrada para registrar un retiro o ingreso
type MovimientoCajaDTO struct {
	Monto  float64 `json:"monto" validate:"required,gt=0"`
	Motivo string  `json:"motivo" validate:"required"`
}

// CierreCajaDTO entrada para finalizar el turno
type CierreCajaDTO struct {
	MontoCierre float64 `json:"monto_cierre" validate:"required,min=0"`
}

// EstadoCajaDTO respuesta con el resumen actual de la caja
type EstadoCajaDTO struct {
	IDSesion      uint      `json:"id_sesion"`
	NombreCaja    string    `json:"nombre_caja"`
	MontoApertura float64   `json:"monto_apertura"`
	TotalIngresos float64   `json:"total_ingresos"`
	TotalEgresos  float64   `json:"total_egresos"`
	SaldoActual   float64   `json:"saldo_actual"`
	FechaApertura time.Time `json:"fecha_apertura"`
}
