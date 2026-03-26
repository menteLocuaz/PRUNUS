package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// migrateFnInventario registra la función almacenada para gestionar movimientos de inventario.
func migrateFnInventario(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("migrateFnInventario: error al iniciar transacción: %w", err)
	}
	defer tx.Rollback()

	statements := []struct {
		name string
		sql  string
	}{
		// ── 1. Asegurar esquema de movimientos ──────────────────────────────
		{
			name: "ALTER TABLE movimientos_inventario ADD COLUMN id_sucursal",
			sql: `
			ALTER TABLE movimientos_inventario 
			ADD COLUMN IF NOT EXISTS id_sucursal UUID;`,
		},
		{
			name: "BACKFILL id_sucursal en movimientos",
			sql: `
			UPDATE movimientos_inventario m
			SET id_sucursal = i.id_sucursal
			FROM inventario i
			WHERE m.id_producto = i.id_producto
			  AND m.id_sucursal IS NULL;`,
		},
		{
			name: "SET NOT NULL id_sucursal en movimientos",
			sql: `
			-- Solo aplicar si se pudo rellenar o si no hay datos.
			-- En un entorno nuevo será obligatorio desde el inicio.
			DO $$ 
			BEGIN 
				IF NOT EXISTS (SELECT 1 FROM movimientos_inventario WHERE id_sucursal IS NULL) THEN
					ALTER TABLE movimientos_inventario ALTER COLUMN id_sucursal SET NOT NULL;
				END IF;
			END $$;`,
		},
		{
			name: "ADD FK id_sucursal en movimientos",
			sql: `
			DO $$ 
			BEGIN 
				IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints 
							  WHERE constraint_name='fk_movimientos_sucursal') THEN
					ALTER TABLE movimientos_inventario 
					ADD CONSTRAINT fk_movimientos_sucursal 
					FOREIGN KEY (id_sucursal) REFERENCES sucursal(id_sucursal);
				END IF;
			END $$;`,
		},

		// ── 2. Función Principal ─────────────────────────────────────────────
		{
			name: "CREATE FUNCTION inventario_ia_movimiento",
			sql: `
			CREATE OR REPLACE FUNCTION inventario_ia_movimiento(
				p_id_sucursal       UUID,
				p_id_usuario        UUID,
				p_tipo_movimiento   VARCHAR(50), -- 'VENTA', 'COMPRA', 'AJUSTE_ENTRADA', 'AJUSTE_SALIDA'
				p_referencia        VARCHAR(255),
				p_items_json        JSONB        -- [{"id_producto": "uuid", "cantidad": 10.50}, ...]
			)
			RETURNS TABLE (
				id_movimiento   UUID,
				id_producto     UUID,
				stock_anterior  NUMERIC(12,2),
				cantidad        NUMERIC(12,2),
				stock_posterior NUMERIC(12,2)
			)
			LANGUAGE plpgsql
			AS $$
			DECLARE
				v_item           RECORD;
				v_stock_actual   NUMERIC(12,2);
				v_es_salida      BOOLEAN;
				v_ip_origen      VARCHAR(45);
				v_error_number   TEXT;
				v_error_msg      TEXT;
			BEGIN
				-- ── 1. Identificar si el tipo de movimiento es una salida ────────
				v_es_salida := p_tipo_movimiento IN ('VENTA', 'AJUSTE_SALIDA', 'SALIDA', 'DEVOLUCION_PROVEEDOR');

				-- ── 2. Iterar items para validación preventiva de stock ─────────
				-- Esto asegura que NINGÚN item se procese si uno solo falla el stock.
				IF v_es_salida THEN
					FOR v_item IN SELECT * FROM jsonb_to_recordset(p_items_json) AS x(id_producto UUID, cantidad NUMERIC)
					LOOP
						SELECT stock_actual INTO v_stock_actual
						FROM inventario
						WHERE id_producto = v_item.id_producto
						  AND id_sucursal = p_id_sucursal
						  AND deleted_at IS NULL
						FOR UPDATE; -- Bloqueo preventivo de fila para evitar condiciones de carrera

						IF v_stock_actual IS NULL OR v_stock_actual < v_item.cantidad THEN
							RAISE EXCEPTION 'Stock insuficiente para el producto %. Actual: %, Requerido: %', 
								v_item.id_producto, COALESCE(v_stock_actual, 0), v_item.cantidad;
						END IF;
					END LOOP;
				END IF;

				-- ── 3. Procesar inserción de movimientos ─────────────────────────
				-- El disparador (trigger) trg_actualizar_stock_movimiento se encargará 
				-- automáticamente de actualizar la tabla 'inventario' y 'producto'.
				FOR v_item IN SELECT * FROM jsonb_to_recordset(p_items_json) AS x(id_producto UUID, cantidad NUMERIC)
				LOOP
					INSERT INTO movimientos_inventario (
						id_producto,
						id_sucursal,
						id_usuario,
						tipo_movimiento,
						cantidad,
						referencia,
						fecha
					) VALUES (
						v_item.id_producto,
						p_id_sucursal,
						p_id_usuario,
						p_tipo_movimiento,
						v_item.cantidad,
						p_referencia,
						CURRENT_TIMESTAMP
					)
					RETURNING 
						movimientos_inventario.id_movimiento, 
						movimientos_inventario.id_producto, 
						movimientos_inventario.stock_anterior, 
						movimientos_inventario.cantidad, 
						movimientos_inventario.stock_posterior
					INTO 
						id_movimiento, 
						id_producto, 
						stock_anterior, 
						cantidad, 
						stock_posterior;
					
					RETURN NEXT;
				END LOOP;

			EXCEPTION
				WHEN OTHERS THEN
					GET STACKED DIAGNOSTICS
						v_error_msg    = MESSAGE_TEXT,
						v_error_number = RETURNED_SQLSTATE;

					v_ip_origen := COALESCE(inet_client_addr()::VARCHAR, 'local/socket');

					RAISE LOG 'inventario_ia_movimiento error: %', json_build_object(
						'fecha',     CURRENT_TIMESTAMP,
						'usuario',   p_id_usuario,
						'sucursal',  p_id_sucursal,
						'tipo',      p_tipo_movimiento,
						'sqlstate',  v_error_number,
						'mensaje',   v_error_msg
					);

					RAISE EXCEPTION 'inventario_ia_movimiento falló [%]: %', v_error_number, v_error_msg;
			END;
			$$;`,
		},
	}

	for _, stmt := range statements {
		if _, err := tx.ExecContext(ctx, stmt.sql); err != nil {
			return fmt.Errorf("migrateFnInventario: fallo en [%s]: %w", stmt.name, err)
		}
	}

	return tx.Commit()
}
