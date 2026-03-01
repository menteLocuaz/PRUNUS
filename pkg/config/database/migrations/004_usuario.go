package migrations

import "database/sql"

func migrateUsuario(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS usuario (
		id_usuario    SERIAL PRIMARY KEY,
		id_sucursal   INTEGER      NOT NULL,
		id_rol        INTEGER      NOT NULL,

		email         VARCHAR(150) NOT NULL UNIQUE,
		usu_nombre    VARCHAR(150) NOT NULL,
		usu_dni       VARCHAR(30)  NOT NULL UNIQUE,
		usu_telefono  VARCHAR(30),
		password      TEXT         NOT NULL,
		estado        INTEGER      NOT NULL DEFAULT 1,

		created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at    TIMESTAMP NULL,

		CONSTRAINT fk_usuario_sucursal
			FOREIGN KEY (id_sucursal)
			REFERENCES sucursal(id_sucursal)
			ON UPDATE CASCADE
			ON DELETE RESTRICT,

		CONSTRAINT fk_usuario_rol
			FOREIGN KEY (id_rol)
			REFERENCES rol(id_rol)
			ON UPDATE CASCADE
			ON DELETE RESTRICT
	);

	CREATE INDEX IF NOT EXISTS idx_usuario_id_sucursal ON usuario(id_sucursal);
	CREATE INDEX IF NOT EXISTS idx_usuario_id_rol      ON usuario(id_rol);
	CREATE INDEX IF NOT EXISTS idx_usuario_estado      ON usuario(estado);
	CREATE INDEX IF NOT EXISTS idx_usuario_deleted_at  ON usuario(deleted_at);
	`
	_, err := db.Exec(query)
	return err
}
