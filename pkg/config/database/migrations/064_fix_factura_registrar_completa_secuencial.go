package migrations

import "database/sql"

func migrateFixFacturacionSecuencial(db *sql.DB) error {
	query := `
	CREATE OR REPLACE FUNCTION factura_registrar_completa(
		p_cabecera_json     JSONB, -- Datos de la tabla 'factura'
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
		-- ── 1. Resolver estados y parámetros base ───────────────────────
		-- Estado Activo para la factura
		SELECT id_status INTO v_id_status_act FROM estatus 
		WHERE std_descripcion = 'Activo' AND mdl_id = 8 LIMIT 1;
		
		-- Estado Pagado para los registros de pago
		SELECT id_status INTO v_id_status_pag FROM estatus 
		WHERE std_descripcion = 'Activo' AND mdl_id = 8 LIMIT 1;

		v_id_estacion := (p_cabecera_json->>'id_estacion')::UUID;
		v_fac_numero := p_cabecera_json->>'fac_numero';

		-- ── 2. Generación automática de secuencial (si es necesario) ────
		IF v_fac_numero IS NULL OR v_fac_numero = '' OR v_fac_numero = 'AUTO' THEN
			v_fac_numero := fn_get_next_secuencial(v_id_estacion, 'FACTURA');
		END IF;

		-- ── 3. Insertar Cabecera de Factura ─────────────────────────────
		INSERT INTO factura (
			fac_numero,
			cfac_subtotal,
			cfac_iva,
			cfac_total,
			cfac_observacion,
			id_user_pos,
			id_estacion,
			id_orden_pedido,
			id_cliente,
			id_periodo,
			id_control_estacion,
			id_status,
			base_impuesto,
			impuesto,
			valor_impuesto,
			metadata
		) VALUES (
			v_fac_numero,
			(p_cabecera_json->>'subtotal')::DECIMAL,
			(p_cabecera_json->>'iva')::DECIMAL,
			(p_cabecera_json->>'total')::DECIMAL,
			p_cabecera_json->>'observacion',
			p_id_usuario,
			v_id_estacion,
			(p_cabecera_json->>'id_orden_pedido')::UUID,
			(p_cabecera_json->>'id_cliente')::UUID,
			(p_cabecera_json->>'id_periodo')::UUID,
			(p_cabecera_json->>'id_control_estacion')::UUID,
			v_id_status_act,
			(p_cabecera_json->>'base_impuesto')::DECIMAL,
			(p_cabecera_json->>'impuesto')::DECIMAL,
			(p_cabecera_json->>'valor_impuesto')::DECIMAL,
			p_cabecera_json->'metadata'
		)
		RETURNING factura.id_factura, factura.fac_numero, factura.cfac_total 
		INTO v_id_factura, v_fac_numero, v_total_detalles;

		-- ── 4. Insertar Detalles ────────────────────────────────────────
		FOR v_item IN SELECT * FROM jsonb_to_recordset(p_detalles_json) 
			AS x(id_producto UUID, cantidad NUMERIC, precio NUMERIC, subtotal NUMERIC, impuesto NUMERIC, total NUMERIC)
		LOOP
			INSERT INTO detalle_factura (
				id_factura,
				id_producto,
				cantidad,
				precio,
				subtotal,
				impuesto,
				total
			) VALUES (
				v_id_factura,
				v_item.id_producto,
				v_item.cantidad,
				v_item.precio,
				v_item.subtotal,
				v_item.impuesto,
				v_item.total
			);
		END LOOP;

		-- ── 5. Registrar Formas de Pago ─────────────────────────────────
		FOR v_pago IN SELECT * FROM jsonb_to_recordset(p_pagos_json) 
			AS x(id_forma_pago UUID, valor_billete DECIMAL, total_pagar DECIMAL)
		LOOP
			v_total_pagos := v_total_pagos + v_pago.total_pagar;
			
			INSERT INTO forma_pago_factura (
				id_factura,
				id_forma_pago,
				valor_billete,
				total_pagar,
				id_usuario,
				id_status
			) VALUES (
				v_id_factura,
				v_pago.id_forma_pago,
				v_pago.valor_billete,
				v_pago.total_pagar,
				p_id_usuario,
				v_id_status_pag
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
	$$;`
	_, err := db.Exec(query)
	return err
}
