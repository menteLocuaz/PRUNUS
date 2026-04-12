-- 1. TABLA DE ACCESOS MULTI-SUCURSAL PARA USUARIOS
CREATE TABLE IF NOT EXISTS usuario_sucursal_acceso (
    id_usuario    UUID NOT NULL REFERENCES usuario(id_usuario) ON DELETE CASCADE,
    id_sucursal   UUID NOT NULL REFERENCES sucursal(id_sucursal) ON DELETE CASCADE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id_usuario, id_sucursal)
);

-- 2. ÍNDICES DE RENDIMIENTO
CREATE INDEX IF NOT EXISTS idx_usa_id_usuario ON usuario_sucursal_acceso(id_usuario);
CREATE INDEX IF NOT EXISTS idx_usa_id_sucursal ON usuario_sucursal_acceso(id_sucursal);

-- 3. APLICAR FRAMEWORK CORE (Auditoría JSONB)
CALL sp_core_setup_table('usuario_sucursal_acceso');
