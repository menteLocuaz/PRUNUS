package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/prunus/pkg/config/database"
	"github.com/spf13/cobra"
)

// Constantes para configuración de migraciones
const (
	migrationsPath = "file://database/migrations"
	defaultSchema  = "public"
)

// migrateCmd representa el comando base de migraciones
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Gestión de migraciones de base de datos",
	Long:  `Permite ejecutar, revertir y consultar el estado de las migraciones SQL.`,
}

// upMigrateCmd ejecuta todas las migraciones pendientes
var upMigrateCmd = &cobra.Command{
	Use:   "up",
	Short: "Ejecuta todas las migraciones pendientes",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := database.InitDB()
		if err != nil {
			log.Fatalf("❌ Error conectando a DB: %v", err)
		}
		defer db.Close()

		m, err := newMigrator(db)
		if err != nil {
			log.Fatalf("❌ Error al inicializar migrador: %v", err)
		}
		defer m.migrate.Close()

		log.Println("🚀 Ejecutando migraciones pendientes...")
		if err := m.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				log.Println("✅ La base de datos ya está actualizada.")
				return
			}
			log.Fatalf("❌ Error ejecutando migraciones: %v", err)
		}
		log.Println("✅ Migraciones completadas exitosamente.")
	},
}

// downMigrateCmd revierte la última migración o N migraciones
var downMigrateCmd = &cobra.Command{
	Use:   "down [n]",
	Short: "Revierte la(s) última(s) migración(es)",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		steps := 1
		if len(args) > 0 {
			var err error
			steps, err = strconv.Atoi(args[0])
			if err != nil {
				log.Fatalf("❌ El número de pasos debe ser un entero: %v", err)
			}
		}

		db, err := database.InitDB()
		if err != nil {
			log.Fatalf("❌ Error conectando a DB: %v", err)
		}
		defer db.Close()

		m, err := newMigrator(db)
		if err != nil {
			log.Fatalf("❌ Error al inicializar migrador: %v", err)
		}
		defer m.migrate.Close()

		log.Printf("⏪ Revirtiendo %d migración(es)...", steps)
		if err := m.Steps(-steps); err != nil {
			log.Fatalf("❌ Error revirtiendo migraciones: %v", err)
		}
		log.Println("✅ Reversión completada exitosamente.")
	},
}

// versionMigrateCmd muestra la versión actual de la base de datos
var versionMigrateCmd = &cobra.Command{
	Use:   "version",
	Short: "Muestra la versión actual de la migración",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := database.InitDB()
		if err != nil {
			log.Fatalf("❌ Error conectando a DB: %v", err)
		}
		defer db.Close()

		m, err := newMigrator(db)
		if err != nil {
			log.Fatalf("❌ Error al inicializar migrador: %v", err)
		}
		defer m.migrate.Close()

		version, dirty, err := m.Version()
		if err != nil {
			if errors.Is(err, migrate.ErrNilVersion) {
				log.Println("ℹ️ No hay migraciones aplicadas.")
				return
			}
			log.Fatalf("❌ Error obteniendo versión: %v", err)
		}

		status := "limpio"
		if dirty {
			status = "DIRTY (requiere intervención manual)"
		}
		log.Printf("📊 Versión actual: %d (%s)", version, status)
	},
}

// forceMigrateCmd limpia el estado dirty forzando una versión
var forceMigrateCmd = &cobra.Command{
	Use:   "force <version>",
	Short: "Fuerza una versión específica (limpia estado dirty)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		version, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalf("❌ La versión debe ser un número entero: %v", err)
		}

		db, err := database.InitDB()
		if err != nil {
			log.Fatalf("❌ Error conectando a DB: %v", err)
		}
		defer db.Close()

		m, err := newMigrator(db)
		if err != nil {
			log.Fatalf("❌ Error al inicializar migrador: %v", err)
		}
		defer m.migrate.Close()

		log.Printf("🛠️ Forzando versión %d para limpiar estado...", version)
		if err := m.Force(version); err != nil {
			log.Fatalf("❌ Error al forzar versión: %v", err)
		}
		log.Println("✅ Estado limpiado correctamente.")
	},
}

// RunMigrationsIfNeeded es una función auxiliar para ejecutar migraciones al arranque si se desea
func RunMigrationsIfNeeded(db *sql.DB) error {
	m, err := newMigrator(db)
	if err != nil {
		return fmt.Errorf("error al inicializar migrador: %w", err)
	}
	// NOTA: No llamamos a m.migrate.Close() aquí porque cerraría la conexión 'db' compartida
	// que el servidor necesita para seguir funcionando.

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		return fmt.Errorf("error ejecutando migraciones: %w", err)
	}
	return nil
}

// MigratorWrapper envuelve la lógica de golang-migrate
type MigratorWrapper struct {
	db      *sql.DB
	migrate *migrate.Migrate
}

// newMigrator inicializa el motor de migraciones con una conexión existente
func newMigrator(db *sql.DB) (*MigratorWrapper, error) {
	// Asegurar esquema base
	_, err := db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s;", defaultSchema))
	if err != nil {
		log.Printf("⚠️ Advertencia al preparar esquema: %v", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{
		SchemaName: defaultSchema,
	})
	if err != nil {
		return nil, fmt.Errorf("error creando driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("error inicializando migrate: %w", err)
	}

	return &MigratorWrapper{
		db:      db,
		migrate: m,
	}, nil
}

// Up ejecuta migraciones pendientes
func (mw *MigratorWrapper) Up() error {
	return mw.migrate.Up()
}

// Steps ejecuta N pasos (positivo hacia adelante, negativo hacia atrás)
func (mw *MigratorWrapper) Steps(n int) error {
	return mw.migrate.Steps(n)
}

// Version obtiene la versión actual
func (mw *MigratorWrapper) Version() (uint, bool, error) {
	return mw.migrate.Version()
}

// Force establece una versión manualmente
func (mw *MigratorWrapper) Force(v int) error {
	return mw.migrate.Force(v)
}

// Close cierra la conexión a la base de datos
func (mw *MigratorWrapper) Close() error {
	sourceErr, dbErr := mw.migrate.Close()
	if sourceErr != nil {
		return sourceErr
	}
	if dbErr != nil {
		return dbErr
	}
	return mw.db.Close()
}

func init() {
	// Registro de subcomandos
	migrateCmd.AddCommand(upMigrateCmd)
	migrateCmd.AddCommand(downMigrateCmd)
	migrateCmd.AddCommand(versionMigrateCmd)
	migrateCmd.AddCommand(forceMigrateCmd)

	rootCmd.AddCommand(migrateCmd)
}
