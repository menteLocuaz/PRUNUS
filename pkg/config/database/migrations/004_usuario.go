package migrations

import "database/sql"

func migrateUsuario(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS usuario (
		id_usuario    UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Cambio: SERIAL -> UUID
		id_sucursal   UUID         NOT NULL,
		id_rol        UUID         NOT NULL,
		id_status     UUID         NOT NULL,

		email         VARCHAR(150) NOT NULL UNIQUE,
		usu_nombre    VARCHAR(150) NOT NULL,
		usu_dni       VARCHAR(30)  NOT NULL UNIQUE,
		usu_telefono  VARCHAR(30),
		password      TEXT         NOT NULL,

		-- Campos Supermercado / POS
		usu_tarjeta_nfc VARCHAR(100) UNIQUE,
		usu_pin_pos     VARCHAR(100),
		nombre_ticket   VARCHAR(50),

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
			ON DELETE RESTRICT,

		CONSTRAINT fk_usuario_status
			FOREIGN KEY (id_status)
			REFERENCES estatus(id_status)
			ON UPDATE CASCADE
			ON DELETE RESTRICT
	);

	CREATE INDEX IF NOT EXISTS idx_usuario_id_sucursal ON usuario(id_sucursal);
	CREATE INDEX IF NOT EXISTS idx_usuario_id_rol      ON usuario(id_rol);
	CREATE INDEX IF NOT EXISTS idx_usuario_id_status   ON usuario(id_status);
	CREATE INDEX IF NOT EXISTS idx_usuario_deleted_at  ON usuario(deleted_at);
	`
	_, err := db.Exec(query)
	return err
}
