-- 1. AÑADIR IDENTIFICADORES A PRODUCTOS
DO $$ 
BEGIN
    -- Código de Barras (EAN/UPC)
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='producto' AND column_name='codigo_barras') THEN
        ALTER TABLE producto ADD COLUMN codigo_barras VARCHAR(50);
    END IF;

    -- SKU (Stock Keeping Unit - Código Interno)
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='producto' AND column_name='sku') THEN
        ALTER TABLE producto ADD COLUMN sku VARCHAR(50);
    END IF;
END $$;

-- 2. ÍNDICES DE BÚSQUEDA RÁPIDA (PostgreSQL Optimization)
-- Creamos índices parciales para ignorar registros eliminados, reduciendo el tamaño del índice.
CREATE INDEX IF NOT EXISTS idx_producto_codigo_barras ON producto(codigo_barras) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_producto_sku ON producto(sku) WHERE deleted_at IS NULL;

-- 3. NORMALIZACIÓN DE INVENTARIO (Ubicación física)
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='inventario' AND column_name='ubicacion') THEN
        ALTER TABLE inventario ADD COLUMN ubicacion VARCHAR(100);
    END IF;
END $$;
