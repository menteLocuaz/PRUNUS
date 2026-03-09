package migrations

import "database/sql"

func migrateDispositivosPos(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS dispositivos_pos (
		id_dispositivo UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		nombre         VARCHAR(150) NOT NULL,
		tipo           VARCHAR(50)  NOT NULL, -- KIOSKO, IMPRESORA, DATÁFONO
		ip             VARCHAR(50),
		id_estacion    UUID      NOT NULL,

		created_at     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at     TIMESTAMP    NULL,

		CONSTRAINT fk_dispositivos_estacion FOREIGN KEY (id_estacion) REFERENCES estaciones_pos(id_estacion)
	);

	CREATE INDEX IF NOT EXISTS idx_dispositivos_estacion   ON dispositivos_pos(id_estacion);
	CREATE INDEX IF NOT EXISTS idx_dispositivos_tipo       ON dispositivos_pos(tipo);
	CREATE INDEX IF NOT EXISTS idx_dispositivos_deleted_at ON dispositivos_pos(deleted_at);
	`
	_, err := db.Exec(query)
	return err
}
