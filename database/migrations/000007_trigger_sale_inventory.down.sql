-- Revertir: Eliminar trigger de vinculación venta-inventario
DROP TRIGGER IF EXISTS trg_registrar_movimiento_venta ON detalle_factura;
DROP FUNCTION IF EXISTS fn_registrar_movimiento_venta();
