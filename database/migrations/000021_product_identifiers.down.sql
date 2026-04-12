DROP INDEX IF EXISTS idx_producto_sku;
DROP INDEX IF EXISTS idx_producto_codigo_barras;
ALTER TABLE producto DROP COLUMN IF EXISTS sku;
ALTER TABLE producto DROP COLUMN IF EXISTS codigo_barras;
