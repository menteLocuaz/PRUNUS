-- Agregar columnas de precio, costo, fecha y referencia a movimientos_inventario
-- Requeridas por el modelo MovimientoInventario y las queries del dashboard.

ALTER TABLE movimientos_inventario
    ADD COLUMN IF NOT EXISTS precio_unitario NUMERIC(18,4) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS costo_unitario  NUMERIC(18,4) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS fecha           TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ADD COLUMN IF NOT EXISTS referencia      VARCHAR(255);

-- Poblar fecha desde created_at en filas históricas
UPDATE movimientos_inventario SET fecha = created_at WHERE fecha = CURRENT_TIMESTAMP AND created_at < CURRENT_TIMESTAMP;

CREATE INDEX IF NOT EXISTS idx_movimientos_fecha_sucursal
    ON movimientos_inventario(id_sucursal, fecha DESC)
    WHERE deleted_at IS NULL;

-- Actualizar la función para que acepte precio_unitario y costo_unitario opcionales del JSON
CREATE OR REPLACE FUNCTION inventario_ia_movimiento(
    p_id_sucursal       UUID,
    p_id_usuario        UUID,
    p_tipo_movimiento   VARCHAR(50),
    p_referencia        VARCHAR(255),
    p_items_json        JSONB
)
RETURNS TABLE (
    id_movimiento   UUID,
    id_producto     UUID,
    stock_anterior  NUMERIC(12,2),
    cantidad        NUMERIC(12,2),
    stock_posterior NUMERIC(12,2)
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_item           RECORD;
    v_stock_actual   NUMERIC(12,2);
    v_es_salida      BOOLEAN;
    v_error_number   TEXT;
    v_error_msg      TEXT;
BEGIN
    v_es_salida := p_tipo_movimiento IN ('VENTA', 'AJUSTE_SALIDA', 'SALIDA', 'DEVOLUCION_PROVEEDOR');

    IF v_es_salida THEN
        FOR v_item IN SELECT * FROM jsonb_to_recordset(p_items_json)
            AS x(id_producto UUID, cantidad NUMERIC, precio_unitario NUMERIC, costo_unitario NUMERIC)
        LOOP
            SELECT stock_actual INTO v_stock_actual
            FROM inventario
            WHERE id_producto = v_item.id_producto
              AND id_sucursal = p_id_sucursal
              AND deleted_at IS NULL
            FOR UPDATE;

            IF v_stock_actual IS NULL OR v_stock_actual < v_item.cantidad THEN
                RAISE EXCEPTION 'Stock insuficiente para el producto %. Disponible: %, Requerido: %',
                    v_item.id_producto, COALESCE(v_stock_actual, 0), v_item.cantidad;
            END IF;
        END LOOP;
    END IF;

    FOR v_item IN SELECT * FROM jsonb_to_recordset(p_items_json)
        AS x(id_producto UUID, cantidad NUMERIC, precio_unitario NUMERIC, costo_unitario NUMERIC)
    LOOP
        INSERT INTO movimientos_inventario (
            id_producto,
            id_sucursal,
            id_usuario,
            tipo_movimiento,
            cantidad,
            precio_unitario,
            costo_unitario,
            referencia,
            observacion,
            fecha,
            created_at
        ) VALUES (
            v_item.id_producto,
            p_id_sucursal,
            p_id_usuario,
            p_tipo_movimiento,
            v_item.cantidad,
            COALESCE(v_item.precio_unitario, 0),
            COALESCE(v_item.costo_unitario, 0),
            p_referencia,
            p_referencia,
            CURRENT_TIMESTAMP,
            CURRENT_TIMESTAMP
        )
        RETURNING
            movimientos_inventario.id_movimiento,
            movimientos_inventario.id_producto,
            movimientos_inventario.stock_anterior,
            movimientos_inventario.cantidad,
            movimientos_inventario.stock_posterior
        INTO
            id_movimiento,
            id_producto,
            stock_anterior,
            cantidad,
            stock_posterior;

        RETURN NEXT;
    END LOOP;

EXCEPTION
    WHEN OTHERS THEN
        GET STACKED DIAGNOSTICS
            v_error_msg    = MESSAGE_TEXT,
            v_error_number = RETURNED_SQLSTATE;

        RAISE LOG 'inventario_ia_movimiento error: %', json_build_object(
            'fecha',     CURRENT_TIMESTAMP,
            'usuario',   p_id_usuario,
            'sucursal',  p_id_sucursal,
            'sqlstate',  v_error_number,
            'mensaje',   v_error_msg
        );

        RAISE EXCEPTION 'inventario_ia_movimiento falló [%]: %', v_error_number, v_error_msg;
END;
$$;
