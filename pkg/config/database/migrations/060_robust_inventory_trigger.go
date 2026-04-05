package migrations

import "database/sql"

// migrateRobustInventoryTrigger implementa un trigger avanzado para movimientos_inventario
// que maneja INSERT, UPDATE y DELETE, asegurando la integridad del stock incluso en
// cambios manuales o correcciones de movimientos.
func migrateRobustInventoryTrigger(db *sql.DB) error {
	query := `
	-- 1. Función para determinar el signo según tipo de movimiento
	CREATE OR REPLACE FUNCTION fn_get_movimiento_signo(p_tipo TEXT, p_cantidad NUMERIC)
	RETURNS INTEGER AS $$
	BEGIN
		IF p_tipo IN ('VENTA', 'AJUSTE_SALIDA', 'SALIDA') THEN
			RETURN -1;
		ELSIF p_tipo IN ('COMPRA', 'DEVOLUCION', 'ENTRADA', 'AJUSTE_ENTRADA') THEN
			RETURN 1;
		ELSE
			-- Caso para AJUSTE genérico (si cantidad < 0, se asume que ya trae el signo, pero normalizamos)
			RETURN 1;
		END IF;
	END;
	$$ LANGUAGE plpgsql;

	-- 2. Función principal del trigger (AHORA AFTER para manejar DELETE y UPDATE complejos)
	CREATE OR REPLACE FUNCTION fn_actualizar_stock_movimiento_v2()
	RETURNS TRIGGER AS $$
	DECLARE
		v_signo_old INTEGER;
		v_signo_new INTEGER;
		v_delta_old NUMERIC(12,2) := 0;
		v_delta_new NUMERIC(12,2) := 0;
	BEGIN
		-- A. REVERTIR IMPACTO ANTERIOR (Si es UPDATE o DELETE)
		IF (TG_OP = 'UPDATE' OR TG_OP = 'DELETE') THEN
			-- Solo revertimos si el registro anterior estaba "activo" (no borrado lógicamente)
			IF OLD.deleted_at IS NULL THEN
				v_signo_old := fn_get_movimiento_signo(OLD.tipo_movimiento, OLD.cantidad);
				v_delta_old := OLD.cantidad * v_signo_old;
				
				-- Restar del inventario lo que el movimiento anterior sumó/restó
				UPDATE inventario 
				SET stock_actual = stock_actual - v_delta_old,
				    updated_at = CURRENT_TIMESTAMP
				WHERE id_producto = OLD.id_producto AND id_sucursal = OLD.id_sucursal;
			END IF;
		END IF;

		-- B. APLICAR IMPACTO NUEVO (Si es INSERT o UPDATE)
		IF (TG_OP = 'INSERT' OR TG_OP = 'UPDATE') THEN
			-- Solo aplicamos si el nuevo registro no está marcado como eliminado
			IF NEW.deleted_at IS NULL THEN
				v_signo_new := fn_get_movimiento_signo(NEW.tipo_movimiento, NEW.cantidad);
				v_delta_new := NEW.cantidad * v_signo_new;

				-- Asegurar que exista el registro en inventario (Upsert manual)
				INSERT INTO inventario (id_producto, id_sucursal, stock_actual, created_at, updated_at)
				SELECT NEW.id_producto, NEW.id_sucursal, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
				WHERE NOT EXISTS (
					SELECT 1 FROM inventario 
					WHERE id_producto = NEW.id_producto AND id_sucursal = NEW.id_sucursal
				);

				-- Actualizar el stock actual sumando el nuevo delta
				UPDATE inventario 
				SET stock_actual = stock_actual + v_delta_new,
				    updated_at = CURRENT_TIMESTAMP
				WHERE id_producto = NEW.id_producto AND id_sucursal = NEW.id_sucursal;
			END IF;
			RETURN NEW;
		END IF;

		IF (TG_OP = 'DELETE') THEN
			RETURN OLD;
		END IF;
		
		RETURN NULL;
	END;
	$$ LANGUAGE plpgsql;

	-- 3. Recrear el Trigger
	-- Eliminamos el anterior (que era BEFORE INSERT)
	DROP TRIGGER IF EXISTS trg_actualizar_stock_movimiento ON movimientos_inventario;
	
	-- Creamos el nuevo como AFTER para capturar todos los cambios de estado
	CREATE TRIGGER trg_actualizar_stock_movimiento
	AFTER INSERT OR UPDATE OR DELETE ON movimientos_inventario
	FOR EACH ROW
	EXECUTE FUNCTION fn_actualizar_stock_movimiento_v2();
	`
	_, err := db.Exec(query)
	return err
}
