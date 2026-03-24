package migrations

import "database/sql"

func migrateTriggerStockSync(db *sql.DB) error {
	query := `
	CREATE OR REPLACE FUNCTION fn_actualizar_stock_movimiento()
	RETURNS TRIGGER AS $$
	DECLARE
		v_stock_actual NUMERIC(12,2) := 0;
		v_signo INTEGER := 1;
	BEGIN
		-- Determinar el signo según el tipo de movimiento
		IF NEW.tipo_movimiento IN ('VENTA', 'AJUSTE_SALIDA', 'SALIDA') THEN
			v_signo := -1;
		ELSIF NEW.tipo_movimiento IN ('COMPRA', 'DEVOLUCION', 'ENTRADA', 'AJUSTE_ENTRADA') THEN
			v_signo := 1;
		ELSE
			-- Si es AJUSTE genérico, usamos el signo de la cantidad si es posible o asumimos entrada
			IF NEW.cantidad < 0 THEN
				v_signo := 1; -- Si la cantidad ya es negativa, no invertimos más, pero es raro
			END IF;
		END IF;

		-- Obtener stock actual o inicializar si no existe
		SELECT stock_actual INTO v_stock_actual 
		FROM inventario 
		WHERE id_producto = NEW.id_producto AND id_sucursal = NEW.id_sucursal;

		IF NOT FOUND THEN
			v_stock_actual := 0;
			-- Crear registro de inventario base si no existe
			INSERT INTO inventario (id_producto, id_sucursal, stock_actual)
			VALUES (NEW.id_producto, NEW.id_sucursal, 0);
		END IF;

		-- Asignar valores de auditoría al movimiento
		NEW.stock_anterior := v_stock_actual;
		NEW.stock_posterior := v_stock_actual + (NEW.cantidad * v_signo);

		-- Actualizar la tabla principal de inventario
		UPDATE inventario 
		SET stock_actual = NEW.stock_posterior,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id_producto = NEW.id_producto AND id_sucursal = NEW.id_sucursal;

		-- Sincronizar stock GLOBAL en la tabla producto (Suma de todas las sucursales)
		UPDATE producto 
		SET stock = (SELECT COALESCE(SUM(stock_actual), 0) FROM inventario WHERE id_producto = NEW.id_producto AND deleted_at IS NULL),
		    updated_at = CURRENT_TIMESTAMP
		WHERE id_producto = NEW.id_producto;

		RETURN NEW;
	END;
	$$ language 'plpgsql';

	DROP TRIGGER IF EXISTS trg_actualizar_stock_movimiento ON movimientos_inventario;
	CREATE TRIGGER trg_actualizar_stock_movimiento
	BEFORE INSERT ON movimientos_inventario
	FOR EACH ROW
	EXECUTE FUNCTION fn_actualizar_stock_movimiento();
	`
	_, err := db.Exec(query)
	return err
}
