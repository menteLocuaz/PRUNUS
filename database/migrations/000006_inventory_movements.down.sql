-- Revertir: Eliminar trigger y tabla de movimientos
DROP TRIGGER IF EXISTS trg_actualizar_stock_movimiento ON movimientos_inventario;
DROP FUNCTION IF EXISTS fn_actualizar_stock_movimiento();
DROP TABLE IF EXISTS movimientos_inventario;

-- Nota: No eliminamos la columna 'stock' de producto por precaución de datos,
-- pero el trigger que la mantenía ya no existirá.
