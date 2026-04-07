package dto

type DashboardResumen struct {
	ValorInventarioTotal float64            `json:"valor_inventario_total"`
	ProductosBajoStock   int                `json:"productos_bajo_stock"`
	VentasMesActual      float64            `json:"ventas_mes_actual"`
	CuentasPorCobrar     float64            `json:"cuentas_por_cobrar"`
	CuentasPorPagar      float64            `json:"cuentas_por_pagar"`
	GastosMensuales      float64            `json:"gastos_mensuales"`
	PuntoEquilibrio      float64            `json:"punto_equilibrio"`
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
