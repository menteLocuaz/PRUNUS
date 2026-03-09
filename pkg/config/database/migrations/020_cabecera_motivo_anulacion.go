package migrations

import "database/sql"

func migrateCabeceraMotivoAnulacion(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS cabecera_motivo_anulacion (
		id_cabecera_motivo UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		descripcion        VARCHAR(255) NOT NULL,
		estado             INTEGER      NOT NULL DEFAULT 1,

		created_at         TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at         TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at         TIMESTAMP    NULL
	);

	CREATE INDEX IF NOT EXISTS idx_cabecera_motivo_estado     ON cabecera_motivo_anulacion(estado);
	CREATE INDEX IF NOT EXISTS idx_cabecera_motivo_deleted_at ON cabecera_motivo_anulacion(deleted_at);
	`
	_, err := db.Exec(query)
	return err
}
