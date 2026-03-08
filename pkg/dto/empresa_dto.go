package dto

type EmpresaResponse struct {
	ID     uint   `json:"id"`
	Nombre string `json:"nombre"`
	RUT    string `json:"rut"`
	Estado int    `json:"estado"`
}

type EmpresaCreateRequest struct {
	Nombre string `json:"nombre" validate:"required,min=3,max=150"`
	RUT    string `json:"rut" validate:"required,max=20"`
	Estado int    `json:"estado" validate:"required,oneof=0 1"`
}

type EmpresaUpdateRequest struct {
	Nombre string `json:"nombre" validate:"required,min=3,max=150"`
	RUT    string `json:"rut" validate:"required,max=20"`
	Estado int    `json:"estado" validate:"required,oneof=0 1"`
}
