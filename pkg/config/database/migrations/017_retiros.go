package migrations

import "database/sql"

func migrateRetiros(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS retiros (
		id_retiro                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		arc_valor                DECIMAL(18,2) NOT NULL DEFAULT 0,
		arc_numero_transacciones INTEGER       NOT NULL DEFAULT 0,
		id_control_estacion      UUID       NOT NULL,
		id_forma_pago            UUID       NOT NULL,
		id_user_pos              UUID       NOT NULL,
		usuario_inicia           UUID       NOT NULL,
		usuario_finaliza         UUID       NULL,
		fecha_inicio             TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
		fecha_finaliza           TIMESTAMP     NULL,
		id_status                UUID       NOT NULL,
		pos_calculado            DECIMAL(18,2) NOT NULL DEFAULT 0,
		diferencia_valor         DECIMAL(18,2) NOT NULL DEFAULT 0,
		retiro_valor             DECIMAL(18,2) NOT NULL DEFAULT 0,
		tpenv_id                 INTEGER       NOT NULL DEFAULT -1,

		created_at               TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at               TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at               TIMESTAMP     NULL,

		CONSTRAINT fk_retiros_control_estacion FOREIGN KEY (id_control_estacion) REFERENCES control_estacion(id_control_estacion),
		CONSTRAINT fk_retiros_forma_pago       FOREIGN KEY (id_forma_pago)       REFERENCES forma_pago(id_forma_pago),
		CONSTRAINT fk_retiros_user_pos         FOREIGN KEY (id_user_pos)         REFERENCES usuario(id_usuario),
		CONSTRAINT fk_retiros_usuario_inicia   FOREIGN KEY (usuario_inicia)      REFERENCES usuario(id_usuario),
		CONSTRAINT fk_retiros_usuario_finaliza FOREIGN KEY (usuario_finaliza)    REFERENCES usuario(id_usuario),
		CONSTRAINT fk_retiros_status           FOREIGN KEY (id_status)           REFERENCES estatus(id_status)
	);

	CREATE INDEX IF NOT EXISTS idx_retiros_control_estacion ON retiros(id_control_estacion);
	CREATE INDEX IF NOT EXISTS idx_retiros_forma_pago       ON retiros(id_forma_pago);
	CREATE INDEX IF NOT EXISTS idx_retiros_status           ON retiros(id_status);
	CREATE INDEX IF NOT EXISTS idx_retiros_deleted_at       ON retiros(deleted_at);
	`
	_, err := db.Exec(query)
	return err
}
