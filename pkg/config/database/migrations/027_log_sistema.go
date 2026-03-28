package migrations

import "database/sql"

func migrateLogSistema(db *sql.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS log_sistema (
			id_log      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			id_usuario  UUID      NOT NULL,
			id_modulo   UUID,
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
		);`,
		// Asegurar que la columna id_modulo exista si la tabla ya fue creada previamente
		`ALTER TABLE log_sistema ADD COLUMN IF NOT EXISTS id_modulo UUID;`,
		// Asegurar columnas de auditoría por si la tabla es antigua
		`ALTER TABLE log_sistema ADD COLUMN IF NOT EXISTS created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;`,
		`ALTER TABLE log_sistema ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;`,
		`ALTER TABLE log_sistema ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP NULL;`,

		// Re-intentar crear la FK si se añadió la columna recién
		`DO $$ 
		BEGIN 
			IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints 
						  WHERE constraint_name='fk_log_sistema_modulo') THEN
				ALTER TABLE log_sistema ADD CONSTRAINT fk_log_sistema_modulo 
				FOREIGN KEY (id_modulo) REFERENCES modulo(id_modulo);
			END IF;
		END $$;`,

		`CREATE INDEX IF NOT EXISTS idx_log_sistema_usuario ON log_sistema(id_usuario);`,
		`CREATE INDEX IF NOT EXISTS idx_log_sistema_modulo ON log_sistema(id_modulo);`,
		`CREATE INDEX IF NOT EXISTS idx_log_sistema_fecha ON log_sistema(fecha);`,
		`CREATE INDEX IF NOT EXISTS idx_log_sistema_tabla ON log_sistema(tabla);`,
		`CREATE INDEX IF NOT EXISTS idx_log_sistema_deleted_at ON log_sistema(deleted_at);`,
	}

	for _, stmt := range statements {
		if _, err := db.Exec(stmt); err != nil {
			return err
		}
	}

	return nil
}
