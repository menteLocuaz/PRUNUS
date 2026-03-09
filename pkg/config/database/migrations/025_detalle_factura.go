package migrations

import "database/sql"

func migrateDetalleFactura(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS detalle_factura (
		id_detalle_factura UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		id_factura         UUID       NOT NULL,
		id_producto        UUID       NOT NULL,
		cantidad           NUMERIC(12,2) NOT NULL DEFAULT 0,
		precio             NUMERIC(12,2) NOT NULL DEFAULT 0,
		subtotal           NUMERIC(12,2) NOT NULL DEFAULT 0,
		impuesto           NUMERIC(12,2) NOT NULL DEFAULT 0,
		total              NUMERIC(12,2) NOT NULL DEFAULT 0,

		created_at         TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at         TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at         TIMESTAMP     NULL,

		CONSTRAINT fk_detalle_factura_factura  FOREIGN KEY (id_factura)  REFERENCES factura(id_factura),
		CONSTRAINT fk_detalle_factura_producto FOREIGN KEY (id_producto) REFERENCES producto(id_producto)
	);

	CREATE INDEX IF NOT EXISTS idx_detalle_factura_factura  ON detalle_factura(id_factura);
	CREATE INDEX IF NOT EXISTS idx_detalle_factura_producto ON detalle_factura(id_producto);
	CREATE INDEX IF NOT EXISTS idx_detalle_factura_deleted_at ON detalle_factura(deleted_at);
	`
	_, err := db.Exec(query)
	return err
}
