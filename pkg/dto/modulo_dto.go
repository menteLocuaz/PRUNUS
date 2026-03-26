package dto

import (
	"github.com/google/uuid"
)

type ModuloDTO struct {
	IDModulo       uuid.UUID   `json:"id_modulo"`
	MdlID          int         `json:"mdl_id"`
	MdlDescripcion string      `json:"mdl_descripcion"`
	Abreviatura    string      `json:"abreviatura,omitempty"`
	IDPadre        *uuid.UUID  `json:"id_padre,omitempty"`
	Nivel          int         `json:"nivel"`
	Orden          int         `json:"orden"`
	Ruta           string      `json:"ruta,omitempty"`
	Icono          string      `json:"icono,omitempty"`
	IsActive       bool        `json:"is_active"`
	IDStatus       *uuid.UUID  `json:"id_status,omitempty"`
	SubModulos     []ModuloDTO `json:"sub_modulos,omitempty"`
}

type CreateModuloDTO struct {
	MdlDescripcion string     `json:"mdl_descripcion" validate:"required"`
	Abreviatura    string     `json:"abreviatura"`
	IDPadre        *uuid.UUID `json:"id_padre"`
	Nivel          int        `json:"nivel"`
	Orden          int        `json:"orden"`
	Ruta           string     `json:"ruta"`
	Icono          string     `json:"icono"`
	IsActive       bool       `json:"is_active"`
	IDStatus       *uuid.UUID `json:"id_status"`
}

type PermisoRolDTO struct {
	IDPermiso uuid.UUID  `json:"id_permiso"`
	IDRol     uuid.UUID  `json:"id_rol"`
	IDModulo  uuid.UUID  `json:"id_modulo"`
	Modulo    *ModuloDTO `json:"modulo,omitempty"`
	CanRead   bool       `json:"can_read"`
	CanWrite  bool       `json:"can_write"`
	CanUpdate bool       `json:"can_update"`
	CanDelete bool       `json:"can_delete"`
	CanImport bool       `json:"can_import"`
	CanExport bool       `json:"can_export"`
}

type UpdatePermisoDTO struct {
	CanRead   bool `json:"can_read"`
	CanWrite  bool `json:"can_write"`
	CanUpdate bool `json:"can_update"`
	CanDelete bool `json:"can_delete"`
	CanImport bool `json:"can_import"`
	CanExport bool `json:"can_export"`
}
