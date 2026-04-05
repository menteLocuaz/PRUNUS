package migrations

import "database/sql"

func migrateCascadingSoftDelete(db *sql.DB) error {
	query := `
	-- 1. Mejorar el trigger de Detalle de Factura para manejar borrados
	-- Esto asegura que si se borra un detalle, se borre el movimiento de inventario asociado.
	CREATE OR REPLACE FUNCTION fn_sincronizar_detalle_movimiento()
	RETURNS TRIGGER AS $$
	BEGIN
		-- Si es INSERT: Crear el movimiento (Lógica existente en 036)
		IF (TG_OP = 'INSERT') THEN
			INSERT INTO movimientos_inventario (
				id_producto, id_sucursal, id_usuario, tipo_movimiento, cantidad, precio_unitario, referencia
			)
			SELECT 
				NEW.id_producto, e.id_sucursal, f.id_user_pos, 'VENTA', NEW.cantidad, NEW.precio, 'Factura #' || f.fac_numero
			FROM factura f
			JOIN estaciones_pos e ON f.id_estacion = e.id_estacion
			WHERE f.id_factura = NEW.id_factura;
			
			RETURN NEW;
		END IF;

		-- Si es UPDATE y se hace Soft Delete del detalle
		IF (TG_OP = 'UPDATE') THEN
			IF (OLD.deleted_at IS NULL AND NEW.deleted_at IS NOT NULL) THEN
				UPDATE movimientos_inventario 
				SET deleted_at = NEW.deleted_at 
				WHERE id_producto = NEW.id_producto 
				  AND referencia LIKE 'Factura #' || (SELECT fac_numero FROM factura WHERE id_factura = NEW.id_factura);
			END IF;
			RETURN NEW;
		END IF;

		-- Si es DELETE físico
		IF (TG_OP = 'DELETE') THEN
			DELETE FROM movimientos_inventario 
			WHERE id_producto = OLD.id_producto 
			  AND referencia LIKE 'Factura #' || (SELECT fac_numero FROM factura WHERE id_factura = OLD.id_factura);
			RETURN OLD;
		END IF;

		RETURN NULL;
	END;
	$$ LANGUAGE plpgsql;

	DROP TRIGGER IF EXISTS trg_registrar_movimiento_venta ON detalle_factura;
	CREATE TRIGGER trg_sincronizar_detalle_movimiento
	AFTER INSERT OR UPDATE OR DELETE ON detalle_factura
	FOR EACH ROW EXECUTE FUNCTION fn_sincronizar_detalle_movimiento();

	-- 2. Trigger de Cascada para Factura
	-- Si se borra la factura, borramos detalles y pagos.
	CREATE OR REPLACE FUNCTION fn_factura_soft_delete_cascade()
	RETURNS TRIGGER AS $$
	BEGIN
		-- Detectar Soft Delete de Factura
		IF (OLD.deleted_at IS NULL AND NEW.deleted_at IS NOT NULL) THEN
			-- Cascada a detalles (esto activará a su vez el trigger de movimientos e inventario)
			UPDATE detalle_factura SET deleted_at = NEW.deleted_at WHERE id_factura = NEW.id_factura;
			-- Cascada a formas de pago
			UPDATE forma_pago_factura SET deleted_at = NEW.deleted_at WHERE id_factura = NEW.id_factura;
		END IF;
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

	DROP TRIGGER IF EXISTS trg_factura_soft_delete_cascade ON factura;
	CREATE TRIGGER trg_factura_soft_delete_cascade
	AFTER UPDATE ON factura
	FOR EACH ROW EXECUTE FUNCTION fn_factura_soft_delete_cascade();
	`
	_, err := db.Exec(query)
	return err
}
