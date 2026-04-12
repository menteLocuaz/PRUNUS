-- 1. CONFIGURACIÓN DE ESTATUS INICIALES
INSERT INTO estatus (id_status, std_descripcion, mdl_id, is_active)
VALUES 
    ('7f7b0e11-1234-4a21-9591-316279f06742', 'Activo', -1, true),
    ('892340e0-4328-491d-9102-80550bb6aac4', 'Pendiente', -1, true),
    ('62ed7d82-0c81-4511-8f02-e7fd140018d8', 'Anulado', -1, true)
ON CONFLICT (id_status) DO NOTHING;

-- 2. ENTIDADES DEMO
INSERT INTO empresa (id_empresa, nombre, rut, id_status)
VALUES ('11111111-1111-4111-a111-111111111111', 'Prunus Business Demo', '12345678-9', '7f7b0e11-1234-4a21-9591-316279f06742')
ON CONFLICT (id_empresa) DO NOTHING;

INSERT INTO sucursal (id_sucursal, id_empresa, nombre_sucursal, id_status)
VALUES ('22222222-2222-4222-a222-222222222222', '11111111-1111-4111-a111-111111111111', 'Sucursal Central', '7f7b0e11-1234-4a21-9591-316279f06742')
ON CONFLICT (id_sucursal) DO NOTHING;

-- 3. MÓDULOS DEL SISTEMA
INSERT INTO modulo (mdl_id, mdl_descripcion, abreviatura, ruta, icono, id_status, is_active, orden)
VALUES 
    (1, 'Configuración Empresa', 'EMP',  '/config/empresa',    'settings',      '7f7b0e11-1234-4a21-9591-316279f06742', true, 1),
    (2, 'Gestión Sucursales',    'SUC',  '/config/sucursales', 'store',         '7f7b0e11-1234-4a21-9591-316279f06742', true, 2),
    (3, 'Seguridad y Usuarios',  'USR',  '/config/usuarios',   'users',         '7f7b0e11-1234-4a21-9591-316279f06742', true, 3),
    (4, 'Catálogo Productos',    'PROD', '/productos',         'package',       '7f7b0e11-1234-4a21-9591-316279f06742', true, 4),
    (5, 'Ventas y POS',          'VENT', '/ventas',            'shopping-cart', '7f7b0e11-1234-4a21-9591-316279f06742', true, 5),
    (6, 'Control de Inventario', 'INV',  '/inventario',        'archive',       '7f7b0e11-1234-4a21-9591-316279f06742', true, 6),
    (7, 'Reportes y Dashboards', 'DASH', '/dashboard',         'pie-chart',     '7f7b0e11-1234-4a21-9591-316279f06742', true, 7),
    (8, 'Gastos Operativos',     'GAST', '/gastos',            'dollar-sign',   '7f7b0e11-1234-4a21-9591-316279f06742', true, 8)
ON CONFLICT (mdl_id) WHERE deleted_at IS NULL DO UPDATE SET
    mdl_descripcion = EXCLUDED.mdl_descripcion,
    ruta = EXCLUDED.ruta,
    icono = EXCLUDED.icono,
    orden = EXCLUDED.orden;

-- 4. ROL ADMINISTRADOR GLOBAL
INSERT INTO rol (id_rol, nombre_rol, id_sucursal, id_status)
VALUES ('7d7b0e11-1234-4a21-9591-316279f06742', 'Administrador Global', '22222222-2222-4222-a222-222222222222', '7f7b0e11-1234-4a21-9591-316279f06742')
ON CONFLICT (id_rol) DO NOTHING;

-- 5. PERMISOS TOTALES ADMIN
CREATE TABLE IF NOT EXISTS permiso_rol (
    id_permiso UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_rol     UUID NOT NULL REFERENCES rol(id_rol),
    id_modulo  UUID NOT NULL REFERENCES modulo(id_modulo),
    can_read   BOOLEAN DEFAULT FALSE,
    can_write  BOOLEAN DEFAULT FALSE,
    can_update BOOLEAN DEFAULT FALSE,
    can_delete BOOLEAN DEFAULT FALSE,
    CONSTRAINT uq_permiso_rol_modulo UNIQUE (id_rol, id_modulo)
);
