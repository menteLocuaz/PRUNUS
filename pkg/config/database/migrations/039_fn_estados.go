// Package migrations contiene las funciones de migración de base de datos.
// Cada función es responsable de crear o modificar objetos en PostgreSQL
// (funciones, índices, tablas, etc.) de forma idempotente.
package migrations

import (
	"context" // Mejora 5: contexto para cancelación y timeout
	"database/sql"
	"fmt"
	"time"
)

// migrateFnEstados registra en PostgreSQL la función almacenada `modulos_ia_estados`
// y su índice de unicidad asociado.
//
// Esta migración es idempotente:
//   - Usa CREATE OR REPLACE FUNCTION → no falla si la función ya existe.
//   - Usa CREATE UNIQUE INDEX IF NOT EXISTS → no falla si el índice ya existe.
//
// Responsabilidades de la función SQL creada:
//   - Opción 0: Actualiza un estado existente en la tabla `estatus`.
//   - Opción 1: Inserta un nuevo estado en la tabla `estatus`.
//   - Ambas opciones retornan la lista actualizada de estados del módulo afectado.
//
// Parámetros de la función SQL:
//   - p_opcion       : 0 = actualizar, 1 = insertar.
//   - p_id_estado    : UUID del estado a actualizar (solo opción 0).
//   - p_descripcion  : Nombre/descripción del estado (no puede ser vacío).
//   - p_factor       : Factor numérico del estado (ej: "+1", "-1", "0").
//   - p_nivel        : Nivel jerárquico del estado dentro del módulo.
//   - p_id_modulo    : ID del módulo al que pertenece el estado.
//   - p_id_cadena    : Reservado para filtro por cadena (uso futuro).
//   - p_id_users_pos : ID del usuario que ejecuta la operación (auditoría).
//
// Retorna un error si la ejecución del DDL falla en la base de datos.
func migrateFnEstados(db *sql.DB) error {
	// Mejora 5: contexto con timeout de 30s para evitar migraciones colgadas.
	// Sin contexto, db.Exec puede bloquearse indefinidamente si el servidor
	// no responde, dejando el proceso de migración sin posibilidad de cancelación.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Mejora 6: DDL separado en statements independientes.
	// Razón: CREATE OR REPLACE FUNCTION y CREATE UNIQUE INDEX son DDL distintos.
	// Mezclarlos en un solo string hace que un fallo en el índice sea difícil
	// de distinguir de un fallo en la función. Separados, el error es preciso.
	statements := []struct {
		name string
		sql  string
	}{
		{
			// Statement 1: Función principal de administración de estados.
			name: "CREATE FUNCTION modulos_ia_estados",
			sql: `
			CREATE OR REPLACE FUNCTION modulos_ia_estados(
				p_opcion        INTEGER,
				p_id_estado     UUID,
				p_descripcion   VARCHAR(100),
				p_factor        VARCHAR(10),
				p_nivel         INTEGER,
				p_id_modulo     INTEGER,
				p_id_cadena     INTEGER,    -- Reservado para filtro futuro
				p_id_users_pos  VARCHAR(40) -- Reservado para auditoría
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
				-- Variables para captura de contexto de error en el bloque EXCEPTION
				v_ip_origen     VARCHAR(45);  -- IP del cliente (IPv4/IPv6)
				v_error_number  TEXT;         -- Código SQLSTATE (ej: '23505' = unique violation)
				v_error_proc    TEXT;         -- PG_EXCEPTION_DETAIL: detalle extendido
				v_error_line    TEXT;         -- PG_EXCEPTION_HINT: sugerencia del error
				v_error_msg     TEXT;         -- MESSAGE_TEXT: mensaje principal

				-- Mejora 4: constante documentada en lugar de magic string '1'.
				-- Representa el tipo de estado estándar creado por esta función.
				-- Si en el futuro se agregan tipos (ej: '2' = temporal), se cambia aquí.
				c_tipo_estado_default CONSTANT VARCHAR(10) := '1';
			BEGIN

				-- Mejora 7: Validación de descripción vacía o solo espacios en blanco.
				-- Sin esta guarda, se podrían insertar estados con nombre '' o '   ',
				-- que son inválidos para el negocio y difíciles de detectar en UI.
				IF TRIM(p_descripcion) = '' OR p_descripcion IS NULL THEN
					RAISE EXCEPTION 'La descripción no puede estar vacía o contener solo espacios';
				END IF;

				-- Normaliza la descripción eliminando espacios extremos.
				-- Evita duplicados por diferencias de espaciado (ej: 'Activo' vs ' Activo').
				p_descripcion := TRIM(p_descripcion);

				-- Validación temprana de opción: fail-fast antes de cualquier I/O.
				IF p_opcion NOT IN (0, 1) THEN
					RAISE EXCEPTION 'Opción inválida: %. Valores permitidos: 0 (actualizar), 1 (insertar)', p_opcion;
				END IF;

				-- ── Opción 0: Actualizar estado existente ───────────────────────
				IF p_opcion = 0 THEN

					-- Verifica unicidad de descripción dentro del módulo,
					-- excluyendo el propio registro para permitir guardar sin cambios.
					IF EXISTS (
						SELECT 1 FROM estatus
						WHERE std_descripcion = p_descripcion
						  AND id_status      <> p_id_estado
						  AND mdl_id          = p_id_modulo
						  AND deleted_at     IS NULL
					) THEN
						RAISE EXCEPTION 'Ya existe un estado con descripción "%" en el módulo %',
							p_descripcion, p_id_modulo;
					END IF;

					UPDATE estatus SET
						std_descripcion = p_descripcion,
						factor          = p_factor,
						nivel           = p_nivel,
						updated_at      = CURRENT_TIMESTAMP
						-- last_user    = p_id_users_pos  ← descomentar cuando exista la columna
					WHERE id_status  = p_id_estado
					  AND deleted_at IS NULL;

					-- Detecta UPDATE sin efecto: id inexistente o ya eliminado.
					-- Sin esta verificación el caller recibe éxito falso.
					IF NOT FOUND THEN
						RAISE EXCEPTION 'Estado con id_status % no encontrado o está eliminado', p_id_estado;
					END IF;

				-- ── Opción 1: Insertar nuevo estado ─────────────────────────────
				ELSIF p_opcion = 1 THEN

					-- Verifica unicidad sin excluir ningún ID (el registro aún no existe).
					IF EXISTS (
						SELECT 1 FROM estatus
						WHERE std_descripcion = p_descripcion
						  AND mdl_id          = p_id_modulo
						  AND deleted_at     IS NULL
					) THEN
						RAISE EXCEPTION 'Ya existe un estado con descripción "%" en el módulo %',
							p_descripcion, p_id_modulo;
					END IF;

					-- Mejora 8: Se poblan created_at y updated_at en el INSERT.
					-- En la versión original solo se guardaba updated_at, dejando
					-- created_at sin valor y perdiendo la fecha real de creación.
					-- created_at nunca debe modificarse después del INSERT.
					INSERT INTO estatus (
						std_descripcion,
						factor,
						nivel,
						mdl_id,
						std_tipo_estado,
						is_active,
						created_at,      -- ← nuevo: fecha de creación inmutable
						updated_at
						-- last_user    ← descomentar cuando exista la columna
					) VALUES (
						p_descripcion,
						p_factor,
						p_nivel,
						p_id_modulo,
						c_tipo_estado_default,  -- Mejora 4: constante en lugar de '1'
						TRUE,                   -- todo estado nuevo nace activo
						CURRENT_TIMESTAMP,      -- created_at: se setea una sola vez
						CURRENT_TIMESTAMP       -- updated_at: se actualizará en futuros UPDATEs
					);

				END IF;

				-- ── Resultado: lista actualizada del módulo ──────────────────────
				-- Único RETURN QUERY al final (DRY): ambas opciones retornan
				-- el mismo conjunto de datos, ordenado para presentación en UI.
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

			-- ── Manejo centralizado de errores ──────────────────────────────────
			EXCEPTION
				WHEN OTHERS THEN
					GET STACKED DIAGNOSTICS
						v_error_msg  = MESSAGE_TEXT,
						v_error_proc = PG_EXCEPTION_DETAIL,
						v_error_line = PG_EXCEPTION_HINT;

					-- Mejora 2: COALESCE previene NULL en conexiones locales (socket Unix).
					-- inet_client_addr() retorna NULL cuando la conexión es local,
					-- lo que causaría que el log muestre "IP: <NULL>" sin contexto.
					v_ip_origen    := COALESCE(inet_client_addr()::VARCHAR, 'local/socket');
					v_error_number := SQLSTATE;

					-- Mejora 3: Un solo RAISE LOG con formato JSON estructurado.
					-- La versión original generaba 7 entradas de log separadas,
					-- dificultando correlacionar todos los campos de un mismo error.
					-- Con JSON, herramientas como Datadog, Loki o pgBadger pueden
					-- parsear y agrupar el evento completo en una sola entrada.
					RAISE LOG 'modulos_ia_estados error: %', json_build_object(
						'fecha',     CURRENT_TIMESTAMP,
						'usuario',   p_id_users_pos,
						'ip',        v_ip_origen,
						'sqlstate',  v_error_number,
						'detalle',   v_error_proc,
						'hint',      v_error_line,
						'mensaje',   v_error_msg
					);

					-- Re-lanza para que Go reciba un error no-nil en db.ExecContext.
					RAISE EXCEPTION 'modulos_ia_estados falló [%]: %', v_error_number, v_error_msg;
			END;
			$$;`,
		},
		{
			// Statement 2: Índice UNIQUE parcial separado de la función.
			// Mejora 6: separado del DDL de la función para aislar errores.
			// Si el índice falla (ej: datos duplicados existentes), el error
			// indica exactamente qué statement falló, no el bloque completo.
			name: "CREATE UNIQUE INDEX idx_estatus_desc_modulo",
			sql: `
			CREATE UNIQUE INDEX IF NOT EXISTS idx_estatus_desc_modulo
				ON estatus (mdl_id, std_descripcion)
				WHERE deleted_at IS NULL;`,
		},
	}

	// Ejecuta cada statement de forma independiente.
	// Mejora 5 + 6: ExecContext con timeout y statements separados permiten
	// identificar exactamente cuál DDL falló y cancelar si el servidor no responde.
	for _, stmt := range statements {
		if _, err := db.ExecContext(ctx, stmt.sql); err != nil {
			// Wrapping del error con el nombre del statement para trazabilidad.
			// Sin esto, el caller solo recibe el error de PostgreSQL sin saber
			// qué parte de la migración falló.
			return fmt.Errorf("migrateFnEstados: fallo en [%s]: %w", stmt.name, err)
		}
	}

	return nil
}
