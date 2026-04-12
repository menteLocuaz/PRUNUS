-- 1. TABLA DE CATEGORÍAS DE GASTOS
CREATE TABLE IF NOT EXISTS gasto_categoria (
    id_categoria_gasto UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nombre             VARCHAR(100) NOT NULL UNIQUE,
    descripcion        TEXT,
    id_status          UUID NOT NULL REFERENCES estatus(id_status),
    created_at         TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at         TIMESTAMPTZ NULL
);

-- 2. TABLA DE GASTOS OPERATIVOS
CREATE TABLE IF NOT EXISTS gasto_operativo (
    id_gasto           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_sucursal        UUID NOT NULL REFERENCES sucursal(id_sucursal),
    id_categoria_gasto UUID NOT NULL REFERENCES gasto_categoria(id_categoria_gasto),
    id_usuario         UUID NOT NULL REFERENCES usuario(id_usuario),
    monto              NUMERIC(18,2) NOT NULL DEFAULT 0,
    referencia         VARCHAR(100),
    observacion        TEXT,
    fecha_gasto        TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    id_status          UUID NOT NULL REFERENCES estatus(id_status),
    metadata           JSONB DEFAULT '{}',
    created_at         TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at         TIMESTAMPTZ NULL
);

-- 3. ÍNDICES DE RENDIMIENTO (PostgreSQL Optimization)
CREATE INDEX IF NOT EXISTS idx_gasto_sucursal_fecha ON gasto_operativo(id_sucursal, fecha_gasto DESC);
CREATE INDEX IF NOT EXISTS idx_gasto_categoria ON gasto_operativo(id_categoria_gasto);

-- 4. APLICAR FRAMEWORK CORE
CALL sp_core_setup_table('gasto_categoria');
CALL sp_core_setup_table('gasto_operativo');
