package dto

type SucursalResponse struct {
	ID             uint   `json:"id_sucursal"`
	IDEmpresa      uint   `json:"id_empresa"`
	NombreSucursal string `json:"nombre_sucursal"`
	Estado         int    `json:"estado"`
}

type SucursalCreateRequest struct {
	IDEmpresa      uint   `json:"id_empresa" validate:"required"`
	NombreSucursal string `json:"nombre_sucursal" validate:"required,min=3,max=150"`
	Estado         int    `json:"estado" validate:"required,oneof=0 1"`
}

type SucursalUpdateRequest struct {
	IDEmpresa      uint   `json:"id_empresa" validate:"required"`
	NombreSucursal string `json:"nombre_sucursal" validate:"required,min=3,max=150"`
	Estado         int    `json:"estado" validate:"required,oneof=0 1"`
}
