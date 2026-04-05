package migrations

import "database/sql"

func migrateLotes(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS lotes (
		id_lote          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		id_producto      UUID       NOT NULL,
		id_sucursal      UUID       NOT NULL,
		codigo_lote      VARCHAR(100) NOT NULL,
		cantidad_inicial NUMERIC(12,2) NOT NULL DEFAULT 0,
		cantidad_actual  NUMERIC(12,2) NOT NULL DEFAULT 0,
		costo_compra     NUMERIC(12,2) NOT NULL DEFAULT 0,
		fecha_vencimiento TIMESTAMP     NULL,
		fecha_recepcion   TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,

		created_at       TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at       TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at       TIMESTAMP     NULL,

		CONSTRAINT fk_lotes_producto FOREIGN KEY (id_producto) REFERENCES producto(id_producto),
		CONSTRAINT fk_lotes_sucursal FOREIGN KEY (id_sucursal) REFERENCES sucursal(id_sucursal)
	);

	-- Índices de rendimiento para búsqueda por producto y sucursal (PEPS/UEPS)
	CREATE INDEX IF NOT EXISTS idx_lotes_producto_sucursal ON lotes(id_producto, id_sucursal, fecha_recepcion) WHERE deleted_at IS NULL AND cantidad_actual > 0;
	CREATE INDEX IF NOT EXISTS idx_lotes_vencimiento      ON lotes(fecha_vencimiento)                     WHERE deleted_at IS NULL AND cantidad_actual > 0;
	`
	_, err := db.Exec(query)
	return err
}
