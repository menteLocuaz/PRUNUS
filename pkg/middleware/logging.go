package middleware

import (
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
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

	// TimeFormat es el formato de tiempo a usar (mantenido por compatibilidad)
	TimeFormat string

	// Output es donde se escribirán los logs (default: os.Stdout)
	Output io.Writer

	// UseJSON indica si los logs deben ser en formato JSON estructurado
	UseJSON bool

	// ColorOutput indica si usar colores (mantenido por compatibilidad, slog no lo soporta nativamente en TextHandler)
	ColorOutput bool
}

// DefaultLogConfig retorna una configuración por defecto para el middleware.
func DefaultLogConfig() *LogConfig {
	return &LogConfig{
		SkipPaths:      []string{},
		LogHeaders:     false,
		LogQueryParams: true,
		LogUserAgent:   true,
		TimeFormat:     time.RFC3339,
		Output:         os.Stdout,
		UseJSON:        true,
		ColorOutput:    false,
	}
}

// SimpleLogConfig retorna una configuración simple para logs en texto plano.
func SimpleLogConfig() *LogConfig {
	return &LogConfig{
		SkipPaths:      []string{},
		LogHeaders:     false,
		LogQueryParams: false,
		LogUserAgent:   false,
		TimeFormat:     "2006-01-02 15:04:05",
		Output:         os.Stdout,
		UseJSON:        false,
		ColorOutput:    true,
	}
}

// ProductionLogConfig retorna una configuración optimizada para producción.
func ProductionLogConfig() *LogConfig {
	return &LogConfig{
		SkipPaths:      []string{"/health", "/healthz", "/metrics"},
		LogHeaders:     false,
		LogQueryParams: true,
		LogUserAgent:   true,
		TimeFormat:     time.RFC3339,
		Output:         os.Stdout,
		UseJSON:        true,
		ColorOutput:    false,
	}
}

// VerboseLogConfig retorna una configuración con máximo detalle.
func VerboseLogConfig() *LogConfig {
	return &LogConfig{
		SkipPaths:      []string{},
		LogHeaders:     true,
		LogQueryParams: true,
		LogUserAgent:   true,
		TimeFormat:     time.RFC3339Nano,
		Output:         os.Stdout,
		UseJSON:        true,
		ColorOutput:    false,
	}
}

// DevelopmentLogConfig retorna una configuración para desarrollo.
func DevelopmentLogConfig() *LogConfig {
	return &LogConfig{
		SkipPaths:      []string{},
		LogHeaders:     false,
		LogQueryParams: true,
		LogUserAgent:   false,
		TimeFormat:     "15:04:05",
		Output:         os.Stdout,
		UseJSON:        false,
		ColorOutput:    true,
	}
}

// Logger es el middleware de logging optimizado que registra información de cada petición HTTP usando slog.
func Logger(config *LogConfig) func(http.Handler) http.Handler {
	if config == nil {
		config = DefaultLogConfig()
	}

	if config.Output == nil {
		config.Output = os.Stdout
	}

	var slogHandler slog.Handler
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	if config.UseJSON {
		slogHandler = slog.NewJSONHandler(config.Output, opts)
	} else {
		slogHandler = slog.NewTextHandler(config.Output, opts)
	}

	logger := slog.New(slogHandler)

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
			ctx := r.Context()

			// Pre-asignar slice de atributos para evitar re-allocations constantes (capacidad estimada de 16)
			attrs := make([]slog.Attr, 0, 16)
			
			// Extraer Request ID del contexto si existe
			if reqID := GetRequestID(ctx); reqID != "" {
				attrs = append(attrs, slog.String("request_id", reqID))
			}

			attrs = append(attrs,
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", lrw.statusCode),
				slog.Duration("latency", latency),
				slog.Int64("latency_ms", latency.Milliseconds()),
				slog.String("client_ip", getClientIP(r)),
				slog.Int64("bytes_out", lrw.written),
			)

			if config.LogQueryParams && r.URL.RawQuery != "" {
				attrs = append(attrs, slog.String("query", r.URL.RawQuery))
			}

			if config.LogUserAgent {
				attrs = append(attrs, slog.String("user_agent", r.UserAgent()))
			}

			if config.LogHeaders {
				headerAttrs := make([]slog.Attr, 0, len(r.Header))
				for key, values := range r.Header {
					headerAttrs = append(headerAttrs, slog.String(key, strings.Join(values, ", ")))
				}
				attrs = append(attrs, slog.Attr{
					Key:   "headers",
					Value: slog.GroupValue(headerAttrs...),
				})
			}

			// Determinar el nivel de log basado en el código de estado
			level := slog.LevelInfo
			if lrw.statusCode >= 500 {
				level = slog.LevelError
			} else if lrw.statusCode >= 400 {
				level = slog.LevelWarn
			}

			// Usar el contexto de la petición para permitir tracing y propagación de datos
			logger.LogAttrs(r.Context(), level, "Petición HTTP completada", attrs...)
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
