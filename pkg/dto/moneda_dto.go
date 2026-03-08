package dto

type MonedaCreateRequest struct {
	Nombre     string `json:"nombre" validate:"required,min=1,max=100"`
	IDSucursal uint   `json:"id_sucursal" validate:"required"`
	Estado     int    `json:"estado" validate:"required,oneof=0 1"`
}

type MonedaUpdateRequest struct {
	Nombre     string `json:"nombre" validate:"required,min=1,max=100"`
	IDSucursal uint   `json:"id_sucursal" validate:"required"`
	Estado     int    `json:"estado" validate:"required,oneof=0 1"`
}
