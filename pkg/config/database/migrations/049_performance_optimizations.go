package migrations

import "database/sql"

func migratePerformanceOptimizations(db *sql.DB) error {
	query := `
	-- 1. Optimizaciones para Factura (Arqueos y Reportes)
	CREATE INDEX IF NOT EXISTS idx_factura_periodo ON factura(id_periodo);
	CREATE INDEX IF NOT EXISTS idx_factura_control_estacion ON factura(id_control_estacion);
	
	-- Índice Compuesto + Cobertura para Arqueo de Caja
	-- Acelera: SELECT SUM(cfac_total) FROM factura WHERE id_control_estacion = ? AND id_status = ?
	CREATE INDEX IF NOT EXISTS idx_factura_arqueo_performance 
	ON factura(id_control_estacion, id_status) 
	INCLUDE (cfac_total, fecha_operacion);

	-- 2. Optimizaciones para Inventario (Consulta rápida en Caja)
	-- Acelera: SELECT precio_venta, stock_actual FROM inventario WHERE id_sucursal = ? AND id_producto = ?
	CREATE INDEX IF NOT EXISTS idx_inventario_lookup_performance 
	ON inventario(id_sucursal, id_producto) 
	INCLUDE (stock_actual, precio_venta);

	-- 3. Optimizaciones para Movimientos (Kardex y Trazabilidad)
	-- Acelera: consultas de historial de producto ordenadas por fecha
	CREATE INDEX IF NOT EXISTS idx_movimientos_kardex_performance 
	ON movimientos_inventario(id_producto, fecha DESC);

	-- 4. Optimizaciones para Detalle Factura (Análisis de Ventas)
	-- Acelera: TOP productos vendidos
	CREATE INDEX IF NOT EXISTS idx_detalle_factura_ventas_performance 
	ON detalle_factura(id_producto, cantidad);
	`
	_, err := db.Exec(query)
	return err
}
