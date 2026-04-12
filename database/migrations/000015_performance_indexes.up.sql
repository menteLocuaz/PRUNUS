-- 1. OPTIMIZACIONES PARA FACTURA (Arqueos y Reportes Financieros)
-- Acelera la suma de totales por control de caja y estado.
-- Uso de Cobertura (Covering Index) para evitar lecturas de tabla (Heap Fetches).
CREATE INDEX IF NOT EXISTS idx_factura_arqueo_performance 
ON factura(id_sucursal, id_status) 
INCLUDE (total, created_at);

-- 2. OPTIMIZACIONES PARA INVENTARIO (Consulta Rápida en Punto de Venta)
-- Acelera la búsqueda de stock disponible por sucursal y producto.
CREATE INDEX IF NOT EXISTS idx_inventario_lookup_performance 
ON inventario(id_sucursal, id_producto) 
INCLUDE (stock_actual);

-- 3. OPTIMIZACIONES PARA MOVIMIENTOS (Kardex / Historial)
-- Acelera consultas de historial de producto ordenadas por fecha descendente.
CREATE INDEX IF NOT EXISTS idx_movimientos_kardex_performance 
ON movimientos_inventario(id_producto, created_at DESC);

-- 4. OPTIMIZACIONES PARA DETALLE FACTURA (Análisis de Ventas)
-- Acelera el cálculo de productos más vendidos (TOP Sales).
CREATE INDEX IF NOT EXISTS idx_detalle_factura_ventas_performance 
ON detalle_factura(id_producto, cantidad);
