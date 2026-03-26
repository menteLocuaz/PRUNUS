package models

import (
	"time"

	"github.com/google/uuid"
)

type Modulo struct {
	IDModulo       uuid.UUID  `json:"id_modulo"`
	MdlID          int        `json:"mdl_id"`
	MdlDescripcion string     `json:"mdl_descripcion"`
	Abreviatura    string     `json:"abreviatura,omitempty"`
	IDPadre        *uuid.UUID `json:"id_padre,omitempty"`
	Nivel          int        `json:"nivel"`
	Orden          int        `json:"orden"`
	Ruta           string     `json:"ruta,omitempty"`
	Icono          string     `json:"icono,omitempty"`
	IsActive       bool       `json:"is_active"`
	IDStatus       *uuid.UUID `json:"id_status,omitempty"`

	SubModulos []Modulo `json:"sub_modulos,omitempty"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type PermisoRol struct {
	IDPermiso uuid.UUID `json:"id_permiso"`
	IDRol     uuid.UUID `json:"id_rol"`
	IDModulo  uuid.UUID `json:"id_modulo"`
	CanRead   bool      `json:"can_read"`
	CanWrite  bool      `json:"can_write"`
	CanUpdate bool      `json:"can_update"`
	CanDelete bool      `json:"can_delete"`
	CanImport bool      `json:"can_import"`
	CanExport bool      `json:"can_export"`

	Modulo *Modulo `json:"modulo,omitempty"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
