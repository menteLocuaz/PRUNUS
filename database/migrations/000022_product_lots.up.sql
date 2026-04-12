-- 1. TABLA DE LOTES (Trazabilidad y PEPS)
CREATE TABLE IF NOT EXISTS lotes (
    id_lote          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_producto      UUID NOT NULL REFERENCES producto(id_producto),
    id_sucursal      UUID NOT NULL REFERENCES sucursal(id_sucursal),
    codigo_lote      VARCHAR(100) NOT NULL,
    cantidad_inicial NUMERIC(12,2) NOT NULL DEFAULT 0,
    cantidad_actual  NUMERIC(12,2) NOT NULL DEFAULT 0,
    costo_compra     NUMERIC(18,2) NOT NULL DEFAULT 0,
    fecha_vencimiento TIMESTAMPTZ NULL,
    fecha_recepcion   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    id_status        UUID REFERENCES estatus(id_status),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at       TIMESTAMPTZ NULL
);

-- 2. ÍNDICES DE RENDIMIENTO (PostgreSQL Optimization)
-- Optimiza la búsqueda del primer lote disponible (PEPS) ignorando lotes vacíos o eliminados.
CREATE INDEX IF NOT EXISTS idx_lotes_fifo_performance 
ON lotes(id_producto, id_sucursal, fecha_recepcion ASC) 
WHERE deleted_at IS NULL AND cantidad_actual > 0;

-- Optimiza alertas de vencimiento.
CREATE INDEX IF NOT EXISTS idx_lotes_vencimiento_performance 
ON lotes(fecha_vencimiento) 
WHERE deleted_at IS NULL AND cantidad_actual > 0;

-- 3. APLICAR FRAMEWORK CORE
CALL sp_core_setup_table('lotes');
