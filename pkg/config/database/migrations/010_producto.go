package migrations

import "database/sql"

// migrateProducto crea la tabla maestra de productos.
// Consolidado para la versión normalizada: El stock y precios residen en 'inventario'.
func migrateProducto(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS producto (
		id_producto       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		nombre            VARCHAR(150)   NOT NULL,
		descripcion       TEXT,
		fecha_vencimiento DATE,
		imagen            TEXT,
		id_status         UUID           NOT NULL,
		id_categoria      UUID           NOT NULL,
		id_moneda         UUID           NOT NULL,
		id_unidad         UUID           NOT NULL,

		created_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at        TIMESTAMP NULL,

		CONSTRAINT fk_producto_status
			FOREIGN KEY (id_status)
			REFERENCES estatus(id_status),

		CONSTRAINT fk_producto_categoria
			FOREIGN KEY (id_categoria)
			REFERENCES categoria(id_categoria)
			ON UPDATE CASCADE
			ON DELETE RESTRICT,

		CONSTRAINT fk_producto_moneda
			FOREIGN KEY (id_moneda)
			REFERENCES moneda(id_moneda)
			ON UPDATE CASCADE
			ON DELETE RESTRICT,

		CONSTRAINT fk_producto_unidad
			FOREIGN KEY (id_unidad)
			REFERENCES unidad(id_unidad)
			ON UPDATE CASCADE
			ON DELETE RESTRICT
	);

	-- Índices de rendimiento
	CREATE INDEX IF NOT EXISTS idx_producto_id_categoria ON producto(id_categoria);
	CREATE INDEX IF NOT EXISTS idx_producto_id_moneda    ON producto(id_moneda);
	CREATE INDEX IF NOT EXISTS idx_producto_id_unidad    ON producto(id_unidad);
	CREATE INDEX IF NOT EXISTS idx_producto_status       ON producto(id_status);
	
	-- Índice parcial para optimizar búsquedas de productos activos (Soft Delete)
	CREATE INDEX IF NOT EXISTS idx_producto_active ON producto(id_producto) WHERE deleted_at IS NULL;
	
	-- Índice compuesto para filtros comunes
	CREATE INDEX IF NOT EXISTS idx_producto_cat_active ON producto(id_categoria) WHERE deleted_at IS NULL;
	
	-- Índice para búsquedas por nombre
	CREATE INDEX IF NOT EXISTS idx_producto_nombre ON producto(nombre);

	COMMENT ON TABLE producto IS 'Catálogo maestro de productos (Atómico/Global)';
	`
	_, err := db.Exec(query)
	return err
}
