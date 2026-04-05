package migrations

import "database/sql"

func migrateMasterAuditing(db *sql.DB) error {
	query := `
	-- 1. Tabla de Auditoría General (JSONB)
	CREATE TABLE IF NOT EXISTS auditoria_maestra (
		id_auditoria    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		tabla_nombre    VARCHAR(100) NOT NULL,
		registro_id     UUID NOT NULL,
		accion          VARCHAR(20) NOT NULL, -- INSERT, UPDATE, DELETE
		old_data        JSONB,
		new_data        JSONB,
		id_usuario      UUID,
		fecha           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		ip_address      VARCHAR(45)
	);

	CREATE INDEX IF NOT EXISTS idx_audit_maestra_tabla ON auditoria_maestra(tabla_nombre);
	CREATE INDEX IF NOT EXISTS idx_audit_maestra_registro ON auditoria_maestra(registro_id);
	CREATE INDEX IF NOT EXISTS idx_audit_maestra_fecha ON auditoria_maestra(fecha);

	-- 2. Función de Auditoría Genérica
	CREATE OR REPLACE FUNCTION fn_audit_generic()
	RETURNS TRIGGER AS $$
	DECLARE
		v_old_json JSONB := NULL;
		v_new_json JSONB := NULL;
		v_user_id  UUID;
		v_pk_id    UUID;
	BEGIN
		-- Intentar obtener el usuario del contexto de sesión
		v_user_id := NULLIF(current_setting('app.current_user_id', true), '')::UUID;

		-- Determinar el ID del registro (Asumimos que la PK se llama id_... o usamos row_to_json)
		IF (TG_OP = 'DELETE' OR TG_OP = 'UPDATE') THEN
			v_old_json := to_jsonb(OLD);
		END IF;
		
		IF (TG_OP = 'INSERT' OR TG_OP = 'UPDATE') THEN
			v_new_json := to_jsonb(NEW);
		END IF;

		-- Extraer ID de la PK dinámicamente (asumiendo convenciones del proyecto id_nombre_tabla)
		-- Si no se puede extraer, se usa el registro completo como referencia
		IF (TG_OP = 'DELETE') THEN
			v_pk_id := (v_old_json->>(TG_TABLE_NAME || '_id'))::UUID; -- Caso tabla_id
			IF v_pk_id IS NULL THEN v_pk_id := (v_old_json->>'id_' || TG_TABLE_NAME)::UUID; END IF; -- Caso id_tabla
		ELSE
			v_pk_id := (v_new_json->>(TG_TABLE_NAME || '_id'))::UUID;
			IF v_pk_id IS NULL THEN v_pk_id := (v_new_json->>'id_' || TG_TABLE_NAME)::UUID; END IF;
		END IF;

		-- Solo insertar si hay cambios reales en UPDATE
		IF (TG_OP = 'UPDATE' AND v_old_json = v_new_json) THEN
			RETURN NEW;
		END IF;

		INSERT INTO auditoria_maestra (
			tabla_nombre, registro_id, accion, old_data, new_data, id_usuario, ip_address
		) VALUES (
			TG_TABLE_NAME, 
			COALESCE(v_pk_id, gen_random_uuid()), -- Fallback si no hay UUID
			TG_OP, v_old_json, v_new_json, v_user_id, 
			COALESCE(inet_client_addr()::VARCHAR, 'local/socket')
		);

		IF (TG_OP = 'DELETE') THEN RETURN OLD; END IF;
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

	-- 3. Procedimiento para aplicar auditoría a una tabla
	CREATE OR REPLACE PROCEDURE sp_habilitar_auditoria(p_tabla TEXT)
	LANGUAGE plpgsql AS $$
	BEGIN
		EXECUTE format(
			'DROP TRIGGER IF EXISTS tr_audit_maestra_%I ON %I;
			 CREATE TRIGGER tr_audit_maestra_%I
			 AFTER INSERT OR UPDATE OR DELETE ON %I
			 FOR EACH ROW EXECUTE FUNCTION fn_audit_generic();',
			p_tabla, p_tabla, p_tabla, p_tabla
		);
	END;
	$$;

	-- 4. Habilitar auditoría para tablas maestras
	CALL sp_habilitar_auditoria('usuario');
	CALL sp_habilitar_auditoria('empresa');
	CALL sp_habilitar_auditoria('sucursal');
	CALL sp_habilitar_auditoria('producto');
	CALL sp_habilitar_auditoria('cliente');
	CALL sp_habilitar_auditoria('proveedor');
	`
	_, err := db.Exec(query)
	return err
}
