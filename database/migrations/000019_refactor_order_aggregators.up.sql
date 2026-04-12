-- 1. TABLA DE AGREGADORES (Maestra)
CREATE TABLE IF NOT EXISTS agregadores (
    id_agregador    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nombre          VARCHAR(100) NOT NULL UNIQUE,
    id_status       UUID NOT NULL REFERENCES estatus(id_status),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMPTZ NULL
);

-- 2. ACTUALIZAR ORDEN_PEDIDO (Cabecera Maestra)
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='orden_pedido' AND column_name='total') THEN
        ALTER TABLE orden_pedido ADD COLUMN total NUMERIC(18,2) NOT NULL DEFAULT 0;
    END IF;
END $$;

-- 3. TABLA DE REFERENCIA ORDEN_AGREGADOR
CREATE TABLE IF NOT EXISTS orden_agregador (
    id_orden_agregador UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_orden_pedido    UUID NOT NULL REFERENCES orden_pedido(id_orden_pedido),
    id_agregador       UUID NOT NULL REFERENCES agregadores(id_agregador),
    referencia_externa VARCHAR(100),
    comision_agregador NUMERIC(18,2) NOT NULL DEFAULT 0,
    metadata           JSONB DEFAULT '{}',
    created_at         TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at         TIMESTAMPTZ NULL
);

-- 4. LIMPIEZA Y NORMALIZACIÓN
DO $$ 
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='orden_agregador' AND column_name='fecha') THEN
        ALTER TABLE orden_agregador DROP COLUMN fecha;
    END IF;
END $$;

-- 5. APLICAR FRAMEWORK CORE
CALL sp_core_setup_table('agregadores');
CALL sp_core_setup_table('orden_agregador');
