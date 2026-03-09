package migrations

import "database/sql"

func migrateMotivoAnulacion(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS motivo_anulacion (
		id_motivo_anulacion UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		id_cabecera_motivo  UUID      NOT NULL,
		descripcion         VARCHAR(255) NOT NULL,
		id_status           UUID      NOT NULL,

		created_at          TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at          TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at          TIMESTAMP    NULL,

		CONSTRAINT fk_motivo_anulacion_cabecera FOREIGN KEY (id_cabecera_motivo) REFERENCES cabecera_motivo_anulacion(id_cabecera_motivo),
		CONSTRAINT fk_motivo_anulacion_status   FOREIGN KEY (id_status)          REFERENCES estatus(id_status)
	);

	CREATE INDEX IF NOT EXISTS idx_motivo_anulacion_cabecera  ON motivo_anulacion(id_cabecera_motivo);
	CREATE INDEX IF NOT EXISTS idx_motivo_anulacion_status    ON motivo_anulacion(id_status);
	CREATE INDEX IF NOT EXISTS idx_motivo_anulacion_deleted_at ON motivo_anulacion(deleted_at);
	`
	_, err := db.Exec(query)
	return err
}
