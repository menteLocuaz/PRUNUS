-- 000047_periodo_snapshot.up.sql
-- Tabla para almacenar el resumen histórico al cerrar un periodo

CREATE TABLE IF NOT EXISTS periodo_snapshot (
    id_snapshot         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_periodo          UUID NOT NULL REFERENCES periodo(id_periodo),
    fecha_cierre        TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    total_ventas        NUMERIC(18,2) DEFAULT 0,
    total_efectivo      NUMERIC(18,2) DEFAULT 0,
    total_tarjeta       NUMERIC(18,2) DEFAULT 0,
    total_otros         NUMERIC(18,2) DEFAULT 0,
    total_diferencias   NUMERIC(18,2) DEFAULT 0, -- Suma de faltantes/sobrantes de todas las cajas
    total_operaciones   INTEGER DEFAULT 0,
    data_json           JSONB, -- Snapshot completo en formato JSON para reportes detallados
    id_usuario_cierre   UUID NOT NULL,
    created_at          TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_periodo_snapshot_periodo ON periodo_snapshot(id_periodo);
