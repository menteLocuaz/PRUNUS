-- 000049_seed_master_estatus.up.sql
-- Carga de estados maestros del sistema para todos los módulos operativos

INSERT INTO estatus (id_status, std_descripcion, std_tipo_estado, mdl_id) VALUES
-- Empresa
('fc273a6a-ab7b-4453-a560-ac62fa64348b', 'Activa', 'GENERAL', 1),
('b4d0544d-1778-4560-a170-a681ab3399bd', 'Suspendida', 'RESTRICCION', 1),
('c3e12abc-0011-4abc-b123-000000000001', 'Inactiva', 'GENERAL', 1),

-- Sucursal
('6cf06fbe-b21c-46e3-a34b-b24f5167cd9a', 'Abierta', 'OPERATIVO', 2),
('34be5a4c-ab4c-4afd-9a3c-98dfeb500fbc', 'Cerrada', 'OPERATIVO', 2),
('dba9516b-257e-46c9-8372-e1bb58e999f1', 'En Mantenimiento', 'TECNICO', 2),
('dba9516b-257e-46c9-8372-e1bb58e000f2', 'Suspendida Temporalmente', 'RESTRICCION', 2),

-- Usuario
('3a99d245-b34f-48a5-ac08-a5a010c5822f', 'Activo', 'ACCESO', 3),
('3ae0a60b-8abc-43ea-95da-c1300245a327', 'Pendiente de Activación', 'REGISTRO', 3),
('f210ac4b-81dd-4e52-a247-79a93b8b1a68', 'Bloqueado', 'SEGURIDAD', 3),
('f210ac4b-81dd-4e52-a247-79a93b8b0000', 'Inactivo', 'ACCESO', 3),
('f210ac4b-81dd-4e52-a247-79a93b8b0001', 'Vacaciones', 'RRHH', 3),

-- Producto
('31f4e127-e7e1-414d-aaef-6e92e4c5d970', 'Disponible', 'STOCK', 4),
('073f4513-a88e-4278-af95-bd9cde61bdbd', 'Agotado', 'STOCK', 4),
('073f4513-a88e-4278-af95-bd9cde610001', 'Stock Bajo', 'STOCK', 4),
('073f4513-a88e-4278-af95-bd9cde610002', 'En Oferta', 'COMERCIAL', 4),
('073f4513-a88e-4278-af95-bd9cde610003', 'Próximo a Vencer', 'COMERCIAL', 4),
('f566215a-a391-4904-b736-26c7f258aa00', 'Descontinuado', 'COMERCIAL', 4),
('f566215a-a391-4904-b736-26c7f258aa01', 'Retirado del Mercado', 'COMERCIAL', 4),

-- Venta
('892340e0-4328-491d-9102-80550bb6aac4', 'Pendiente de Pago', 'TRANSACCION', 5),
('0f447fd7-9849-4a68-b82f-c69297e7a924', 'Pagada', 'TRANSACCION', 5),
('0f447fd7-9849-4a68-b82f-c69297e70001', 'Pago Parcial', 'TRANSACCION', 5),
('62ed7d82-0c81-4511-8f02-e7fd140018d8', 'Anulada', 'TRANSACCION', 5),
('62ed7d82-0c81-4511-8f02-e7fd14000001', 'Devuelta', 'TRANSACCION', 5),
('62ed7d82-0c81-4511-8f02-e7fd14000002', 'En Proceso de Devolución', 'TRANSACCION', 5),

-- Compra / Abastecimiento
('0cd9aa6e-5768-45d2-a66d-12758a3bd0cc', 'Solicitada', 'FLUJO', 6),
('0cd9aa6e-5768-45d2-a66d-12758a3b0001', 'Aprobada', 'FLUJO', 6),
('0cd9aa6e-5768-45d2-a66d-12758a3b0002', 'En Tránsito', 'FLUJO', 6),
('00363491-8508-4220-9661-e99f05b0d545', 'Recibida Parcialmente', 'FLUJO', 6),
('00363491-8508-4220-9661-e99f05b00001', 'Recibida Completa', 'FLUJO', 6),
('9e370f25-d719-46ec-ade7-9fb476c96f44', 'Cancelada', 'FLUJO', 6),

-- Finanzas / Caja
('877e725b-ae57-4501-b1cd-1158fe2df087', 'Pendiente de Conciliación', 'CONTABLE', 7),
('e69d8b1d-d267-47ab-b0ef-507ed6382cd3', 'Conciliado', 'CONTABLE', 7),
('e69d8b1d-d267-47ab-b0ef-507ed6380001', 'Con Diferencia', 'CONTABLE', 7),
('e69d8b1d-d267-47ab-b0ef-507ed6380002', 'Cerrado', 'CONTABLE', 7),
('e69d8b1d-d267-47ab-b0ef-507ed6380003', 'Auditado', 'CONTABLE', 7),

-- Caja / POS
('59039503-85cf-e511-80c1-000c29c9e0e0', 'Activo', 'SESION', 8),
('5a039503-85cf-e511-80c1-000c29c9e0e0', 'Inactivo', 'SESION', 8),
('99039503-85cf-e511-80c1-000c29c9e0e0', 'Fondo Asignado', 'APERTURA', 8),
('5e8dd0fb-5550-e711-80c1-000c29c9e0e0', 'Fondo Por Confirmar', 'APERTURA', 8),
('a864475f-0d34-e711-80c1-000c29c9e0e0', 'Fondo Activo', 'APERTURA', 8),
('2160b065-0d34-e711-80c1-000c29c9e0e0', 'Fondo Retirado', 'CIERRE', 8),
('e8297cfa-630e-e611-80c1-000c29c9e0e0', 'Retiro Efectivo', 'RETIRO', 8),
('84920103-640e-e611-80c1-000c29c9e0e0', 'Retiro Total', 'RETIRO', 8),
('159e3fe6-630e-e611-80c1-000c29c9e0e0', 'Arqueo de Retiros', 'ARQUEO', 8),
('0d4515fe-c907-e611-a6b8-000c29c9e0e0', 'Arqueo Final', 'ARQUEO', 8),
('9b039503-85cf-e511-80c1-000c29c9e0e0', 'Ingreso Administrador', 'ADMIN', 8),
('9c039503-85cf-e511-80c1-000c29c9e0e0', 'Salida Administrador', 'ADMIN', 8),
('9a039503-85cf-e511-80c1-000c29c9e0e0', 'Desmontado', 'CIERRE', 8)
ON CONFLICT (id_status) DO UPDATE SET 
    std_descripcion = EXCLUDED.std_descripcion,
    std_tipo_estado = EXCLUDED.std_tipo_estado,
    mdl_id = EXCLUDED.mdl_id;
