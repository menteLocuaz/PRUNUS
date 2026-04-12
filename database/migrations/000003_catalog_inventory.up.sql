-- 1. CATEGORÍA
CREATE TABLE IF NOT EXISTS categoria (
    id_categoria    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cat_nombre      VARCHAR(150) NOT NULL,
    id_status       UUID NOT NULL REFERENCES estatus(id_status),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMPTZ NULL
);

-- 2. PRODUCTO CON BÚSQUEDA VECTORIAL
CREATE TABLE IF NOT EXISTS producto (
    id_producto     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_categoria    UUID NOT NULL REFERENCES categoria(id_categoria),
    pro_nombre      VARCHAR(255) NOT NULL,
    pro_descripcion TEXT,
    pro_codigo      VARCHAR(100) UNIQUE,
    precio_venta    NUMERIC(18,2) DEFAULT 0,
    precio_compra   NUMERIC(18,2) DEFAULT 0,
    id_status       UUID NOT NULL REFERENCES estatus(id_status),
    metadata        JSONB DEFAULT '{}',
    search_vector   tsvector,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMPTZ NULL
);

-- 3. INVENTARIO
CREATE TABLE IF NOT EXISTS inventario (
    id_inventario   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_producto     UUID NOT NULL REFERENCES producto(id_producto),
    id_sucursal     UUID NOT NULL REFERENCES sucursal(id_sucursal),
    stock_actual    NUMERIC(12,2) DEFAULT 0,
    stock_minimo    NUMERIC(12,2) DEFAULT 0,
    ubicacion       VARCHAR(100),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMPTZ NULL,
    UNIQUE(id_producto, id_sucursal)
);

-- NORMALIZACIÓN
DO $$ 
BEGIN
    -- Categoría
    ALTER TABLE categoria ALTER COLUMN created_at TYPE TIMESTAMPTZ;
    ALTER TABLE categoria ALTER COLUMN updated_at TYPE TIMESTAMPTZ;
    ALTER TABLE categoria ALTER COLUMN deleted_at TYPE TIMESTAMPTZ;

    -- Producto
    ALTER TABLE producto ALTER COLUMN created_at TYPE TIMESTAMPTZ;
    ALTER TABLE producto ALTER COLUMN updated_at TYPE TIMESTAMPTZ;
    ALTER TABLE producto ALTER COLUMN deleted_at TYPE TIMESTAMPTZ;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='producto' AND column_name='search_vector') THEN
        ALTER TABLE producto ADD COLUMN search_vector tsvector;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='producto' AND column_name='metadata') THEN
        ALTER TABLE producto ADD COLUMN metadata JSONB DEFAULT '{}';
    END IF;

    -- Inventario
    ALTER TABLE inventario ALTER COLUMN created_at TYPE TIMESTAMPTZ;
    ALTER TABLE inventario ALTER COLUMN updated_at TYPE TIMESTAMPTZ;
    ALTER TABLE inventario ALTER COLUMN deleted_at TYPE TIMESTAMPTZ;
END $$;

-- Trigger e Índices
CREATE INDEX IF NOT EXISTS idx_producto_search ON producto USING gin(search_vector);
CREATE INDEX IF NOT EXISTS idx_producto_metadata ON producto USING gin(metadata);

CREATE OR REPLACE FUNCTION fn_producto_search_sync() RETURNS TRIGGER AS $$
BEGIN
    NEW.search_vector := to_tsvector('spanish', COALESCE(NEW.pro_nombre, '') || ' ' || COALESCE(NEW.pro_descripcion, ''));
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS tr_producto_search_sync ON producto;
CREATE TRIGGER tr_producto_search_sync BEFORE INSERT OR UPDATE ON producto
FOR EACH ROW EXECUTE FUNCTION fn_producto_search_sync();

-- Aplicar Framework Core
CALL sp_core_setup_table('categoria');
CALL sp_core_setup_table('producto');
CALL sp_core_setup_table('inventario');
