package migrations

import "database/sql"

// migrateNormalizeInventory consolida la gestión de precios y stock en la tabla inventario.
// Diseñada para ser segura tanto en instalaciones nuevas como en migraciones de datos legacy.
func migrateNormalizeInventory(db *sql.DB) error {
	query := `
	-- 1. Asegurar restricción de unicidad en inventario para evitar duplicados por sucursal
	DO $$ 
	BEGIN 
		IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'uq_inventario_producto_sucursal') THEN
			ALTER TABLE inventario ADD CONSTRAINT uq_inventario_producto_sucursal UNIQUE (id_producto, id_sucursal);
		END IF;
	END $$;

	-- 2. Migrar datos operativos de producto a inventario (Solo si las columnas legacy existen)
	-- Esto permite que la migración corra sin errores en bases de datos ya normalizadas.
	DO $$ 
	BEGIN 
		IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='producto' AND column_name='precio_venta') THEN
			INSERT INTO inventario (
				id_producto, 
				id_sucursal, 
				stock_actual, 
				precio_compra, 
				precio_venta, 
				created_at, 
				updated_at
			)
			SELECT 
				id_producto, 
				id_sucursal, 
				stock, 
				precio_compra, 
				precio_venta, 
				created_at, 
				updated_at
			FROM producto
			ON CONFLICT (id_producto, id_sucursal) DO UPDATE SET 
				stock_actual = EXCLUDED.stock_actual,
				precio_compra = EXCLUDED.precio_compra,
				precio_venta = EXCLUDED.precio_venta,
				updated_at = CURRENT_TIMESTAMP;

			-- 3. Limpieza de columnas redundantes en la tabla producto
			ALTER TABLE producto DROP COLUMN IF EXISTS precio_compra;
			ALTER TABLE producto DROP COLUMN IF EXISTS precio_venta;
			ALTER TABLE producto DROP COLUMN IF EXISTS stock;
			ALTER TABLE producto DROP COLUMN IF EXISTS id_sucursal;
			
			DROP INDEX IF EXISTS idx_producto_id_sucursal;
		END IF;
	END $$;

	-- 4. Documentación de esquema
	COMMENT ON TABLE producto IS 'Catálogo maestro de productos (Atómico/Global)';
	COMMENT ON TABLE inventario IS 'Gestión operativa de productos por sucursal (Precios y Existencias)';
	`
	_, err := db.Exec(query)
	return err
}
