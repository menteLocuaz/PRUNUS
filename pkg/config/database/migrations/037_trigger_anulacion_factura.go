package migrations

import "database/sql"

func migrateTriggerAnulacionFactura(db *sql.DB) error {
	query := `
	CREATE OR REPLACE FUNCTION fn_revertir_stock_anulacion()
	RETURNS TRIGGER AS $$
	DECLARE
		v_detalle RECORD;
		v_id_sucursal UUID;
		v_id_anulado UUID := '62ed7d82-0c81-4511-8f02-e7fd140018d8';
	BEGIN
		-- Solo actuar si el estado cambia a "Anulada"
		IF NEW.id_status = v_id_anulado AND OLD.id_status != v_id_anulado THEN
			
			-- Obtener la sucursal de la factura
			SELECT id_sucursal INTO v_id_sucursal 
			FROM estaciones_pos 
			WHERE id_estacion = NEW.id_estacion;

			-- Recorrer el detalle de la factura para revertir stock
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

	DROP TRIGGER IF EXISTS trg_revertir_stock_anulacion ON factura;
	CREATE TRIGGER trg_revertir_stock_anulacion
	AFTER UPDATE ON factura
	FOR EACH ROW
	EXECUTE FUNCTION fn_revertir_stock_anulacion();
	`
	_, err := db.Exec(query)
	return err
}
