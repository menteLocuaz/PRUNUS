-- 1. REFACTOR: FACTURACIÓN ATÓMICA CON SECUENCIALES AUTOMÁTICOS
-- Esta función integra fn_get_next_secuencial para automatizar la numeración.
CREATE OR REPLACE FUNCTION factura_registrar_completa(
    p_cabecera_json     JSONB, -- Datos para la tabla 'factura'
    p_detalles_json     JSONB, -- Array de items para 'detalle_factura'
    p_pagos_json        JSONB, -- Array de pagos para 'forma_pago_factura'
    p_id_usuario        UUID
)
RETURNS TABLE (
    id_factura      UUID,
    fac_numero      VARCHAR(50),
    total           DECIMAL(18,2),
    status_msg      TEXT
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_id_factura        UUID;
    v_fac_numero        VARCHAR(50);
    v_id_status_act     UUID;
    v_id_status_pag     UUID;
    v_id_estacion       UUID;
    v_item              RECORD;
    v_pago              RECORD;
    v_total_detalles    DECIMAL(18,2) := 0;
    v_total_pagos       DECIMAL(18,2) := 0;
    v_ip_origen         VARCHAR(45);
    v_error_number      TEXT;
    v_error_msg         TEXT;
BEGIN
    -- ── 1. Resolver estados y parámetros base ────────────────────────
    SELECT id_status INTO v_id_status_act FROM estatus 
    WHERE std_descripcion = 'Activo' AND mdl_id IN (5, 8) LIMIT 1;
    
    SELECT id_status INTO v_id_status_pag FROM estatus 
    WHERE std_descripcion = 'Activo' AND mdl_id IN (5, 8) LIMIT 1;

    v_id_estacion := (p_cabecera_json->>'id_estacion')::UUID;
    v_fac_numero  := p_cabecera_json->>'fac_numero';

    -- ── 2. Generación automática de secuencial (si es necesario) ─────
    IF v_fac_numero IS NULL OR v_fac_numero = '' OR v_fac_numero = 'AUTO' THEN
        v_fac_numero := fn_get_next_secuencial(v_id_estacion, 'FACTURA');
    END IF;

    -- ── 3. Insertar Cabecera de Factura ─────────────────────────────
    INSERT INTO factura (
        fac_numero,
        id_sucursal,
        id_usuario,
        id_cliente,
        subtotal,
        impuesto,
        total,
        id_status,
        metadata
    ) VALUES (
        v_fac_numero,
        (p_cabecera_json->>'id_sucursal')::UUID,
        p_id_usuario,
        (p_cabecera_json->>'id_cliente')::UUID,
        (p_cabecera_json->>'subtotal')::DECIMAL,
        (p_cabecera_json->>'impuesto')::DECIMAL,
        (p_cabecera_json->>'total')::DECIMAL,
        v_id_status_act,
        p_cabecera_json->'metadata'
    )
    RETURNING factura.id_factura, factura.fac_numero, factura.total 
    INTO v_id_factura, v_fac_numero, v_total_detalles;

    -- ── 4. Insertar Detalles ────────────────────────────────────────
    FOR v_item IN SELECT * FROM jsonb_to_recordset(p_detalles_json) 
        AS x(id_producto UUID, cantidad NUMERIC, precio_unitario NUMERIC, subtotal NUMERIC, impuesto NUMERIC, total NUMERIC)
    LOOP
        INSERT INTO detalle_factura (
            id_factura,
            id_producto,
            cantidad,
            precio_unitario,
            subtotal,
            impuesto,
            total
        ) VALUES (
            v_id_factura,
            v_item.id_producto,
            v_item.cantidad,
            v_item.precio_unitario,
            v_item.subtotal,
            v_item.impuesto,
            v_item.total
        );
    END LOOP;

    -- ── 5. Registrar Formas de Pago ─────────────────────────────────
    FOR v_pago IN SELECT * FROM jsonb_to_recordset(p_pagos_json) 
        AS x(metodo_pago VARCHAR, monto DECIMAL, referencia VARCHAR)
    LOOP
        v_total_pagos := v_total_pagos + v_pago.monto;
        
        INSERT INTO forma_pago_factura (
            id_factura,
            metodo_pago,
            monto,
            referencia
        ) VALUES (
            v_id_factura,
            v_pago.metodo_pago,
            v_pago.monto,
            v_pago.referencia
        );
    END LOOP;

    -- ── 6. Validación final de cuadre ───────────────────────────────
    IF ABS(v_total_pagos - v_total_detalles) > 0.01 THEN
        RAISE EXCEPTION 'El total pagado (%) no coincide con el total de la factura (%)', 
            v_total_pagos, v_total_detalles;
    END IF;

    -- Retorno exitoso
    id_factura := v_id_factura;
    fac_numero := v_fac_numero;
    total      := v_total_detalles;
    status_msg := 'Factura procesada correctamente';
    RETURN NEXT;

EXCEPTION
    WHEN OTHERS THEN
        GET STACKED DIAGNOSTICS
            v_error_msg    = MESSAGE_TEXT,
            v_error_number = RETURNED_SQLSTATE;

        v_ip_origen := COALESCE(inet_client_addr()::VARCHAR, 'local/socket');

        RAISE LOG 'factura_registrar_completa error: %', json_build_object(
            'fecha',      CURRENT_TIMESTAMP,
            'usuario',    p_id_usuario,
            'sqlstate',   v_error_number,
            'mensaje',    v_error_msg
        );

        RAISE EXCEPTION 'factura_registrar_completa falló [%]: %', v_error_number, v_error_msg;
END;
$$;
