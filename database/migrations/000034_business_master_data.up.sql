-- 1. MONEDAS
CREATE TABLE IF NOT EXISTS moneda (
    id_moneda       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nombre          VARCHAR(50) NOT NULL,
    codigo_iso      VARCHAR(3) NOT NULL UNIQUE, -- USD, EUR, COP, etc.
    simbolo         VARCHAR(5) NOT NULL,
    id_status       UUID NOT NULL REFERENCES estatus(id_status),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMPTZ NULL
);

-- 2. UNIDADES DE MEDIDA
CREATE TABLE IF NOT EXISTS unidad_medida (
    id_unidad       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nombre          VARCHAR(50) NOT NULL,
    abreviatura     VARCHAR(10) NOT NULL UNIQUE, -- KG, UND, LTS, etc.
    id_status       UUID NOT NULL REFERENCES estatus(id_status),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMPTZ NULL
);

-- 3. CLIENTES
CREATE TABLE IF NOT EXISTS cliente (
    id_cliente      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nombre_completo VARCHAR(255) NOT NULL,
    tipo_documento  VARCHAR(20), -- CEDULA, RUC, NIT, etc.
    documento       VARCHAR(50) UNIQUE,
    email           VARCHAR(255),
    telefono        VARCHAR(50),
    direccion       TEXT,
    id_status       UUID NOT NULL REFERENCES estatus(id_status),
    metadata        JSONB DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMPTZ NULL
);

-- 4. PROVEEDORES
CREATE TABLE IF NOT EXISTS proveedor (
    id_proveedor    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    razon_social    VARCHAR(255) NOT NULL,
    nit_rut         VARCHAR(50) UNIQUE NOT NULL,
    contacto_nombre VARCHAR(255),
    telefono        VARCHAR(50),
    email           VARCHAR(255),
    direccion       TEXT,
    id_status       UUID NOT NULL REFERENCES estatus(id_status),
    metadata        JSONB DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMPTZ NULL
);

-- Índices de búsqueda rápida
CREATE INDEX IF NOT EXISTS idx_cliente_documento ON cliente(documento) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_proveedor_nit ON proveedor(nit_rut) WHERE deleted_at IS NULL;

-- Aplicar Framework Core
CALL sp_core_setup_table('moneda');
CALL sp_core_setup_table('unidad_medida');
CALL sp_core_setup_table('cliente');
CALL sp_core_setup_table('proveedor');
