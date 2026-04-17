-- 000042_fix_periodo_columns.down.sql
-- Revertir cambios de la tabla 'periodo'

ALTER TABLE periodo RENAME COLUMN prd_fecha_apertura TO fecha_inicio;
ALTER TABLE periodo RENAME COLUMN prd_fecha_cierre TO fecha_fin;

ALTER TABLE periodo DROP COLUMN IF EXISTS prd_usuario_apertura;
ALTER TABLE periodo DROP COLUMN IF EXISTS prd_usuario_cierre;

ALTER TABLE periodo ADD COLUMN IF NOT EXISTS is_cerrado BOOLEAN DEFAULT FALSE;
