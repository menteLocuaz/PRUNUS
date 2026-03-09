package models

import (
	"time"

	"github.com/google/uuid"
)

type Producto struct {
	IDProducto       uuid.UUID  `json:"id_producto"`
	Nombre           string     `json:"nombre"`
	Descripcion      string     `json:"descripcion"`
	PrecioCompra     float64    `json:"precio_compra"`
	PrecioVenta      float64    `json:"precio_venta"`
	Stock            uint       `json:"stock"`
	FechaVencimiento *time.Time `json:"fecha_vencimiento,omitempty"`
	Imagen           string     `json:"imagen"`
	IDStatus         uuid.UUID  `json:"id_status"`
	IDSucursal       uuid.UUID  `json:"id_sucursal"`
	IDCategoria      uuid.UUID  `json:"id_categoria"`
	IDMoneda         uuid.UUID  `json:"id_moneda"`
	IDUnidad         uuid.UUID  `json:"id_unidad"`

	// Relaciones
	Sucursal  *Sucursal  `json:"sucursal,omitempty"`
	Categoria *Categoria `json:"categoria,omitempty"`
	Moneda    *Moneda    `json:"moneda,omitempty"`
	Unidad    *Unidad    `json:"unidad,omitempty"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
