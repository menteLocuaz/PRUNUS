package migrations

import "database/sql"

func migrateInventario(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS inventario (
		id_inventario   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		id_producto     UUID       NOT NULL,
		id_sucursal     UUID       NOT NULL,
		stock_actual    NUMERIC(12,2) NOT NULL DEFAULT 0,
		stock_minimo    NUMERIC(12,2) NOT NULL DEFAULT 0,
		stock_maximo    NUMERIC(12,2) NOT NULL DEFAULT 0,
		precio_compra   NUMERIC(12,2) NOT NULL DEFAULT 0,
		precio_venta    NUMERIC(12,2) NOT NULL DEFAULT 0,

		created_at      TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at      TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at      TIMESTAMP     NULL,

		CONSTRAINT fk_inventario_producto FOREIGN KEY (id_producto) REFERENCES producto(id_producto),
		CONSTRAINT fk_inventario_sucursal FOREIGN KEY (id_sucursal) REFERENCES sucursal(id_sucursal)
	);

	CREATE INDEX IF NOT EXISTS idx_inventario_producto   ON inventario(id_producto);
	CREATE INDEX IF NOT EXISTS idx_inventario_sucursal   ON inventario(id_sucursal);
	CREATE INDEX IF NOT EXISTS idx_inventario_deleted_at ON inventario(deleted_at);
	`
	_, err := db.Exec(query)
	return err
}
