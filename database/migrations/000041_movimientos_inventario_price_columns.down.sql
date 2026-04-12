DROP INDEX IF EXISTS idx_movimientos_fecha_sucursal;

ALTER TABLE movimientos_inventario
    DROP COLUMN IF EXISTS precio_unitario,
    DROP COLUMN IF EXISTS costo_unitario,
    DROP COLUMN IF EXISTS fecha,
    DROP COLUMN IF EXISTS referencia;
