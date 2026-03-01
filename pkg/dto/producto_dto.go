package dto

import "time"

type ProductoCreateRequest struct {
	Nombre           string    `json:"nombre"`
	Descripcion      string    `json:"descripcion"`
	PrecioCompra     float64   `json:"precio_compra"`
	PrecioVenta      float64   `json:"precio_venta"`
	Stock            uint      `json:"stock"`
	FechaVencimiento time.Time `json:"fecha_vencimiento"`
	Imagen           string    `json:"imagen"`
	Estado           int       `json:"estado"`
	IDSucursal       uint      `json:"id_sucursal"`
	IDCategoria      uint      `json:"id_categoria"`
	IDMoneda         uint      `json:"id_moneda"`
	IDUnidad         uint      `json:"id_unidad"`
}

type ProductoUpdateRequest struct {
	Nombre           string    `json:"nombre"`
	Descripcion      string    `json:"descripcion"`
	PrecioCompra     float64   `json:"precio_compra"`
	PrecioVenta      float64   `json:"precio_venta"`
	Stock            uint      `json:"stock"`
	FechaVencimiento time.Time `json:"fecha_vencimiento"`
	Imagen           string    `json:"imagen"`
	Estado           int       `json:"estado"`
	IDSucursal       uint      `json:"id_sucursal"`
	IDCategoria      uint      `json:"id_categoria"`
	IDMoneda         uint      `json:"id_moneda"`
	IDUnidad         uint      `json:"id_unidad"`
}
