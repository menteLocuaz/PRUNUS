-- 1. EXTENDER FACTURAS Y COMPRAS PARA GESTIÓN DE CRÉDITO
DO $$ 
BEGIN
    -- Fecha de vencimiento para Cuentas por Cobrar (Facturas)
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='factura' AND column_name='fecha_vencimiento') THEN
        ALTER TABLE factura ADD COLUMN fecha_vencimiento TIMESTAMPTZ;
    END IF;

    -- Fecha de vencimiento para Cuentas por Pagar (Orden de Compra)
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name='orden_compra') THEN
        IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='orden_compra' AND column_name='fecha_vencimiento') THEN
            ALTER TABLE orden_compra ADD COLUMN fecha_vencimiento TIMESTAMPTZ;
        END IF;
    END IF;
END $$;

-- 2. ÍNDICES PARA REPORTES DE ANTIGÜEDAD DE DEUDA (PostgreSQL Optimization)
-- Utilizamos una subconsulta para obtener el ID del estado 'Pendiente' de forma dinámica.
DO $$ 
DECLARE 
    v_id_pendiente UUID;
BEGIN
    SELECT id_status INTO v_id_pendiente FROM estatus WHERE std_descripcion = 'Pendiente' LIMIT 1;

    IF v_id_pendiente IS NOT NULL THEN
        EXECUTE format('CREATE INDEX IF NOT EXISTS idx_factura_vencimiento_pend ON factura(fecha_vencimiento) 
                        WHERE id_status = %L AND deleted_at IS NULL', v_id_pendiente);
    END IF;
END $$;

-- 3. OPTIMIZACIÓN DE REPORTES DE RENTABILIDAD (Pareto / Análisis de Ventas)
CREATE INDEX IF NOT EXISTS idx_movimientos_fecha_tipo_perf 
ON movimientos_inventario(created_at DESC, tipo_movimiento)
WHERE deleted_at IS NULL;
