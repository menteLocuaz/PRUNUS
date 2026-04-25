-- 000048_fix_periodo_nullability.up.sql
-- Corrige la restricción NOT NULL en prd_fecha_cierre para permitir la apertura de periodos

ALTER TABLE periodo ALTER COLUMN prd_fecha_cierre DROP NOT NULL;

-- Aprovechamos para asegurar que el campo nombre también sea opcional si se desea, 
-- pero prd_fecha_cierre es el error crítico actual.
COMMENT ON COLUMN periodo.prd_fecha_cierre IS 'Fecha y hora de cierre del periodo. Debe ser NULL mientras el periodo esté abierto.';
