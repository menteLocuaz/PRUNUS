// Package migrations contiene las funciones de migración de base de datos.
// Cada función aplica cambios DDL de forma idempotente sobre PostgreSQL.
package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// migrateEstatus crea y normaliza la tabla `estatus` en PostgreSQL.
//
// Es idempotente: puede ejecutarse múltiples veces sin efectos secundarios.
//
// Responsabilidades:
//  1. Crear la tabla `estatus` si no existe.
//  2. Normalizar el nombre de columna `stp_tipo_estado` → `std_tipo_estado`
//     en bases de datos migradas desde versiones anteriores.
//  3. Asegurar que todas las columnas requeridas existan (migraciones incrementales).
//  4. Crear índices de rendimiento.
//  5. Insertar los estados semilla del sistema si no existen.
func migrateEstatus(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Iniciamos una transacción para asegurar que SET search_path persista
	// en toda la sesión y que la migración sea atómica.
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("migrateEstatus: error al iniciar transacción: %w", err)
	}
	// Aseguramos el rollback si algo sale mal
	defer tx.Rollback()

	statements := []struct {
		name string
		sql  string
	}{
		// ── 0. Esquema ──────────────────────────────────────────────────────
		{
			name: "CREATE SCHEMA public",
			sql:  "CREATE SCHEMA IF NOT EXISTS public;",
		},
		{
			name: "SET search_path",
			sql:  "SET search_path TO public, pg_catalog;",
		},
		// ── 1. Tabla principal ───────────────────────────────────────────────
		{
			name: "CREATE TABLE estatus",
			sql: `
			CREATE TABLE IF NOT EXISTS public.estatus (
				id_status       UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
				std_descripcion VARCHAR(255) NOT NULL,
				std_tipo_estado VARCHAR(255) NOT NULL DEFAULT '1',
				factor          VARCHAR(10),
				nivel           INTEGER      NOT NULL DEFAULT 0,
				mdl_id          INTEGER      NOT NULL,
				is_active       BOOLEAN      NOT NULL DEFAULT TRUE,
				created_at      TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
				updated_at      TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
				deleted_at      TIMESTAMPTZ  NULL
			);`,
		},

		// ── 2. Normalización de nombre de columna (typo histórico) ──────────
		{
			name: "RENAME stp_tipo_estado → std_tipo_estado",
			sql: `
			DO $$
			BEGIN
				IF EXISTS (
					SELECT 1 FROM information_schema.columns
					WHERE table_schema = 'public'
					  AND table_name   = 'estatus'
					  AND column_name  = 'stp_tipo_estado'
				) AND NOT EXISTS (
					SELECT 1 FROM information_schema.columns
					WHERE table_schema = 'public'
					  AND table_name   = 'estatus'
					  AND column_name  = 'std_tipo_estado'
				) THEN
					ALTER TABLE estatus RENAME COLUMN stp_tipo_estado TO std_tipo_estado;
				END IF;
			END $$;`,
		},

		// ── 3. Columna std_tipo_estado (si el rename no ocurrió) ────────────
		{
			name: "ADD COLUMN std_tipo_estado",
			sql:  `ALTER TABLE estatus ADD COLUMN IF NOT EXISTS std_tipo_estado VARCHAR(255);`,
		},
		{
			name: "BACKFILL std_tipo_estado NULLs",
			sql:  `UPDATE estatus SET std_tipo_estado = '1' WHERE std_tipo_estado IS NULL;`,
		},
		{
			name: "SET NOT NULL std_tipo_estado",
			sql:  `ALTER TABLE estatus ALTER COLUMN std_tipo_estado SET NOT NULL;`,
		},

		// ── 4. Columna stp_tipo_estado legacy: quitar NOT NULL si aún existe ─
		{
			name: "DROP NOT NULL stp_tipo_estado (legacy)",
			sql: `
			DO $$
			BEGIN
				IF EXISTS (
					SELECT 1 FROM information_schema.columns
					WHERE table_schema = 'public'
					  AND table_name   = 'estatus'
					  AND column_name  = 'stp_tipo_estado'
				) THEN
					ALTER TABLE estatus ALTER COLUMN stp_tipo_estado DROP NOT NULL;
				END IF;
			END $$;`,
		},

		// ── 5. Columnas opcionales (migraciones incrementales) ───────────────
		{
			name: "ADD COLUMN mdl_id",
			sql:  `ALTER TABLE estatus ADD COLUMN IF NOT EXISTS mdl_id INTEGER;`,
		},
		{
			name: "BACKFILL mdl_id NULLs",
			sql:  `UPDATE estatus SET mdl_id = -1 WHERE mdl_id IS NULL;`,
		},
		{
			name: "SET NOT NULL mdl_id",
			sql:  `ALTER TABLE estatus ALTER COLUMN mdl_id SET NOT NULL;`,
		},
		{
			name: "ADD COLUMN factor",
			sql:  `ALTER TABLE estatus ADD COLUMN IF NOT EXISTS factor VARCHAR(10);`,
		},
		{
			name: "ADD COLUMN nivel",
			sql:  `ALTER TABLE estatus ADD COLUMN IF NOT EXISTS nivel INTEGER DEFAULT 0;`,
		},
		{
			name: "ADD COLUMN is_active",
			sql:  `ALTER TABLE estatus ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT TRUE;`,
		},

		// ── 6. Columnas de auditoría ─────────────────────────────────────────
		{
			name: "ADD COLUMN created_at (nullable primero)",
			sql:  `ALTER TABLE estatus ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ;`,
		},
		{
			name: "BACKFILL created_at",
			sql:  `UPDATE estatus SET created_at = CURRENT_TIMESTAMP WHERE created_at IS NULL;`,
		},
		{
			name: "SET NOT NULL created_at",
			sql: `
			ALTER TABLE estatus
				ALTER COLUMN created_at SET NOT NULL,
				ALTER COLUMN created_at SET DEFAULT CURRENT_TIMESTAMP;`,
		},
		{
			name: "ADD COLUMN updated_at (nullable primero)",
			sql:  `ALTER TABLE estatus ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ;`,
		},
		{
			name: "BACKFILL updated_at",
			sql:  `UPDATE estatus SET updated_at = CURRENT_TIMESTAMP WHERE updated_at IS NULL;`,
		},
		{
			name: "SET NOT NULL updated_at",
			sql: `
			ALTER TABLE estatus
				ALTER COLUMN updated_at SET NOT NULL,
				ALTER COLUMN updated_at SET DEFAULT CURRENT_TIMESTAMP;`,
		},
		{
			name: "ADD COLUMN deleted_at",
			sql: `ALTER TABLE estatus ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ NULL;`,
		},

		// ── 7. Índices ───────────────────────────────────────────────────────
		{
			name: "CREATE INDEX idx_estatus_tipo",
			sql:  `CREATE INDEX IF NOT EXISTS idx_estatus_tipo ON estatus(std_tipo_estado);`,
		},
		{
			name: "CREATE INDEX idx_estatus_modulo",
			sql:  `CREATE INDEX IF NOT EXISTS idx_estatus_modulo ON estatus(mdl_id);`,
		},
		{
			name: "CREATE PARTIAL INDEX idx_estatus_deleted_at",
			sql: `
			CREATE INDEX IF NOT EXISTS idx_estatus_deleted_at
				ON estatus(deleted_at)
				WHERE deleted_at IS NOT NULL;`,
		},

		// ── 8. Datos semilla ─────────────────────────────────────────────────
		{
			name: "SEED estatus",
			sql: `
			INSERT INTO estatus (id_status, std_descripcion, std_tipo_estado, mdl_id) VALUES
				('59039503-85cf-e511-80c1-000c29c9e0e0', 'Activo',             '1', 8),
				('5a039503-85cf-e511-80c1-000c29c9e0e0', 'Inactivo',           '1', 8),
				('99039503-85cf-e511-80c1-000c29c9e0e0', 'Fondo Asignado',     '1', 8),
				('9a039503-85cf-e511-80c1-000c29c9e0e0', 'Desmontado',         '1', 8),
				('9b039503-85cf-e511-80c1-000c29c9e0e0', 'Ingreso Admin',      '1', 8),
				('9c039503-85cf-e511-80c1-000c29c9e0e0', 'Salir Admin',        '1', 8),
				('0d4515fe-c907-e611-a6b8-000c29c9e0e0', 'Arqueo',             '1', 8),
				('159e3fe6-630e-e611-80c1-000c29c9e0e0', 'Arqueo Retiros',     '1', 8),
				('e8297cfa-630e-e611-80c1-000c29c9e0e0', 'Retiro Efectivo',    '1', 8),
				('84920103-640e-e611-80c1-000c29c9e0e0', 'Retiro Total',       '1', 8),
				('a864475f-0d34-e711-80c1-000c29c9e0e0', 'Fondo Activo',       '1', 8),
				('2160b065-0d34-e711-80c1-000c29c9e0e0', 'Fondo Retirado',     '1', 8),
				('5e8dd0fb-5550-e711-80c1-000c29c9e0e0', 'Fondo Por Confirmar','1', 8)
			ON CONFLICT (id_status) DO NOTHING;`,
		},
	}

	// Ejecutamos cada statement dentro de la transacción.
	for _, stmt := range statements {
		if _, err := tx.ExecContext(ctx, stmt.sql); err != nil {
			return fmt.Errorf("migrateEstatus: fallo en [%s]: %w", stmt.name, err)
		}
	}

	// Si todo salió bien, confirmamos los cambios.
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("migrateEstatus: error al hacer commit: %w", err)
	}

	return nil
}
