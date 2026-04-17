-- ============================================================
-- CORRECCIÓN DE DISCREPANCIAS ENTRE STORE Y ESQUEMA SQL
-- ============================================================

-- 1. ESTATUS: agregar columnas que el store usa en SELECT/INSERT/UPDATE
ALTER TABLE estatus
    ADD COLUMN IF NOT EXISTS std_tipo_estado VARCHAR(100),
    ADD COLUMN IF NOT EXISTS factor VARCHAR(100),
    ADD COLUMN IF NOT EXISTS nivel INTEGER DEFAULT 0;

-- 2. CATEGORIA: renombrar cat_nombre → nombre y agregar id_sucursal
DO $$ 
BEGIN 
    IF EXISTS (SELECT 1 FROM information_schema.columns 
               WHERE table_name = 'categoria' AND column_name = 'cat_nombre') THEN
        ALTER TABLE categoria RENAME COLUMN cat_nombre TO nombre;
    END IF;
END $$;

ALTER TABLE categoria
    ADD COLUMN IF NOT EXISTS id_sucursal UUID REFERENCES sucursal(id_sucursal);

-- 3. MONEDA: agregar id_sucursal
ALTER TABLE moneda
    ADD COLUMN IF NOT EXISTS id_sucursal UUID REFERENCES sucursal(id_sucursal);

-- 4. UNIDAD: el store consulta la tabla como "unidad" pero la migración la creó como "unidad_medida"
--    Renombrar + agregar id_sucursal
DO $$ 
BEGIN 
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'unidad_medida') THEN
        ALTER TABLE unidad_medida RENAME TO unidad;
    END IF;
END $$;

ALTER TABLE unidad
    ADD COLUMN IF NOT EXISTS id_sucursal UUID REFERENCES sucursal(id_sucursal);

-- 5. CONTROL_ESTACION: renombrar columnas para alinear con el store Go
DO $$ 
BEGIN 
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'control_estacion' AND column_name = 'id_control') THEN
        ALTER TABLE control_estacion RENAME COLUMN id_control TO id_control_estacion;
    END IF;
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'control_estacion' AND column_name = 'id_usuario') THEN
        ALTER TABLE control_estacion RENAME COLUMN id_usuario TO usuario_asignado;
    END IF;
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'control_estacion' AND column_name = 'fecha_apertura') THEN
        ALTER TABLE control_estacion RENAME COLUMN fecha_apertura TO fecha_inicio;
    END IF;
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'control_estacion' AND column_name = 'fecha_cierre') THEN
        ALTER TABLE control_estacion RENAME COLUMN fecha_cierre TO fecha_salida;
    END IF;
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'control_estacion' AND column_name = 'monto_apertura') THEN
        ALTER TABLE control_estacion RENAME COLUMN monto_apertura TO fondo_base;
    END IF;
END $$;

ALTER TABLE control_estacion
    ADD COLUMN IF NOT EXISTS id_user_pos          UUID REFERENCES usuario(id_usuario),
    ADD COLUMN IF NOT EXISTS id_periodo           UUID,
    ADD COLUMN IF NOT EXISTS fondo_retirado       NUMERIC(18,2),
    ADD COLUMN IF NOT EXISTS usuario_retiro_fondo UUID,
    ADD COLUMN IF NOT EXISTS ctrc_motivo_descuadre TEXT;

-- 6. RETIROS: renombrar FK y agregar columnas de arqueo que usa el store Go
DO $$ 
BEGIN 
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'retiros' AND column_name = 'id_control') THEN
        ALTER TABLE retiros RENAME COLUMN id_control TO id_control_estacion;
    END IF;
END $$;

ALTER TABLE retiros
    ADD COLUMN IF NOT EXISTS arc_valor         NUMERIC(18,2),
    ADD COLUMN IF NOT EXISTS usuario_inicia    UUID,
    ADD COLUMN IF NOT EXISTS usuario_finaliza  UUID,
    ADD COLUMN IF NOT EXISTS fecha_inicio      TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS fecha_finaliza    TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS retiro_valor      NUMERIC(18,2),
    ADD COLUMN IF NOT EXISTS diferencia_valor  NUMERIC(18,2),
    ADD COLUMN IF NOT EXISTS pos_calculado     NUMERIC(18,2),
    ADD COLUMN IF NOT EXISTS id_status         UUID REFERENCES estatus(id_status),
    ADD COLUMN IF NOT EXISTS id_forma_pago     UUID REFERENCES forma_pago(id_forma_pago),
    ADD COLUMN IF NOT EXISTS tpenv_id          INT DEFAULT -1;

-- 7. AUDITORIA_CAJA: renombrar FK para alinear con el store Go
DO $$ 
BEGIN 
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'auditoria_caja' AND column_name = 'id_control') THEN
        ALTER TABLE auditoria_caja RENAME COLUMN id_control TO id_control_estacion;
    END IF;
END $$;

-- 8. FUNCIÓN fn_get_gastos_mensuales usada por el dashboard
CREATE OR REPLACE FUNCTION fn_get_gastos_mensuales(
    p_id_sucursal UUID,
    p_fecha       DATE
) RETURNS NUMERIC AS $$
    SELECT COALESCE(
        SUM(monto),
        0
    )
    FROM gasto_operativo
    WHERE id_sucursal = p_id_sucursal
      AND date_trunc('month', fecha_gasto) = date_trunc('month', p_fecha::TIMESTAMPTZ)
      AND deleted_at IS NULL;
$$ LANGUAGE SQL STABLE;
