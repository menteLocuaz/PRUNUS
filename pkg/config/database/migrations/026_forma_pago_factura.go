package migrations

import "database/sql"

func migrateFormaPagoFactura(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS forma_pago_factura (
		id_pago_factura UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		id_factura      UUID       NOT NULL,
		id_forma_pago   UUID       NOT NULL,
		valor_billete   DECIMAL(18,2) NOT NULL DEFAULT 0,
		total_pagar     DECIMAL(18,2) NOT NULL DEFAULT 0,
		fecha           TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
		id_usuario      UUID       NOT NULL,
		id_status       UUID       NOT NULL,

		created_at      TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at      TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at      TIMESTAMP     NULL,

		CONSTRAINT fk_pago_factura_factura    FOREIGN KEY (id_factura)    REFERENCES factura(id_factura),
		CONSTRAINT fk_pago_factura_forma_pago FOREIGN KEY (id_forma_pago) REFERENCES forma_pago(id_forma_pago),
		CONSTRAINT fk_pago_factura_usuario    FOREIGN KEY (id_usuario)    REFERENCES usuario(id_usuario),
		CONSTRAINT fk_pago_factura_status     FOREIGN KEY (id_status)     REFERENCES estatus(id_status)
	);

	CREATE INDEX IF NOT EXISTS idx_pago_factura_factura    ON forma_pago_factura(id_factura);
	CREATE INDEX IF NOT EXISTS idx_pago_factura_forma_pago ON forma_pago_factura(id_forma_pago);
	CREATE INDEX IF NOT EXISTS idx_pago_factura_deleted_at ON forma_pago_factura(deleted_at);
	`
	_, err := db.Exec(query)
	return err
}
