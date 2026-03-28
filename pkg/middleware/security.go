package middleware

import (
	"context"
	"net"
	"net/http"
	"strings"
)

const (
	// ClientIPKey es la clave usada para almacenar la IP del cliente en el contexto.
	ClientIPKey contextKey = "client_ip"
)

// MaxPayloadSize limita el tamaño de la solicitud entrante (1MB por defecto)
func MaxPayloadSize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Limitar el cuerpo de la solicitud a 1MB
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MB
		next.ServeHTTP(w, r)
	})
}

// SecureHeaders añade cabeceras de seguridad para mitigar ataques comunes
func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		next.ServeHTTP(w, r)
	})
}

// ClientIP captura la IP real del cliente y la inyecta en el contexto
func ClientIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := extractIP(r)
		ctx := context.WithValue(r.Context(), ClientIPKey, ip)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetClientIP extrae la IP del cliente del contexto
func GetClientIP(ctx context.Context) string {
	if ctx == nil {
		return "unknown"
	}
	if ip, ok := ctx.Value(ClientIPKey).(string); ok {
		return ip
	}
	return "unknown"
}

func extractIP(r *http.Request) string {
	// 1. X-Forwarded-For (Proxies)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}

	// 2. X-Real-IP
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}

	// 3. RemoteAddr
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
