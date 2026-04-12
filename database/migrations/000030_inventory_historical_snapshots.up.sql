-- 1. TABLA DE SNAPSHOTS DIARIOS DE INVENTARIO
CREATE TABLE IF NOT EXISTS inventario_historico (
    id_historico   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_sucursal    UUID NOT NULL REFERENCES sucursal(id_sucursal),
    fecha_snapshot DATE NOT NULL,
    valor_total    NUMERIC(18,4) NOT NULL DEFAULT 0,
    cantidad_total NUMERIC(18,4) NOT NULL DEFAULT 0,
    num_productos  INTEGER NOT NULL DEFAULT 0,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uk_inv_hist_sucursal_fecha UNIQUE (id_sucursal, fecha_snapshot)
);

CREATE INDEX IF NOT EXISTS idx_inv_hist_sucursal_fecha ON inventario_historico(id_sucursal, fecha_snapshot DESC);

-- 2. FUNCIÓN DE CAPTURA DE SNAPSHOT
CREATE OR REPLACE FUNCTION fn_snapshot_inventario(p_sucursal_id UUID)
RETURNS VOID AS $$
BEGIN
    INSERT INTO inventario_historico (
        id_sucursal, fecha_snapshot, valor_total, cantidad_total, num_productos
    )
    SELECT
        p_sucursal_id,
        CURRENT_DATE,
        COALESCE(SUM(stock_actual * precio_compra), 0),
        COALESCE(SUM(stock_actual), 0),
        COUNT(*)::INTEGER
    FROM inventario
    WHERE id_sucursal = p_sucursal_id
      AND deleted_at IS NULL
    ON CONFLICT (id_sucursal, fecha_snapshot) DO UPDATE
        SET valor_total    = EXCLUDED.valor_total,
            cantidad_total = EXCLUDED.cantidad_total,
            num_productos  = EXCLUDED.num_productos,
            created_at     = CURRENT_TIMESTAMP;
END;
$$ LANGUAGE plpgsql;

-- 3. ACTUALIZAR VALIDACIÓN DE TIPOS DE MOVIMIENTO
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'movimientos_inventario_tipo_movimiento_check') THEN
        ALTER TABLE movimientos_inventario DROP CONSTRAINT movimientos_inventario_tipo_movimiento_check;
    END IF;

    ALTER TABLE movimientos_inventario
        ADD CONSTRAINT movimientos_inventario_tipo_movimiento_check
        CHECK (tipo_movimiento IN (
            'ENTRADA','SALIDA','AJUSTE','DEVOLUCION','TRASLADO',
            'COMPRA','VENTA','MERMA','CADUCADO', 'AJUSTE_ENTRADA', 'AJUSTE_SALIDA'
        ));
END $$;

-- 4. APLICAR FRAMEWORK CORE
CALL sp_core_setup_table('inventario_historico');
