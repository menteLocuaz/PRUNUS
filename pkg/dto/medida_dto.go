package dto

type MedidaCreateRequest struct {
	Nombre     string `json:"nombre" validate:"required,min=1,max=100"`
	IDSucursal uint   `json:"id_sucursal" validate:"required"`
}

type MedidaUpdateRequest struct {
	Nombre     string `json:"nombre" validate:"required,min=1,max=100"`
	IDSucursal uint   `json:"id_sucursal" validate:"required"`
}
