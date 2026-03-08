package dto

type ClienteCreateRequest struct {
	EmpresaCliente string `json:"empresa_cliente" validate:"required,max=200"`
	Nombre         string `json:"nombre" validate:"required,max=200"`
	RUC            string `json:"ruc" validate:"required,max=20"`
	Direccion      string `json:"direccion" validate:"required,max=255"`
	Telefono       string `json:"telefono" validate:"required,max=20"`
	Email          string `json:"email" validate:"required,email,max=150"`
	Estado         int    `json:"estado" validate:"required,oneof=0 1"`
}

type ClienteUpdateRequest struct {
	EmpresaCliente string `json:"empresa_cliente" validate:"required,max=200"`
	Nombre         string `json:"nombre" validate:"required,max=200"`
	RUC            string `json:"ruc" validate:"required,max=20"`
	Direccion      string `json:"direccion" validate:"required,max=255"`
	Telefono       string `json:"telefono" validate:"required,max=20"`
	Email          string `json:"email" validate:"required,email,max=150"`
	Estado         int    `json:"estado" validate:"required,oneof=0 1"`
}
