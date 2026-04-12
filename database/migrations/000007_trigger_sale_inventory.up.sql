-- 1. ASEGURAR COLUMNAS ADICIONALES EN MOVIMIENTOS (Normalización)
DO $$ 
BEGIN
    -- Añadir id_usuario a movimientos para trazabilidad si no existe
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='movimientos_inventario' AND column_name='id_usuario') THEN
        ALTER TABLE movimientos_inventario ADD COLUMN id_usuario UUID REFERENCES usuario(id_usuario);
    END IF;
END $$;

-- 2. FUNCIÓN DEL TRIGGER
CREATE OR REPLACE FUNCTION fn_registrar_movimiento_venta()
RETURNS TRIGGER AS $$
DECLARE
    v_id_sucursal UUID;
    v_id_usuario  UUID;
    v_fac_numero  VARCHAR(50);
BEGIN
    -- Obtener datos directamente de la cabecera de la factura
    -- En nuestra versión normalizada (000004), la factura ya tiene id_sucursal e id_usuario
    SELECT id_sucursal, id_usuario, fac_numero
    INTO v_id_sucursal, v_id_usuario, v_fac_numero
    FROM factura
    WHERE id_factura = NEW.id_factura;

    -- Insertar el movimiento de inventario
    -- El trigger 'trg_actualizar_stock_movimiento' (creado en 000006) 
    -- se disparará automáticamente al insertar aquí, actualizando el stock real.
    INSERT INTO movimientos_inventario (
        id_producto,
        id_sucursal,
        id_usuario,
        tipo_movimiento,
        cantidad,
        id_referencia,
        observacion
    ) VALUES (
        NEW.id_producto,
        v_id_sucursal,
        v_id_usuario,
        'VENTA',
        NEW.cantidad,
        NEW.id_factura,
        'Factura #' || v_fac_numero
    );

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 3. CREACIÓN DEL TRIGGER
-- Se ejecuta AFTER INSERT para asegurar que el detalle ya existe y la factura es válida
DROP TRIGGER IF EXISTS trg_registrar_movimiento_venta ON detalle_factura;
CREATE TRIGGER trg_registrar_movimiento_venta
AFTER INSERT ON detalle_factura
FOR EACH ROW
EXECUTE FUNCTION fn_registrar_movimiento_venta();
