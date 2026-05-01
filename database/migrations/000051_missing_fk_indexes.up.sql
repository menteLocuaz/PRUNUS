-- Índices sobre FK faltantes en tablas de catálogo y transaccionales.
-- Se omite CONCURRENTLY para compatibilidad con golang-migrate (requiere transacción).
-- Para entornos productivos con tablas grandes, ejecutar manualmente con CONCURRENTLY.

-- categoria
CREATE INDEX IF NOT EXISTS idx_categoria_status
    ON categoria(id_status) WHERE deleted_at IS NULL;

-- moneda
CREATE INDEX IF NOT EXISTS idx_moneda_sucursal
    ON moneda(id_sucursal) WHERE deleted_at IS NULL;

-- medida / unidad
CREATE INDEX IF NOT EXISTS idx_medida_sucursal
    ON unidad(id_sucursal) WHERE deleted_at IS NULL;

-- estaciones_pos
CREATE INDEX IF NOT EXISTS idx_estaciones_pos_sucursal
    ON estaciones_pos(id_sucursal) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_estaciones_pos_status
    ON estaciones_pos(id_status) WHERE deleted_at IS NULL;

-- clientes
CREATE INDEX IF NOT EXISTS idx_clientes_status
    ON cliente(id_status) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_clientes_sucursal
    ON cliente(id_sucursal) WHERE deleted_at IS NULL;

-- proveedores
CREATE INDEX IF NOT EXISTS idx_proveedores_status
    ON proveedor(id_status) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_proveedores_sucursal
    ON proveedor(id_sucursal) WHERE deleted_at IS NULL;

-- detalle_factura — clave de join para cada lectura de factura
CREATE INDEX IF NOT EXISTS idx_detalle_factura_factura
    ON detalle_factura(id_factura);

-- forma_pago_factura — clave de join para pagos de factura
CREATE INDEX IF NOT EXISTS idx_forma_pago_factura_factura
    ON forma_pago_factura(id_factura);

-- movimientos_inventario — tabla de alta cardinalidad usada en reportes analíticos
CREATE INDEX IF NOT EXISTS idx_movimientos_inv_sucursal
    ON movimientos_inventario(id_sucursal, fecha DESC) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_movimientos_inv_tipo
    ON movimientos_inventario(tipo_movimiento, id_sucursal) WHERE deleted_at IS NULL;
