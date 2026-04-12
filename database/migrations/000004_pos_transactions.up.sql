-- 1. CABECERA DE FACTURA (Idempotente)
CREATE TABLE IF NOT EXISTS factura (
    id_factura      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fac_numero      VARCHAR(50) UNIQUE NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMPTZ NULL
);

-- 2. ASEGURAR COLUMNAS Y RENOMBRAR LEGACY
DO $$ 
BEGIN
    -- Asegurar id_sucursal (Nueva para optimización de reportes)
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='factura' AND column_name='id_sucursal') THEN
        ALTER TABLE factura ADD COLUMN id_sucursal UUID;
        
        -- Solo intentar poblar si la tabla de estaciones ya existe (Evita error de relación no existente)
        IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name='estaciones_pos') THEN
            EXECUTE 'UPDATE factura f SET id_sucursal = e.id_sucursal 
                     FROM estaciones_pos e WHERE f.id_estacion = e.id_estacion';
        END IF;
    END IF;

    -- Asegurar id_usuario (Mapping de id_user_pos)
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='factura' AND column_name='id_usuario') THEN
        ALTER TABLE factura ADD COLUMN id_usuario UUID;
        -- Solo mapear si la columna legacy existe
        IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='factura' AND column_name='id_user_pos') THEN
            EXECUTE 'UPDATE factura SET id_usuario = id_user_pos WHERE id_user_pos IS NOT NULL';
        END IF;
    END IF;

    -- Normalizar Nombres de Montos (Si existen los viejos, creamos los nuevos)
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='factura' AND column_name='total') THEN
        ALTER TABLE factura ADD COLUMN total NUMERIC(18,2) DEFAULT 0;
        IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='factura' AND column_name='cfac_total') THEN
            EXECUTE 'UPDATE factura SET total = cfac_total WHERE cfac_total IS NOT NULL';
        END IF;
    END IF;

    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='factura' AND column_name='subtotal') THEN
        ALTER TABLE factura ADD COLUMN subtotal NUMERIC(18,2) DEFAULT 0;
        IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='factura' AND column_name='cfac_subtotal') THEN
            EXECUTE 'UPDATE factura SET subtotal = cfac_subtotal WHERE cfac_subtotal IS NOT NULL';
        END IF;
    END IF;

    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='factura' AND column_name='impuesto') THEN
        ALTER TABLE factura ADD COLUMN impuesto NUMERIC(18,2) DEFAULT 0;
        IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='factura' AND column_name='cfac_iva') THEN
            EXECUTE 'UPDATE factura SET impuesto = cfac_iva WHERE cfac_iva IS NOT NULL';
        END IF;
    END IF;

    -- Columnas de Estado y Relaciones faltantes
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='factura' AND column_name='id_status') THEN
        ALTER TABLE factura ADD COLUMN id_status UUID;
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='factura' AND column_name='id_cliente') THEN
        ALTER TABLE factura ADD COLUMN id_cliente UUID;
    END IF;

    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='factura' AND column_name='metadata') THEN
        ALTER TABLE factura ADD COLUMN metadata JSONB DEFAULT '{}';
    END IF;
END $$;

-- 3. DETALLE DE FACTURA (Idempotente)
CREATE TABLE IF NOT EXISTS detalle_factura (
    id_detalle      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_factura      UUID NOT NULL REFERENCES factura(id_factura),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMPTZ NULL
);

DO $$ 
BEGIN
    -- Asegurar columnas detalle
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='detalle_factura' AND column_name='id_producto') THEN
        ALTER TABLE detalle_factura ADD COLUMN id_producto UUID NOT NULL;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='detalle_factura' AND column_name='cantidad') THEN
        ALTER TABLE detalle_factura ADD COLUMN cantidad NUMERIC(12,2) NOT NULL DEFAULT 0;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='detalle_factura' AND column_name='total') THEN
        ALTER TABLE detalle_factura ADD COLUMN total NUMERIC(18,2) NOT NULL DEFAULT 0;
    END IF;
END $$;

-- 4. FORMAS DE PAGO (Idempotente)
CREATE TABLE IF NOT EXISTS forma_pago_factura (
    id_pago_fac     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_factura      UUID NOT NULL REFERENCES factura(id_factura),
    metodo_pago     VARCHAR(50),
    monto           NUMERIC(18,2),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMPTZ NULL
);

-- 5. NORMALIZACIÓN DE TIEMPOS
DO $$ 
BEGIN
    ALTER TABLE factura ALTER COLUMN created_at TYPE TIMESTAMPTZ;
    ALTER TABLE factura ALTER COLUMN updated_at TYPE TIMESTAMPTZ;
    ALTER TABLE factura ALTER COLUMN deleted_at TYPE TIMESTAMPTZ;
    
    ALTER TABLE detalle_factura ALTER COLUMN created_at TYPE TIMESTAMPTZ;
    ALTER TABLE detalle_factura ALTER COLUMN deleted_at TYPE TIMESTAMPTZ;
    
    ALTER TABLE forma_pago_factura ALTER COLUMN created_at TYPE TIMESTAMPTZ;
    ALTER TABLE forma_pago_factura ALTER COLUMN deleted_at TYPE TIMESTAMPTZ;
END $$;

-- 6. ÍNDICES DE ALTO RENDIMIENTO (PostgreSQL Optimization)
-- Solo se crean si existe la columna para evitar errores
DO $$ 
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='factura' AND column_name='id_sucursal') THEN
        CREATE INDEX IF NOT EXISTS idx_factura_sucursal_fecha ON factura(id_sucursal, created_at DESC);
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='detalle_factura' AND column_name='id_producto') THEN
        CREATE INDEX IF NOT EXISTS idx_detalle_factura_prod ON detalle_factura(id_producto);
    END IF;
END $$;

-- 7. APLICAR FRAMEWORK CORE
CALL sp_core_setup_table('factura'::text);
CALL sp_core_setup_table('detalle_factura'::text);
CALL sp_core_setup_table('forma_pago_factura'::text);
