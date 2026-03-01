package dto

type CategoriaCreateRequest struct {
	Nombre     string `json:"nombre"`
	IDSucursal uint   `json:"id_sucursal"`
}

type CategoriaUpdateRequest struct {
	Nombre     string `json:"nombre"`
	IDSucursal uint   `json:"id_sucursal"`
}
