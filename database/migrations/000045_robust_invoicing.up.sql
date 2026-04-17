-- 000045_robust_invoicing.up.sql
-- Implementación robusta de facturación atómica con control de stock, lotes e idempotencia

CREATE OR REPLACE FUNCTION factura_registrar_completa(
    p_cabecera_json     JSONB,
    p_detalles_json     JSONB,
    p_pagos_json        JSONB,
    p_id_usuario        UUID,
    p_id_operacion      UUID DEFAULT NULL -- ID para idempotencia
)
RETURNS TABLE (
    id_factura UUID, 
    fac_numero VARCHAR(50), 
    total DECIMAL(18,2), 
    status_msg TEXT
)
LANGUAGE plpgsql AS $$
DECLARE
    v_id_factura        UUID;
    v_fac_numero        VARCHAR(50);
    v_id_status_act     UUID;
    v_id_estacion       UUID;
    v_id_sucursal       UUID;
    v_item              RECORD;
    v_pago              RECORD;
    v_total_detalles    DECIMAL(18,2) := 0;
    v_total_pagos       DECIMAL(18,2) := 0;
    v_stock_actual      NUMERIC(12,2);
    v_error_msg         TEXT;
BEGIN
    -- 1. CONTROL DE IDEMPOTENCIA (Si se provee p_id_operacion)
    -- Asumimos que guardamos el id_operacion en el campo 'metadata' de la factura para este ejemplo
    IF p_id_operacion IS NOT NULL THEN
        SELECT f.id_factura, f.fac_numero, f.total INTO v_id_factura, v_fac_numero, v_total_detalles
        FROM factura f 
        WHERE (f.metadata->>'id_operacion')::UUID = p_id_operacion
        LIMIT 1;

        IF FOUND THEN
            id_factura := v_id_factura; fac_numero := v_fac_numero; total := v_total_detalles; status_msg := 'Operación recuperada (Idempotencia)';
            RETURN NEXT;
            RETURN;
        END IF;
    END IF;

    -- 2. RESOLVER PARÁMETROS
    SELECT id_status INTO v_id_status_act FROM estatus WHERE std_descripcion = 'Activo' AND mdl_id IN (5, 8) LIMIT 1;
    v_id_estacion := (p_cabecera_json->>'id_estacion')::UUID;
    v_id_sucursal := (p_cabecera_json->>'id_sucursal')::UUID;

    -- 3. BLOQUEO TRANSACCIONAL Y VALIDACIÓN DE STOCK
    FOR v_item IN SELECT * FROM jsonb_to_recordset(p_detalles_json) 
        AS x(id_producto UUID, cantidad NUMERIC, id_lote UUID)
    LOOP
        -- Bloqueamos el registro de inventario
        SELECT stock_actual INTO v_stock_actual 
        FROM inventario 
        WHERE id_producto = v_item.id_producto AND id_sucursal = v_id_sucursal
        FOR UPDATE;

        IF v_stock_actual < v_item.cantidad THEN
            RAISE EXCEPTION 'Stock insuficiente para producto %: solicitado %, disponible %', v_item.id_producto, v_item.cantidad, v_stock_actual;
        END IF;

        -- Si usa lotes, validar el stock del lote también
        IF v_item.id_lote IS NOT NULL THEN
            IF NOT EXISTS (SELECT 1 FROM lote WHERE id_lote = v_item.id_lote AND cantidad_actual >= v_item.cantidad) THEN
                RAISE EXCEPTION 'Stock insuficiente o vencido en lote % para producto %', v_item.id_lote, v_item.id_producto;
            END IF;
        END IF;
    END LOOP;

    -- 4. GENERAR CORRELATIVO
    v_fac_numero := p_cabecera_json->>'fac_numero';
    IF v_fac_numero IS NULL OR v_fac_numero = '' OR v_fac_numero = 'AUTO' THEN
        v_fac_numero := fn_get_next_secuencial(v_id_estacion, 'FACTURA');
    END IF;

    -- 5. INSERTAR CABECERA (Incluyendo metadatos de idempotencia)
    INSERT INTO factura (
        fac_numero, id_sucursal, id_usuario, id_cliente, subtotal, impuesto, total, id_status, 
        id_estacion, id_periodo, id_control_estacion, metadata
    ) VALUES (
        v_fac_numero, v_id_sucursal, p_id_usuario, (p_cabecera_json->>'id_cliente')::UUID, 
        (p_cabecera_json->>'subtotal')::DECIMAL, (p_cabecera_json->>'impuesto')::DECIMAL, 
        (p_cabecera_json->>'total')::DECIMAL, v_id_status_act,
        v_id_estacion, (p_cabecera_json->>'id_periodo')::UUID, (p_cabecera_json->>'id_control_estacion')::UUID,
        jsonb_set(COALESCE(p_cabecera_json->'metadata', '{}'::JSONB), '{id_operacion}', to_jsonb(p_id_operacion))
    )
    RETURNING factura.id_factura INTO v_id_factura;

    -- 6. INSERTAR DETALLES Y REGISTRAR MOVIMIENTOS
    FOR v_item IN SELECT * FROM jsonb_to_recordset(p_detalles_json) 
        AS x(id_producto UUID, cantidad NUMERIC, precio_unitario NUMERIC, subtotal NUMERIC, impuesto NUMERIC, total NUMERIC, id_lote UUID)
    LOOP
        v_total_detalles := v_total_detalles + v_item.total;

        INSERT INTO detalle_factura (id_factura, id_producto, cantidad, precio_unitario, subtotal, impuesto, total, id_lote)
        VALUES (v_id_factura, v_item.id_producto, v_item.cantidad, v_item.precio_unitario, v_item.subtotal, v_item.impuesto, v_item.total, v_item.id_lote);

        -- Insertar movimiento de inventario (esto activará el trigger de stock_actual)
        INSERT INTO movimientos_inventario (id_producto, id_sucursal, cantidad, tipo_movimiento, referencia, id_usuario, id_lote)
        VALUES (v_item.id_producto, v_id_sucursal, v_item.cantidad, 'VENTA', 'FACTURA: ' || v_fac_numero, p_id_usuario, v_item.id_lote);

        -- Si hay lote, restar del lote directamente (o vía trigger en movimientos_inventario)
        IF v_item.id_lote IS NOT NULL THEN
            UPDATE lote SET cantidad_actual = cantidad_actual - v_item.cantidad WHERE id_lote = v_item.id_lote;
        END IF;
    END LOOP;

    -- 7. REGISTRAR PAGOS
    FOR v_pago IN SELECT * FROM jsonb_to_recordset(p_pagos_json) 
        AS x(metodo_pago VARCHAR, monto DECIMAL, referencia VARCHAR)
    LOOP
        v_total_pagos := v_total_pagos + v_pago.monto;
        INSERT INTO forma_pago_factura (id_factura, metodo_pago, monto, referencia)
        VALUES (v_id_factura, v_pago.metodo_pago, v_pago.monto, v_pago.referencia);
    END LOOP;

    -- 8. VALIDACIÓN FINAL DE CUADRE (Tolerancia de redondeo 0.01)
    IF ABS(v_total_pagos - (p_cabecera_json->>'total')::DECIMAL) > 0.01 THEN
        RAISE EXCEPTION 'El total pagado (%) no coincide con el total esperado (%)', v_total_pagos, p_cabecera_json->>'total';
    END IF;

    id_factura := v_id_factura; fac_numero := v_fac_numero; total := (p_cabecera_json->>'total')::DECIMAL; status_msg := 'Factura registrada con éxito';
    RETURN NEXT;

EXCEPTION WHEN OTHERS THEN
    GET STACKED DIAGNOSTICS v_error_msg = MESSAGE_TEXT;
    RAISE EXCEPTION 'Error en factura_registrar_completa: %', v_error_msg;
END;
$$;
