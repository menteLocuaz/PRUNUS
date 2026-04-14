package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prunus/pkg/config"
	"github.com/prunus/pkg/config/database"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/routers"
	"github.com/prunus/pkg/store"
	zaplogger "github.com/prunus/pkg/utils/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var port string

// serveCmd representa el comando para iniciar el servidor
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Inicia el servidor API REST",
	Long:  `Levanta el servidor HTTP de Prunus e inyecta todas las dependencias necesarias con soporte para Graceful Shutdown.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Validar variables de configuración críticas antes de arrancar
		if err := config.Validate(
			"DB_HOST", "DB_USER", "DB_PASSWORD", "DB_NAME",
			"JWT_SECRET",
		); err != nil {
			log.Fatalf("❌ Error de configuración: %v", err)
		}

		// 2. Conexión a Base de Datos
		db, err := database.InitDB()
		if err != nil {
			log.Fatalf("❌ Error crítico conectando a la base de datos: %v", err)
		}
		defer db.Close()

		// 3. Conexión a Redis
		rdb, err := database.InitRedis()
		var cacheStore models.CacheStore
		if err != nil {
			fmt.Printf("⚠️ Aviso: No se pudo conectar a Redis: %v. Cache desactivado.\n", err)
		} else {
			cacheStore = store.NewRedisStore(rdb)
		}

		// 4. Logger
		logger := zaplogger.New()
		defer logger.Sync() //nolint:errcheck
		zap.ReplaceGlobals(logger)

		// 5. Handlers y Router
		h := RegisterHandlers(db, cacheStore, logger)
		router := routers.NewMainRouter(h, logger)

		// 5b. Worker de Segundo Plano para Snapshots de Inventario
		go StartBackgroundWorker(db)

		// 6. Configuración del Servidor
		finalPort := config.GetDefault("PORT", "9090")
		if port != "9090" {
			finalPort = port
		}

		srv := &http.Server{
			Addr:         ":" + finalPort,
			Handler:      router,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		}

		// 7. Graceful Shutdown
		go func() {
			fmt.Printf("🚀 Servidor Prunus iniciado en el puerto %s\n", finalPort)
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("❌ Error iniciando el servidor: %v", err)
			}
		}()

		// Canal para escuchar señales de interrupción (Ctrl+C, SIGTERM)
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

		<-stop // Esperar señal

		fmt.Println("\n🛑 Apagando el servidor de forma segura...")

		// Tiempo de gracia para cerrar conexiones activas (10 segundos)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("❌ Error durante el apagado del servidor: %v", err)
		}

		fmt.Println("✅ Servidor detenido correctamente.")
	},
}

func init() {
	serveCmd.Flags().StringVarP(&port, "port", "p", "9090", "Puerto en el que escuchará el servidor")
	rootCmd.AddCommand(serveCmd)
}

// StartBackgroundWorker ejecuta tareas periódicas como snapshots de inventario
func StartBackgroundWorker(db *sql.DB) {
	fmt.Println("⏲️ Worker de segundo plano iniciado (Snapshots diarios)")
	
	// Ejecutar inmediatamente al arrancar para asegurar el dato del día
	takeInventorySnapshots(db)

	ticker := time.NewTicker(24 * time.Hour)
	for range ticker.C {
		takeInventorySnapshots(db)
	}
}

func takeInventorySnapshots(db *sql.DB) {
	// Obtener todas las sucursales activas
	rows, err := db.Query("SELECT id_sucursal FROM sucursal WHERE deleted_at IS NULL")
	if err != nil {
		fmt.Printf("❌ Error obteniendo sucursales para snapshot: %v\n", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var sucursalID string
		if err := rows.Scan(&sucursalID); err != nil {
			continue
		}

		// Llamar a la función de Postgres fn_snapshot_inventario(uuid)
		_, err := db.Exec("SELECT fn_snapshot_inventario($1)", sucursalID)
		if err != nil {
			fmt.Printf("❌ Error generando snapshot para sucursal %s: %v\n", sucursalID, err)
		}
	}
	fmt.Printf("📸 Snapshots de inventario generados correctamente (%s)\n", time.Now().Format("2006-01-02 15:04:05"))
}
