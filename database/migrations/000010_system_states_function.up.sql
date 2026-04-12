-- 1. FUNCIÓN MAESTRA DE ADMINISTRACIÓN DE ESTADOS
-- Refactor: Validaciones fail-fast, normalización de texto y logs JSON
CREATE OR REPLACE FUNCTION modulos_ia_estados(
    p_opcion        INTEGER,
    p_id_estado     UUID,
    p_descripcion   VARCHAR(100),
    p_factor        VARCHAR(10),
    p_nivel         INTEGER,
    p_id_modulo     INTEGER,
    p_id_cadena     INTEGER,    -- Reservado
    p_id_users_pos  VARCHAR(40) -- Auditoría
)
RETURNS TABLE (
    id_status       UUID,
    std_descripcion VARCHAR(255),
    factor          VARCHAR(10),
    nivel           INTEGER,
    mdl_id          INTEGER,
    std_tipo_estado VARCHAR(255),
    is_active       BOOLEAN
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_ip_origen     VARCHAR(45);
    v_error_number  TEXT;
    v_error_msg     TEXT;
    c_tipo_estado_default CONSTANT VARCHAR(10) := '1';
BEGIN
    -- Validación de entrada
    IF TRIM(p_descripcion) = '' OR p_descripcion IS NULL THEN
        RAISE EXCEPTION 'La descripción no puede estar vacía o contener solo espacios';
    END IF;

    p_descripcion := TRIM(p_descripcion);

    IF p_opcion NOT IN (0, 1) THEN
        RAISE EXCEPTION 'Opción inválida: %. Valores permitidos: 0 (actualizar), 1 (insertar)', p_opcion;
    END IF;

    -- ── Opción 0: Actualizar ──────────────────────────────────────────
    IF p_opcion = 0 THEN
        IF EXISTS (
            SELECT 1 FROM estatus
            WHERE std_descripcion = p_descripcion
              AND id_status      <> p_id_estado
              AND mdl_id          = p_id_modulo
              AND deleted_at     IS NULL
        ) THEN
            RAISE EXCEPTION 'Ya existe un estado con descripción "%" en el módulo %', p_descripcion, p_id_modulo;
        END IF;

        UPDATE estatus SET
            std_descripcion = p_descripcion,
            factor          = p_factor,
            nivel           = p_nivel,
            updated_at      = CURRENT_TIMESTAMP
        WHERE id_status  = p_id_estado
          AND deleted_at IS NULL;

        IF NOT FOUND THEN
            RAISE EXCEPTION 'Estado no encontrado o ya eliminado';
        END IF;

    -- ── Opción 1: Insertar ────────────────────────────────────────────
    ELSIF p_opcion = 1 THEN
        IF EXISTS (
            SELECT 1 FROM estatus
            WHERE std_descripcion = p_descripcion
              AND mdl_id          = p_id_modulo
              AND deleted_at     IS NULL
        ) THEN
            RAISE EXCEPTION 'Ya existe un estado con descripción "%" en el módulo %', p_descripcion, p_id_modulo;
        END IF;

        INSERT INTO estatus (
            std_descripcion,
            factor,
            nivel,
            mdl_id,
            std_tipo_estado,
            is_active,
            created_at,
            updated_at
        ) VALUES (
            p_descripcion,
            p_factor,
            p_nivel,
            p_id_modulo,
            c_tipo_estado_default,
            TRUE,
            CURRENT_TIMESTAMP,
            CURRENT_TIMESTAMP
        );
    END IF;

    -- Resultado unificado
    RETURN QUERY
        SELECT
            e.id_status,
            e.std_descripcion,
            e.factor,
            e.nivel,
            e.mdl_id,
            e.std_tipo_estado,
            e.is_active
        FROM estatus e
        WHERE e.mdl_id     = p_id_modulo
          AND e.deleted_at IS NULL
        ORDER BY e.nivel, e.std_descripcion;

EXCEPTION
    WHEN OTHERS THEN
        v_ip_origen    := COALESCE(inet_client_addr()::VARCHAR, 'local/socket');
        v_error_number := SQLSTATE;
        v_error_msg    := MESSAGE_TEXT;

        RAISE LOG 'modulos_ia_estados error: %', json_build_object(
            'fecha',     CURRENT_TIMESTAMP,
            'usuario',   p_id_users_pos,
            'ip',        v_ip_origen,
            'sqlstate',  v_error_number,
            'mensaje',   v_error_msg
        );

        RAISE EXCEPTION 'modulos_ia_estados falló [%]: %', v_error_number, v_error_msg;
END;
$$;

-- 2. ÍNDICE DE UNICIDAD POR MÓDULO (SQL Optimization)
CREATE UNIQUE INDEX IF NOT EXISTS idx_estatus_desc_modulo
    ON estatus (mdl_id, std_descripcion)
    WHERE deleted_at IS NULL;
