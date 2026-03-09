package migrations

import "database/sql"

func migrateOrdenAgregador(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS orden_agregador (
		id_orden_agregador UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		id_orden_pedido    UUID      NOT NULL,
		id_agregador       UUID      NOT NULL,
		codigo_externo     VARCHAR(100),
		datos_agregador    JSONB,
		fecha              TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,

		created_at         TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at         TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at         TIMESTAMP    NULL,

		CONSTRAINT fk_orden_agregador_pedido    FOREIGN KEY (id_orden_pedido) REFERENCES orden_pedido(id_orden_pedido),
		CONSTRAINT fk_orden_agregador_plataforma FOREIGN KEY (id_agregador)    REFERENCES agregadores(id_agregador)
	);

	CREATE INDEX IF NOT EXISTS idx_orden_agregador_pedido    ON orden_agregador(id_orden_pedido);
	CREATE INDEX IF NOT EXISTS idx_orden_agregador_plataforma ON orden_agregador(id_agregador);
	CREATE INDEX IF NOT EXISTS idx_orden_agregador_externo   ON orden_agregador(codigo_externo);
	CREATE INDEX IF NOT EXISTS idx_orden_agregador_deleted_at ON orden_agregador(deleted_at);
	`
	_, err := db.Exec(query)
	return err
}
