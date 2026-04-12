package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/prunus/pkg/config/database"
	"github.com/spf13/cobra"
)

// migrateCmd representa el comando para ejecutar migraciones
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Ejecuta las migraciones de la base de datos (SQL)",
	Long:  `Este comando utiliza golang-migrate para ejecutar los archivos SQL en /database/migrations.`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := database.InitDB()
		if err != nil {
			log.Fatalf("❌ Error conectando a la base de datos: %v", err)
		}
		defer db.Close()

		log.Println("🚀 Iniciando ejecución de migraciones SQL...")
		if err := runSqlMigrations(db); err != nil {
			log.Fatalf("❌ Error en migraciones SQL: %v", err)
		}

		log.Println("✅ Base de datos actualizada correctamente")
	},
}

// runSqlMigrations inicializa y ejecuta el motor de golang-migrate
func runSqlMigrations(db *sql.DB) error {
	// 1. Asegurar esquema public y search_path
	_, err := db.Exec("CREATE SCHEMA IF NOT EXISTS public; SET search_path TO public;")
	if err != nil {
		log.Printf("⚠️ Advertencia al preparar esquema: %v", err)
	}

	// 2. Usar ruta relativa estándar para el origen de archivos.
	// golang-migrate interpreta file:// como relativo a la raíz del proceso.
	srcURL := "file://database/migrations"

	// 3. Configurar driver con esquema explícito
	driver, err := postgres.WithInstance(db, &postgres.Config{
		SchemaName: "public",
	})
	if err != nil {
		return fmt.Errorf("error al crear driver de postgres: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(srcURL, "postgres", driver)
	if err != nil {
		return fmt.Errorf("error al inicializar migrate: %w", err)
	}

	// 4. Ejecutar migraciones
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error al ejecutar migraciones: %w", err)
	}

	return nil
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

// forceMigrateCmd permite limpiar el estado dirty forzando una versión específica
var forceMigrateCmd = &cobra.Command{
	Use:   "force [version]",
	Short: "Fuerza una versión específica de la migración (limpia estado dirty)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		version, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalf("❌ La versión debe ser un número entero: %v", err)
		}

		db, err := database.InitDB()
		if err != nil {
			log.Fatalf("❌ Error conectando a la base de datos: %v", err)
		}
		defer db.Close()

		// Asegurar esquema public y search_path
		_, _ = db.Exec("CREATE SCHEMA IF NOT EXISTS public; SET search_path TO public;")

		srcURL := "file://database/migrations"
		driver, err := postgres.WithInstance(db, &postgres.Config{
			SchemaName: "public",
		})
		if err != nil {
			log.Fatalf("❌ Error al crear driver: %v", err)
		}

		m, err := migrate.NewWithDatabaseInstance(srcURL, "postgres", driver)
		if err != nil {
			log.Fatalf("❌ Error al inicializar migrate: %v", err)
		}

		if err := m.Force(version); err != nil {
			log.Fatalf("❌ Error al forzar versión: %v", err)
		}

		log.Printf("✅ Versión forzada a %d correctamente", version)
	},
}

func init() {
	// Registrar el comando en el root
	migrateCmd.AddCommand(forceMigrateCmd)
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(seedCmd)
}
