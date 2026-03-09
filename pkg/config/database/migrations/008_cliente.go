package migrations

import "database/sql"

func migrateCliente(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS cliente (
		id_cliente      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		empresa_cliente VARCHAR(150) NOT NULL,
		nombre          VARCHAR(150) NOT NULL,
		ruc             VARCHAR(20),
		direccion       VARCHAR(255),
		telefono        VARCHAR(30),
		email           VARCHAR(150),
		id_status       UUID         NOT NULL,

		created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at      TIMESTAMP NULL,

		CONSTRAINT fk_cliente_status
			FOREIGN KEY (id_status)
			REFERENCES estatus(id_status)
			ON UPDATE CASCADE
			ON DELETE RESTRICT
	);

	CREATE INDEX IF NOT EXISTS idx_cliente_id_status  ON cliente(id_status);
	CREATE INDEX IF NOT EXISTS idx_cliente_deleted_at ON cliente(deleted_at);
	`
	_, err := db.Exec(query)
	return err
}
