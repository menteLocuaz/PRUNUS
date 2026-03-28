package migrations

import "database/sql"

// migrateStatusConsistencyTriggers establece validaciones automáticas para asegurar que
// cada tabla use únicamente los estatus que corresponden a su módulo de negocio.
func migrateStatusConsistencyTriggers(db *sql.DB) error {
	query := `
	-- 1. Asegurar esquema de configuración para metadatos de base de datos
	CREATE SCHEMA IF NOT EXISTS config;

	-- 2. Tabla de configuración dinámica de módulos por tabla
	CREATE TABLE IF NOT EXISTS config.trigger_module_map (
		table_name  TEXT    NOT NULL,
		mdl_id      INTEGER NOT NULL,
		descripcion TEXT,                          
		activo      BOOLEAN NOT NULL DEFAULT TRUE, 
		lastupdate  TIMESTAMP DEFAULT NOW(),
		lastuser    TEXT,
		CONSTRAINT pk_trigger_module_map PRIMARY KEY (table_name)
	);

	COMMENT ON TABLE config.trigger_module_map IS 'Mapa dinámico de módulos por tabla para validación de estatus.';

	-- 3. Sincronización de IDs según pkg/config/database/migrations/012_estatus.go
	-- Mapeo consolidado: 1:Empresa, 2:Sucursal, 3:Usuario, 4:Producto, 5:Venta, 6:Compra, 7:Finanzas, 8:POS
	INSERT INTO config.trigger_module_map (table_name, mdl_id, descripcion, lastuser)
	VALUES
		('empresa',      1, 'Módulo Estructura/Empresa', 'system'),
		('sucursal',     2, 'Módulo Estructura/Sedes',   'system'),
		('usuario',      3, 'Módulo Seguridad/Usuarios', 'system'),
		('producto',     4, 'Módulo Catálogo/Productos', 'system'),
		('factura',      5, 'Módulo Ventas/Facturas',    'system'),
		('orden_pedido', 6, 'Módulo Operaciones/Pedidos','system'),
		('categoria',    4, 'Módulo Catálogo (Asociado a Producto)', 'system'),
		('moneda',      -1, 'Módulo Global/Transversal', 'system'),
		('medida',      -1, 'Módulo Global/Transversal', 'system')
	ON CONFLICT (table_name) DO UPDATE 
	SET mdl_id = EXCLUDED.mdl_id, 
	    descripcion = EXCLUDED.descripcion,
	    lastupdate = NOW();

	-- 4. Función genérica de validación
	CREATE OR REPLACE FUNCTION config.fn_validate_status_module()
	RETURNS TRIGGER AS $$
	DECLARE
		v_expected_mdl_id INTEGER;
		v_actual_mdl_id   INTEGER;
		v_table_name      TEXT;
		v_activo          BOOLEAN;
	BEGIN
		v_table_name := TG_TABLE_NAME;

		-- Buscar configuración para la tabla actual
		SELECT mdl_id, activo INTO v_expected_mdl_id, v_activo
		FROM config.trigger_module_map
		WHERE table_name = v_table_name;

		-- Si la tabla no está mapeada o la validación está desactivada, permitir
		IF NOT FOUND OR v_activo = FALSE THEN
			RETURN NEW;
		END IF;

		-- Obtener el módulo asignado al estatus en la tabla maestra
		SELECT mdl_id INTO v_actual_mdl_id
		FROM public.estatus
		WHERE id_status = NEW.id_status;

		IF NOT FOUND THEN
			RAISE EXCEPTION '[Integridad] El estatus con ID "%" no existe en la tabla maestra.', NEW.id_status;
		END IF;

		-- mdl_id = -1 es un comodín para estatus globales/transversales
		IF v_actual_mdl_id = -1 THEN
			RETURN NEW;
		END IF;

		-- Validar coherencia: El estatus debe pertenecer al módulo de la tabla
		IF v_actual_mdl_id != v_expected_mdl_id THEN
			RAISE EXCEPTION '[Incoherencia de Estatus] La tabla "%" (Módulo %) no puede usar el estatus ID "%" que pertenece al Módulo %.',
				v_table_name, v_expected_mdl_id, NEW.id_status, v_actual_mdl_id;
		END IF;

		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

	-- 5. Procedimiento Helper para aplicar triggers de forma masiva
	CREATE OR REPLACE PROCEDURE config.sp_aplicar_trigger_validacion(p_table_name TEXT)
	LANGUAGE plpgsql AS $$
	DECLARE
		v_trigger_name TEXT;
	BEGIN
		v_trigger_name := 'tr_validate_status_' || p_table_name;
		EXECUTE format(
			'DROP TRIGGER IF EXISTS %I ON %I;
			 CREATE TRIGGER %I
			 BEFORE INSERT OR UPDATE OF id_status ON %I
			 FOR EACH ROW EXECUTE FUNCTION config.fn_validate_status_module();',
			v_trigger_name, p_table_name, v_trigger_name, p_table_name
		);
	END;
	$$;

	-- 6. Ejecución de la aplicación de triggers
	CALL config.sp_aplicar_trigger_validacion('empresa');
	CALL config.sp_aplicar_trigger_validacion('sucursal');
	CALL config.sp_aplicar_trigger_validacion('usuario');
	CALL config.sp_aplicar_trigger_validacion('producto');
	CALL config.sp_aplicar_trigger_validacion('factura');
	CALL config.sp_aplicar_trigger_validacion('orden_pedido');
	CALL config.sp_aplicar_trigger_validacion('categoria');
	CALL config.sp_aplicar_trigger_validacion('moneda');
	CALL config.sp_aplicar_trigger_validacion('medida');
	`
	_, err := db.Exec(query)
	return err
}
