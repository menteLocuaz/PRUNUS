package migrations

import "database/sql"

func migrateOrdenPedido(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS orden_pedido (
		id_orden_pedido    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		odp_fecha_creacion TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
		odp_observacion    VARCHAR(255),
		id_user_pos        UUID      NOT NULL,
		id_periodo         UUID      NOT NULL,
		id_estacion        UUID      NOT NULL,
		id_status          UUID      NOT NULL,
		direccion          VARCHAR(255),
		canal              VARCHAR(50)  NOT NULL,

		created_at         TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at         TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at         TIMESTAMP    NULL,

		CONSTRAINT fk_orden_pedido_user_pos  FOREIGN KEY (id_user_pos)  REFERENCES usuario(id_usuario),
		CONSTRAINT fk_orden_pedido_periodo   FOREIGN KEY (id_periodo)   REFERENCES periodo(id_periodo),
		CONSTRAINT fk_orden_pedido_estacion  FOREIGN KEY (id_estacion)  REFERENCES estaciones_pos(id_estacion),
		CONSTRAINT fk_orden_pedido_status    FOREIGN KEY (id_status)    REFERENCES estatus(id_status)
	);

	CREATE INDEX IF NOT EXISTS idx_orden_pedido_user_pos  ON orden_pedido(id_user_pos);
	CREATE INDEX IF NOT EXISTS idx_orden_pedido_periodo   ON orden_pedido(id_periodo);
	CREATE INDEX IF NOT EXISTS idx_orden_pedido_estacion  ON orden_pedido(id_estacion);
	CREATE INDEX IF NOT EXISTS idx_orden_pedido_status    ON orden_pedido(id_status);
	CREATE INDEX IF NOT EXISTS idx_orden_pedido_canal     ON orden_pedido(canal);
	CREATE INDEX IF NOT EXISTS idx_orden_pedido_deleted_at ON orden_pedido(deleted_at);
	`
	_, err := db.Exec(query)
	return err
}
