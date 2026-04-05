package models

import (
	"time"

	"github.com/google/uuid"
)

type Producto struct {
	IDProducto       uuid.UUID  `json:"id_producto"`
	Nombre           string     `json:"nombre"`
	Descripcion      string     `json:"descripcion"`
	CodigoBarras     string     `json:"codigo_barras"`
	SKU              string     `json:"sku"`
	FechaVencimiento *time.Time `json:"fecha_vencimiento,omitempty"`
	Imagen           string     `json:"imagen"`
	IDStatus         uuid.UUID  `json:"id_status"`
	IDCategoria      uuid.UUID  `json:"id_categoria"`
	IDMoneda         uuid.UUID  `json:"id_moneda"`
	IDUnidad         uuid.UUID  `json:"id_unidad"`

	// Relaciones
	Categoria *Categoria `json:"categoria,omitempty"`
	Moneda    *Moneda    `json:"moneda,omitempty"`
	Unidad    *Unidad    `json:"unidad,omitempty"`

	// Información de Inventario (Opcional, poblada vía JOIN)
	Inventario []*Inventario `json:"inventario,omitempty"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
