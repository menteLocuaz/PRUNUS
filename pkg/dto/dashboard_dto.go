package dto

type DashboardResumen struct {
	ValorInventarioTotal float64            `json:"valor_inventario_total"`
	ProductosBajoStock   int                `json:"productos_bajo_stock"`
	VentasMesActual      float64            `json:"ventas_mes_actual"`
	CuentasPorCobrar     float64            `json:"cuentas_por_cobrar"`
	CuentasPorPagar      float64            `json:"cuentas_por_pagar"`
	GastosMensuales      float64            `json:"gastos_mensuales"`
	PuntoEquilibrio      float64            `json:"punto_equilibrio"`
	CicloConversionEfectivo float64         `json:"ciclo_conversion_efectivo"`
	DIO                  float64            `json:"dio"` // Días de Inventario
	DSO                  float64            `json:"dso"` // Días de Cobro
	DPO                  float64            `json:"dpo"` // Días de Pago
	TopProductos         []TopProductoDTO   `json:"top_productos"`
	VentasVsCompras      []VentasComprasDTO `json:"ventas_vs_compras"`
}

type TopProductoDTO struct {
	Nombre     string  `json:"nombre"`
	Cantidad   float64 `json:"cantidad"`
	Rentabilidad float64 `json:"rentabilidad"`
}

type VentasComprasDTO struct {
	Mes     string  `json:"mes"`
	Ventas  float64 `json:"ventas"`
	Compras float64 `json:"compras"`
}

type InventarioCategoriaDTO struct {
	Categoria string  `json:"categoria"`
	Valor     float64 `json:"valor"`
	Porcentaje float64 `json:"porcentaje"`
}

type AntiguedadDeudaDTO struct {
	Rango string  `json:"rango"` // 0-30, 31-60, 61-90, 90+
	Monto float64 `json:"monto"`
}

type MermaItemDTO struct {
	IDProducto    string  `json:"id_producto"`
	ProNombre     string  `json:"pro_nombre"`
	ProCodigo     string  `json:"pro_codigo"`
	CantidadMerma float64 `json:"cantidad_merma"`
	Motivo        string  `json:"motivo"`
	CostoTotal    float64 `json:"costo_total"`
	Fecha         string  `json:"fecha"`
}

type MermasResponseDTO struct {
	TotalMermas float64        `json:"total_mermas"`
	Moneda      string         `json:"moneda"`
	Items       []MermaItemDTO `json:"items"`
}
