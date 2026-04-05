package migrations

import "database/sql"

func migrateConcurrencyControl(db *sql.DB) error {
	query := `
	-- 1. Actualizar el trigger de stock con bloqueos explícitos
	CREATE OR REPLACE FUNCTION fn_actualizar_stock_movimiento_v2()
	RETURNS TRIGGER AS $$
	DECLARE
		v_delta_old NUMERIC(12,2) := 0;
		v_delta_new NUMERIC(12,2) := 0;
		v_lock      RECORD;
	BEGIN
		-- A. REVERTIR IMPACTO ANTERIOR (Si es UPDATE o DELETE)
		IF (TG_OP = 'UPDATE' OR TG_OP = 'DELETE') THEN
			IF OLD.deleted_at IS NULL THEN
				-- Bloqueo pestimista de la fila de inventario
				SELECT * INTO v_lock FROM inventario 
				WHERE id_producto = OLD.id_producto AND id_sucursal = OLD.id_sucursal 
				FOR UPDATE;

				v_delta_old := OLD.cantidad * fn_get_movimiento_signo(OLD.tipo_movimiento, OLD.cantidad);
				
				UPDATE inventario 
				SET stock_actual = stock_actual - v_delta_old,
				    updated_at = CURRENT_TIMESTAMP
				WHERE id_producto = OLD.id_producto AND id_sucursal = OLD.id_sucursal;
			END IF;
		END IF;

		-- B. APLICAR IMPACTO NUEVO (Si es INSERT o UPDATE)
		IF (TG_OP = 'INSERT' OR TG_OP = 'UPDATE') THEN
			IF NEW.deleted_at IS NULL THEN
				-- Asegurar que exista el registro de forma atómica
				INSERT INTO inventario (id_producto, id_sucursal, stock_actual, created_at, updated_at)
				VALUES (NEW.id_producto, NEW.id_sucursal, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
				ON CONFLICT (id_producto, id_sucursal) DO NOTHING;

				-- Bloqueo pestimista de la fila de inventario
				SELECT * INTO v_lock FROM inventario 
				WHERE id_producto = NEW.id_producto AND id_sucursal = NEW.id_sucursal 
				FOR UPDATE;

				v_delta_new := NEW.cantidad * fn_get_movimiento_signo(NEW.tipo_movimiento, NEW.cantidad);

				UPDATE inventario 
				SET stock_actual = stock_actual + v_delta_new,
				    updated_at = CURRENT_TIMESTAMP
				WHERE id_producto = NEW.id_producto AND id_sucursal = NEW.id_sucursal;
			END IF;
			RETURN NEW;
		END IF;

		IF (TG_OP = 'DELETE') THEN RETURN OLD; END IF;
		RETURN NULL;
	END;
	$$ LANGUAGE plpgsql;

	-- 2. Actualizar función de facturación con bloqueo preventivo
	CREATE OR REPLACE FUNCTION factura_registrar_completa(
		p_cabecera_json     JSONB,
		p_detalles_json     JSONB,
		p_pagos_json        JSONB,
		p_id_usuario        UUID
	)
	RETURNS TABLE (id_factura UUID, fac_numero VARCHAR(50), total DECIMAL(18,2), status_msg TEXT)
	LANGUAGE plpgsql AS $$
	DECLARE
		v_id_factura        UUID;
		v_fac_numero        VARCHAR(50);
		v_id_status_act     UUID;
		v_id_status_pag     UUID;
		v_id_estacion       UUID;
		v_id_sucursal       UUID;
		v_item              RECORD;
		v_pago              RECORD;
		v_total_detalles    DECIMAL(18,2) := 0;
		v_total_pagos       DECIMAL(18,2) := 0;
		v_ip_origen         VARCHAR(45);
		v_error_number      TEXT;
		v_error_msg         TEXT;
	BEGIN
		-- Resolvemos IDs básicos
		SELECT id_status INTO v_id_status_act FROM estatus WHERE std_descripcion = 'Activo' AND mdl_id = 8 LIMIT 1;
		SELECT id_status INTO v_id_status_pag FROM estatus WHERE std_descripcion = 'Activo' AND mdl_id = 8 LIMIT 1;
		
		v_id_estacion := (p_cabecera_json->>'id_estacion')::UUID;
		SELECT id_sucursal INTO v_id_sucursal FROM estaciones_pos WHERE id_estacion = v_id_estacion;

		-- BLOQUEO PREVENTIVO: Bloquear todos los productos de la venta en orden para evitar Deadlocks
		-- Esto asegura que ninguna otra transacción toque estos stocks hasta que esta termine.
		PERFORM 1 FROM inventario 
		WHERE id_sucursal = v_id_sucursal 
		  AND id_producto IN (SELECT (value->>'id_producto')::UUID FROM jsonb_array_elements(p_detalles_json))
		ORDER BY id_producto
		FOR UPDATE;

		-- Generar correlativo
		v_fac_numero := p_cabecera_json->>'fac_numero';
		IF v_fac_numero IS NULL OR v_fac_numero = '' OR v_fac_numero = 'AUTO' THEN
			v_fac_numero := fn_get_next_secuencial(v_id_estacion, 'FACTURA');
		END IF;

		-- Insertar Cabecera
		INSERT INTO factura (fac_numero, cfac_subtotal, cfac_iva, cfac_total, cfac_observacion, id_user_pos, id_estacion, id_orden_pedido, id_cliente, id_periodo, id_control_estacion, id_status, base_impuesto, impuesto, valor_impuesto, metadata)
		VALUES (v_fac_numero, (p_cabecera_json->>'subtotal')::DECIMAL, (p_cabecera_json->>'iva')::DECIMAL, (p_cabecera_json->>'total')::DECIMAL, p_cabecera_json->>'observacion', p_id_usuario, v_id_estacion, (p_cabecera_json->>'id_orden_pedido')::UUID, (p_cabecera_json->>'id_cliente')::UUID, (p_cabecera_json->>'id_periodo')::UUID, (p_cabecera_json->>'id_control_estacion')::UUID, v_id_status_act, (p_cabecera_json->>'base_impuesto')::DECIMAL, (p_cabecera_json->>'impuesto')::DECIMAL, (p_cabecera_json->>'valor_impuesto')::DECIMAL, p_cabecera_json->'metadata')
		RETURNING factura.id_factura, factura.fac_numero, factura.cfac_total INTO v_id_factura, v_fac_numero, v_total_detalles;

		-- Insertar Detalles (El trigger AFTER insertará movimientos y actualizará stock bajo el bloqueo ya obtenido)
		FOR v_item IN SELECT * FROM jsonb_to_recordset(p_detalles_json) AS x(id_producto UUID, cantidad NUMERIC, precio NUMERIC, subtotal NUMERIC, impuesto NUMERIC, total NUMERIC)
		LOOP
			INSERT INTO detalle_factura (id_factura, id_producto, cantidad, precio, subtotal, impuesto, total)
			VALUES (v_id_factura, v_item.id_producto, v_item.cantidad, v_item.precio, v_item.subtotal, v_item.impuesto, v_item.total);
		END LOOP;

		-- Registrar Pagos
		FOR v_pago IN SELECT * FROM jsonb_to_recordset(p_pagos_json) AS x(id_forma_pago UUID, valor_billete DECIMAL, total_pagar DECIMAL)
		LOOP
			v_total_pagos := v_total_pagos + v_pago.total_pagar;
			INSERT INTO forma_pago_factura (id_factura, id_forma_pago, valor_billete, total_pagar, id_usuario, id_status)
			VALUES (v_id_factura, v_pago.id_forma_pago, v_pago.valor_billete, v_pago.total_pagar, p_id_usuario, v_id_status_pag);
		END LOOP;

		-- Validación final de cuadre
		IF ABS(v_total_pagos - v_total_detalles) > 0.01 THEN
			RAISE EXCEPTION 'El total pagado (%) no coincide con el total de la factura (%)', v_total_pagos, v_total_detalles;
		END IF;

		id_factura := v_id_factura; fac_numero := v_fac_numero; total := v_total_detalles; status_msg := 'Factura procesada correctamente';
		RETURN NEXT;
	EXCEPTION WHEN OTHERS THEN
		GET STACKED DIAGNOSTICS v_error_msg = MESSAGE_TEXT, v_error_number = RETURNED_SQLSTATE;
		RAISE EXCEPTION 'factura_registrar_completa falló [%]: %', v_error_number, v_error_msg;
	END;
	$$;
	`
	_, err := db.Exec(query)
	return err
}
