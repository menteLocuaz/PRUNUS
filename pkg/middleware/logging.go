package middleware

import (
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	zaplogger "github.com/prunus/pkg/utils/logger"
	"go.uber.org/zap"
)

// loggingResponseWriter es un wrapper de http.ResponseWriter que captura
// el código de estado HTTP y el tamaño de la respuesta.
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	written    int64
}

// rwPool permite reutilizar estructuras loggingResponseWriter para reducir la presión sobre el GC.
var rwPool = sync.Pool{
	New: func() any {
		return &loggingResponseWriter{}
	},
}

// WriteHeader captura el código de estado antes de escribirlo.
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// Write captura el tamaño de la respuesta.
func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	n, err := lrw.ResponseWriter.Write(b)
	lrw.written += int64(n)
	return n, err
}

// LogConfig contiene la configuración para el middleware de logging.
type LogConfig struct {
	// SkipPaths son las rutas que no se loguearán (ej: /health, /metrics)
	SkipPaths []string

	// LogHeaders indica si se deben loguear los headers de la petición
	LogHeaders bool

	// LogQueryParams indica si se deben loguear los query parameters
	LogQueryParams bool

	// LogUserAgent indica si se debe loguear el User-Agent
	LogUserAgent bool

	// Output es donde se escribirán los logs (default: os.Stdout)
	Output io.Writer
}

// DefaultLogConfig retorna una configuración por defecto para el middleware.
func DefaultLogConfig() *LogConfig {
	return &LogConfig{
		SkipPaths:      []string{},
		LogHeaders:     false,
		LogQueryParams: true,
		LogUserAgent:   true,
		Output:         os.Stdout,
	}
}

// SimpleLogConfig retorna una configuración simple.
func SimpleLogConfig() *LogConfig {
	return &LogConfig{
		SkipPaths:      []string{},
		LogHeaders:     false,
		LogQueryParams: false,
		LogUserAgent:   false,
		Output:         os.Stdout,
	}
}

// ProductionLogConfig retorna una configuración optimizada para producción.
func ProductionLogConfig() *LogConfig {
	return &LogConfig{
		SkipPaths:      []string{"/health", "/healthz", "/metrics"},
		LogHeaders:     false,
		LogQueryParams: true,
		LogUserAgent:   true,
		Output:         os.Stdout,
	}
}

// VerboseLogConfig retorna una configuración con máximo detalle.
func VerboseLogConfig() *LogConfig {
	return &LogConfig{
		SkipPaths:      []string{},
		LogHeaders:     true,
		LogQueryParams: true,
		LogUserAgent:   true,
		Output:         os.Stdout,
	}
}

// DevelopmentLogConfig retorna una configuración para desarrollo.
func DevelopmentLogConfig() *LogConfig {
	return &LogConfig{
		SkipPaths:      []string{},
		LogHeaders:     false,
		LogQueryParams: true,
		LogUserAgent:   false,
		Output:         os.Stdout,
	}
}

// Logger es el middleware de logging que registra cada petición HTTP en formato JSON
// usando zap. Incluye automáticamente request_id y user_id del contexto en cada log,
// lo que permite correlación de trazas en ELK / Grafana.
func Logger(config *LogConfig, log *zap.Logger) func(http.Handler) http.Handler {
	if config == nil {
		config = DefaultLogConfig()
	}
	if config.Output == nil {
		config.Output = os.Stdout
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verificar si esta ruta debe ser ignorada
			for _, skipPath := range config.SkipPaths {
				if strings.HasPrefix(r.URL.Path, skipPath) {
					next.ServeHTTP(w, r)
					return
				}
			}

			start := time.Now()

			// Obtener wrapper del pool y resetear sus valores
			lrw := rwPool.Get().(*loggingResponseWriter)
			lrw.ResponseWriter = w
			lrw.statusCode = http.StatusOK
			lrw.written = 0
			defer rwPool.Put(lrw)

			// Ejecutar el siguiente handler
			next.ServeHTTP(lrw, r)

			latency := time.Since(start)

			// Logger enriquecido con request_id y user_id del contexto
			entry := zaplogger.WithContext(r.Context(), log).With(
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", lrw.statusCode),
				zap.Duration("latency", latency),
				zap.Int64("latency_ms", latency.Milliseconds()),
				zap.String("client_ip", getClientIP(r)),
				zap.Int64("bytes_out", lrw.written),
			)

			if config.LogQueryParams && r.URL.RawQuery != "" {
				entry = entry.With(zap.String("query", r.URL.RawQuery))
			}

			if config.LogUserAgent {
				entry = entry.With(zap.String("user_agent", r.UserAgent()))
			}

			if config.LogHeaders {
				headerFields := make([]zap.Field, 0, len(r.Header))
				for key, values := range r.Header {
					headerFields = append(headerFields, zap.String(key, strings.Join(values, ", ")))
				}
				entry = entry.With(zap.Namespace("headers"))
				entry = entry.With(headerFields...)
			}

			// Nivel de log según código de estado
			switch {
			case lrw.statusCode >= 500:
				entry.Error("Petición HTTP completada")
			case lrw.statusCode >= 400:
				entry.Warn("Petición HTTP completada")
			default:
				entry.Info("Petición HTTP completada")
			}
		})
	}
}

// getClientIP extrae la IP del cliente de la petición de forma optimizada.
func getClientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		if idx := strings.Index(ip, ","); idx != -1 {
			return strings.TrimSpace(ip[:idx])
		}
		return strings.TrimSpace(ip)
	}

	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return strings.TrimSpace(ip)
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
