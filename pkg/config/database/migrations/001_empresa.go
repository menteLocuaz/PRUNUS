package migrations

import "database/sql"

func migrateEmpresa(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS empresa (
		id_empresa  UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Cambio: SERIAL -> UUID
		nombre      VARCHAR(255) NOT NULL,
		rut         VARCHAR(20)  NOT NULL UNIQUE,
		id_status   UUID         NOT NULL,

		created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at  TIMESTAMP NULL,

		CONSTRAINT fk_empresa_status FOREIGN KEY (id_status) REFERENCES estatus(id_status)
	);

	CREATE INDEX IF NOT EXISTS idx_empresa_id_status   ON empresa(id_status);
	CREATE INDEX IF NOT EXISTS idx_empresa_deleted_at ON empresa(deleted_at);
	`
	_, err := db.Exec(query)
	return err
}
