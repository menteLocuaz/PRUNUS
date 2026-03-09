package migrations

import "database/sql"

func migrateAuditoriaCaja(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS auditoria_caja (
		id_auditoria         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		id_control_estacion  UUID       NOT NULL,
		tipo_movimiento      VARCHAR(50)   NOT NULL, -- AJUSTE, SOBRANTE, FALTANTE, etc.
		valor                DECIMAL(18,2) NOT NULL DEFAULT 0,
		fecha                TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
		id_usuario           UUID       NOT NULL,
		descripcion          TEXT,

		created_at           TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at           TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at           TIMESTAMP     NULL,

		CONSTRAINT fk_auditoria_caja_control FOREIGN KEY (id_control_estacion) REFERENCES control_estacion(id_control_estacion),
		CONSTRAINT fk_auditoria_caja_usuario FOREIGN KEY (id_usuario)          REFERENCES usuario(id_usuario)
	);

	CREATE INDEX IF NOT EXISTS idx_auditoria_caja_control    ON auditoria_caja(id_control_estacion);
	CREATE INDEX IF NOT EXISTS idx_auditoria_caja_usuario    ON auditoria_caja(id_usuario);
	CREATE INDEX IF NOT EXISTS idx_auditoria_caja_fecha      ON auditoria_caja(fecha);
	CREATE INDEX IF NOT EXISTS idx_auditoria_caja_deleted_at ON auditoria_caja(deleted_at);
	`
	_, err := db.Exec(query)
	return err
}
