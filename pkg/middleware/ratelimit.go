package middleware

import (
	"net"
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

// IPVisitor representa un visitante identificado por su IP con su propio limitador
type IPVisitor struct {
	limiter  *rate.Limiter
	lastSeen int64
}

var (
	visitors = make(map[string]*rate.Limiter)
	mu       sync.Mutex
)

// getLimiter devuelve el limitador para una IP específica, creándolo si no existe o actualizándolo si los parámetros cambian
func getLimiter(ip string, r rate.Limit, b int) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(r, b)
		visitors[ip] = limiter
		return limiter
	}

	// Actualizar el límite y ráfaga por si cambiaron en la configuración
	limiter.SetLimit(r)
	limiter.SetBurst(b)

	return limiter
}

// RateLimit configura el middleware para limitar la tasa de peticiones por IP usando Token Bucket
func RateLimit(requestsPerSecond float64, burst int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			}

			limiter := getLimiter(ip, rate.Limit(requestsPerSecond), burst)
			if !limiter.Allow() {
				http.Error(w, "Demasiadas peticiones", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
