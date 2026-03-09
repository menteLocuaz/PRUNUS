package migrations

import "database/sql"

func migrateRol(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS rol (
		id_rol      UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Cambio: SERIAL -> UUID
		nombre_rol  VARCHAR(100) NOT NULL,
		id_sucursal UUID         NOT NULL,
		id_status   UUID         NOT NULL,

		created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at  TIMESTAMP NULL,

		CONSTRAINT fk_rol_sucursal
			FOREIGN KEY (id_sucursal)
			REFERENCES sucursal(id_sucursal)
			ON UPDATE CASCADE
			ON DELETE RESTRICT,

		CONSTRAINT fk_rol_status
			FOREIGN KEY (id_status)
			REFERENCES estatus(id_status)
			ON UPDATE CASCADE
			ON DELETE RESTRICT
	);

	CREATE INDEX IF NOT EXISTS idx_rol_id_status  ON rol(id_status);
	CREATE INDEX IF NOT EXISTS idx_rol_deleted_at ON rol(deleted_at);
	`
	_, err := db.Exec(query)
	return err
}
