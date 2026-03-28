package migrations

import "database/sql"

// migrateRefactorOrdenAgregador normaliza la relación entre orden_pedido y orden_agregador.
// 1. Agrega el total a la orden de pedido (Cabecera única de verdad).
// 2. Elimina redundancias en orden_agregador (fecha).
// 3. Agrega metadatos específicos (comisión).
func migrateRefactorOrdenAgregador(db *sql.DB) error {
	query := `
	-- 1. Actualizar orden_pedido
	ALTER TABLE orden_pedido ADD COLUMN IF NOT EXISTS odp_total DECIMAL(18,2) NOT NULL DEFAULT 0;
	COMMENT ON COLUMN orden_pedido.odp_total IS 'Monto total de la orden (Subtotal + Impuestos).';

	-- 2. Actualizar orden_agregador
	-- Eliminar fecha redundante (ya existe en orden_pedido.odp_fecha_creacion)
	ALTER TABLE orden_agregador DROP COLUMN IF EXISTS fecha;

	-- Agregar comisión específica del agregador
	ALTER TABLE orden_agregador ADD COLUMN IF NOT EXISTS comision_agregador DECIMAL(18,2) NOT NULL DEFAULT 0;
	COMMENT ON COLUMN orden_agregador.comision_agregador IS 'Monto de comisión cobrado por la plataforma externa.';

	-- 3. Documentación
	COMMENT ON TABLE orden_agregador IS 'Metadatos y referencias externas de pedidos provenientes de agregadores (UberEats, Rappi, etc).';
	`
	_, err := db.Exec(query)
	return err
}
