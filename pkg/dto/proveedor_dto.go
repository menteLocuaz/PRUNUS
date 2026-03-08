package dto

type ProveedorCreateRequest struct {
	Nombre     string `json:"nombre" validate:"required,min=3,max=150"`
	RUC        string `json:"ruc" validate:"required,max=20"`
	Telefono   string `json:"telefono" validate:"omitempty,max=20"`
	Direccion  string `json:"direccion" validate:"omitempty,max=255"`
	Email      string `json:"email" validate:"omitempty,email,max=150"`
	Estado     int    `json:"estado" validate:"required,oneof=0 1"`
	IDSucursal uint   `json:"id_sucursal" validate:"required"`
	IDEmpresa  uint   `json:"id_empresa" validate:"required"`
}

type ProveedorUpdateRequest struct {
	Nombre     string `json:"nombre" validate:"required,min=3,max=150"`
	RUC        string `json:"ruc" validate:"required,max=20"`
	Telefono   string `json:"telefono" validate:"omitempty,max=20"`
	Direccion  string `json:"direccion" validate:"omitempty,max=255"`
	Email      string `json:"email" validate:"omitempty,email,max=150"`
	Estado     int    `json:"estado" validate:"required,oneof=0 1"`
	IDSucursal uint   `json:"id_sucursal" validate:"required"`
	IDEmpresa  uint   `json:"id_empresa" validate:"required"`
}
