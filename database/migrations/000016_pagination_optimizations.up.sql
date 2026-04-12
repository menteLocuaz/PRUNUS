-- 1. ÍNDICES PARA PAGINACIÓN POR CURSOR (Keyset Pagination)
-- Estos índices permiten recorrer grandes volúmenes de datos de forma ultra-rápida
-- sin el costo incremental de OFFSET.

-- Factura: Optimiza el listado histórico de ventas en el POS y Administración
CREATE INDEX IF NOT EXISTS idx_factura_created_at_pagination ON factura(created_at DESC);

-- Producto: Optimiza el listado del catálogo de productos (Grid de ventas)
CREATE INDEX IF NOT EXISTS idx_producto_created_at_pagination ON producto(created_at DESC);

-- Inventario: Optimiza el listado de stock por sucursal
CREATE INDEX IF NOT EXISTS idx_inventario_created_at_pagination ON inventario(created_at DESC);

-- Auditoría: Optimiza el visor de logs del sistema
CREATE INDEX IF NOT EXISTS idx_auditoria_maestra_fecha_pagination ON auditoria_maestra(fecha DESC);
