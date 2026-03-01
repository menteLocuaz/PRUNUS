package dto

type MonedaCreateRequest struct {
	Nombre     string `json:"nombre"`
	IDSucursal uint   `json:"id_sucursal"`
	Estado     int    `json:"estado"`
}

type MonedaUpdateRequest struct {
	Nombre     string `json:"nombre"`
	IDSucursal uint   `json:"id_sucursal"`
	Estado     int    `json:"estado"`
}
