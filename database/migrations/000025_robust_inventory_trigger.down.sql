DROP TRIGGER IF EXISTS trg_actualizar_stock_movimiento ON movimientos_inventario;
DROP FUNCTION IF EXISTS fn_actualizar_stock_movimiento_v2();
DROP FUNCTION IF EXISTS fn_get_movimiento_signo(TEXT, NUMERIC);
