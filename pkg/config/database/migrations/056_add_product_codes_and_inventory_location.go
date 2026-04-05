package migrations

import "database/sql"

func migrateAddProductCodesAndInventoryLocation(db *sql.DB) error {
	query := `
	-- Añadir campos de identificación a productos
	ALTER TABLE producto ADD COLUMN IF NOT EXISTS codigo_barras VARCHAR(50);
	ALTER TABLE producto ADD COLUMN IF NOT EXISTS sku VARCHAR(50);
	
	-- Añadir ubicación física al inventario
	ALTER TABLE inventario ADD COLUMN IF NOT EXISTS ubicacion VARCHAR(100);

	-- Índices para búsquedas rápidas por código
	CREATE INDEX IF NOT EXISTS idx_producto_codigo_barras ON producto(codigo_barras) WHERE deleted_at IS NULL;
	CREATE INDEX IF NOT EXISTS idx_producto_sku           ON producto(sku)           WHERE deleted_at IS NULL;
	`
	_, err := db.Exec(query)
	return err
}
