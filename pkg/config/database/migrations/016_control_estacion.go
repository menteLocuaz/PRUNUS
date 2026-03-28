package migrations

import "database/sql"

func migrateControlEstacion(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS control_estacion (
		id_control_estacion  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		id_estacion          UUID      NOT NULL,
		fecha_inicio         TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
		fecha_salida         TIMESTAMP    NULL,
		fondo_base           DECIMAL(18,2) NOT NULL DEFAULT 0,
		usuario_asignado     UUID      NOT NULL,
		fecha_fondo_aceptado TIMESTAMP    NULL,
		usuario_retiro_fondo UUID      NULL,
		fondo_retirado       DECIMAL(18,2) NULL,
		id_status            UUID      NOT NULL,
		id_user_pos          UUID      NOT NULL,
		id_periodo           UUID      NOT NULL,
		ctrc_motivo_descuadre TEXT      NULL,

		created_at           TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at           TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at           TIMESTAMP    NULL,

		CONSTRAINT fk_control_estacion_estacion FOREIGN KEY (id_estacion)      REFERENCES estaciones_pos(id_estacion),
		CONSTRAINT fk_control_estacion_asignado FOREIGN KEY (usuario_asignado) REFERENCES usuario(id_usuario),
		CONSTRAINT fk_control_estacion_retiro   FOREIGN KEY (usuario_retiro_fondo) REFERENCES usuario(id_usuario),
		CONSTRAINT fk_control_estacion_status   FOREIGN KEY (id_status)        REFERENCES estatus(id_status),
		CONSTRAINT fk_control_estacion_user_pos FOREIGN KEY (id_user_pos)      REFERENCES usuario(id_usuario),
		CONSTRAINT fk_control_estacion_periodo  FOREIGN KEY (id_periodo)       REFERENCES periodo(id_periodo)
	);

	CREATE INDEX IF NOT EXISTS idx_control_estacion_estacion ON control_estacion(id_estacion);
	CREATE INDEX IF NOT EXISTS idx_control_estacion_status   ON control_estacion(id_status);
	CREATE INDEX IF NOT EXISTS idx_control_estacion_periodo  ON control_estacion(id_periodo);
	CREATE INDEX IF NOT EXISTS idx_control_estacion_deleted_at ON control_estacion(deleted_at);
	`
	_, err := db.Exec(query)
	return err
}
