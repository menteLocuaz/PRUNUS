package models

import (
	"time"

	"github.com/google/uuid"
)

type Usuario struct {
	IDUsuario     uuid.UUID  `json:"id_usuario"`
	IDSucursal    uuid.UUID  `json:"id_sucursal"`
	IDRol         uuid.UUID  `json:"id_rol"`
	Username      string     `json:"username"`
	Email         string     `json:"email"`
	UsuNombre     string     `json:"usu_nombre"`
	UsuDNI        string     `json:"usu_dni"`
	UsuTelefono   string     `json:"usu_telefono"`
	UsuTarjetaNFC string     `json:"usu_tarjeta_nfc,omitempty"`
	UsuPinPOS     string     `json:"-"` // PIN cifrado para acceso rápido
	NombreTicket  string     `json:"nombre_ticket,omitempty"`
	Password      string     `json:"-"` // Ocultar password de respuestas JSON
	IDStatus      uuid.UUID  `json:"id_status"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`

	// Relaciones de navegación
	Rol        *Rol        `json:"rol,omitempty"`
	Sucursal   *Sucursal   `json:"sucursal,omitempty"`
	Sucursales []uuid.UUID `json:"sucursales_acceso,omitempty"` // IDs de sucursales habilitadas
	Permisos      []string   `json:"permisos"`            // Slugs o rutas de módulos permitidos
	EnTurno       bool       `json:"en_turno"`            // Indica si el usuario puede operar caja
}
