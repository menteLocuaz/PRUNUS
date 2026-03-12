package performance

import (
	"context"
	"log/slog"
	"time"
)

// Thresholds predefinidos
const (
	DbThreshold       = 200 * time.Millisecond
	RedisThreshold    = 50 * time.Millisecond
	ExternalThreshold = 500 * time.Millisecond
)

// Trace mide el tiempo de ejecución y loguea un WARN si excede el threshold indicado.
func Trace(ctx context.Context, layer, operation string, threshold time.Duration, start time.Time) {
	elapsed := time.Since(start)
	if elapsed > threshold {
		slog.WarnContext(ctx, "Operación lenta detectada",
			slog.String("capa", layer),
			slog.String("operacion", operation),
			slog.Duration("latencia", elapsed),
			slog.Int64("latencia_ms", elapsed.Milliseconds()),
			slog.Duration("threshold", threshold),
		)
	}
}
