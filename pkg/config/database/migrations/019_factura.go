package migrations

import "database/sql"

func migrateFactura(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS factura (
		id_factura          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		cfac_fecha_creacion TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
		fac_numero          VARCHAR(50)   NOT NULL UNIQUE,
		cfac_subtotal       DECIMAL(18,2) NOT NULL DEFAULT 0,
		cfac_iva            DECIMAL(18,2) NOT NULL DEFAULT 0,
		cfac_total          DECIMAL(18,2) NOT NULL DEFAULT 0,
		cfac_observacion    VARCHAR(255),
		id_user_pos         UUID       NOT NULL,
		id_estacion         UUID       NOT NULL,
		id_orden_pedido     UUID       NOT NULL,
		id_cliente          UUID       NOT NULL,
		id_motivo_anulacion UUID,
		id_periodo          UUID       NOT NULL,
		id_control_estacion UUID       NOT NULL,
		id_status           UUID       NOT NULL,
		fecha_operacion     TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
		base_impuesto       DECIMAL(18,2) NOT NULL DEFAULT 0,
		impuesto            DECIMAL(18,2) NOT NULL DEFAULT 0,
		valor_impuesto      DECIMAL(18,2) NOT NULL DEFAULT 0,
		metadata            JSONB,

		created_at          TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at          TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at          TIMESTAMP     NULL,

		CONSTRAINT fk_factura_user_pos         FOREIGN KEY (id_user_pos)         REFERENCES usuario(id_usuario),
		CONSTRAINT fk_factura_estacion         FOREIGN KEY (id_estacion)         REFERENCES estaciones_pos(id_estacion),
		CONSTRAINT fk_factura_orden_pedido     FOREIGN KEY (id_orden_pedido)     REFERENCES orden_pedido(id_orden_pedido),
		CONSTRAINT fk_factura_cliente          FOREIGN KEY (id_cliente)          REFERENCES cliente(id_cliente),
		CONSTRAINT fk_factura_periodo          FOREIGN KEY (id_periodo)          REFERENCES periodo(id_periodo),
		CONSTRAINT fk_factura_control_estacion FOREIGN KEY (id_control_estacion) REFERENCES control_estacion(id_control_estacion),
		CONSTRAINT fk_factura_status           FOREIGN KEY (id_status)           REFERENCES estatus(id_status)
	);

	CREATE INDEX IF NOT EXISTS idx_factura_fac_numero      ON factura(fac_numero);
	CREATE INDEX IF NOT EXISTS idx_factura_id_user_pos     ON factura(id_user_pos);
	CREATE INDEX IF NOT EXISTS idx_factura_id_estacion     ON factura(id_estacion);
	CREATE INDEX IF NOT EXISTS idx_factura_id_orden_pedido ON factura(id_orden_pedido);
	CREATE INDEX IF NOT EXISTS idx_factura_id_cliente      ON factura(id_cliente);
	CREATE INDEX IF NOT EXISTS idx_factura_id_status       ON factura(id_status);
	CREATE INDEX IF NOT EXISTS idx_factura_deleted_at      ON factura(deleted_at);
	`
	_, err := db.Exec(query)
	return err
}
