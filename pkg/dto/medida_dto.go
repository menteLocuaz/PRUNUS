package dto

type MedidaCreateRequest struct {
	Nombre     string `json:"nombre"`
	IDSucursal uint   `json:"id_sucursal"`
}

type MedidaUpdateRequest struct {
	Nombre     string `json:"nombre"`
	IDSucursal uint   `json:"id_sucursal"`
}
