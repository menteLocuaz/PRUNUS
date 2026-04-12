DROP TABLE IF EXISTS puertos;
DROP TABLE IF EXISTS impresora;
DROP TABLE IF EXISTS canal_impresion;

ALTER TABLE orden_compra
    DROP COLUMN IF EXISTS fecha_vencimiento,
    DROP COLUMN IF EXISTS observaciones,
    DROP COLUMN IF EXISTS impuesto,
    DROP COLUMN IF EXISTS subtotal,
    DROP COLUMN IF EXISTS id_moneda,
    DROP COLUMN IF EXISTS numero_orden;

DROP TABLE IF EXISTS detalle_orden_compra;
DROP TABLE IF EXISTS movimiento_caja;
DROP TABLE IF EXISTS sesion_caja;
DROP TABLE IF EXISTS caja;
DROP TABLE IF EXISTS log_sistema;
