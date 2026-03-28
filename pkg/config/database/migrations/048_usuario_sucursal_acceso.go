package migrations

import "database/sql"

func migrateUsuarioSucursalAcceso(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS usuario_sucursal_acceso (
		id_usuario    UUID NOT NULL,
		id_sucursal   UUID NOT NULL,
		created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		
		PRIMARY KEY (id_usuario, id_sucursal),
		CONSTRAINT fk_usa_usuario FOREIGN KEY (id_usuario) REFERENCES usuario(id_usuario) ON DELETE CASCADE,
		CONSTRAINT fk_usa_sucursal FOREIGN KEY (id_sucursal) REFERENCES sucursal(id_sucursal) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_usa_id_usuario ON usuario_sucursal_acceso(id_usuario);
	CREATE INDEX IF NOT EXISTS idx_usa_id_sucursal ON usuario_sucursal_acceso(id_sucursal);
	`
	_, err := db.Exec(query)
	return err
}
