-- 1. ASEGURAR COLUMNAS OPERATIVAS EN INVENTARIO
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='inventario' AND column_name='precio_venta') THEN
        ALTER TABLE inventario ADD COLUMN precio_venta NUMERIC(18,2) DEFAULT 0;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='inventario' AND column_name='precio_compra') THEN
        ALTER TABLE inventario ADD COLUMN precio_compra NUMERIC(18,2) DEFAULT 0;
    END IF;
    
    -- Unicidad obligatoria para evitar duplicados en la misma sede
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'uq_inventario_producto_sucursal') THEN
        ALTER TABLE inventario ADD CONSTRAINT uq_inventario_producto_sucursal UNIQUE (id_producto, id_sucursal);
    END IF;
END $$;

-- 2. MIGRACIÓN DE DATOS (Legacy to Normalized)
-- Solo se ejecuta si la tabla producto aún tiene los campos operativos.
DO $$ 
DECLARE
    v_query TEXT;
    v_has_sucursal BOOLEAN;
    v_has_stock BOOLEAN;
BEGIN 
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='producto' AND column_name='precio_venta') THEN
        -- Verificar qué columnas existen para construir el query dinámico
        SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='producto' AND column_name='id_sucursal') INTO v_has_sucursal;
        SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='producto' AND column_name='stock') INTO v_has_stock;

        v_query := 'INSERT INTO inventario (id_producto, id_sucursal, stock_actual, precio_compra, precio_venta, created_at, updated_at) SELECT id_producto, ';
        
        -- Selección dinámica de sucursal
        IF v_has_sucursal THEN
            v_query := v_query || 'COALESCE(id_sucursal, (SELECT id_sucursal FROM sucursal LIMIT 1)), ';
        ELSE
            v_query := v_query || '(SELECT id_sucursal FROM sucursal LIMIT 1), ';
        END IF;

        -- Selección dinámica de stock
        IF v_has_stock THEN
            v_query := v_query || 'COALESCE(stock, 0), ';
        ELSE
            v_query := v_query || '0, ';
        END IF;

        v_query := v_query || 'COALESCE(precio_compra, 0), COALESCE(precio_venta, 0), created_at, updated_at FROM producto ON CONFLICT (id_producto, id_sucursal) DO UPDATE SET precio_compra = EXCLUDED.precio_compra, precio_venta = EXCLUDED.precio_venta, updated_at = CURRENT_TIMESTAMP;';

        -- Ejecutar la migración de datos
        EXECUTE v_query;

        -- 3. LIMPIEZA DE COLUMNAS REDUNDANTES EN PRODUCTO (SQL Optimization)
        ALTER TABLE producto DROP COLUMN IF EXISTS precio_compra;
        ALTER TABLE producto DROP COLUMN IF EXISTS precio_venta;
        ALTER TABLE producto DROP COLUMN IF EXISTS stock;
        ALTER TABLE producto DROP COLUMN IF EXISTS id_sucursal;
        
        DROP INDEX IF EXISTS idx_producto_id_sucursal;
    END IF;
END $$;

-- 4. DOCUMENTACIÓN Y SEMÁNTICA
COMMENT ON TABLE producto IS 'Catálogo maestro de productos: Definición global (Nombre, Código, IVA, etc.)';
COMMENT ON TABLE inventario IS 'Gestión operativa: Precios, Costos y Stock específico por Sucursal';
