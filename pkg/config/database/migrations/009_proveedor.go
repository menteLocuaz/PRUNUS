package migrations

import "database/sql"

func migrateProveedor(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS proveedor (
		id_proveedor  SERIAL PRIMARY KEY,
		nombre        VARCHAR(150) NOT NULL,
		ruc           VARCHAR(20),
		telefono      VARCHAR(30),
		direccion     VARCHAR(255),
		email         VARCHAR(150),
		estado        INTEGER      NOT NULL DEFAULT 1,
		id_sucursal   INTEGER      NOT NULL,
		id_empresa    INTEGER      NOT NULL,

		created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at    TIMESTAMP NULL,

		CONSTRAINT fk_proveedor_sucursal
			FOREIGN KEY (id_sucursal)
			REFERENCES sucursal(id_sucursal)
			ON UPDATE CASCADE
			ON DELETE RESTRICT,

		CONSTRAINT fk_proveedor_empresa
			FOREIGN KEY (id_empresa)
			REFERENCES empresa(id_empresa)
			ON UPDATE CASCADE
			ON DELETE RESTRICT
	);

	CREATE INDEX IF NOT EXISTS idx_proveedor_id_sucursal ON proveedor(id_sucursal);
	CREATE INDEX IF NOT EXISTS idx_proveedor_id_empresa  ON proveedor(id_empresa);
	CREATE INDEX IF NOT EXISTS idx_proveedor_estado      ON proveedor(estado);
	CREATE INDEX IF NOT EXISTS idx_proveedor_deleted_at  ON proveedor(deleted_at);
	`
	_, err := db.Exec(query)
	return err
}
