package dto

type SucursalResponse struct {
	ID             uint   `json:"id_sucursal"`
	IDEmpresa      uint   `json:"id_empresa"`
	NombreSucursal string `json:"nombre_sucursal"`
	Estado         int    `json:"estado"`
}

type SucursalCreateRequest struct {
	IDEmpresa      uint   `json:"id_empresa"`
	NombreSucursal string `json:"nombre_sucursal"`
	Estado         int    `json:"estado"`
}

type SucursalUpdateRequest struct {
	IDEmpresa      uint   `json:"id_empresa"`
	NombreSucursal string `json:"nombre_sucursal"`
	Estado         int    `json:"estado"`
}
