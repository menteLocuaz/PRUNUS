package middleware

import (
	"net/http"

	"github.com/rs/cors"
)

// CORS configura el middleware para el manejo de Cross-Origin Resource Sharing
func CORS() func(http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // En producción, especificar dominios reales
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // 5 minutos
	})

	return c.Handler
}
