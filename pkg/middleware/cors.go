package middleware

import (
	"net/http"
	"strings"

	"github.com/prunus/pkg/config"
	"github.com/rs/cors"
)

// CORS configura el middleware para el manejo de Cross-Origin Resource Sharing.
// Orígenes permitidos se leen de CORS_ALLOWED_ORIGINS (coma-separados).
// Fallback: http://localhost:3000 para desarrollo local.
func CORS() func(http.Handler) http.Handler {
	raw := config.GetDefault("CORS_ALLOWED_ORIGINS", "http://localhost:3000")
	origins := strings.Split(raw, ",")
	for i, o := range origins {
		origins[i] = strings.TrimSpace(o)
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	return c.Handler
}
