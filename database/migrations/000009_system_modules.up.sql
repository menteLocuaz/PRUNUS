-- 1. TABLA DE MÓDULOS DEL SISTEMA (Consolidada y Normalizada)
CREATE TABLE IF NOT EXISTS modulo (
    id_modulo       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    mdl_id          SERIAL NOT NULL, -- ID numérico para legibilidad/orden
    mdl_descripcion VARCHAR(150) NOT NULL,
    abreviatura     VARCHAR(50),
    ruta            VARCHAR(255),
    icono           VARCHAR(100),
    nivel           INTEGER NOT NULL DEFAULT 0,
    orden           INTEGER NOT NULL DEFAULT 0,
    id_padre        UUID REFERENCES modulo(id_modulo),
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    id_status       UUID REFERENCES estatus(id_status),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMPTZ NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_modulo_mdl_id ON modulo(mdl_id) WHERE deleted_at IS NULL;

-- 2. FUNCIÓN DE ADMINISTRACIÓN DE MÓDULOS (Refactored)
-- Opción 0 → Actualizar, Opción 1 → Insertar
CREATE OR REPLACE FUNCTION modulos_ia_modulo(
    p_opcion        INTEGER,
    p_id_modulo     INTEGER,        -- mdl_id
    p_descripcion   VARCHAR(100),
    p_abreviatura   VARCHAR(50),
    p_nivel         INTEGER,
    p_estado        BOOLEAN,
    p_id_cadena     INTEGER,        -- Reservado
    p_id_users_pos  VARCHAR(40)     -- Auditoría
)
RETURNS TABLE (
    id_modulo       UUID,
    mdl_id          INTEGER,
    mdl_descripcion VARCHAR(100),
    abreviatura     VARCHAR(50),
    nivel           INTEGER,
    is_active       BOOLEAN,
    id_status       UUID
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_next_id       INTEGER;
    v_ip_origen     VARCHAR(45);
    v_error_number  TEXT;
    v_error_msg     TEXT;
BEGIN
    -- Validación de duplicados
    IF EXISTS (
        SELECT 1 FROM modulo 
        WHERE mdl_descripcion = TRIM(p_descripcion) 
          AND mdl_id <> p_id_modulo 
          AND deleted_at IS NULL
    ) THEN
        RAISE EXCEPTION 'Ya existe un módulo con la descripción "%"', p_descripcion;
    END IF;

    -- ── Opción 0: Actualizar ──────────────────────────────────────────
    IF p_opcion = 0 THEN
        UPDATE modulo SET
            mdl_descripcion = TRIM(p_descripcion),
            abreviatura     = p_abreviatura,
            nivel           = p_nivel,
            is_active       = p_estado,
            updated_at      = CURRENT_TIMESTAMP
        WHERE mdl_id = p_id_modulo;

        IF NOT FOUND THEN
            RAISE EXCEPTION 'Módulo con ID % no encontrado', p_id_modulo;
        END IF;

    -- ── Opción 1: Insertar ────────────────────────────────────────────
    ELSIF p_opcion = 1 THEN
        INSERT INTO modulo (
            mdl_descripcion,
            abreviatura,
            nivel,
            is_active,
            updated_at
        )
        VALUES (
            TRIM(p_descripcion),
            p_abreviatura,
            p_nivel,
            p_estado,
            CURRENT_TIMESTAMP
        )
        RETURNING modulo.mdl_id INTO v_next_id;

        -- Sincronizar secuencia si es necesario
        PERFORM setval(pg_get_serial_sequence('modulo', 'mdl_id'), (SELECT MAX(mdl_id) FROM modulo));

    ELSE
        RAISE EXCEPTION 'Opción % no soportada. Use 0 para actualizar o 1 para insertar.', p_opcion;
    END IF;

    -- Retornar lista actualizada
    RETURN QUERY
        SELECT
            m.id_modulo,
            m.mdl_id,
            m.mdl_descripcion,
            m.abreviatura,
            m.nivel,
            m.is_active,
            m.id_status
        FROM modulo m
        WHERE m.deleted_at IS NULL
        ORDER BY m.mdl_id;

EXCEPTION
    WHEN OTHERS THEN
        GET STACKED DIAGNOSTICS
            v_error_msg    = MESSAGE_TEXT,
            v_error_number = RETURNED_SQLSTATE;

        v_ip_origen := COALESCE(inet_client_addr()::VARCHAR, 'local/socket');

        RAISE LOG 'modulos_ia_modulo error: %', json_build_object(
            'fecha',      CURRENT_TIMESTAMP,
            'usuario',    p_id_users_pos,
            'ip',         v_ip_origen,
            'sqlstate',   v_error_number,
            'mensaje',    v_error_msg
        );

        RAISE EXCEPTION 'modulos_ia_modulo falló [%]: %', v_error_number, v_error_msg;
END;
$$;

-- 3. APLICAR FRAMEWORK CORE
CALL sp_core_setup_table('modulo');
