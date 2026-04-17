package models

import (
	"time"

	"github.com/google/uuid"
)

// Producto representa el catálogo maestro de artículos en el sistema.
// Sigue la estructura definida en la base de datos (tabla: producto).
type Producto struct {
	IDProducto  uuid.UUID `json:"id_producto"`
	IDCategoria uuid.UUID `json:"id_categoria"`
	IDStatus    uuid.UUID `json:"id_status"`

	// Campos con discrepancia corregida según el esquema DB
	Nombre       string `json:"pro_nombre"`      // DB: pro_nombre
	Descripcion  string `json:"pro_descripcion"` // DB: pro_descripcion
	CodigoBarras string `json:"pro_codigo"`      // DB: pro_codigo (Código principal)
	SKU          string `json:"sku"`             // DB: sku (Añadido en migración 000021)

	// Identificadores adicionales (Migración 000021)
	// Nota: Si se requiere usar el campo codigo_barras de la DB además de pro_codigo,
	// se debería añadir un campo extra aquí.

	// Campos que no están en la tabla 'producto' base pero se usan en la lógica
	// TODO: Verificar si estos campos deben ser añadidos a la tabla via migración
	FechaVencimiento *time.Time `json:"fecha_vencimiento,omitempty"`
	Imagen           string     `json:"imagen,omitempty"`
	IDMoneda         uuid.UUID  `json:"id_moneda,omitempty"`
	IDUnidad         uuid.UUID  `json:"id_unidad,omitempty"`

	// Metadatos y búsqueda (PostgreSQL Optimization)
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	SearchVector string                 `json:"-"` // tsvector no se expone en JSON

	// Relaciones (Navegación)
	Categoria *Categoria `json:"categoria,omitempty"`
	Moneda    *Moneda    `json:"moneda,omitempty"`
	Unidad    *Unidad    `json:"unidad,omitempty"`

	// Información de Inventario (Poblada vía JOIN o servicios externos)
	Inventario []*Inventario `json:"inventario,omitempty"`

	// Auditoría (sp_core_setup_table)
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
