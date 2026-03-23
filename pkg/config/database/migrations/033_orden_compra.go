package migrations

import "database/sql"

func migrateOrdenCompra(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS orden_compra (
		id_orden_compra  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		numero_orden     VARCHAR(50)    NOT NULL UNIQUE,
		id_proveedor     UUID           NOT NULL,
		id_sucursal      UUID           NOT NULL,
		id_usuario       UUID           NOT NULL,
		id_moneda        UUID           NOT NULL,
		id_status        UUID           NOT NULL,
		fecha_emision    TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
		fecha_recepcion  TIMESTAMP,
		subtotal         NUMERIC(12,2)  NOT NULL DEFAULT 0,
		impuesto         NUMERIC(12,2)  NOT NULL DEFAULT 0,
		total            NUMERIC(12,2)  NOT NULL DEFAULT 0,
		observaciones    TEXT,

		created_at       TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at       TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at       TIMESTAMP      NULL,

		CONSTRAINT fk_compra_proveedor FOREIGN KEY (id_proveedor) REFERENCES proveedor(id_proveedor),
		CONSTRAINT fk_compra_sucursal  FOREIGN KEY (id_sucursal)  REFERENCES sucursal(id_sucursal),
		CONSTRAINT fk_compra_usuario   FOREIGN KEY (id_usuario)   REFERENCES usuario(id_usuario),
		CONSTRAINT fk_compra_moneda    FOREIGN KEY (id_moneda)    REFERENCES moneda(id_moneda),
		CONSTRAINT fk_compra_status    FOREIGN KEY (id_status)    REFERENCES estatus(id_status)
	);

	CREATE TABLE IF NOT EXISTS detalle_orden_compra (
		id_detalle_compra UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		id_orden_compra   UUID           NOT NULL,
		id_producto       UUID           NOT NULL,
		cantidad_pedida   NUMERIC(12,2)  NOT NULL DEFAULT 0,
		cantidad_recibida NUMERIC(12,2)  NOT NULL DEFAULT 0,
		precio_unitario   NUMERIC(12,2)  NOT NULL DEFAULT 0,
		impuesto          NUMERIC(12,2)  NOT NULL DEFAULT 0,
		total             NUMERIC(12,2)  NOT NULL DEFAULT 0,

		CONSTRAINT fk_detalle_compra_cabecera FOREIGN KEY (id_orden_compra) REFERENCES orden_compra(id_orden_compra) ON DELETE CASCADE,
		CONSTRAINT fk_detalle_compra_producto FOREIGN KEY (id_producto) REFERENCES producto(id_producto)
	);

	CREATE INDEX IF NOT EXISTS idx_compra_proveedor ON orden_compra(id_proveedor);
	CREATE INDEX IF NOT EXISTS idx_compra_sucursal  ON orden_compra(id_sucursal);
	CREATE INDEX IF NOT EXISTS idx_compra_status    ON orden_compra(id_status);
	`
	_, err := db.Exec(query)
	return err
}
