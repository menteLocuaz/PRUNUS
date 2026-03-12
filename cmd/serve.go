package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prunus/pkg/config/database"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/routers"
	"github.com/prunus/pkg/store"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var port string

// serveCmd representa el comando para iniciar el servidor
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Inicia el servidor API REST",
	Long:  `Levanta el servidor HTTP de Prunus e inyecta todas las dependencias necesarias con soporte para Graceful Shutdown.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Conexión a Base de Datos
		db, err := database.InitDB()
		if err != nil {
			log.Fatalf("❌ Error crítico conectando a la base de datos: %v", err)
		}
		defer db.Close()

		// 2. Conexión a Redis
		rdb, err := database.InitRedis()
		var cacheStore models.CacheStore
		if err != nil {
			fmt.Printf("⚠️ Aviso: No se pudo conectar a Redis: %v. Cache desactivado.\n", err)
		} else {
			cacheStore = store.NewRedisStore(rdb)
		}

		// 3. Logger
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		slog.SetDefault(logger)

		// 4. Handlers y Router
		h := RegisterHandlers(db, cacheStore, logger)
		router := routers.NewMainRouter(h)

		// 5. Configuración del Servidor
		finalPort := viper.GetString("PORT")
		if port != "9090" {
			finalPort = port
		}
		if finalPort == "" {
			finalPort = "9090"
		}

		srv := &http.Server{
			Addr:    ":" + finalPort,
			Handler: router,
		}

		// 6. Graceful Shutdown
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
