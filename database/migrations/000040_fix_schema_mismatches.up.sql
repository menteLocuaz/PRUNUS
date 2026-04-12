-- ============================================================
-- CORRECCIÓN DE DISCREPANCIAS ENTRE STORE Y ESQUEMA SQL
-- ============================================================

-- 1. ESTATUS: agregar columna std_tipo_estado que el store usa en SELECT/INSERT/UPDATE
ALTER TABLE estatus
    ADD COLUMN IF NOT EXISTS std_tipo_estado VARCHAR(100);

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

-- 5. FUNCIÓN fn_get_gastos_mensuales usada por el dashboard
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
