package migrations

import "database/sql"

func migrateLogSistema(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS log_sistema (
		id_log      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		id_usuario  UUID      NOT NULL,
		id_modulo   UUID,          -- El módulo desde donde se realiza la acción
		accion      VARCHAR(100) NOT NULL,
		tabla       VARCHAR(100) NOT NULL,
		registro_id UUID      NOT NULL,
		fecha       TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
		ip          VARCHAR(50),

		created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at  TIMESTAMP    NULL,

		CONSTRAINT fk_log_sistema_usuario FOREIGN KEY (id_usuario) REFERENCES usuario(id_usuario),
		CONSTRAINT fk_log_sistema_modulo  FOREIGN KEY (id_modulo)  REFERENCES modulo(id_modulo)
	);

	CREATE INDEX IF NOT EXISTS idx_log_sistema_usuario    ON log_sistema(id_usuario);
	CREATE INDEX IF NOT EXISTS idx_log_sistema_modulo     ON log_sistema(id_modulo);
	CREATE INDEX IF NOT EXISTS idx_log_sistema_fecha      ON log_sistema(fecha);
	CREATE INDEX IF NOT EXISTS idx_log_sistema_tabla      ON log_sistema(tabla);
	CREATE INDEX IF NOT EXISTS idx_log_sistema_deleted_at ON log_sistema(deleted_at);
	`
	_, err := db.Exec(query)
	return err
}
