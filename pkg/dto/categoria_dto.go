package dto

type CategoriaCreateRequest struct {
	Nombre     string `json:"nombre" validate:"required,min=3,max=100"`
	IDSucursal uint   `json:"id_sucursal" validate:"required"`
}

type CategoriaUpdateRequest struct {
	Nombre     string `json:"nombre" validate:"required,min=3,max=100"`
	IDSucursal uint   `json:"id_sucursal" validate:"required"`
}
