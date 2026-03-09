package migrations

import "database/sql"

func migrateMoneda(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS moneda (
		id_moneda    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		nombre       VARCHAR(100) NOT NULL,
		id_sucursal  UUID         NOT NULL,
		id_status    UUID         NOT NULL,

		created_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at   TIMESTAMP NULL,

		CONSTRAINT fk_moneda_sucursal
			FOREIGN KEY (id_sucursal)
			REFERENCES sucursal(id_sucursal)
			ON UPDATE CASCADE
			ON DELETE RESTRICT,

		CONSTRAINT fk_moneda_status
			FOREIGN KEY (id_status)
			REFERENCES estatus(id_status)
			ON UPDATE CASCADE
			ON DELETE RESTRICT
	);

	CREATE INDEX IF NOT EXISTS idx_moneda_id_sucursal ON moneda(id_sucursal);
	CREATE INDEX IF NOT EXISTS idx_moneda_id_status   ON moneda(id_status);
	CREATE INDEX IF NOT EXISTS idx_moneda_deleted_at  ON moneda(deleted_at);
	`
	_, err := db.Exec(query)
	return err
}
