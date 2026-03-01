package models

import "time"

// Producto representa un producto asociado a una sucursal.
type Producto struct {
	IDProducto       uint      `json:"id_producto"`
	Nombre           string    `json:"nombre"`
	Descripcion      string    `json:"descripcion"`
	PrecioCompra     float64   `json:"precio_compra"`
	PrecioVenta      float64   `json:"precio_venta"`
	Stock            uint      `json:"stock"`
	FechaVencimiento time.Time `json:"fecha_vencimiento"`
	Imagen           string    `json:"imagen"`
	Estado           int       `json:"estado"`

	// Claves foráneas
	IDSucursal  uint `json:"id_sucursal"`
	IDCategoria uint `json:"id_categoria"`
	IDMoneda    uint `json:"id_moneda"`
	IDUnidad    uint `json:"id_unidad"`

	// Relaciones de navegación (no se persisten directamente en la tabla)
	Sucursal  *Sucursal  `json:"sucursal,omitempty"`
	Categoria *Categoria `json:"categoria,omitempty"`
	Moneda    *Moneda    `json:"moneda,omitempty"`
	Unidad    *Unidad    `json:"unidad,omitempty"`

	// Campos para auditoría y eliminación lógica
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
