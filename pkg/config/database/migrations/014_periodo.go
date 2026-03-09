package migrations

import "database/sql"

func migratePeriodo(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS periodo (
		id_periodo           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		prd_fecha_apertura   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		prd_fecha_cierre     TIMESTAMP NULL,
		prd_usuario_apertura UUID   NOT NULL,
		prd_usuario_cierre   UUID   NULL,
		id_status            UUID   NOT NULL,

		created_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at           TIMESTAMP NULL,

		CONSTRAINT fk_periodo_usuario_apertura FOREIGN KEY (prd_usuario_apertura) REFERENCES usuario(id_usuario),
		CONSTRAINT fk_periodo_usuario_cierre   FOREIGN KEY (prd_usuario_cierre)   REFERENCES usuario(id_usuario),
		CONSTRAINT fk_periodo_status           FOREIGN KEY (id_status)            REFERENCES estatus(id_status)
	);

	CREATE INDEX IF NOT EXISTS idx_periodo_usuario_apertura ON periodo(prd_usuario_apertura);
	CREATE INDEX IF NOT EXISTS idx_periodo_usuario_cierre   ON periodo(prd_usuario_cierre);
	CREATE INDEX IF NOT EXISTS idx_periodo_status           ON periodo(id_status);
	CREATE INDEX IF NOT EXISTS idx_periodo_deleted_at       ON periodo(deleted_at);
	`
	_, err := db.Exec(query)
	return err
}
