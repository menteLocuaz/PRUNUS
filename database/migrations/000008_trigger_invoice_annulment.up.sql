-- 1. FUNCIÓN DEL TRIGGER PARA REVERTIR STOCK
CREATE OR REPLACE FUNCTION fn_revertir_stock_anulacion()
RETURNS TRIGGER AS $$
DECLARE
    v_detalle       RECORD;
    v_id_anulado    UUID;
BEGIN
    -- Obtener dinámicamente el ID del estatus 'Anulado' para el módulo de facturación (mdl_id = 8 o similar)
    -- Esto es más robusto que un UUID hardcodeado
    SELECT id_status INTO v_id_anulado 
    FROM estatus 
    WHERE std_descripcion IN ('Anulado', 'Anulada') 
      AND mdl_id IN (5, 8) -- Módulos de Ventas/POS
    LIMIT 1;

    -- Solo actuar si el estado cambia a "Anulado"
    IF NEW.id_status = v_id_anulado AND OLD.id_status != v_id_anulado THEN
        
        -- Recorrer el detalle de la factura para generar movimientos de devolución
        -- Usamos los campos normalizados: id_sucursal e id_usuario ya están en la cabecera (NEW)
        FOR v_detalle IN (
            SELECT id_producto, cantidad, precio_unitario 
            FROM detalle_factura 
            WHERE id_factura = NEW.id_factura 
              AND deleted_at IS NULL
        ) LOOP
            -- Insertar el movimiento de inventario (el trigger de movimientos actualizará el stock)
            INSERT INTO movimientos_inventario (
                id_producto,
                id_sucursal,
                id_usuario,
                tipo_movimiento,
                cantidad,
                id_referencia,
                observacion
            ) VALUES (
                v_detalle.id_producto,
                NEW.id_sucursal,
                NEW.id_usuario,
                'DEVOLUCION',
                v_detalle.cantidad,
                NEW.id_factura,
                'ANULACION FACTURA #' || NEW.fac_numero
            );
        END LOOP;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 2. CREACIÓN DEL TRIGGER
-- Se ejecuta AFTER UPDATE porque necesitamos comparar el estado anterior y el nuevo
DROP TRIGGER IF EXISTS trg_revertir_stock_anulacion ON factura;
CREATE TRIGGER trg_revertir_stock_anulacion
AFTER UPDATE ON factura
FOR EACH ROW
EXECUTE FUNCTION fn_revertir_stock_anulacion();
