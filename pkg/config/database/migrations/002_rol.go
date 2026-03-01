package migrations

import "database/sql"

func migrateRol(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS rol (
		id_rol      SERIAL PRIMARY KEY,
		nombre_rol  VARCHAR(100) NOT NULL,
		estado      INTEGER      NOT NULL DEFAULT 1,

		created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at  TIMESTAMP NULL
	);

	CREATE INDEX IF NOT EXISTS idx_rol_estado     ON rol(estado);
	CREATE INDEX IF NOT EXISTS idx_rol_deleted_at ON rol(deleted_at);
	`
	_, err := db.Exec(query)
	return err
}
