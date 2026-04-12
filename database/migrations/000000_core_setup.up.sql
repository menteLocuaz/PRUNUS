-- 1. Esquema de control para lógica interna
CREATE SCHEMA IF NOT EXISTS core;

-- 2. Función genérica para actualizar automáticamente el campo updated_at
-- Evita tener que escribir la lógica de tiempo en cada UPDATE desde la aplicación.
CREATE OR REPLACE FUNCTION fn_update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 3. Auditoría Avanzada con Soporte JSONB (Skill: Refactor & Security)
-- Captura cambios detallados, usuario que realiza la acción e IP, enmascarando campos sensibles.
CREATE OR REPLACE FUNCTION fn_audit_generic()
RETURNS TRIGGER AS $$
DECLARE
    v_user_id UUID := NULLIF(current_setting('app.current_user_id', true), '')::UUID;
    v_ip TEXT := COALESCE(NULLIF(current_setting('app.ip_address', true), ''), inet_client_addr()::TEXT);
    v_pk_col TEXT;
    v_pk_val TEXT;
    v_old_json JSONB := NULL;
    v_new_json JSONB := NULL;
BEGIN
    -- Identificar la Columna de Clave Primaria dinámicamente
    SELECT a.attname INTO v_pk_col 
    FROM pg_index i
    JOIN pg_attribute a ON a.attrelid = i.indrelid AND a.attnum = ANY(i.indkey)
    WHERE i.indrelid = TG_RELID AND i.indisprimary;

    -- Extraer el valor de la PK para el registro de auditoría
    IF (TG_OP = 'DELETE') THEN
        EXECUTE format('SELECT ($1).%I::text', v_pk_col) INTO v_pk_val USING OLD;
        v_old_json := to_jsonb(OLD);
    ELSIF (TG_OP = 'UPDATE') THEN
        EXECUTE format('SELECT ($1).%I::text', v_pk_col) INTO v_pk_val USING NEW;
        v_old_json := to_jsonb(OLD);
        v_new_json := to_jsonb(NEW);
    ELSIF (TG_OP = 'INSERT') THEN
        EXECUTE format('SELECT ($1).%I::text', v_pk_col) INTO v_pk_val USING NEW;
        v_new_json := to_jsonb(NEW);
    END IF;

    -- Enmascaramiento de campos sensibles (Security First)
    IF v_old_json ? 'password' THEN v_old_json := v_old_json || '{"password": "[REDACTED]"}'; END IF;
    IF v_new_json ? 'password' THEN v_new_json := v_new_json || '{"password": "[REDACTED]"}'; END IF;

    -- Insertar en la tabla de auditoría (se asume creada en migraciones posteriores o existente)
    -- Si la tabla no existe aún, este trigger fallará en el primer registro, pero asegura consistencia.
    INSERT INTO auditoria_maestra (tabla_nombre, registro_id, accion, old_data, new_data, id_usuario, ip_address)
    VALUES (TG_TABLE_NAME, v_pk_val, TG_OP, v_old_json, v_new_json, v_user_id, v_ip);

    RETURN COALESCE(NEW, OLD);
EXCEPTION WHEN OTHERS THEN
    -- En auditoría, preferimos que la operación principal continúe si falla el log, 
    -- o podrías cambiarlo para que falle si el cumplimiento es estricto.
    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;
