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
	}

	for _, s := range steps {
		if err := s.fn(db); err != nil {
			return fmt.Errorf("migración %s: %w", s.name, err)
		}
	}

	return nil
}
