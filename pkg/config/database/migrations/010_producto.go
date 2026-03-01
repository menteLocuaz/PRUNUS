package migrations

import "database/sql"

func migrateProducto(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS producto (
		id_producto       SERIAL PRIMARY KEY,
		nombre            VARCHAR(150)   NOT NULL,
		descripcion       TEXT,
		precio_compra     NUMERIC(12,2)  NOT NULL DEFAULT 0,
		precio_venta      NUMERIC(12,2)  NOT NULL DEFAULT 0,
		stock             INTEGER        NOT NULL DEFAULT 0,
		fecha_vencimiento DATE,
		imagen            TEXT,
		estado            INTEGER        NOT NULL DEFAULT 1,
		id_sucursal       INTEGER        NOT NULL,
		id_categoria      INTEGER        NOT NULL,
		id_moneda         INTEGER        NOT NULL,
		id_unidad         INTEGER        NOT NULL,

		created_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at        TIMESTAMP NULL,

		CONSTRAINT fk_producto_sucursal
			FOREIGN KEY (id_sucursal)
			REFERENCES sucursal(id_sucursal)
			ON UPDATE CASCADE
			ON DELETE RESTRICT,

		CONSTRAINT fk_producto_categoria
			FOREIGN KEY (id_categoria)
			REFERENCES categoria(id_categoria)
			ON UPDATE CASCADE
			ON DELETE RESTRICT,

		CONSTRAINT fk_producto_moneda
			FOREIGN KEY (id_moneda)
			REFERENCES moneda(id_moneda)
			ON UPDATE CASCADE
			ON DELETE RESTRICT,

		CONSTRAINT fk_producto_unidad
			FOREIGN KEY (id_unidad)
			REFERENCES unidad(id_unidad)
			ON UPDATE CASCADE
			ON DELETE RESTRICT
	);

	CREATE INDEX IF NOT EXISTS idx_producto_id_sucursal  ON producto(id_sucursal);
	CREATE INDEX IF NOT EXISTS idx_producto_id_categoria ON producto(id_categoria);
	CREATE INDEX IF NOT EXISTS idx_producto_id_moneda    ON producto(id_moneda);
	CREATE INDEX IF NOT EXISTS idx_producto_id_unidad    ON producto(id_unidad);
	CREATE INDEX IF NOT EXISTS idx_producto_estado       ON producto(estado);
	CREATE INDEX IF NOT EXISTS idx_producto_deleted_at   ON producto(deleted_at);
	`
	_, err := db.Exec(query)
	return err
}
