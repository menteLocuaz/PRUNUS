package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/httprate"
)

// RateLimit configura el middleware para limitar la tasa de peticiones por IP
func RateLimit(requests int, window time.Duration) func(http.Handler) http.Handler {
	return httprate.LimitByIP(requests, window)
}
