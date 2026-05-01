package models

import (
	"time"

	"github.com/google/uuid" // Nueva dependencia
)

// Caja representa un punto de venta o caja física
type Caja struct {
	IDCaja     uuid.UUID `json:"id_caja"`
	Nombre     string    `json:"nombre"`
	IDSucursal uuid.UUID `json:"id_sucursal"`
	Estado     int       `json:"estado"` // 1: Activa, 0: Inactiva

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// SesionCaja representa el turno de un cajero
type SesionCaja struct {
	IDSesion      uuid.UUID  `json:"id_sesion"`
	IDCaja        uuid.UUID  `json:"id_caja"`
	IDUsuario     uuid.UUID  `json:"id_usuario"`
	MontoApertura float64    `json:"monto_apertura"` // Base de la caja
	MontoCierre   *float64   `json:"monto_cierre,omitempty"`
	FechaApertura time.Time  `json:"fecha_apertura"`
	FechaCierre   *time.Time `json:"fecha_cierre,omitempty"`
	Estado        string     `json:"estado"` // ABIERTA, CERRADA

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// MovimientoCaja representa un ingreso o egreso (retiro) de efectivo
type MovimientoCaja struct {
	IDMovimiento uuid.UUID `json:"id_movimiento"`
	IDSesion     uuid.UUID `json:"id_sesion"`
	Tipo         string    `json:"tipo"` // INGRESO, EGRESO
	Monto        float64   `json:"monto"`
	Motivo       string    `json:"motivo"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
