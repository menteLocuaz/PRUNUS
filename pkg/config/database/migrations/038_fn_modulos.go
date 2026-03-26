package migrations

import "database/sql"

func migrateFnModulos(db *sql.DB) error {
	query := `
	-- ============================================================
	-- Función PostgreSQL equivalente a [config].[MODULOS_IA_modulos]
	-- Opción 0 → Actualizar módulo existente
	-- Opción 1 → Insertar nuevo módulo
	-- Retorna: TABLE con la lista de módulos
	-- ============================================================
	CREATE OR REPLACE FUNCTION modulos_ia_modulo(
		p_opcion        INTEGER,
		p_id_modulo     INTEGER,        -- mdl_id (serial legible, no UUID)
		p_descripcion   VARCHAR(100),
		p_abreviatura   VARCHAR(50),
		p_nivel         INTEGER,
		p_estado        BOOLEAN,
		p_id_cadena     INTEGER,        -- Reservado para uso futuro
		p_id_users_pos  VARCHAR(40)     -- Reservado para auditoría
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
		v_valido        BOOLEAN := TRUE;
		v_next_id       INTEGER;
		v_ip_origen     VARCHAR(45);
		v_error_number  TEXT;
		v_error_proc    TEXT;
		v_error_line    TEXT;
		v_error_msg     TEXT;
	BEGIN

		-- ── Validación: descripción no repetida en otro módulo ──────────────
		SELECT FALSE INTO v_valido
		FROM modulo
		WHERE mdl_descripcion = p_descripcion
		  AND mdl_id <> p_id_modulo
		  AND deleted_at IS NULL
		LIMIT 1;

		IF v_valido IS NULL THEN
			v_valido := TRUE;
		END IF;

		IF v_valido THEN

			-- ── Opción 0: Actualizar módulo existente ───────────────────────
			IF p_opcion = 0 THEN
				UPDATE modulo SET
					mdl_descripcion = p_descripcion,
					abreviatura     = p_abreviatura,
					nivel           = p_nivel,
					is_active       = p_estado,
					updated_at      = CURRENT_TIMESTAMP
				WHERE mdl_id = p_id_modulo;

			-- ── Opción 1: Insertar nuevo módulo ─────────────────────────────
			ELSIF p_opcion = 1 THEN

				SELECT COALESCE(MAX(mdl_id), 0)
				INTO v_next_id
				FROM modulo;

				v_next_id := v_next_id + 1;

				INSERT INTO modulo (
					mdl_id,
					mdl_descripcion,
					abreviatura,
					nivel,
					is_active,
					updated_at
				)
				VALUES (
					v_next_id,
					p_descripcion,
					p_abreviatura,
					p_nivel,
					p_estado,
					CURRENT_TIMESTAMP
				);

				-- Resincronizar secuencia SERIAL tras inserción manual si es necesario
				PERFORM setval(
					pg_get_serial_sequence('modulo', 'mdl_id'),
					v_next_id
				);

			END IF;

			-- ── Lista de módulos ──
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

		END IF;

	EXCEPTION
		WHEN OTHERS THEN
			GET STACKED DIAGNOSTICS
				v_error_msg  = MESSAGE_TEXT,
				v_error_proc = PG_EXCEPTION_DETAIL,
				v_error_line = PG_EXCEPTION_HINT;

			v_ip_origen    := inet_client_addr()::VARCHAR;
			v_error_number := SQLSTATE;

			RAISE LOG 'Fecha: %',           CURRENT_TIMESTAMP;
			RAISE LOG 'Usuario: %',         p_id_users_pos;
			RAISE LOG 'IP: %',              v_ip_origen;
			RAISE LOG 'ErrorNumber: %',     v_error_number;
			RAISE LOG 'ErrorProcedure: %',  v_error_proc;
			RAISE LOG 'ErrorLine: %',       v_error_line;
			RAISE LOG 'ErrorMessage: %',    v_error_msg;

			RAISE EXCEPTION 'modulos_ia_modulo falló [%]: %', v_error_number, v_error_msg;

	END;
	$$;
	`
	_, err := db.Exec(query)
	return err
}
