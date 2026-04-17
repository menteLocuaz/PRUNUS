-- 000043_periodo_operational_controls.up.sql
-- Añade controles de concurrencia, auditoría y soporte multi-sucursal a la tabla 'periodo'

-- 1. Añadir columnas de auditoría y sucursal
ALTER TABLE periodo 
    ADD COLUMN IF NOT EXISTS id_sucursal         UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000',
    ADD COLUMN IF NOT EXISTS prd_ip_apertura     VARCHAR(45),
    ADD COLUMN IF NOT EXISTS prd_motivo_apertura TEXT,
    ADD COLUMN IF NOT EXISTS prd_ip_cierre       VARCHAR(45);

-- 2. Garantizar que id_sucursal esté correctamente referenciado (si existe la tabla)
DO $$ 
BEGIN 
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'sucursal') THEN
        -- Ajustamos el DEFAULT si es necesario o eliminamos si ya hay datos
        ALTER TABLE periodo ALTER COLUMN id_sucursal DROP DEFAULT;
        ALTER TABLE periodo ADD CONSTRAINT fk_periodo_sucursal FOREIGN KEY (id_sucursal) REFERENCES sucursal(id_sucursal);
    END IF;
END $$;

-- 3. BLOQUEO DE CONCURRENCIA: Restricción de unicidad parcial
-- Esto evita que dos procesos inserten un periodo abierto para la misma sucursal simultáneamente.
CREATE UNIQUE INDEX IF NOT EXISTS idx_periodo_unico_activo_sucursal 
ON periodo (id_sucursal) 
WHERE (prd_fecha_cierre IS NULL AND deleted_at IS NULL);

-- 4. Comentarios
COMMENT ON COLUMN periodo.prd_ip_apertura IS 'IP desde la cual se abrió el periodo';
COMMENT ON COLUMN periodo.prd_motivo_apertura IS 'Motivo opcional de la apertura';
