package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

func migrateSeedModulosPermissions(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	const idStatusActivo = "7f7b0e11-1234-4a21-9591-316279f06742"
	const idRolAdmin = "7d7b0e11-1234-4a21-9591-316279f06742"
	const idEmpresaDefault = "11111111-1111-4111-a111-111111111111"
	const idSucursalDefault = "22222222-2222-4222-a222-222222222222"

	// 0. Asegurar que las columnas existan en la tabla modulo
	alterTableQueries := []string{
		`ALTER TABLE modulo ADD COLUMN IF NOT EXISTS id_padre UUID NULL;`,
		`ALTER TABLE modulo ADD COLUMN IF NOT EXISTS orden INTEGER NOT NULL DEFAULT 0;`,
		`ALTER TABLE modulo ADD COLUMN IF NOT EXISTS ruta VARCHAR(255) NULL;`,
		`ALTER TABLE modulo ADD COLUMN IF NOT EXISTS icono VARCHAR(100) NULL;`,
	}

	for _, query := range alterTableQueries {
		if _, err := db.ExecContext(ctx, query); err != nil {
			return fmt.Errorf("error al alterar tabla modulo: %w", err)
		}
	}

	// 0.1 Crear la tabla permiso_rol si no existe
	createPermisoRolQuery := `
	CREATE TABLE IF NOT EXISTS permiso_rol (
		id_permiso  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		id_rol      UUID NOT NULL,
		id_modulo   UUID NOT NULL,
		can_read    BOOLEAN NOT NULL DEFAULT FALSE,
		can_write   BOOLEAN NOT NULL DEFAULT FALSE,
		can_update  BOOLEAN NOT NULL DEFAULT FALSE,
		can_delete  BOOLEAN NOT NULL DEFAULT FALSE,
		can_import  BOOLEAN NOT NULL DEFAULT FALSE,
		can_export  BOOLEAN NOT NULL DEFAULT FALSE,
		created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at  TIMESTAMP NULL,
		CONSTRAINT fk_permiso_rol_rol    FOREIGN KEY (id_rol)    REFERENCES rol(id_rol),
		CONSTRAINT fk_permiso_rol_modulo FOREIGN KEY (id_modulo) REFERENCES modulo(id_modulo),
		CONSTRAINT uk_rol_modulo UNIQUE (id_rol, id_modulo)
	);`

	if _, err := db.ExecContext(ctx, createPermisoRolQuery); err != nil {
		return fmt.Errorf("error al crear tabla permiso_rol: %w", err)
	}

	// 0.2 Asegurar Empresa
	insertEmpresa := `INSERT INTO empresa (id_empresa, nombre, rut, id_status)
	VALUES ($1, 'Empresa Demo', '12345678-9', $2) ON CONFLICT (id_empresa) DO NOTHING;`
	if _, err := db.ExecContext(ctx, insertEmpresa, idEmpresaDefault, idStatusActivo); err != nil {
		return fmt.Errorf("error al asegurar empresa: %w", err)
	}

	// 0.3 Asegurar Sucursal
	insertSucursal := `INSERT INTO sucursal (id_sucursal, id_empresa, nombre_sucursal, id_status)
	VALUES ($1, $2, 'Sucursal Central', $3) ON CONFLICT (id_sucursal) DO NOTHING;`
	if _, err := db.ExecContext(ctx, insertSucursal, idSucursalDefault, idEmpresaDefault, idStatusActivo); err != nil {
		return fmt.Errorf("error al asegurar sucursal: %w", err)
	}

	// 1. Módulos Principales
	queryModulosPrincipales := `
	INSERT INTO modulo (mdl_id, mdl_descripcion, id_status, is_active, abreviatura, nivel, orden, ruta, icono)
	VALUES 
		(1, 'Configuración Empresa', $1, true, 'EMP',  0, 1, '/config/empresa', 'settings'),
		(2, 'Gestión Sucursales',    $1, true, 'SUC',  0, 2, '/config/sucursales', 'store'),
		(3, 'Usuarios y Roles',      $1, true, 'USR',  0, 3, '/config/usuarios', 'users'),
		(4, 'Catálogo Productos',    $1, true, 'PROD', 0, 4, '/productos', 'package'),
		(5, 'Ventas y POS',          $1, true, 'VENT', 0, 5, '/ventas', 'shopping-cart'),
		(8, 'Control de Caja',       $1, true, 'CAJA', 0, 6, '/caja', 'monitor')
	ON CONFLICT (mdl_id) DO UPDATE SET 
		mdl_descripcion = EXCLUDED.mdl_descripcion,
		ruta = EXCLUDED.ruta,
		icono = EXCLUDED.icono,
		is_active = EXCLUDED.is_active,
		orden = EXCLUDED.orden;`

	if _, err := db.ExecContext(ctx, queryModulosPrincipales, idStatusActivo); err != nil {
		return fmt.Errorf("error módulos principales: %w", err)
	}

	// 2. Submódulos
	submodulosQueries := []struct {
		name string
		sql  string
	}{
		{
			name: "Apertura de Caja",
			sql: `INSERT INTO modulo (mdl_id, mdl_descripcion, id_status, is_active, abreviatura, nivel, orden, ruta, icono, id_padre)
				  SELECT 801, 'Apertura de Caja', $1, true, 'OPEN', 1, 1, '/caja/apertura', 'lock-open', id_modulo
				  FROM modulo WHERE mdl_id = 8
				  ON CONFLICT (mdl_id) DO UPDATE SET ruta = EXCLUDED.ruta, icono = EXCLUDED.icono, id_padre = EXCLUDED.id_padre;`,
		},
		{
			name: "Arqueo de Caja",
			sql: `INSERT INTO modulo (mdl_id, mdl_descripcion, id_status, is_active, abreviatura, nivel, orden, ruta, icono, id_padre)
				  SELECT 802, 'Arqueo de Caja', $1, true, 'ARQ', 1, 2, '/caja/arqueo', 'calculator', id_modulo
				  FROM modulo WHERE mdl_id = 8
				  ON CONFLICT (mdl_id) DO UPDATE SET ruta = EXCLUDED.ruta, icono = EXCLUDED.icono, id_padre = EXCLUDED.id_padre;`,
		},
		{
			name: "Nueva Venta",
			sql: `INSERT INTO modulo (mdl_id, mdl_descripcion, id_status, is_active, abreviatura, nivel, orden, ruta, icono, id_padre)
				  SELECT 501, 'Nueva Venta', $1, true, 'SALE', 1, 1, '/pos/ventas', 'plus-circle', id_modulo
				  FROM modulo WHERE mdl_id = 5
				  ON CONFLICT (mdl_id) DO UPDATE SET ruta = EXCLUDED.ruta, icono = EXCLUDED.icono, id_padre = EXCLUDED.id_padre;`,
		},
	}

	for _, sm := range submodulosQueries {
		if _, err := db.ExecContext(ctx, sm.sql, idStatusActivo); err != nil {
			fmt.Printf("Aviso: Submódulo %s no se pudo crear: %v\n", sm.name, err)
		}
	}

	// 4. Asegurar Rol Administrador
	rolQuery := `
	INSERT INTO rol (id_rol, nombre_rol, id_sucursal, id_status)
	SELECT $1, 'Administrador', id_sucursal, $2
	FROM sucursal WHERE id_sucursal = $3
	ON CONFLICT (id_rol) DO NOTHING;`
	
	if _, err := db.ExecContext(ctx, rolQuery, idRolAdmin, idStatusActivo, idSucursalDefault); err != nil {
		return fmt.Errorf("error al asegurar el rol Admin: %w", err)
	}

	// 5. Permisos Totales Admin
	permisosQuery := `
	INSERT INTO permiso_rol (id_rol, id_modulo, can_read, can_write, can_update, can_delete, can_import, can_export)
	SELECT $1, m.id_modulo, true, true, true, true, true, true
	FROM modulo m
	ON CONFLICT (id_rol, id_modulo) DO UPDATE SET 
		can_read = true, can_write = true, can_update = true, can_delete = true;`

	if _, err := db.ExecContext(ctx, permisosQuery, idRolAdmin); err != nil {
		return fmt.Errorf("error permisos admin: %w", err)
	}

	// Asegurar que los permisos de import/export también se actualicen (ON CONFLICT anterior los omitía)
	updateImportExport := `
	UPDATE permiso_rol SET can_import = true, can_export = true
	WHERE id_rol = $1 AND (can_import = false OR can_export = false);`

	if _, err := db.ExecContext(ctx, updateImportExport, idRolAdmin); err != nil {
		return fmt.Errorf("error actualizando import/export admin: %w", err)
	}

	return nil
}
