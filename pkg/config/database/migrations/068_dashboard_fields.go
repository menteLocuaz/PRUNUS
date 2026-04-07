package migrations

import "database/sql"

func migrateDashboardFields(db *sql.DB) error {
	query := `
	-- Agregar fecha de vencimiento a Facturas para Cuentas por Cobrar
	ALTER TABLE factura 
	ADD COLUMN IF NOT EXISTS fecha_vencimiento TIMESTAMP;

	-- Agregar fecha de vencimiento a OrdenCompra para Cuentas por Pagar
	ALTER TABLE orden_compra 
	ADD COLUMN IF NOT EXISTS fecha_vencimiento TIMESTAMP;

	-- Asegurar que existen índices para reportes de antigüedad de deuda
	CREATE INDEX IF NOT EXISTS idx_factura_vencimiento ON factura(fecha_vencimiento) 
	WHERE id_status = '892340e0-4328-491d-9102-80550bb6aac4'; -- Solo pendientes de pago

	CREATE INDEX IF NOT EXISTS idx_orden_compra_vencimiento ON orden_compra(fecha_vencimiento);

	-- Optimizar reportes de rentabilidad (Pareto)
	CREATE INDEX IF NOT EXISTS idx_movimientos_fecha_tipo ON movimientos_inventario(fecha DESC, tipo_movimiento)
	WHERE deleted_at IS NULL;
	`
	_, err := db.Exec(query)
	return err
}
