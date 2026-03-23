package migrations

import "database/sql"

func migrateUpdateMovimientosInventario(db *sql.DB) error {
	query := `
	ALTER TABLE movimientos_inventario 
	ADD COLUMN IF NOT EXISTS id_sucursal    UUID,
	ADD COLUMN IF NOT EXISTS costo_unitario  NUMERIC(12,2) DEFAULT 0,
	ADD COLUMN IF NOT EXISTS precio_unitario NUMERIC(12,2) DEFAULT 0,
	ADD COLUMN IF NOT EXISTS stock_anterior  NUMERIC(12,2) DEFAULT 0,
	ADD COLUMN IF NOT EXISTS stock_posterior NUMERIC(12,2) DEFAULT 0;

	-- Añadir llave foránea si no existe
	DO $$ 
	BEGIN 
		IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_movimientos_sucursal') THEN
			ALTER TABLE movimientos_inventario 
			ADD CONSTRAINT fk_movimientos_sucursal 
			FOREIGN KEY (id_sucursal) REFERENCES sucursal(id_sucursal);
		END IF;
	END $$;

	CREATE INDEX IF NOT EXISTS idx_movimientos_sucursal ON movimientos_inventario(id_sucursal);
	`
	_, err := db.Exec(query)
	return err
}
