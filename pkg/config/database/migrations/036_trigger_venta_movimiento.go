package migrations

import "database/sql"

func migrateTriggerVentaMovimiento(db *sql.DB) error {
	query := `
	CREATE OR REPLACE FUNCTION fn_registrar_movimiento_venta()
	RETURNS TRIGGER AS $$
	DECLARE
		v_id_sucursal UUID;
		v_id_usuario  UUID;
		v_referencia  VARCHAR(255);
	BEGIN
		-- Obtener datos de la factura
		SELECT f.id_user_pos, e.id_sucursal, 'Factura #' || f.fac_numero
		INTO v_id_usuario, v_id_sucursal, v_referencia
		FROM factura f
		JOIN estaciones_pos e ON f.id_estacion = e.id_estacion
		WHERE f.id_factura = NEW.id_factura;

		-- Insertar el movimiento de inventario (el trigger de movimientos se encargará del resto)
		INSERT INTO movimientos_inventario (
			id_producto,
			id_sucursal,
			id_usuario,
			tipo_movimiento,
			cantidad,
			precio_unitario,
			referencia
		) VALUES (
			NEW.id_producto,
			v_id_sucursal,
			v_id_usuario,
			'VENTA',
			NEW.cantidad,
			NEW.precio,
			v_referencia
		);

		RETURN NEW;
	END;
	$$ language 'plpgsql';

	DROP TRIGGER IF EXISTS trg_registrar_movimiento_venta ON detalle_factura;
	CREATE TRIGGER trg_registrar_movimiento_venta
	AFTER INSERT ON detalle_factura
	FOR EACH ROW
	EXECUTE FUNCTION fn_registrar_movimiento_venta();
	`
	_, err := db.Exec(query)
	return err
}
