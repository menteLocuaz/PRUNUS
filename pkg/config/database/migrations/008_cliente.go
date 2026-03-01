package migrations

import "database/sql"

func migrateCliente(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS cliente (
		id_cliente      SERIAL PRIMARY KEY,
		empresa_cliente VARCHAR(150) NOT NULL,
		nombre          VARCHAR(150) NOT NULL,
		ruc             VARCHAR(20),
		direccion       VARCHAR(255),
		telefono        VARCHAR(30),
		email           VARCHAR(150),
		estado          INTEGER      NOT NULL DEFAULT 1,

		created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at      TIMESTAMP NULL
	);

	CREATE INDEX IF NOT EXISTS idx_cliente_estado     ON cliente(estado);
	CREATE INDEX IF NOT EXISTS idx_cliente_deleted_at ON cliente(deleted_at);
	`
	_, err := db.Exec(query)
	return err
}
