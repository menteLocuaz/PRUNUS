-- 1. LOG DE SISTEMA
CREATE TABLE IF NOT EXISTS log_sistema (
    id_log      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_usuario  UUID REFERENCES usuario(id_usuario),
    accion      VARCHAR(100) NOT NULL,
    tabla       VARCHAR(100),
    registro_id TEXT,
    ip          VARCHAR(45),
    fecha       TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_log_sistema_fecha     ON log_sistema(fecha DESC);
CREATE INDEX IF NOT EXISTS idx_log_sistema_id_usuario ON log_sistema(id_usuario);

-- 2. CAJA (Punto de Venta físico)
CREATE TABLE IF NOT EXISTS caja (
    id_caja     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nombre      VARCHAR(100) NOT NULL,
    id_sucursal UUID NOT NULL REFERENCES sucursal(id_sucursal),
    estado      SMALLINT NOT NULL DEFAULT 1, -- 1: Activa, 0: Inactiva
    created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMPTZ NULL
);

-- 3. SESIÓN DE CAJA (Turno del cajero)
CREATE TABLE IF NOT EXISTS sesion_caja (
    id_sesion      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_caja        UUID NOT NULL REFERENCES caja(id_caja),
    id_usuario     UUID NOT NULL REFERENCES usuario(id_usuario),
    monto_apertura NUMERIC(18,2) NOT NULL DEFAULT 0,
    monto_cierre   NUMERIC(18,2),
    fecha_apertura TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    fecha_cierre   TIMESTAMPTZ,
    estado         VARCHAR(20) NOT NULL DEFAULT 'ABIERTA', -- ABIERTA, CERRADA
    created_at     TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at     TIMESTAMPTZ NULL
);

CREATE INDEX IF NOT EXISTS idx_sesion_caja_id_caja    ON sesion_caja(id_caja)    WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_sesion_caja_id_usuario ON sesion_caja(id_usuario) WHERE deleted_at IS NULL;

-- 4. MOVIMIENTO DE CAJA (Ingresos y egresos de efectivo)
CREATE TABLE IF NOT EXISTS movimiento_caja (
    id_movimiento UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_sesion     UUID NOT NULL REFERENCES sesion_caja(id_sesion),
    tipo          VARCHAR(20) NOT NULL, -- INGRESO, EGRESO
    monto         NUMERIC(18,2) NOT NULL,
    motivo        TEXT,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at    TIMESTAMPTZ NULL
);

CREATE INDEX IF NOT EXISTS idx_movimiento_caja_id_sesion ON movimiento_caja(id_sesion) WHERE deleted_at IS NULL;

-- 5. DETALLE ORDEN DE COMPRA
CREATE TABLE IF NOT EXISTS detalle_orden_compra (
    id_detalle_compra  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_orden_compra    UUID NOT NULL REFERENCES orden_compra(id_orden_compra),
    id_producto        UUID NOT NULL REFERENCES producto(id_producto),
    cantidad_pedida    NUMERIC(18,4) NOT NULL,
    cantidad_recibida  NUMERIC(18,4) NOT NULL DEFAULT 0,
    precio_unitario    NUMERIC(18,2) NOT NULL,
    impuesto           NUMERIC(18,2) NOT NULL DEFAULT 0,
    total              NUMERIC(18,2) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_detalle_orden_compra_orden   ON detalle_orden_compra(id_orden_compra);
CREATE INDEX IF NOT EXISTS idx_detalle_orden_compra_producto ON detalle_orden_compra(id_producto);

-- 6. COLUMNAS FALTANTES EN ORDEN_COMPRA
ALTER TABLE orden_compra
    ADD COLUMN IF NOT EXISTS numero_orden      VARCHAR(50),
    ADD COLUMN IF NOT EXISTS id_moneda         UUID REFERENCES moneda(id_moneda),
    ADD COLUMN IF NOT EXISTS subtotal          NUMERIC(18,2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS impuesto          NUMERIC(18,2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS observaciones     TEXT,
    ADD COLUMN IF NOT EXISTS fecha_vencimiento TIMESTAMPTZ;

-- 7. CANAL DE IMPRESIÓN (Configuración POS)
CREATE TABLE IF NOT EXISTS canal_impresion (
    id_canal_impresion UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    descripcion        VARCHAR(100) NOT NULL,
    cdn_id             INTEGER NOT NULL,       -- ID de cadena/empresa legado
    id_status          UUID NOT NULL REFERENCES estatus(id_status),
    created_at         TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at         TIMESTAMPTZ NULL
);

-- 8. IMPRESORA (Configuración POS)
CREATE TABLE IF NOT EXISTS impresora (
    id_impresora UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nombre       VARCHAR(100) NOT NULL,
    rst_id       INTEGER NOT NULL,             -- ID de restaurante/sucursal legado
    id_status    UUID NOT NULL REFERENCES estatus(id_status),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at   TIMESTAMPTZ NULL
);

-- 9. PUERTOS (Configuración POS)
CREATE TABLE IF NOT EXISTS puertos (
    id_puertos  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    descripcion VARCHAR(20) NOT NULL,          -- Ej: COM1, COM2, USB1
    id_status   UUID NOT NULL REFERENCES estatus(id_status),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMPTZ NULL
);
