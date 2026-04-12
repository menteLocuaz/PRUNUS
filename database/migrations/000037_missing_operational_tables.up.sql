-- 1. PERIODOS CONTABLES
CREATE TABLE IF NOT EXISTS periodo (
    id_periodo      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nombre          VARCHAR(100) NOT NULL,
    fecha_inicio    DATE NOT NULL,
    fecha_fin       DATE NOT NULL,
    is_cerrado      BOOLEAN DEFAULT FALSE,
    id_status       UUID NOT NULL REFERENCES estatus(id_status),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMPTZ NULL
);

-- 2. DISPOSITIVOS POS (Hardware)
CREATE TABLE IF NOT EXISTS dispositivos_pos (
    id_dispositivo  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_estacion     UUID NOT NULL REFERENCES estaciones_pos(id_estacion),
    nombre          VARCHAR(100) NOT NULL,
    tipo_dispositivo VARCHAR(50), -- IMPRESORA, SCANNER, CAJON, etc.
    configuracion   JSONB DEFAULT '{}',
    id_status       UUID NOT NULL REFERENCES estatus(id_status),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMPTZ NULL
);

-- 3. AUDITORÍA DE CAJA (Event Log para Arqueos)
CREATE TABLE IF NOT EXISTS auditoria_caja (
    id_audit_caja   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_control      UUID NOT NULL REFERENCES control_estacion(id_control),
    id_usuario      UUID NOT NULL REFERENCES usuario(id_usuario),
    tipo_movimiento VARCHAR(50) NOT NULL, -- APERTURA, CIERRE, AJUSTE, SOBRANTE, FALTANTE
    valor           NUMERIC(18,2) NOT NULL,
    descripcion     TEXT,
    fecha           TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMPTZ NULL
);

-- Aplicar Framework Core
CALL sp_core_setup_table('periodo');
CALL sp_core_setup_table('dispositivos_pos');
CALL sp_core_setup_table('auditoria_caja');
