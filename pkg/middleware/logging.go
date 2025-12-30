package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// loggingResponseWriter es un wrapper de http.ResponseWriter que captura
// el código de estado HTTP y el tamaño de la respuesta
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	written    int64
}

// newLoggingResponseWriter crea una nueva instancia del wrapper
func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK, // Por defecto 200 OK
		written:        0,
	}
}

// WriteHeader captura el código de estado antes de escribirlo
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// Write captura el tamaño de la respuesta
func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	n, err := lrw.ResponseWriter.Write(b)
	lrw.written += int64(n)
	return n, err
}

// LogConfig contiene la configuración para el middleware de logging
type LogConfig struct {
	// SkipPaths son las rutas que no se loguearán (ej: /health, /metrics)
	SkipPaths []string

	// LogHeaders indica si se deben loguear los headers de la petición
	LogHeaders bool

	// LogQueryParams indica si se deben loguear los query parameters
	LogQueryParams bool

	// LogUserAgent indica si se debe loguear el User-Agent
	LogUserAgent bool

	// TimeFormat es el formato de tiempo a usar (default: RFC3339)
	TimeFormat string

	// Output es donde se escribirán los logs (default: os.Stdout)
	Output io.Writer

	// UseJSON indica si los logs deben ser en formato JSON estructurado
	UseJSON bool

	// ColorOutput indica si usar colores en la salida (solo para logs no-JSON)
	ColorOutput bool
}

// DefaultLogConfig retorna una configuración por defecto para el middleware
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

// SimpleLogConfig retorna una configuración simple para logs en texto plano
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

// LogEntry representa una entrada de log estructurada
type LogEntry struct {
	Time       string            `json:"time"`
	Method     string            `json:"method"`
	Path       string            `json:"path"`
	Query      string            `json:"query,omitempty"`
	Status     int               `json:"status"`
	LatencyMs  int64             `json:"latency_ms"`
	ClientIP   string            `json:"client_ip"`
	UserAgent  string            `json:"user_agent,omitempty"`
	BytesOut   int64             `json:"bytes_out"`
	Headers    map[string]string `json:"headers,omitempty"`
}

// Logger es el middleware de logging que registra información de cada petición HTTP
func Logger(config *LogConfig) func(http.Handler) http.Handler {
	// Si no se proporciona configuración, usar la por defecto
	if config == nil {
		config = DefaultLogConfig()
	}

	// Si no se especifica output, usar stdout
	if config.Output == nil {
		config.Output = os.Stdout
	}

	// Si no se especifica formato de tiempo, usar RFC3339
	if config.TimeFormat == "" {
		config.TimeFormat = time.RFC3339
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

			// Registrar tiempo de inicio
			start := time.Now()

			// Crear wrapper del ResponseWriter
			lrw := newLoggingResponseWriter(w)

			// Ejecutar el siguiente handler
			next.ServeHTTP(lrw, r)

			// Calcular latencia
			latency := time.Since(start)

			// Crear entrada de log
			entry := LogEntry{
				Time:      start.Format(config.TimeFormat),
				Method:    r.Method,
				Path:      r.URL.Path,
				Status:    lrw.statusCode,
				LatencyMs: latency.Milliseconds(),
				ClientIP:  getClientIP(r),
				BytesOut:  lrw.written,
			}

			// Agregar query params si está configurado
			if config.LogQueryParams && r.URL.RawQuery != "" {
				entry.Query = r.URL.RawQuery
			}

			// Agregar User-Agent si está configurado
			if config.LogUserAgent {
				entry.UserAgent = r.UserAgent()
			}

			// Agregar headers si está configurado
			if config.LogHeaders {
				entry.Headers = make(map[string]string)
				for key, values := range r.Header {
					entry.Headers[key] = strings.Join(values, ", ")
				}
			}

			// Escribir el log
			if config.UseJSON {
				writeJSONLog(config.Output, entry)
			} else {
				writeTextLog(config.Output, entry, config.ColorOutput)
			}
		})
	}
}

// getClientIP extrae la IP del cliente de la petición
func getClientIP(r *http.Request) string {
	// Intentar obtener de headers comunes de proxies
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		// X-Forwarded-For puede contener múltiples IPs, tomar la primera
		if idx := strings.Index(ip, ","); idx != -1 {
			return strings.TrimSpace(ip[:idx])
		}
		return strings.TrimSpace(ip)
	}

	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return strings.TrimSpace(ip)
	}

	// Fallback a RemoteAddr
	if idx := strings.LastIndex(r.RemoteAddr, ":"); idx != -1 {
		return r.RemoteAddr[:idx]
	}

	return r.RemoteAddr
}

// writeJSONLog escribe el log en formato JSON
func writeJSONLog(output io.Writer, entry LogEntry) {
	encoder := json.NewEncoder(output)
	if err := encoder.Encode(entry); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding log entry: %v\n", err)
	}
}

// writeTextLog escribe el log en formato de texto plano
func writeTextLog(output io.Writer, entry LogEntry, useColor bool) {
	// Códigos de color ANSI
	var (
		reset  = ""
		cyan   = ""
		yellow = ""
		red    = ""
		green  = ""
		blue   = ""
	)

	if useColor {
		reset = "\033[0m"
		cyan = "\033[36m"
		yellow = "\033[33m"
		red = "\033[31m"
		green = "\033[32m"
		blue = "\033[34m"
	}

	// Colorear según el código de estado
	statusColor := green
	if entry.Status >= 400 && entry.Status < 500 {
		statusColor = yellow
	} else if entry.Status >= 500 {
		statusColor = red
	}

	// Formatear el log
	logLine := fmt.Sprintf("%s[%s]%s %s%-6s%s %s%3d%s %s%-50s%s %s%5dms%s %s%s%s\n",
		cyan, entry.Time, reset,
		blue, entry.Method, reset,
		statusColor, entry.Status, reset,
		"", entry.Path, "",
		yellow, entry.LatencyMs, reset,
		"", entry.ClientIP, "",
	)

	fmt.Fprint(output, logLine)
}
