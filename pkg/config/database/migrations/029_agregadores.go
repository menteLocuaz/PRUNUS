package migrations

import "database/sql"

func migrateAgregadores(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS agregadores (
		id_agregador UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		nombre       VARCHAR(100) NOT NULL UNIQUE,
		descripcion  VARCHAR(255),

		created_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at   TIMESTAMP NULL
	);

	CREATE INDEX IF NOT EXISTS idx_agregadores_nombre     ON agregadores(nombre);
	CREATE INDEX IF NOT EXISTS idx_agregadores_deleted_at ON agregadores(deleted_at);
	`
	_, err := db.Exec(query)
	return err
}
