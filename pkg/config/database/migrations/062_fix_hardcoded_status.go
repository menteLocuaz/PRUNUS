package migrations

import "database/sql"

func migrateFixHardcodedStatus(db *sql.DB) error {
	query := `
	CREATE OR REPLACE FUNCTION fn_revertir_stock_anulacion()
	RETURNS TRIGGER AS $$
	DECLARE
		v_detalle RECORD;
		v_id_sucursal UUID;
		v_id_anulado UUID;
	BEGIN
		-- Obtener dinámicamente el ID del estatus 'Anulada' para el módulo de Ventas (5)
		SELECT id_status INTO v_id_anulado 
		FROM estatus 
		WHERE std_descripcion = 'Anulada' AND mdl_id = 5 
		LIMIT 1;

		-- Si no se encuentra, lanzar una excepción para evitar inconsistencias
		IF v_id_anulado IS NULL THEN
			RAISE EXCEPTION 'No se encontró el estatus [Anulada] para el módulo de Ventas (5). Verifique la tabla estatus.';
		END IF;

		-- Solo actuar si el estado cambia a "Anulada"
		IF NEW.id_status = v_id_anulado AND OLD.id_status != v_id_anulado THEN
			
			-- Obtener la sucursal de la factura
			SELECT id_sucursal INTO v_id_sucursal 
			FROM estaciones_pos 
			WHERE id_estacion = NEW.id_estacion;

			-- Recorrer el detalle de la factura para revertir stock
			-- Insertar en movimientos_inventario (el trigger robusto se encargará del stock)
			FOR v_detalle IN (SELECT id_producto, cantidad, precio FROM detalle_factura WHERE id_factura = NEW.id_factura) LOOP
				INSERT INTO movimientos_inventario (
					id_producto,
					id_sucursal,
					id_usuario,
					tipo_movimiento,
					cantidad,
					precio_unitario,
					referencia
				) VALUES (
					v_detalle.id_producto,
					v_id_sucursal,
					NEW.id_user_pos,
					'DEVOLUCION',
					v_detalle.cantidad,
					v_detalle.precio,
					'ANULACION FACTURA #' || NEW.fac_numero
				);
			END LOOP;
		END IF;

		RETURN NEW;
	END;
	$$ language 'plpgsql';
	`
	_, err := db.Exec(query)
	return err
}
