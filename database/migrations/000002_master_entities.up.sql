-- 1. ESTATUS
CREATE TABLE IF NOT EXISTS estatus (
    id_status       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    std_descripcion VARCHAR(100) NOT NULL,
    mdl_id          INTEGER NOT NULL,
    is_active       BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMPTZ NULL
);

-- 2. EMPRESA Y SUCURSAL
CREATE TABLE IF NOT EXISTS empresa (
    id_empresa  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nombre      VARCHAR(255) NOT NULL,
    rut         VARCHAR(20) UNIQUE NOT NULL,
    id_status   UUID NOT NULL REFERENCES estatus(id_status),
    metadata    JSONB DEFAULT '{}',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMPTZ NULL
);

CREATE TABLE IF NOT EXISTS sucursal (
    id_sucursal     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_empresa      UUID NOT NULL REFERENCES empresa(id_empresa),
    nombre_sucursal VARCHAR(255) NOT NULL,
    direccion       TEXT,
    telefono        VARCHAR(50),
    id_status       UUID NOT NULL REFERENCES estatus(id_status),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMPTZ NULL
);

-- 3. ROLES Y USUARIOS
CREATE TABLE IF NOT EXISTS rol (
    id_rol      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nombre_rol  VARCHAR(100) NOT NULL,
    id_sucursal UUID NOT NULL REFERENCES sucursal(id_sucursal),
    id_status   UUID NOT NULL REFERENCES estatus(id_status),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMPTZ NULL,
    UNIQUE(nombre_rol, id_sucursal)
);

CREATE TABLE IF NOT EXISTS usuario (
    id_usuario      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_sucursal     UUID NOT NULL REFERENCES sucursal(id_sucursal),
    id_rol          UUID NOT NULL REFERENCES rol(id_rol),
    username        VARCHAR(100) UNIQUE NOT NULL,
    email           VARCHAR(255) UNIQUE NOT NULL,
    password        TEXT NOT NULL,
    usu_nombre      VARCHAR(255) NOT NULL,
    usu_dni         VARCHAR(50),
    usu_pin_pos     VARCHAR(10),
    en_turno        BOOLEAN DEFAULT FALSE,
    id_status       UUID NOT NULL REFERENCES estatus(id_status),
    metadata        JSONB DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMPTZ NULL
);

-- NORMALIZACIÓN: Asegurar tipos TIMESTAMPTZ y columnas nuevas si la tabla ya existía
DO $$ 
BEGIN
    -- Estatus
    ALTER TABLE estatus ALTER COLUMN created_at TYPE TIMESTAMPTZ;
    ALTER TABLE estatus ALTER COLUMN updated_at TYPE TIMESTAMPTZ;
    ALTER TABLE estatus ALTER COLUMN deleted_at TYPE TIMESTAMPTZ;
    
    -- Empresa
    ALTER TABLE empresa ALTER COLUMN created_at TYPE TIMESTAMPTZ;
    ALTER TABLE empresa ALTER COLUMN updated_at TYPE TIMESTAMPTZ;
    ALTER TABLE empresa ALTER COLUMN deleted_at TYPE TIMESTAMPTZ;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='empresa' AND column_name='metadata') THEN
        ALTER TABLE empresa ADD COLUMN metadata JSONB DEFAULT '{}';
    END IF;

    -- Sucursal
    ALTER TABLE sucursal ALTER COLUMN created_at TYPE TIMESTAMPTZ;
    ALTER TABLE sucursal ALTER COLUMN updated_at TYPE TIMESTAMPTZ;
    ALTER TABLE sucursal ALTER COLUMN deleted_at TYPE TIMESTAMPTZ;

    -- Usuario
    ALTER TABLE usuario ALTER COLUMN created_at TYPE TIMESTAMPTZ;
    ALTER TABLE usuario ALTER COLUMN updated_at TYPE TIMESTAMPTZ;
    ALTER TABLE usuario ALTER COLUMN deleted_at TYPE TIMESTAMPTZ;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='usuario' AND column_name='en_turno') THEN
        ALTER TABLE usuario ADD COLUMN en_turno BOOLEAN DEFAULT FALSE;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='usuario' AND column_name='metadata') THEN
        ALTER TABLE usuario ADD COLUMN metadata JSONB DEFAULT '{}';
    END IF;
END $$;

-- 4. ORDEN DE PEDIDO (Cabecera Maestra)
CREATE TABLE IF NOT EXISTS orden_pedido (
    id_orden_pedido    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    odp_fecha_creacion TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    odp_observacion    TEXT,
    id_user_pos        UUID NOT NULL REFERENCES usuario(id_usuario),
    id_periodo         UUID, -- Se vinculará formalmente en migraciones posteriores (000037)
    id_estacion        UUID, -- Se vinculará formalmente en migraciones posteriores (000011)
    id_status          UUID NOT NULL REFERENCES estatus(id_status),
    direccion          TEXT,
    canal              VARCHAR(50),
    odp_total          NUMERIC(18,2) NOT NULL DEFAULT 0,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at         TIMESTAMPTZ NULL
);

-- Aplicar Framework Core
CALL sp_core_setup_table('estatus'::text);
CALL sp_core_setup_table('empresa'::text);
CALL sp_core_setup_table('sucursal'::text);
CALL sp_core_setup_table('rol'::text);
CALL sp_core_setup_table('usuario'::text);
CALL sp_core_setup_table('orden_pedido'::text);
