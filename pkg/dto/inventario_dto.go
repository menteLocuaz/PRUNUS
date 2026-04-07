package dto

import (
	"time"

	"github.com/google/uuid"
)

type InventarioCreateRequest struct {
	IDProducto   uuid.UUID `json:"id_producto" validate:"required"`
	IDSucursal   uuid.UUID `json:"id_sucursal" validate:"required"`
	StockActual  float64   `json:"stock_actual" validate:"required,gte=0"`
	StockMinimo  float64   `json:"stock_minimo" validate:"required,gte=0"`
	StockMaximo  float64   `json:"stock_maximo" validate:"required,gte=0"`
	PrecioCompra float64   `json:"precio_compra" validate:"required,gte=0"`
	PrecioVenta  float64   `json:"precio_venta" validate:"required,gte=0"`
}

type InventarioUpdateRequest struct {
	StockActual  float64 `json:"stock_actual" validate:"required,gte=0"`
	StockMinimo  float64 `json:"stock_minimo" validate:"required,gte=0"`
	StockMaximo  float64 `json:"stock_maximo" validate:"required,gte=0"`
	PrecioCompra float64 `json:"precio_compra" validate:"required,gte=0"`
	PrecioVenta  float64 `json:"precio_venta" validate:"required,gte=0"`
}

type MovimientoCreateRequest struct {
	IDProducto     uuid.UUID `json:"id_producto" validate:"required"`
	IDSucursal     uuid.UUID `json:"id_sucursal" validate:"required"`
	TipoMovimiento string    `json:"tipo_movimiento" validate:"required,oneof=ENTRADA SALIDA AJUSTE DEVOLUCION TRASLADO"`
	Cantidad       float64   `json:"cantidad" validate:"required,gt=0"`
	Referencia     string    `json:"referencia"`
}

type MovimientoItemRequest struct {
	IDProducto uuid.UUID `json:"id_producto" validate:"required"`
	Cantidad   float64   `json:"cantidad" validate:"required,gt=0"`
}

type MovimientoMasivoRequest struct {
	IDSucursal     uuid.UUID               `json:"id_sucursal" validate:"required"`
	TipoMovimiento string                  `json:"tipo_movimiento" validate:"required,oneof=ENTRADA SALIDA AJUSTE DEVOLUCION TRASLADO"`
	Referencia     string                  `json:"referencia"`
	Items          []MovimientoItemRequest `json:"items" validate:"required,min=1,dive"`
}

// RotacionProductoResponse representa el índice de rotación de un producto en un periodo.
type RotacionProductoResponse struct {
	IDProducto         uuid.UUID `json:"id_producto"`
	COGS               float64   `json:"cogs"`                // Costo de ventas monetario en el periodo
	UnidadesVendidas   float64   `json:"unidades_vendidas"`   // Total de unidades de salida/venta
	InventarioPromedio float64   `json:"inventario_promedio"` // (stock_inicio + stock_fin) / 2 en unidades
	IndiceRotacion     float64   `json:"indice_rotacion"`     // unidades_vendidas / inventario_promedio
}

// ComposicionCategoriaResponse agrupa el valor y cantidad de stock por categoría.
type ComposicionCategoriaResponse struct {
	IDCategoria     uuid.UUID `json:"id_categoria"`
	NombreCategoria string    `json:"nombre_categoria"`
	NumProductos    int       `json:"num_productos"`
	CantidadTotal   float64   `json:"cantidad_total"`             // Suma de stock_actual
	ValorTotal      float64   `json:"valor_total"`                // Suma de stock_actual * precio_compra
	PorcentajeValor float64   `json:"porcentaje_valor,omitempty"` // % del valor total del inventario
}

// AlertaStockResponse enriquece la alerta de stock bajo con datos del producto.
type AlertaStockResponse struct {
	IDInventario   uuid.UUID `json:"id_inventario"`
	IDProducto     uuid.UUID `json:"id_producto"`
	NombreProducto string    `json:"nombre_producto"`
	SKU            string    `json:"sku,omitempty"`
	IDSucursal     uuid.UUID `json:"id_sucursal"`
	StockActual    float64   `json:"stock_actual"`
	StockMinimo    float64   `json:"stock_minimo"`
	Deficit        float64   `json:"deficit"`       // stock_minimo - stock_actual
	PrecioCompra   float64   `json:"precio_compra"`
}

// RotacionFiltroParams parámetros de filtro para el cálculo de rotación.
type RotacionFiltroParams struct {
	FechaInicio time.Time
	FechaFin    time.Time
}
