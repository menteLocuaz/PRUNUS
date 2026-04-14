package main

import (
	"log"

	"github.com/prunus/pkg/config/database"
	"github.com/spf13/cobra"
)

// seedCmd representa el comando para poblar datos iniciales (Módulos y Permisos)
var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Puebla los módulos y permisos iniciales del sistema",
	Long: `Este comando inserta o actualiza los módulos base del sistema y 
asigna permisos completos al rol 'Administrador'. Es útil para inicializar 
el sistema después de las migraciones SQL.`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := database.InitDB()
		if err != nil {
			log.Fatalf("❌ Error conectando a la base de datos: %v", err)
		}
		defer db.Close()

		log.Println("🌱 Iniciando sembrado de módulos y permisos...")

		// SQL para sembrado de datos maestros
		query := `
		-- Insertar/Actualizar Módulos Base
		INSERT INTO modulo (mdl_id, mdl_descripcion, id_status, is_active, abreviatura, nivel, orden, ruta, icono)
		VALUES 
			(1, 'Configuración Empresa', '7f7b0e11-1234-4a21-9591-316279f06742', true, 'EMP',  0, 1, '/config/empresa', 'settings'),
			(2, 'Gestión Sucursales',    '7f7b0e11-1234-4a21-9591-316279f06742', true, 'SUC',  0, 2, '/config/sucursales', 'store'),
			(3, 'Usuarios y Roles',      '7f7b0e11-1234-4a21-9591-316279f06742', true, 'USR',  0, 3, '/config/usuarios', 'users'),
			(4, 'Catálogo Productos',    '7f7b0e11-1234-4a21-9591-316279f06742', true, 'PROD', 0, 4, '/productos', 'package'),
			(5, 'Ventas y POS',          '7f7b0e11-1234-4a21-9591-316279f06742', true, 'VENT', 0, 5, '/ventas', 'shopping-cart'),
			(8, 'Control de Caja',       '7f7b0e11-1234-4a21-9591-316279f06742', true, 'CAJA', 0, 6, '/caja', 'monitor')
		ON CONFLICT (mdl_id) DO UPDATE SET 
			mdl_descripcion = EXCLUDED.mdl_descripcion, 
			ruta = EXCLUDED.ruta, 
			icono = EXCLUDED.icono;

		-- Asignar permisos completos a cualquier rol con nombre 'Administrador'
		INSERT INTO permiso_rol (id_rol, id_modulo, can_read, can_write, can_update, can_delete)
		SELECT r.id_rol, m.id_modulo, true, true, true, true
		FROM rol r, modulo m
		WHERE r.nombre_rol = 'Administrador'
		ON CONFLICT (id_rol, id_modulo) DO UPDATE SET 
			can_read = true, 
			can_write = true, 
			can_update = true, 
			can_delete = true;
		`

		if _, err := db.Exec(query); err != nil {
			log.Fatalf("❌ Error durante el sembrado de datos: %v", err)
		}
		log.Println("✅ Módulos y permisos sincronizados correctamente.")
	},
}

func init() {
	// Registrar el comando en el root
	rootCmd.AddCommand(seedCmd)
}
