package migrations

import "database/sql"

func migrateSucursal(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS sucursal (
		id_sucursal     SERIAL PRIMARY KEY,
		id_empresa      INTEGER      NOT NULL,
		nombre_sucursal VARCHAR(255) NOT NULL,
		estado          INTEGER      NOT NULL DEFAULT 1,

		created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at  TIMESTAMP NULL,

		CONSTRAINT fk_sucursal_empresa
			FOREIGN KEY (id_empresa)
			REFERENCES empresa(id_empresa)
			ON UPDATE CASCADE
			ON DELETE RESTRICT
	);

	CREATE INDEX IF NOT EXISTS idx_sucursal_id_empresa  ON sucursal(id_empresa);
	CREATE INDEX IF NOT EXISTS idx_sucursal_estado      ON sucursal(estado);
	CREATE INDEX IF NOT EXISTS idx_sucursal_deleted_at  ON sucursal(deleted_at);
	`
	_, err := db.Exec(query)
	return err
}
