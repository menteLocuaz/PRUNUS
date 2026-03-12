package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const (
	// RequestIDKey es la clave usada para almacenar el Request ID en el contexto.
	RequestIDKey contextKey = "request_id"
	// RequestIDHeader es el header HTTP usado para el Request ID.
	RequestIDHeader = "X-Request-ID"
)

// RequestID genera un ID único para cada petición y lo inyecta en el contexto y los headers.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Intentar obtener el ID del header (por si viene de un proxy o gateway)
		requestID := r.Header.Get(RequestIDHeader)
		if requestID == "" {
			// Si no existe, generar uno nuevo (UUID v4)
			requestID = uuid.New().String()
		}

		// Inyectar el ID en los headers de la respuesta para trazabilidad del cliente
		w.Header().Set(RequestIDHeader, requestID)

		// Crear un nuevo contexto con el ID
		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)

		// Continuar con el siguiente handler usando el nuevo contexto
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetRequestID extrae el Request ID del contexto de forma segura.
func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}
