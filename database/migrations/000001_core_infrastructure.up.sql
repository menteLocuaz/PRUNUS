-- 1. EXTENSIONES Y SEGURIDAD
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";
CREATE SCHEMA IF NOT EXISTS core;

-- 2. FUNCIÓN DE ACTUALIZACIÓN DE TIEMPO
CREATE OR REPLACE FUNCTION fn_update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 3. AUDITORÍA MAESTRA AVANZADA
CREATE TABLE IF NOT EXISTS auditoria_maestra (
    id_auditoria    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tabla_nombre    VARCHAR(100) NOT NULL,
    registro_id     UUID NOT NULL,
    accion          VARCHAR(20) NOT NULL,
    old_data        JSONB,
    new_data        JSONB,
    id_usuario      UUID,
    fecha           TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ip_address      VARCHAR(45)
);

CREATE INDEX IF NOT EXISTS idx_audit_maestra_tabla_fecha ON auditoria_maestra(tabla_nombre, fecha DESC);
CREATE INDEX IF NOT EXISTS idx_audit_maestra_registro ON auditoria_maestra(registro_id);

CREATE OR REPLACE FUNCTION fn_audit_generic()
RETURNS TRIGGER AS $$
DECLARE
    v_old_json JSONB := NULL;
    v_new_json JSONB := NULL;
    v_user_id  UUID  := NULLIF(current_setting('app.current_user_id', true), '')::UUID;
    v_ip       TEXT  := COALESCE(NULLIF(current_setting('app.ip_address', true), ''), inet_client_addr()::TEXT, '127.0.0.1');
    v_pk_col   TEXT;
    v_pk_val   UUID;
    v_exclude  TEXT[] := ARRAY['password', 'pin', 'secret_key', 'usu_pin_pos'];
BEGIN
    SELECT a.attname INTO v_pk_col FROM pg_index i 
    JOIN pg_attribute a ON a.attrelid = i.indrelid AND a.attnum = ANY(i.indkey)
    WHERE i.indrelid = TG_RELID AND i.indisprimary;

    IF (TG_OP IN ('UPDATE', 'DELETE')) THEN
        v_old_json := to_jsonb(OLD);
        IF v_old_json IS NOT NULL THEN
            FOR i IN 1..array_length(v_exclude, 1) LOOP v_old_json := v_old_json - v_exclude[i]; END LOOP;
            v_pk_val := (v_old_json->>v_pk_col)::UUID;
        END IF;
    END IF;
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        v_new_json := to_jsonb(NEW);
        IF v_new_json IS NOT NULL THEN
            FOR i IN 1..array_length(v_exclude, 1) LOOP v_new_json := v_new_json - v_exclude[i]; END LOOP;
            v_pk_val := (v_new_json->>v_pk_col)::UUID;
        END IF;
    END IF;

    IF (TG_OP = 'UPDATE' AND v_old_json = v_new_json) THEN RETURN NEW; END IF;

    INSERT INTO auditoria_maestra (tabla_nombre, registro_id, accion, old_data, new_data, id_usuario, ip_address)
    VALUES (TG_TABLE_NAME, COALESCE(v_pk_val, gen_random_uuid()), TG_OP, v_old_json, v_new_json, v_user_id, v_ip);

    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;

-- 4. PROCEDIMIENTO DE CONFIGURACIÓN MAESTRA
CREATE OR REPLACE PROCEDURE sp_core_setup_table(p_table TEXT, p_options JSONB DEFAULT '{}')
LANGUAGE plpgsql AS $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = p_table AND column_name = 'updated_at') THEN
        EXECUTE format('DROP TRIGGER IF EXISTS tr_upd_%I ON %I', p_table, p_table);
        EXECUTE format('CREATE TRIGGER tr_upd_%I BEFORE UPDATE ON %I FOR EACH ROW EXECUTE FUNCTION fn_update_updated_at_column()', p_table, p_table);
    END IF;
    IF COALESCE((p_options->>'audit')::BOOLEAN, TRUE) THEN
        EXECUTE format('DROP TRIGGER IF EXISTS tr_aud_%I ON %I', p_table, p_table);
        EXECUTE format('CREATE TRIGGER tr_aud_%I AFTER INSERT OR UPDATE OR DELETE ON %I FOR EACH ROW EXECUTE FUNCTION fn_audit_generic()', p_table, p_table);
    END IF;
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = p_table AND column_name = 'deleted_at') THEN
        EXECUTE format('CREATE INDEX IF NOT EXISTS %I ON %I (created_at DESC) WHERE deleted_at IS NULL', 'idx_active_' || p_table, p_table);
    END IF;
END;
$$;
