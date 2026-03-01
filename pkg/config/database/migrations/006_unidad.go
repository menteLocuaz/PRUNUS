package migrations

import "database/sql"

func migrateUnidad(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS unidad (
		id_unidad    SERIAL PRIMARY KEY,
		nombre       VARCHAR(100) NOT NULL,
		id_sucursal  INTEGER      NOT NULL,

		created_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at   TIMESTAMP NULL,

		CONSTRAINT fk_unidad_sucursal
			FOREIGN KEY (id_sucursal)
			REFERENCES sucursal(id_sucursal)
			ON UPDATE CASCADE
			ON DELETE RESTRICT
	);

	CREATE INDEX IF NOT EXISTS idx_unidad_id_sucursal ON unidad(id_sucursal);
	CREATE INDEX IF NOT EXISTS idx_unidad_deleted_at  ON unidad(deleted_at);
	`
	_, err := db.Exec(query)
	return err
}
