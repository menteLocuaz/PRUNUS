package migrations

import "database/sql"

func migrateFormaPago(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS forma_pago (
		id_forma_pago   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		fmp_codigo      VARCHAR(50)  NOT NULL,
		fmp_descripcion VARCHAR(255) NOT NULL,
		id_status       UUID      NOT NULL,

		created_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at      TIMESTAMP    NULL,

		CONSTRAINT fk_forma_pago_status FOREIGN KEY (id_status) REFERENCES estatus(id_status)
	);

	CREATE INDEX IF NOT EXISTS idx_forma_pago_status     ON forma_pago(id_status);
	CREATE INDEX IF NOT EXISTS idx_forma_pago_deleted_at ON forma_pago(deleted_at);
	`
	_, err := db.Exec(query)
	return err
}
