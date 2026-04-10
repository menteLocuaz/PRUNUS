package migrations

import (
	"database/sql"
	"fmt"
)

func migrateTriggerAutoAdminPermissions(db *sql.DB) error {
	steps := []struct {
		name string
		sql  string
	}{
		{
			name: "función trigger auto permiso admin",
			sql: `
			CREATE OR REPLACE FUNCTION fn_auto_permiso_admin()
			RETURNS TRIGGER AS $$
			BEGIN
				-- Al insertar un módulo nuevo, otorgar permisos completos a todos
				-- los roles que ya tienen permisos totales sobre algún módulo.
				-- Esto asegura que roles tipo Administrador siempre vean los nuevos módulos.
				INSERT INTO permiso_rol (
					id_rol, id_modulo,
					can_read, can_write, can_update, can_delete, can_import, can_export
				)
				SELECT DISTINCT pr.id_rol, NEW.id_modulo, true, true, true, true, true, true
				FROM permiso_rol pr
				WHERE pr.can_read   = true
				  AND pr.can_write  = true
				  AND pr.can_update = true
				  AND pr.can_delete = true
				  AND pr.deleted_at IS NULL
				ON CONFLICT (id_rol, id_modulo) DO NOTHING;

				RETURN NEW;
			END;
			$$ LANGUAGE plpgsql;`,
		},
		{
			name: "eliminar trigger existente si existe",
			sql:  `DROP TRIGGER IF EXISTS trg_auto_permiso_admin ON modulo;`,
		},
		{
			name: "trigger en tabla modulo",
			sql: `
			CREATE TRIGGER trg_auto_permiso_admin
			AFTER INSERT ON modulo
			FOR EACH ROW
			EXECUTE FUNCTION fn_auto_permiso_admin();`,
		},
		{
			name: "backfill permisos faltantes para roles admin",
			sql: `
			-- Sincronizar módulos que ya existen pero no tienen permiso_rol
			-- para los roles con permisos completos (roles administradores).
			INSERT INTO permiso_rol (
				id_rol, id_modulo,
				can_read, can_write, can_update, can_delete, can_import, can_export
			)
			SELECT DISTINCT roles_admin.id_rol, m.id_modulo, true, true, true, true, true, true
			FROM modulo m
			CROSS JOIN (
				SELECT DISTINCT id_rol
				FROM permiso_rol
				WHERE can_read   = true
				  AND can_write  = true
				  AND can_update = true
				  AND can_delete = true
				  AND deleted_at IS NULL
			) AS roles_admin
			WHERE m.deleted_at IS NULL
			ON CONFLICT (id_rol, id_modulo) DO NOTHING;`,
		},
	}

	for _, step := range steps {
		if _, err := db.Exec(step.sql); err != nil {
			return fmt.Errorf("migración 071 - %s: %w", step.name, err)
		}
	}

	return nil
}
