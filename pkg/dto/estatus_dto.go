package dto

import "github.com/google/uuid"

type CreateEstatusRequest struct {
	StdDescripcion string `json:"std_descripcion" validate:"required"`
	StdTipoEstado  string `json:"std_tipo_estado" validate:"required"`
	MdlID          int    `json:"mdl_id" validate:"required"`
}

type UpdateEstatusRequest struct {
	StdDescripcion string `json:"std_descripcion" validate:"required"`
	StdTipoEstado  string `json:"std_tipo_estado" validate:"required"`
	MdlID          int    `json:"mdl_id" validate:"required"`
}

type EstatusResponse struct {
	IDStatus       uuid.UUID `json:"id_status"`
	StdDescripcion string    `json:"std_descripcion"`
	StdTipoEstado  string    `json:"std_tipo_estado"`
	MdlID          int       `json:"mdl_id"`
	CreatedAt      string    `json:"created_at"`
	UpdatedAt      string    `json:"updated_at"`
}

type EstatusModuleGroup struct {
	Modulo string            `json:"modulo"`
	Items  []EstatusResponse `json:"items"`
}

type EstatusMasterCatalog map[int]EstatusModuleGroup
