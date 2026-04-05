package migrations

import "database/sql"

func migrateInventoryPerformanceIndexes(db *sql.DB) error {
	query := `
	-- Índice para optimizar alertas de stock (productos bajo stock mínimo)
	CREATE INDEX IF NOT EXISTS idx_inventario_alertas ON inventario(id_sucursal, stock_actual, stock_minimo) WHERE deleted_at IS NULL AND stock_actual <= stock_minimo;
	`
	_, err := db.Exec(query)
	return err
}
