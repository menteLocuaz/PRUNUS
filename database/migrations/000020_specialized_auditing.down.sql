-- Revertir: Eliminar auditorías especializadas
DROP TRIGGER IF EXISTS tr_audit_precios ON inventario;
DROP FUNCTION IF EXISTS fn_audit_precios();
DROP TABLE IF EXISTS historial_precios;
DROP TABLE IF EXISTS factura_audit;
