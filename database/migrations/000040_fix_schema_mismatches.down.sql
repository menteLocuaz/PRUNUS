DROP FUNCTION IF EXISTS fn_get_gastos_mensuales(UUID, DATE);

ALTER TABLE unidad DROP COLUMN IF EXISTS id_sucursal;
ALTER TABLE unidad RENAME TO unidad_medida;

ALTER TABLE moneda DROP COLUMN IF EXISTS id_sucursal;

ALTER TABLE categoria DROP COLUMN IF EXISTS id_sucursal;
ALTER TABLE categoria RENAME COLUMN nombre TO cat_nombre;

ALTER TABLE estatus
    DROP COLUMN IF EXISTS std_tipo_estado,
    DROP COLUMN IF EXISTS factor,
    DROP COLUMN IF EXISTS nivel;

ALTER TABLE auditoria_caja RENAME COLUMN id_control_estacion TO id_control;

ALTER TABLE retiros DROP COLUMN IF EXISTS tpenv_id;
ALTER TABLE retiros DROP COLUMN IF EXISTS id_forma_pago;
ALTER TABLE retiros DROP COLUMN IF EXISTS id_status;
ALTER TABLE retiros DROP COLUMN IF EXISTS pos_calculado;
ALTER TABLE retiros DROP COLUMN IF EXISTS diferencia_valor;
ALTER TABLE retiros DROP COLUMN IF EXISTS retiro_valor;
ALTER TABLE retiros DROP COLUMN IF EXISTS fecha_finaliza;
ALTER TABLE retiros DROP COLUMN IF EXISTS fecha_inicio;
ALTER TABLE retiros DROP COLUMN IF EXISTS usuario_finaliza;
ALTER TABLE retiros DROP COLUMN IF EXISTS usuario_inicia;
ALTER TABLE retiros DROP COLUMN IF EXISTS arc_valor;
ALTER TABLE retiros RENAME COLUMN id_control_estacion TO id_control;

ALTER TABLE control_estacion DROP COLUMN IF EXISTS ctrc_motivo_descuadre;
ALTER TABLE control_estacion DROP COLUMN IF EXISTS usuario_retiro_fondo;
ALTER TABLE control_estacion DROP COLUMN IF EXISTS fondo_retirado;
ALTER TABLE control_estacion DROP COLUMN IF EXISTS id_periodo;
ALTER TABLE control_estacion DROP COLUMN IF EXISTS id_user_pos;
ALTER TABLE control_estacion RENAME COLUMN fondo_base TO monto_apertura;
ALTER TABLE control_estacion RENAME COLUMN fecha_salida TO fecha_cierre;
ALTER TABLE control_estacion RENAME COLUMN fecha_inicio TO fecha_apertura;
ALTER TABLE control_estacion RENAME COLUMN usuario_asignado TO id_usuario;
ALTER TABLE control_estacion RENAME COLUMN id_control_estacion TO id_control;
