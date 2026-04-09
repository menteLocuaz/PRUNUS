package main

import (
	"log"

	"github.com/prunus/pkg/config/database"
	"github.com/prunus/pkg/config/database/migrations"
	"github.com/spf13/cobra"
)

// migrateCmd representa el comando para ejecutar migraciones
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Ejecuta las migraciones de la base de datos",
	Long:  `Este comando se conecta a la base de datos y ejecuta todas las migraciones pendientes del proyecto.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Conexión a Base de Datos
		db, err := database.InitDB()
		if err != nil {
			log.Fatalf("❌ Error conectando a la base de datos para migración: %v", err)
		}
		defer db.Close()

		// 2. Ejecutar Migraciones
		log.Println("🛠️ Iniciando migraciones de la base de datos...")
		if err := migrations.RunMigrations(db); err != nil {
			log.Fatalf("❌ Error durante la ejecución de migraciones: %v", err)
		}
		log.Println("✅ Migraciones completadas exitosamente")
	},
}

// seedCmd representa el comando para poblar datos iniciales
var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Puebla los módulos y permisos iniciales",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := database.InitDB()
		if err != nil {
			log.Fatalf("❌ Error conectando a la base de datos: %v", err)
		}
		defer db.Close()

		log.Println("🌱 Iniciando sembrado de módulos y permisos...")
		// Reutilizamos la lógica de la migración 070 pero forzando ejecución
		// Nota: En un entorno real podrías exportar la función o tener un paquete seeds
		// Por simplicidad en este entorno, ejecutaremos un script SQL directo similar
		query := `
		-- Insertar/Actualizar Módulos
		INSERT INTO modulo (mdl_id, mdl_descripcion, id_status, is_active, abreviatura, nivel, orden, ruta, icono)
		VALUES 
			(1, 'Configuración Empresa', '7f7b0e11-1234-4a21-9591-316279f06742', true, 'EMP',  0, 1, '/config/empresa', 'settings'),
			(2, 'Gestión Sucursales',    '7f7b0e11-1234-4a21-9591-316279f06742', true, 'SUC',  0, 2, '/config/sucursales', 'store'),
			(3, 'Usuarios y Roles',      '7f7b0e11-1234-4a21-9591-316279f06742', true, 'USR',  0, 3, '/config/usuarios', 'users'),
			(4, 'Catálogo Productos',    '7f7b0e11-1234-4a21-9591-316279f06742', true, 'PROD', 0, 4, '/productos', 'package'),
			(5, 'Ventas y POS',          '7f7b0e11-1234-4a21-9591-316279f06742', true, 'VENT', 0, 5, '/ventas', 'shopping-cart'),
			(8, 'Control de Caja',       '7f7b0e11-1234-4a21-9591-316279f06742', true, 'CAJA', 0, 6, '/caja', 'monitor')
		ON CONFLICT (mdl_id) DO UPDATE SET 
			mdl_descripcion = EXCLUDED.mdl_descripcion, ruta = EXCLUDED.ruta, icono = EXCLUDED.icono;

		-- Asignar todo a CUALQUIER rol que se llame 'Administrador'
		INSERT INTO permiso_rol (id_rol, id_modulo, can_read, can_write, can_update, can_delete)
		SELECT r.id_rol, m.id_modulo, true, true, true, true
		FROM rol r, modulo m
		WHERE r.nombre_rol = 'Administrador'
		ON CONFLICT (id_rol, id_modulo) DO UPDATE SET can_read = true, can_write = true;
		`
		if _, err := db.Exec(query); err != nil {
			log.Fatalf("❌ Error durante el seed: %v", err)
		}
		log.Println("✅ Módulos y permisos sincronizados.")
	},
}

func init() {
	// Registrar el comando en el root
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(seedCmd)
}
