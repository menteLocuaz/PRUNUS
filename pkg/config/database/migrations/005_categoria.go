package migrations

import "database/sql"

func migrateCategoria(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS categoria (
		id_categoria  SERIAL PRIMARY KEY,
		nombre        VARCHAR(150) NOT NULL,
		id_sucursal   INTEGER      NOT NULL,

		created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at    TIMESTAMP NULL,

		CONSTRAINT fk_categoria_sucursal
			FOREIGN KEY (id_sucursal)
			REFERENCES sucursal(id_sucursal)
			ON UPDATE CASCADE
			ON DELETE RESTRICT
	);

	CREATE INDEX IF NOT EXISTS idx_categoria_id_sucursal ON categoria(id_sucursal);
	CREATE INDEX IF NOT EXISTS idx_categoria_deleted_at  ON categoria(deleted_at);
	`
	_, err := db.Exec(query)
	return err
}
