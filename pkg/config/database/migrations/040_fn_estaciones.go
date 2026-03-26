package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// migrateFnEstaciones registra la función almacenada para administrar estaciones POS.
func migrateFnEstaciones(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("migrateFnEstaciones: error al iniciar transacción: %w", err)
	}
	defer tx.Rollback()

	statements := []struct {
		name string
		sql  string
	}{
		{
			name: "CREATE FUNCTION estaciones_ia_estacion",
			sql: `
			CREATE OR REPLACE FUNCTION estaciones_ia_estacion(
				p_opcion            INTEGER,        -- 0 = actualizar, 1 = insertar
				p_id_estacion       UUID,           -- ID de la estación (solo opción 0)
				p_codigo            VARCHAR(50),    -- Código único de la estación
				p_nombre            VARCHAR(255),   -- Nombre de la estación
				p_ip                VARCHAR(50),    -- Dirección IP
				p_id_sucursal       UUID,           -- Sucursal a la que pertenece
				p_estado_desc       VARCHAR(100),   -- Descripción del estado (ej: 'Activo')
				p_id_user_pos       UUID            -- ID del usuario que realiza la operación (auditoría)
			)
			RETURNS TABLE (
				id_estacion     UUID,
				codigo          VARCHAR(50),
				nombre          VARCHAR(255),
				ip              VARCHAR(50),
				id_sucursal     UUID,
				is_active       BOOLEAN,
				id_status       UUID
			)
			LANGUAGE plpgsql
			AS $$
			DECLARE
				v_id_status     UUID;
				v_ip_origen     VARCHAR(45);
				v_error_number  TEXT;
				v_error_proc    TEXT;
				v_error_line    TEXT;
				v_error_msg     TEXT;
			BEGIN
				-- ── 1. Resolver el UUID del estado ──────────────────────────────
				-- Se busca en la tabla estatus por descripción y módulo 6 (Estaciones)
				SELECT id_status INTO v_id_status
				FROM estatus
				WHERE std_descripcion = p_estado_desc
				  AND mdl_id          = 6  -- Módulo Estaciones
				  AND deleted_at     IS NULL
				LIMIT 1;

				IF v_id_status IS NULL THEN
					RAISE EXCEPTION 'El estado "%" no existe para el módulo de Estaciones', p_estado_desc;
				END IF;

				-- ── 2. Validaciones de unicidad ─────────────────────────────────
				-- Verificar IP única dentro de la sucursal (excluyendo la actual)
				IF EXISTS (
					SELECT 1 FROM estaciones_pos
					WHERE ip          = p_ip
					  AND id_sucursal = p_id_sucursal
					  AND id_estacion <> COALESCE(p_id_estacion, '00000000-0000-0000-0000-000000000000'::UUID)
					  AND deleted_at IS NULL
				) THEN
					RAISE EXCEPTION 'La dirección IP % ya está asignada a otra estación en esta sucursal', p_ip;
				END IF;

				-- Verificar código único
				IF EXISTS (
					SELECT 1 FROM estaciones_pos
					WHERE codigo      = p_codigo
					  AND id_estacion <> COALESCE(p_id_estacion, '00000000-0000-0000-0000-000000000000'::UUID)
					  AND deleted_at IS NULL
				) THEN
					RAISE EXCEPTION 'El código de estación % ya existe', p_codigo;
				END IF;

				-- ── 3. Operación Principal ──────────────────────────────────────
				
				-- Opción 0: Actualizar estación existente
				IF p_opcion = 0 THEN
					UPDATE estaciones_pos SET
						codigo      = p_codigo,
						nombre      = p_nombre,
						ip          = p_ip,
						id_sucursal = p_id_sucursal,
						id_status   = v_id_status,
						updated_at  = CURRENT_TIMESTAMP
					WHERE id_estacion = p_id_estacion
					  AND deleted_at IS NULL;

					IF NOT FOUND THEN
						RAISE EXCEPTION 'Estación con ID % no encontrada o está eliminada', p_id_estacion;
					END IF;

				-- Opción 1: Insertar nueva estación
				ELSIF p_opcion = 1 THEN
					INSERT INTO estaciones_pos (
						codigo,
						nombre,
						ip,
						id_sucursal,
						id_status,
						created_at,
						updated_at
					) VALUES (
						p_codigo,
						p_nombre,
						p_ip,
						p_id_sucursal,
						v_id_status,
						CURRENT_TIMESTAMP,
						CURRENT_TIMESTAMP
					);
				ELSE
					RAISE EXCEPTION 'Opción inválida: %. Use 0 para actualizar o 1 para insertar', p_opcion;
				END IF;

				-- ── 4. Retorno de Resultados ────────────────────────────────────
				RETURN QUERY
					SELECT
						e.id_estacion,
						e.codigo,
						e.nombre,
						e.ip,
						e.id_sucursal,
						s.is_active,
						e.id_status
					FROM estaciones_pos e
					JOIN estatus s ON s.id_status = e.id_status
					WHERE e.id_sucursal = p_id_sucursal
					  AND e.deleted_at IS NULL
					ORDER BY e.nombre;

			EXCEPTION
				WHEN OTHERS THEN
					GET STACKED DIAGNOSTICS
						v_error_msg  = MESSAGE_TEXT,
						v_error_proc = PG_EXCEPTION_DETAIL,
						v_error_line = PG_EXCEPTION_HINT;

					v_ip_origen    := COALESCE(inet_client_addr()::VARCHAR, 'local/socket');
					v_error_number := SQLSTATE;

					RAISE LOG 'estaciones_ia_estacion error: %', json_build_object(
						'fecha',     CURRENT_TIMESTAMP,
						'usuario',   p_id_user_pos,
						'ip',        v_ip_origen,
						'sqlstate',  v_error_number,
						'detalle',   v_error_proc,
						'hint',      v_error_line,
						'mensaje',   v_error_msg
					);

					RAISE EXCEPTION 'estaciones_ia_estacion falló [%]: %', v_error_number, v_error_msg;
			END;
			$$;`,
		},
	}

	for _, stmt := range statements {
		if _, err := tx.ExecContext(ctx, stmt.sql); err != nil {
			return fmt.Errorf("migrateFnEstaciones: fallo en [%s]: %w", stmt.name, err)
		}
	}

	return tx.Commit()
}
