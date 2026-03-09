package models

import (
	"time"

	"github.com/google/uuid" // Nueva dependencia
)

type LogSistema struct {
	IDLog      uuid.UUID  `json:"id_log"`
	IDUsuario  uuid.UUID  `json:"id_usuario"`
	Accion     string     `json:"accion"`
	Tabla      string     `json:"tabla"`
	RegistroID uint       `json:"registro_id"`
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
	IDUsuario         uint       `json:"id_usuario"`
	Descripcion       string     `json:"descripcion"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
}
