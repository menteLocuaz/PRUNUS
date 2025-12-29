package dto

type EmpresaResponse struct {
	ID     uint   `json:"id"`
	Nombre string `json:"nombre"`
	RUT    string `json:"rut"`
	Estado int    `json:"estado"`
}

type EmpresaCreateRequest struct {
	Nombre string `json:"nombre"`
	RUT    string `json:"rut"`
	Estado int    `json:"estado"`
}

type EmpresaUpdateRequest struct {
	Nombre string `json:"nombre"`
	RUT    string `json:"rut"`
	Estado int    `json:"estado"`
}
