-- 1. FUNCIÓN DEL TRIGGER: OTORGAR PERMISOS AUTOMÁTICOS
CREATE OR REPLACE FUNCTION fn_auto_permiso_admin()
RETURNS TRIGGER AS $$
BEGIN
    -- Al insertar un módulo nuevo, otorgar permisos completos a todos
    -- los roles que ya tienen permisos totales sobre algún módulo activo.
    INSERT INTO permiso_rol (
        id_rol, id_modulo,
        can_read, can_write, can_update, can_delete
    )
    SELECT DISTINCT pr.id_rol, NEW.id_modulo, true, true, true, true
    FROM permiso_rol pr
    WHERE pr.can_read   = true
      AND pr.can_write  = true
      AND pr.can_update = true
      AND pr.can_delete = true
    ON CONFLICT DO NOTHING;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 2. ASIGNACIÓN DEL TRIGGER A LA TABLA MODULO
DROP TRIGGER IF EXISTS trg_auto_permiso_admin ON modulo;
CREATE TRIGGER trg_auto_permiso_admin
AFTER INSERT ON modulo
FOR EACH ROW
EXECUTE FUNCTION fn_auto_permiso_admin();

-- 3. SINCRONIZACIÓN INICIAL (Backfill)
INSERT INTO permiso_rol (
    id_rol, id_modulo,
    can_read, can_write, can_update, can_delete
)
SELECT DISTINCT roles_admin.id_rol, m.id_modulo, true, true, true, true
FROM modulo m
CROSS JOIN (
    SELECT DISTINCT id_rol
    FROM permiso_rol
    WHERE can_read   = true
      AND can_write  = true
      AND can_update = true
      AND can_delete = true
) AS roles_admin
WHERE m.deleted_at IS NULL
ON CONFLICT DO NOTHING;
