-- 000046_advanced_cash_control.up.sql
-- Tablas y columnas para control avanzado de arqueo y cierre

-- 1. Tabla para desglose de denominaciones (monedas y billetes)
CREATE TABLE IF NOT EXISTS arqueo_denominaciones (
    id_arqueo_den   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_control_est  UUID NOT NULL REFERENCES control_estacion(id_control_estacion),
    tipo            VARCHAR(20) NOT NULL, -- 'APERTURA', 'CIERRE'
    valor_nominal   NUMERIC(18,2) NOT NULL,
    cantidad        INTEGER NOT NULL,
    subtotal        NUMERIC(18,2) NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- 2. Asegurar que control_estacion tenga campos para el resumen final (si no existen)
ALTER TABLE control_estacion 
    ADD COLUMN IF NOT EXISTS ventas_efectivo    NUMERIC(18,2) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS ventas_tarjeta     NUMERIC(18,2) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS ventas_transferencia NUMERIC(18,2) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS total_retiros      NUMERIC(18,2) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS total_gastos       NUMERIC(18,2) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS saldo_esperado     NUMERIC(18,2) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS saldo_real         NUMERIC(18,2) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS diferencia         NUMERIC(18,2) DEFAULT 0;

-- 3. Índices para reportes
CREATE INDEX IF NOT EXISTS idx_arqueo_den_control ON arqueo_denominaciones(id_control_est, tipo);
