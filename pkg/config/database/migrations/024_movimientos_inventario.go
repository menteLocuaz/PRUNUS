package migrations

import "database/sql"

func migrateMovimientosInventario(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS movimientos_inventario (
		id_movimiento   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		id_producto     UUID          NOT NULL,
		id_sucursal     UUID          NOT NULL,
		id_usuario      UUID          NOT NULL,
		tipo_movimiento VARCHAR(50)   NOT NULL, -- VENTA, COMPRA, AJUSTE_ENTRADA, AJUSTE_SALIDA, DEVOLUCION
		cantidad        NUMERIC(12,2) NOT NULL DEFAULT 0,
		costo_unitario  NUMERIC(12,2) NOT NULL DEFAULT 0,
		precio_unitario NUMERIC(12,2) NOT NULL DEFAULT 0,
		stock_anterior  NUMERIC(12,2) NOT NULL DEFAULT 0,
		stock_posterior NUMERIC(12,2) NOT NULL DEFAULT 0,
		fecha           TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
		referencia      VARCHAR(255),

		created_at      TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at      TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at      TIMESTAMP     NULL,

		CONSTRAINT fk_movimientos_producto FOREIGN KEY (id_producto) REFERENCES producto(id_producto),
		CONSTRAINT fk_movimientos_usuario  FOREIGN KEY (id_usuario)  REFERENCES usuario(id_usuario),
		CONSTRAINT fk_movimientos_sucursal FOREIGN KEY (id_sucursal) REFERENCES sucursal(id_sucursal)
	);

	CREATE INDEX IF NOT EXISTS idx_movimientos_producto   ON movimientos_inventario(id_producto);
	CREATE INDEX IF NOT EXISTS idx_movimientos_usuario    ON movimientos_inventario(id_usuario);
	CREATE INDEX IF NOT EXISTS idx_movimientos_sucursal   ON movimientos_inventario(id_sucursal);
	CREATE INDEX IF NOT EXISTS idx_movimientos_tipo       ON movimientos_inventario(tipo_movimiento);
	CREATE INDEX IF NOT EXISTS idx_movimientos_deleted_at ON movimientos_inventario(deleted_at);
	`
	_, err := db.Exec(query)
	return err
}
