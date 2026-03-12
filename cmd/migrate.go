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

func init() {
	// Registrar el comando en el root
	rootCmd.AddCommand(migrateCmd)
}
