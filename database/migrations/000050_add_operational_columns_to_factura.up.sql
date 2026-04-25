-- 000050_add_operational_columns_to_factura.up.sql
-- Añade columnas para vinculación con estación, periodo y control de caja en la tabla factura

ALTER TABLE factura 
    ADD COLUMN IF NOT EXISTS id_estacion         UUID,
    ADD COLUMN IF NOT EXISTS id_periodo          UUID,
    ADD COLUMN IF NOT EXISTS id_control_estacion  UUID;

-- Comentarios para auditoría
COMMENT ON COLUMN factura.id_estacion IS 'ID de la estación física donde se generó la venta';
COMMENT ON COLUMN factura.id_periodo IS 'ID del periodo contable administrativo';
COMMENT ON COLUMN factura.id_control_estacion IS 'ID de la sesión de caja (control_estacion) activa';

-- Índices para optimización de reportes por turno/estación
CREATE INDEX IF NOT EXISTS idx_factura_periodo ON factura(id_periodo);
CREATE INDEX IF NOT EXISTS idx_factura_control ON factura(id_control_estacion);
