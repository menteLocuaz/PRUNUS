DROP FUNCTION IF EXISTS fn_get_gastos_mensuales(UUID, DATE);

ALTER TABLE unidad DROP COLUMN IF EXISTS id_sucursal;
ALTER TABLE unidad RENAME TO unidad_medida;

ALTER TABLE moneda DROP COLUMN IF EXISTS id_sucursal;

ALTER TABLE categoria DROP COLUMN IF EXISTS id_sucursal;
ALTER TABLE categoria RENAME COLUMN nombre TO cat_nombre;

ALTER TABLE estatus DROP COLUMN IF EXISTS std_tipo_estado;
