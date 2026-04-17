-- 000044_pos_opening_controls.up.sql
-- Controles de consistencia para apertura de caja

-- 1. Asegurar estatus críticos (BLOQUEADA, ANULADA) si no existen
-- Nota: Usamos IDs fijos para consistencia con estatus_constants.go
INSERT INTO estatus (id_status, std_descripcion, mdl_id) 
VALUES 
    ('B1039503-85CF-E511-80C1-000C29C9E0E0', 'Bloqueada', 8),
    ('A1039503-85CF-E511-80C1-000C29C9E0E0', 'Anulada', 8)
ON CONFLICT (id_status) DO NOTHING;

-- 2. BLOQUEO DE CONSISTENCIA: Un usuario no puede tener dos cajas abiertas al mismo tiempo.
-- Creamos un índice de unicidad parcial.
CREATE UNIQUE INDEX IF NOT EXISTS idx_control_estacion_usuario_activo 
ON control_estacion (usuario_asignado) 
WHERE (fecha_salida IS NULL AND deleted_at IS NULL);

-- 3. BLOQUEO DE CONCURRENCIA: Una estación no puede ser abierta dos veces simultáneamente.
-- (Ya debería estar protegido por la lógica de negocio, pero esto lo blinda a nivel DB)
CREATE UNIQUE INDEX IF NOT EXISTS idx_control_estacion_estacion_activa 
ON control_estacion (id_estacion) 
WHERE (fecha_salida IS NULL AND deleted_at IS NULL);
