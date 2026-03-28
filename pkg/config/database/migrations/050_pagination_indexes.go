package migrations

import "database/sql"

func migratePaginationIndexes(db *sql.DB) error {
	query := `
	-- Índices para optimizar la paginación por cursor (Keyset Pagination)
	
	-- Factura: Optimiza el listado histórico de ventas
	CREATE INDEX IF NOT EXISTS idx_factura_created_at ON factura(created_at DESC);
	
	-- Producto: Optimiza el listado del catálogo de productos
	CREATE INDEX IF NOT EXISTS idx_producto_created_at ON producto(created_at DESC);
	
	-- Inventario: Optimiza el listado global de stock
	CREATE INDEX IF NOT EXISTS idx_inventario_created_at ON inventario(created_at DESC);
	
	-- Log Sistema: Asegurar índice para auditoría (aunque ya existía uno en fecha, created_at es más consistente)
	CREATE INDEX IF NOT EXISTS idx_log_sistema_created_at ON log_sistema(created_at DESC);
	`
	_, err := db.Exec(query)
	return err
}
