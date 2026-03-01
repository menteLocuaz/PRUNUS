package dto

type ProveedorCreateRequest struct {
	Nombre    string `json:"nombre"`
	RUC       string `json:"ruc"`
	Telefono  string `json:"telefono"`
	Direccion string `json:"direccion"`
	Email     string `json:"email"`
	Estado    int    `json:"estado"`
	IDSucursal uint  `json:"id_sucursal"`
	IDEmpresa  uint  `json:"id_empresa"`
}

type ProveedorUpdateRequest struct {
	Nombre    string `json:"nombre"`
	RUC       string `json:"ruc"`
	Telefono  string `json:"telefono"`
	Direccion string `json:"direccion"`
	Email     string `json:"email"`
	Estado    int    `json:"estado"`
	IDSucursal uint  `json:"id_sucursal"`
	IDEmpresa  uint  `json:"id_empresa"`
}
