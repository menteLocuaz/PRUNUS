package migrations

import (
	"database/sql"
	"fmt"
)

func RunMigrations(db *sql.DB) error {
	steps := []struct {
		name string
		fn   func(*sql.DB) error
	}{
		{"012_estatus", migrateEstatus},
		{"001_empresa", migrateEmpresa},
		{"003_sucursal", migrateSucursal},
		{"002_rol", migrateRol},
		{"011_modulos", migrateModulos},
		{"004_usuario", migrateUsuario},
		{"005_categoria", migrateCategoria},
		{"006_unidad", migrateUnidad},
		{"007_moneda", migrateMoneda},
		{"008_cliente", migrateCliente},
		{"009_proveedor", migrateProveedor},
		{"010_producto", migrateProducto},
		{"013_forma_pago", migrateFormaPago},
		{"014_periodo", migratePeriodo},
		{"015_estaciones_pos", migrateEstacionesPos},
		{"016_control_estacion", migrateControlEstacion},
		{"017_retiros", migrateRetiros},
		{"018_orden_pedido", migrateOrdenPedido},
		{"019_factura", migrateFactura},
		{"020_cabecera_motivo_anulacion", migrateCabeceraMotivoAnulacion},
		{"021_motivo_anulacion", migrateMotivoAnulacion},
		{"022_impuesto", migrateImpuesto},
		{"023_inventario", migrateInventario},
		{"024_movimientos_inventario", migrateMovimientosInventario},
		{"025_detalle_factura", migrateDetalleFactura},
		{"026_forma_pago_factura", migrateFormaPagoFactura},
		{"027_log_sistema", migrateLogSistema},
		{"028_auditoria_caja", migrateAuditoriaCaja},
		{"029_agregadores", migrateAgregadores},
		{"030_orden_agregador", migrateOrdenAgregador},
		{"031_dispositivos_pos", migrateDispositivosPos},
		{"032_update_movimientos_inventario", migrateUpdateMovimientosInventario},
		{"033_orden_compra", migrateOrdenCompra},
		{"034_trigger_updated_at", migrateTriggerUpdatedAt},
		{"035_trigger_stock_sync", migrateTriggerStockSync},
		{"036_trigger_venta_movimiento", migrateTriggerVentaMovimiento},
		{"037_trigger_anulacion_factura", migrateTriggerAnulacionFactura},
		{"038_fn_modulos", migrateFnModulos},
		{"039_fn_estados", migrateFnEstados},
		{"040_fn_estaciones", migrateFnEstaciones},
		{"041_fn_inventario", migrateFnInventario},
		{"042_fn_facturacion", migrateFnFacturacion},
		{"048_usuario_sucursal_acceso", migrateUsuarioSucursalAcceso},
		{"049_performance_optimizations", migratePerformanceOptimizations},
		{"050_pagination_indexes", migratePaginationIndexes},
		{"051_normalize_inventory", migrateNormalizeInventory},
		{"052_status_consistency_triggers", migrateStatusConsistencyTriggers},
		{"053_refactor_orden_agregador", migrateRefactorOrdenAgregador},
		{"054_specialized_auditing", migrateSpecializedAuditing},
	}

	for _, s := range steps {
		if err := s.fn(db); err != nil {
			return fmt.Errorf("migración %s: %w", s.name, err)
		}
	}

	return nil
}
