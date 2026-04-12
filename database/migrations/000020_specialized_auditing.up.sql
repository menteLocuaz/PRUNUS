-- 1. AUDITORÍA DE FACTURACIÓN (Trazabilidad Crítica)
CREATE TABLE IF NOT EXISTS factura_audit (
    id_audit        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_factura      UUID NOT NULL REFERENCES factura(id_factura),
    id_usuario      UUID REFERENCES usuario(id_usuario),
    accion          VARCHAR(50) NOT NULL, -- UPDATE, DELETE, ANULACION
    estado_anterior UUID REFERENCES estatus(id_status),
    estado_nuevo    UUID REFERENCES estatus(id_status),
    observaciones   TEXT,
    fecha           TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ip_address      VARCHAR(45)
);

CREATE INDEX IF NOT EXISTS idx_factura_audit_id_factura ON factura_audit(id_factura);

-- 2. HISTORIAL DE PRECIOS (Event Sourcing Pattern)
CREATE TABLE IF NOT EXISTS historial_precios (
    id_historial    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_producto     UUID NOT NULL REFERENCES producto(id_producto),
    id_sucursal     UUID NOT NULL REFERENCES sucursal(id_sucursal),
    precio_anterior NUMERIC(18,2) NOT NULL,
    precio_nuevo    NUMERIC(18,2) NOT NULL,
    tipo_precio     VARCHAR(20) NOT NULL, -- VENTA, COMPRA
    id_usuario      UUID REFERENCES usuario(id_usuario),
    fecha           TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_historial_precios_prod_suc ON historial_precios(id_producto, id_sucursal);

-- 3. TRIGGER AUTOMÁTICO PARA CAMBIOS DE PRECIOS
CREATE OR REPLACE FUNCTION fn_audit_precios()
RETURNS TRIGGER AS $$
DECLARE
    v_user_id UUID;
BEGIN
    v_user_id := NULLIF(current_setting('app.current_user_id', true), '')::UUID;

    IF (OLD.precio_venta IS DISTINCT FROM NEW.precio_venta) THEN
        INSERT INTO historial_precios (id_producto, id_sucursal, precio_anterior, precio_nuevo, tipo_precio, id_usuario)
        VALUES (NEW.id_producto, NEW.id_sucursal, OLD.precio_venta, NEW.precio_venta, 'VENTA', v_user_id);
    END IF;

    IF (OLD.precio_compra IS DISTINCT FROM NEW.precio_compra) THEN
        INSERT INTO historial_precios (id_producto, id_sucursal, precio_anterior, precio_nuevo, tipo_precio, id_usuario)
        VALUES (NEW.id_producto, NEW.id_sucursal, OLD.precio_compra, NEW.precio_compra, 'COMPRA', v_user_id);
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS tr_audit_precios ON inventario;
CREATE TRIGGER tr_audit_precios
AFTER UPDATE OF precio_venta, precio_compra ON inventario
FOR EACH ROW EXECUTE FUNCTION fn_audit_precios();

-- 4. APLICAR FRAMEWORK CORE
CALL sp_core_setup_table('factura_audit', '{"audit": false}');
CALL sp_core_setup_table('historial_precios', '{"audit": false}');
