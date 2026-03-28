package models

import (
	"time"

	"github.com/google/uuid" // Nueva dependencia
)

type LogSistema struct {
	IDLog      uuid.UUID  `json:"id_log"`
	IDUsuario  uuid.UUID  `json:"id_usuario"`
	IDModulo   *uuid.UUID `json:"id_modulo,omitempty"`
	Accion     string     `json:"accion"`
	Tabla      string     `json:"tabla"`
	RegistroID uuid.UUID  `json:"registro_id"`
	Fecha      time.Time  `json:"fecha"`
	IP         string     `json:"ip"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}

type AuditoriaCaja struct {
	IDAuditoria       uuid.UUID  `json:"id_auditoria"`
	IDControlEstacion uuid.UUID  `json:"id_control_estacion"`
	TipoMovimiento    uuid.UUID  `json:"tipo_movimiento"`
	Valor             float64    `json:"valor"`
	Fecha             time.Time  `json:"fecha"`
	IDUsuario         uuid.UUID  `json:"id_usuario"`
	Descripcion       string     `json:"descripcion"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
}

type FacturaAudit struct {
	IDAudit        uuid.UUID  `json:"id_audit"`
	IDFactura      uuid.UUID  `json:"id_factura"`
	IDUsuario      *uuid.UUID `json:"id_usuario,omitempty"`
	Accion         string     `json:"accion"`
	EstadoAnterior *uuid.UUID `json:"estado_anterior,omitempty"`
	EstadoNuevo    *uuid.UUID `json:"estado_nuevo,omitempty"`
	Observaciones  string     `json:"observaciones"`
	Fecha          time.Time  `json:"fecha"`
	IPAddress      string     `json:"ip_address"`
}

type HistorialPrecios struct {
	IDHistorial    uuid.UUID  `json:"id_historial"`
	IDProducto     uuid.UUID  `json:"id_producto"`
	IDSucursal     uuid.UUID  `json:"id_sucursal"`
	PrecioAnterior float64    `json:"precio_anterior"`
	PrecioNuevo    float64    `json:"precio_nuevo"`
	TipoPrecio     string     `json:"tipo_precio"` // VENTA, COMPRA
	IDUsuario      *uuid.UUID `json:"id_usuario,omitempty"`
	Fecha          time.Time  `json:"fecha"`
}
