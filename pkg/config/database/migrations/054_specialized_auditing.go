package migrations

import "database/sql"

// migrateSpecializedAuditing implementa tablas de auditoría dedicadas para reducir la carga en log_sistema.
// Sigue el patrón Event Sourcing para cambios en precios y trazabilidad de facturación.
func migrateSpecializedAuditing(db *sql.DB) error {
	query := `
	-- 1. Auditoría de Facturación
	CREATE TABLE IF NOT EXISTS factura_audit (
		id_audit        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		id_factura      UUID         NOT NULL,
		id_usuario      UUID,         -- Quién realizó el cambio
		accion          VARCHAR(50)  NOT NULL, -- UPDATE, DELETE, ANULACION
		estado_anterior UUID,
		estado_nuevo    UUID,
		observaciones   TEXT,
		fecha           TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
		ip_address      VARCHAR(45),
		
		CONSTRAINT fk_audit_factura FOREIGN KEY (id_factura) REFERENCES factura(id_factura)
	);

	-- 2. Historial de Precios (Event Sourcing)
	CREATE TABLE IF NOT EXISTS historial_precios (
		id_historial    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		id_producto     UUID           NOT NULL,
		id_sucursal     UUID           NOT NULL,
		precio_anterior DECIMAL(18,2) NOT NULL,
		precio_nuevo    DECIMAL(18,2) NOT NULL,
		tipo_precio     VARCHAR(20)    NOT NULL, -- VENTA, COMPRA
		id_usuario      UUID,
		fecha           TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
		
		CONSTRAINT fk_historial_producto FOREIGN KEY (id_producto) REFERENCES producto(id_producto),
		CONSTRAINT fk_historial_sucursal FOREIGN KEY (id_sucursal) REFERENCES sucursal(id_sucursal)
	);

	-- 3. Triggers para Automatización
	
	-- Función para auditar precios
	CREATE OR REPLACE FUNCTION fn_audit_precios()
	RETURNS TRIGGER AS $$
	DECLARE
		v_user_id UUID;
	BEGIN
		-- Obtener el usuario del contexto de sesión de PG (seteado por el app)
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

	-- 4. Índices de Rendimiento
	CREATE INDEX IF NOT EXISTS idx_factura_audit_id_factura ON factura_audit(id_factura);
	CREATE INDEX IF NOT EXISTS idx_historial_precios_producto ON historial_precios(id_producto, id_sucursal);
	`
	_, err := db.Exec(query)
	return err
}
