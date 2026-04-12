-- Revertir: Eliminar trigger de reversión de stock por anulación
DROP TRIGGER IF EXISTS trg_revertir_stock_anulacion ON factura;
DROP FUNCTION IF EXISTS fn_revertir_stock_anulacion();
