package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// IPVisitor representa un visitante identificado por su IP con su propio limitador.
type IPVisitor struct {
	limiter  *rate.Limiter
	lastSeen int64 // Unix timestamp de la última petición
}

var (
	visitors = make(map[string]*IPVisitor)
	mu       sync.Mutex
)

func init() {
	// Evicta IPs inactivas cada 5 minutos para evitar memory leak.
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			cutoff := time.Now().Add(-10 * time.Minute).Unix()
			mu.Lock()
			for ip, v := range visitors {
				if v.lastSeen < cutoff {
					delete(visitors, ip)
				}
			}
			mu.Unlock()
		}
	}()
}

// getLimiter devuelve el limitador para una IP específica, creándolo si no existe.
func getLimiter(ip string, r rate.Limit, b int) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	v, exists := visitors[ip]
	if !exists {
		v = &IPVisitor{limiter: rate.NewLimiter(r, b)}
		visitors[ip] = v
	}
	v.lastSeen = time.Now().Unix()
	v.limiter.SetLimit(r)
	v.limiter.SetBurst(b)
	return v.limiter
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
