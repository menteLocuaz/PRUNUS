package dto

type ClienteCreateRequest struct {
	EmpresaCliente string `json:"empresa_cliente"`
	Nombre         string `json:"nombre"`
	RUC            string `json:"ruc"`
	Direccion      string `json:"direccion"`
	Telefono       string `json:"telefono"`
	Email          string `json:"email"`
	Estado         int    `json:"estado"`
}

type ClienteUpdateRequest struct {
	EmpresaCliente string `json:"empresa_cliente"`
	Nombre         string `json:"nombre"`
	RUC            string `json:"ruc"`
	Direccion      string `json:"direccion"`
	Telefono       string `json:"telefono"`
	Email          string `json:"email"`
	Estado         int    `json:"estado"`
}
