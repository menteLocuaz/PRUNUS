package migrations

import "database/sql"

func migrateImpuesto(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS impuesto (
		id_impuesto UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		nombre      VARCHAR(100) NOT NULL,
		porcentaje  DECIMAL(5,2) NOT NULL DEFAULT 0,
		tipo        VARCHAR(50)  NOT NULL,

		created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at  TIMESTAMP    NULL
	);

	CREATE INDEX IF NOT EXISTS idx_impuesto_tipo       ON impuesto(tipo);
	CREATE INDEX IF NOT EXISTS idx_impuesto_deleted_at ON impuesto(deleted_at);
	`
	_, err := db.Exec(query)
	return err
}
