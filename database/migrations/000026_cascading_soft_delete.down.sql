-- Revertir: Eliminar borrado lógico en cascada
DROP TRIGGER IF EXISTS trg_factura_soft_delete_cascade ON factura;
DROP FUNCTION IF EXISTS fn_factura_soft_delete_cascade();
DROP TRIGGER IF EXISTS trg_sincronizar_detalle_movimiento ON detalle_factura;
DROP FUNCTION IF EXISTS fn_sincronizar_detalle_movimiento();
