package main

import (
	"context"
	"fmt"
	"log"

	"github.com/prunus/pkg/config/database"
	"github.com/spf13/cobra"
)

var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Operaciones sobre la caché de Redis",
}

var cacheClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Elimina todas las claves de la base de datos Redis activa",
	Run: func(cmd *cobra.Command, args []string) {
		rdb, err := database.InitRedis()
		if err != nil {
			log.Fatalf("❌ Error conectando a Redis: %v", err)
		}
		defer rdb.Close()

		ctx := context.Background()
		if err := rdb.FlushDB(ctx).Err(); err != nil {
			log.Fatalf("❌ Error al limpiar la caché: %v", err)
		}

		fmt.Println("✅ Caché de Redis limpiada correctamente")
	},
}

func init() {
	cacheCmd.AddCommand(cacheClearCmd)
	rootCmd.AddCommand(cacheCmd)
}
