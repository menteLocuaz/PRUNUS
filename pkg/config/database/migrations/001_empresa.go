package migrations

import "database/sql"

func migrateEmpresa(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS empresa (
		id_empresa  SERIAL PRIMARY KEY,
		nombre      VARCHAR(255) NOT NULL,
		rut         VARCHAR(20)  NOT NULL UNIQUE,
		estado      INTEGER      NOT NULL DEFAULT 1,

		created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at  TIMESTAMP NULL
	);

	CREATE INDEX IF NOT EXISTS idx_empresa_estado     ON empresa(estado);
	CREATE INDEX IF NOT EXISTS idx_empresa_deleted_at ON empresa(deleted_at);
	`
	_, err := db.Exec(query)
	return err
}
