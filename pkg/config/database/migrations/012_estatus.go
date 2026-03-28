// Package migrations contiene las funciones de migración de base de datos.
package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// migrateEstatus crea y normaliza la tabla `estatus` en PostgreSQL.
// Consolidado para migraciones limpias e idempotentes.
func migrateEstatus(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("migrateEstatus: error al iniciar transacción: %w", err)
	}
	defer tx.Rollback()

	// 1. Definición final de la tabla estatus
	// Se asegura la existencia del esquema public y se define la tabla limpia.
	query := `
	CREATE SCHEMA IF NOT EXISTS public;

	CREATE TABLE IF NOT EXISTS public.estatus (
		id_status       UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
		std_descripcion VARCHAR(255) NOT NULL,
		std_tipo_estado VARCHAR(255) NOT NULL DEFAULT 'GENERAL',
		factor          VARCHAR(10),
		nivel           INTEGER      NOT NULL DEFAULT 0,
		mdl_id          INTEGER      NOT NULL DEFAULT -1,
		is_active       BOOLEAN      NOT NULL DEFAULT TRUE,
		created_at      TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at      TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at      TIMESTAMPTZ  NULL
	);

	-- Índices de rendimiento
	CREATE INDEX IF NOT EXISTS idx_estatus_tipo ON public.estatus(std_tipo_estado);
	CREATE INDEX IF NOT EXISTS idx_estatus_modulo ON public.estatus(mdl_id);
	CREATE INDEX IF NOT EXISTS idx_estatus_deleted_at ON public.estatus(deleted_at) WHERE deleted_at IS NOT NULL;

	-- 2. Datos semilla esenciales del sistema
	INSERT INTO public.estatus (id_status, std_descripcion, std_tipo_estado, mdl_id) VALUES
    ('fc273a6a-ab7b-4453-a560-ac62fa64348b', 'Activa',                      'GENERAL',      1),
    ('b4d0544d-1778-4560-a170-a681ab3399bd', 'Suspendida',                  'RESTRICCION',  1),
    ('c3e12abc-0011-4abc-b123-000000000001', 'Inactiva',                    'GENERAL',      1),
    ('6cf06fbe-b21c-46e3-a34b-b24f5167cd9a', 'Abierta',                    'OPERATIVO',    2),
    ('34be5a4c-ab4c-4afd-9a3c-98dfeb500fbc', 'Cerrada',                    'OPERATIVO',    2),
    ('3a99d245-b34f-48a5-ac08-a5a010c5822f', 'Activo',                     'ACCESO',       3),
    ('31f4e127-e7e1-414d-aaef-6e92e4c5d970', 'Disponible',                 'STOCK',        4),
    ('073f4513-a88e-4278-af95-bd9cde61bdbd', 'Agotado',                    'STOCK',        4),
    ('892340e0-4328-491d-9102-80550bb6aac4', 'Pendiente de Pago',          'TRANSACCION',  5),
    ('0f447fd7-9849-4a68-b82f-c69297e7a924', 'Pagada',                     'TRANSACCION',  5),
    ('62ed7d82-0c81-4511-8f02-e7fd140018d8', 'Anulada',                    'TRANSACCION',  5),
    ('0cd9aa6e-5768-45d2-a66d-12758a3bd0cc', 'Solicitada',                 'FLUJO',        6),
    ('877e725b-ae57-4501-b1cd-1158fe2df087', 'Pendiente de Conciliación',  'CONTABLE',     7),
    ('e69d8b1d-d267-47ab-b0ef-507ed6382cd3', 'Conciliado',                 'CONTABLE',     7),
    ('59039503-85cf-e511-80c1-000c29c9e0e0', 'Activo',                     'SESION',       8),
    ('5a039503-85cf-e511-80c1-000c29c9e0e0', 'Inactivo',                   'SESION',       8),
    ('99039503-85cf-e511-80c1-000c29c9e0e0', 'Fondo Asignado',             'APERTURA',     8)
	ON CONFLICT (id_status) DO UPDATE SET 
		std_descripcion = EXCLUDED.std_descripcion,
		std_tipo_estado = EXCLUDED.std_tipo_estado,
		mdl_id = EXCLUDED.mdl_id;
	`

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("migrateEstatus: fallo en ejecución: %w", err)
	}

	return tx.Commit()
}
