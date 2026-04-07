package migrations

import "database/sql"

func migrateInventarioHistorico(db *sql.DB) error {
	query := `
	-- Tabla de snapshots diarios del valor del inventario por sucursal
	CREATE TABLE IF NOT EXISTS inventario_historico (
		id_historico   UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
		id_sucursal    UUID        NOT NULL REFERENCES sucursal(id_sucursal),
		fecha_snapshot DATE        NOT NULL,
		valor_total    NUMERIC(18,4) NOT NULL DEFAULT 0,
		cantidad_total NUMERIC(18,4) NOT NULL DEFAULT 0,
		num_productos  INTEGER     NOT NULL DEFAULT 0,
		created_at     TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
		UNIQUE (id_sucursal, fecha_snapshot)
	);

	CREATE INDEX IF NOT EXISTS idx_inv_hist_sucursal_fecha
		ON inventario_historico(id_sucursal, fecha_snapshot DESC);

	-- Función para capturar el snapshot del valor actual del inventario de una sucursal.
	-- Usa ON CONFLICT para actualizar si ya existe el snapshot del día.
	CREATE OR REPLACE FUNCTION fn_snapshot_inventario(p_sucursal_id UUID)
	RETURNS VOID AS $$
	BEGIN
		INSERT INTO inventario_historico (
			id_sucursal, fecha_snapshot, valor_total, cantidad_total, num_productos
		)
		SELECT
			p_sucursal_id,
			CURRENT_DATE,
			COALESCE(SUM(stock_actual * precio_compra), 0),
			COALESCE(SUM(stock_actual), 0),
			COUNT(*)::INTEGER
		FROM inventario
		WHERE id_sucursal = p_sucursal_id
		  AND deleted_at IS NULL
		ON CONFLICT (id_sucursal, fecha_snapshot) DO UPDATE
			SET valor_total    = EXCLUDED.valor_total,
			    cantidad_total = EXCLUDED.cantidad_total,
			    num_productos  = EXCLUDED.num_productos,
			    created_at     = CURRENT_TIMESTAMP;
	END;
	$$ LANGUAGE plpgsql;

	-- Tipos de movimiento para pérdidas por merma y caducidad.
	-- El check constraint en movimientos_inventario debe permitirlos.
	-- Actualizamos el constraint si existe; si no existe, la tabla lo admite libremente.
	DO $$
	BEGIN
		IF EXISTS (
			SELECT 1 FROM pg_constraint
			WHERE conname = 'movimientos_inventario_tipo_movimiento_check'
		) THEN
			ALTER TABLE movimientos_inventario
				DROP CONSTRAINT movimientos_inventario_tipo_movimiento_check;

			ALTER TABLE movimientos_inventario
				ADD CONSTRAINT movimientos_inventario_tipo_movimiento_check
				CHECK (tipo_movimiento IN (
					'ENTRADA','SALIDA','AJUSTE','DEVOLUCION','TRASLADO',
					'COMPRA','VENTA','MERMA','CADUCADO'
				));
		END IF;
	END;
	$$;
	`
	_, err := db.Exec(query)
	return err
}
