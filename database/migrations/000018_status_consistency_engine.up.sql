-- 1. ESQUEMA DE CONFIGURACIÓN DINÁMICA
CREATE SCHEMA IF NOT EXISTS config;

-- 2. MAPA DE MÓDULOS POR TABLA
CREATE TABLE IF NOT EXISTS config.trigger_module_map (
    table_name  TEXT PRIMARY KEY,
    mdl_id      INTEGER NOT NULL,
    descripcion TEXT,
    activo      BOOLEAN NOT NULL DEFAULT TRUE,
    lastupdate  TIMESTAMPTZ DEFAULT NOW()
);

-- 3. POBLACIÓN DE REGLAS DE CONSISTENCIA
INSERT INTO config.trigger_module_map (table_name, mdl_id, descripcion)
VALUES
    ('empresa',      1, 'Módulo Estructura/Empresa'),
    ('sucursal',     2, 'Módulo Estructura/Sedes'),
    ('usuario',      3, 'Módulo Seguridad/Usuarios'),
    ('producto',     4, 'Módulo Catálogo/Productos'),
    ('factura',      5, 'Módulo Ventas/Facturas'),
    ('orden_pedido', 6, 'Módulo Operaciones/Pedidos'),
    ('moneda',      -1, 'Módulo Global/Transversal')
ON CONFLICT (table_name) DO UPDATE 
SET mdl_id = EXCLUDED.mdl_id, descripcion = EXCLUDED.descripcion, lastupdate = NOW();

-- 4. FUNCIÓN GENÉRICA DE VALIDACIÓN DE ESTATUS
CREATE OR REPLACE FUNCTION config.fn_validate_status_module()
RETURNS TRIGGER AS $$
DECLARE
    v_expected_mdl_id INTEGER;
    v_actual_mdl_id   INTEGER;
    v_activo          BOOLEAN;
BEGIN
    SELECT mdl_id, activo INTO v_expected_mdl_id, v_activo
    FROM config.trigger_module_map
    WHERE table_name = TG_TABLE_NAME;

    IF NOT FOUND OR v_activo = FALSE THEN
        RETURN NEW;
    END IF;

    SELECT mdl_id INTO v_actual_mdl_id
    FROM public.estatus
    WHERE id_status = NEW.id_status;

    IF NOT FOUND THEN
        RAISE EXCEPTION '[Integridad] El estatus ID "%" no existe en la tabla maestra.', NEW.id_status;
    END IF;

    IF v_actual_mdl_id = -1 THEN
        RETURN NEW;
    END IF;

    IF v_actual_mdl_id != v_expected_mdl_id THEN
        RAISE EXCEPTION '[Incoherencia de Estatus] La tabla "%" (Módulo %) no puede usar el estatus ID "%" que pertenece al Módulo %.',
            TG_TABLE_NAME, v_expected_mdl_id, NEW.id_status, v_actual_mdl_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 5. PROCEDIMIENTO PARA APLICAR EL MOTOR DE REGLAS
CREATE OR REPLACE PROCEDURE config.sp_apply_status_rule(p_table_name TEXT)
LANGUAGE plpgsql AS $$
DECLARE
    v_trigger_name TEXT;
BEGIN
    v_trigger_name := 'tr_status_check_' || p_table_name;
    EXECUTE format(
        'DROP TRIGGER IF EXISTS %I ON %I;
         CREATE TRIGGER %I
         BEFORE INSERT OR UPDATE OF id_status ON %I
         FOR EACH ROW EXECUTE FUNCTION config.fn_validate_status_module();',
        v_trigger_name, p_table_name, v_trigger_name, p_table_name
    );
END;
$$;

-- 6. APLICACIÓN INICIAL
DO $$ 
DECLARE 
    r RECORD;
BEGIN
    FOR r IN SELECT table_name FROM config.trigger_module_map LOOP
        -- Verificar existencia de la tabla antes de intentar aplicar el trigger
        IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = r.table_name) THEN
            CALL config.sp_apply_status_rule(r.table_name);
        END IF;
    END LOOP;
END $$;
