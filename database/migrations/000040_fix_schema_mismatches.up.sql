-- ============================================================
-- CORRECCIÓN DE DISCREPANCIAS ENTRE STORE Y ESQUEMA SQL
-- ============================================================

-- 1. ESTATUS: agregar columnas que el store usa en SELECT/INSERT/UPDATE
ALTER TABLE estatus
    ADD COLUMN IF NOT EXISTS std_tipo_estado VARCHAR(100),
    ADD COLUMN IF NOT EXISTS factor VARCHAR(100),
    ADD COLUMN IF NOT EXISTS nivel INTEGER DEFAULT 0;

-- 2. CATEGORIA: renombrar cat_nombre → nombre y agregar id_sucursal
ALTER TABLE categoria
    RENAME COLUMN cat_nombre TO nombre;

ALTER TABLE categoria
    ADD COLUMN IF NOT EXISTS id_sucursal UUID REFERENCES sucursal(id_sucursal);

-- 3. MONEDA: agregar id_sucursal
ALTER TABLE moneda
    ADD COLUMN IF NOT EXISTS id_sucursal UUID REFERENCES sucursal(id_sucursal);

-- 4. UNIDAD: el store consulta la tabla como "unidad" pero la migración la creó como "unidad_medida"
--    Renombrar + agregar id_sucursal
ALTER TABLE unidad_medida RENAME TO unidad;

ALTER TABLE unidad
    ADD COLUMN IF NOT EXISTS id_sucursal UUID REFERENCES sucursal(id_sucursal);

-- 5. CONTROL_ESTACION: renombrar columnas para alinear con el store Go
ALTER TABLE control_estacion RENAME COLUMN id_control TO id_control_estacion;
ALTER TABLE control_estacion RENAME COLUMN id_usuario TO usuario_asignado;
ALTER TABLE control_estacion RENAME COLUMN fecha_apertura TO fecha_inicio;
ALTER TABLE control_estacion RENAME COLUMN fecha_cierre TO fecha_salida;
ALTER TABLE control_estacion RENAME COLUMN monto_apertura TO fondo_base;

ALTER TABLE control_estacion
    ADD COLUMN IF NOT EXISTS id_user_pos          UUID REFERENCES usuario(id_usuario),
    ADD COLUMN IF NOT EXISTS id_periodo           UUID,
    ADD COLUMN IF NOT EXISTS fondo_retirado       NUMERIC(18,2),
    ADD COLUMN IF NOT EXISTS usuario_retiro_fondo UUID,
    ADD COLUMN IF NOT EXISTS ctrc_motivo_descuadre TEXT;

-- 6. RETIROS: renombrar FK y agregar columnas de arqueo que usa el store Go
ALTER TABLE retiros RENAME COLUMN id_control TO id_control_estacion;

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
ALTER TABLE auditoria_caja RENAME COLUMN id_control TO id_control_estacion;

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
