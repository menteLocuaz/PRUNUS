package models

import (
	"time"

	"github.com/google/uuid"
)

// Estatus representa un estado maestro dentro del sistema (ej: Activo, Inactivado, Pagado, etc.)
// Estos estados están segmentados por módulos para mantener la integridad referencial.
type Estatus struct {
	IDStatus       uuid.UUID  `json:"id_status"`
	StdDescripcion string     `json:"std_descripcion"`
	StdTipoEstado  string     `json:"std_tipo_estado"`           // Categorización del estado (ej: '1' para normal)
	Factor         string     `json:"factor,omitempty"`          // Factor multiplicador o código legado si aplica
	Nivel          int        `json:"nivel,omitempty"`           // Nivel de jerarquía o peso del estado
	MdlID          int        `json:"mdl_id"`                    // ID del módulo al que pertenece
	MdlDescripcion string     `json:"mdl_descripcion,omitempty"` // Nombre del módulo (poblado vía JOIN)
	IsActive       bool       `json:"is_active"`                 // Flag de activación
	CreatedAt      time.Time  `json:"created_at"`                // Fecha de creación
	UpdatedAt      time.Time  `json:"updated_at"`                // Fecha de última actualización
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`      // Fecha de eliminación lógica (soft delete)
}
