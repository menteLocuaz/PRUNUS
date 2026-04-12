-- 1. IMPUESTOS
CREATE TABLE IF NOT EXISTS impuesto (
    id_impuesto     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nombre          VARCHAR(50) NOT NULL,
    porcentaje      NUMERIC(5,2) NOT NULL DEFAULT 0,
    id_status       UUID NOT NULL REFERENCES estatus(id_status),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMPTZ NULL
);

-- 2. FORMAS DE PAGO
CREATE TABLE IF NOT EXISTS forma_pago (
    id_forma_pago   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nombre          VARCHAR(100) NOT NULL UNIQUE,
    requiere_ref    BOOLEAN DEFAULT FALSE,
    id_status       UUID NOT NULL REFERENCES estatus(id_status),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMPTZ NULL
);

-- 3. CONTROL DE CAJA / ESTACIÓN
CREATE TABLE IF NOT EXISTS control_estacion (
    id_control      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_estacion     UUID NOT NULL REFERENCES estaciones_pos(id_estacion),
    id_usuario      UUID NOT NULL REFERENCES usuario(id_usuario),
    fecha_apertura  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    fecha_cierre    TIMESTAMPTZ,
    monto_apertura  NUMERIC(18,2) NOT NULL DEFAULT 0,
    monto_cierre    NUMERIC(18,2),
    id_status       UUID NOT NULL REFERENCES estatus(id_status),
    metadata        JSONB DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMPTZ NULL
);

-- 4. RETIROS DE CAJA
CREATE TABLE IF NOT EXISTS retiros (
    id_retiro       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_control      UUID NOT NULL REFERENCES control_estacion(id_control),
    id_usuario      UUID NOT NULL REFERENCES usuario(id_usuario),
    monto           NUMERIC(18,2) NOT NULL,
    motivo          TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMPTZ NULL
);

-- 5. MOTIVOS DE ANULACIÓN
CREATE TABLE IF NOT EXISTS motivo_anulacion (
    id_motivo       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    descripcion     VARCHAR(255) NOT NULL,
    id_status       UUID NOT NULL REFERENCES estatus(id_status),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMPTZ NULL
);

-- 6. ORDEN DE COMPRA (Vinculación con Proveedores)
CREATE TABLE IF NOT EXISTS orden_compra (
    id_orden_compra UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_proveedor    UUID NOT NULL REFERENCES proveedor(id_proveedor),
    id_sucursal     UUID NOT NULL REFERENCES sucursal(id_sucursal),
    id_usuario      UUID NOT NULL REFERENCES usuario(id_usuario),
    total           NUMERIC(18,2) NOT NULL DEFAULT 0,
    fecha_emision   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    fecha_recepcion TIMESTAMPTZ,
    id_status       UUID NOT NULL REFERENCES estatus(id_status),
    metadata        JSONB DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMPTZ NULL
);

-- Aplicar Framework Core
CALL sp_core_setup_table('impuesto');
CALL sp_core_setup_table('forma_pago');
CALL sp_core_setup_table('control_estacion');
CALL sp_core_setup_table('retiros');
CALL sp_core_setup_table('motivo_anulacion');
CALL sp_core_setup_table('orden_compra');
