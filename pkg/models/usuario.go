package models

import (
	"time"

	"github.com/google/uuid"
)

type Usuario struct {
	IDUsuario   uuid.UUID  `json:"id_usuario"`
	IDSucursal  uuid.UUID  `json:"id_sucursal"`
	IDRol       uuid.UUID  `json:"id_rol"`
	Email       string     `json:"email"`
	UsuNombre   string     `json:"usu_nombre"`
	UsuDNI      string     `json:"usu_dni"`
	UsuTelefono string     `json:"usu_telefono"`
	Password    string     `json:"-"` // Ocultar password de respuestas JSON
	IDStatus    uuid.UUID  `json:"id_status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`

	// Relaciones de navegación
	Rol      *Rol      `json:"rol,omitempty"`
	Sucursal *Sucursal `json:"sucursal,omitempty"`
}
