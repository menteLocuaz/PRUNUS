package performance

import (
	"context"
	"time"

	"github.com/prunus/pkg/utils/logger"
	"go.uber.org/zap"
)

// Thresholds predefinidos
const (
	DbThreshold       = 200 * time.Millisecond
	RedisThreshold    = 50 * time.Millisecond
	ExternalThreshold = 500 * time.Millisecond
)

// Trace mide el tiempo de ejecución y loguea un WARN si excede el threshold indicado.
// Usa el logger global (zap.L()) enriquecido con request_id y user_id del contexto.
// El logger global se inicializa en cmd/serve.go mediante zap.ReplaceGlobals.
func Trace(ctx context.Context, layer, operation string, threshold time.Duration, start time.Time) {
	elapsed := time.Since(start)
	if elapsed > threshold {
		logger.WithContext(ctx, zap.L()).Warn("Operación lenta detectada",
			zap.String("capa", layer),
			zap.String("operacion", operation),
			zap.Duration("latencia", elapsed),
			zap.Int64("latencia_ms", elapsed.Milliseconds()),
			zap.Duration("threshold", threshold),
		)
	}
}
