package middleware

import (
	"context"
	"net/http"
	"time"
)

// Timeout es un middleware que cancela el contexto de la petición si excede el tiempo especificado
func Timeout(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			// Sobrescribir el contexto de la petición
			r = r.WithContext(ctx)

			// Canal para detectar si el handler terminó
			done := make(chan struct{})

			go func() {
				next.ServeHTTP(w, r)
				close(done)
			}()

			select {
			case <-done:
				// El handler terminó a tiempo
				return
			case <-ctx.Done():
				// El contexto expiró (timeout) o fue cancelado
				if ctx.Err() == context.DeadlineExceeded {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusGatewayTimeout)
					_, _ = w.Write([]byte(`{"error": "El servidor tardó demasiado en responder (Timeout)", "code": 504}`))
				}
			}
		})
	}
}
