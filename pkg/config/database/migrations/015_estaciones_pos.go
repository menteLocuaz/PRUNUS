package migrations

import "database/sql"

func migrateEstacionesPos(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS estaciones_pos (
		id_estacion     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		codigo          VARCHAR(50)  NOT NULL UNIQUE,
		nombre          VARCHAR(255) NOT NULL,
		ip              VARCHAR(50)  NOT NULL,
		id_sucursal     UUID      NOT NULL,
		id_status       UUID      NOT NULL,

		created_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at      TIMESTAMP    NULL,

		CONSTRAINT fk_estaciones_pos_sucursal FOREIGN KEY (id_sucursal) REFERENCES sucursal(id_sucursal),
		CONSTRAINT fk_estaciones_pos_status   FOREIGN KEY (id_status)   REFERENCES estatus(id_status)
	);

	CREATE INDEX IF NOT EXISTS idx_estaciones_pos_sucursal   ON estaciones_pos(id_sucursal);
	CREATE INDEX IF NOT EXISTS idx_estaciones_pos_status     ON estaciones_pos(id_status);
	CREATE INDEX IF NOT EXISTS idx_estaciones_pos_deleted_at ON estaciones_pos(deleted_at);
	`
	_, err := db.Exec(query)
	return err
}
