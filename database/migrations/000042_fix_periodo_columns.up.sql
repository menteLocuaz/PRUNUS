-- 000042_fix_periodo_columns.up.sql
-- Alinea la tabla 'periodo' con los nombres de columna esperados por el Store en Go

ALTER TABLE periodo RENAME COLUMN fecha_inicio TO prd_fecha_apertura;
ALTER TABLE periodo RENAME COLUMN fecha_fin TO prd_fecha_cierre;

-- Asegurar que las columnas de auditoría existan si el Store las usa
ALTER TABLE periodo ADD COLUMN IF NOT EXISTS prd_usuario_apertura UUID;
ALTER TABLE periodo ADD COLUMN IF NOT EXISTS prd_usuario_cierre   UUID;

-- Eliminar columna redundante is_cerrado ya que se usa prd_fecha_cierre para determinar si está abierto
ALTER TABLE periodo DROP COLUMN IF EXISTS is_cerrado;
