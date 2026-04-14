package dto

import "github.com/google/uuid"

// CreateEstatusRequest representa la carga para crear un nuevo estado.
type CreateEstatusRequest struct {
	StdDescripcion string `json:"std_descripcion" validate:"required"`
	StdTipoEstado  string `json:"std_tipo_estado" validate:"required"`
	Factor         string `json:"factor"`
	Nivel          int    `json:"nivel"`
	MdlID          int    `json:"mdl_id" validate:"required"`
	IsActive       bool   `json:"is_active"`
}

// UpdateEstatusRequest representa la carga para actualizar un estado existente.
type UpdateEstatusRequest struct {
	StdDescripcion string `json:"std_descripcion" validate:"required"`
	StdTipoEstado  string `json:"std_tipo_estado" validate:"required"`
	Factor         string `json:"factor"`
	Nivel          int    `json:"nivel"`
	MdlID          int    `json:"mdl_id" validate:"required"`
	IsActive       bool   `json:"is_active"`
}

// EstatusResponse representa la respuesta estándar para un estado.
type EstatusResponse struct {
	IDStatus       uuid.UUID `json:"id_status"`
	StdDescripcion string    `json:"std_descripcion"`
	StdTipoEstado  string    `json:"std_tipo_estado"`
	Factor         string    `json:"factor,omitempty"`
	Nivel          int       `json:"nivel,omitempty"`
	MdlID          int       `json:"mdl_id"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      string    `json:"created_at"`
	UpdatedAt      string    `json:"updated_at"`
}

// EstatusModuleGroup agrupa los estados por el módulo al que pertenecen.
type EstatusModuleGroup struct {
	Modulo string            `json:"modulo"`
	Items  []EstatusResponse `json:"items"`
}

// EstatusMasterCatalog es un mapa de estados agrupados por ID de módulo.
type EstatusMasterCatalog map[int]EstatusModuleGroup
