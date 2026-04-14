package logger

import (
	"context"
	"os"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New crea un *zap.Logger de producción: JSON estructurado, stdout, nivel Info.
// Incluye caller y stack trace en errores para análisis en ELK/Grafana.
func New() *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.LowercaseLevelEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(os.Stdout),
		zap.InfoLevel,
	)
	return zap.New(core,
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	)
}

// NewDevelopment crea un logger de consola legible para desarrollo local.
func NewDevelopment() *zap.Logger {
	l, _ := zap.NewDevelopment()
	return l
}

// WithContext devuelve el logger enriquecido con request_id y user_id extraídos
// del contexto. Si el contexto no contiene esos valores, devuelve base sin modificar.
func WithContext(ctx context.Context, base *zap.Logger) *zap.Logger {
	if ctx == nil {
		return base
	}

	fields := make([]zap.Field, 0, 2)

	// request_id — clave string (ver middleware/request_id.go)
	if id, ok := ctx.Value("request_id").(string); ok && id != "" {
		fields = append(fields, zap.String("request_id", id))
	}

	// user_id — inyectado por middleware/auth.go tras validar JWT
	if uid, ok := ctx.Value("user_id").(uuid.UUID); ok {
		fields = append(fields, zap.String("user_id", uid.String()))
	}

	if len(fields) == 0 {
		return base
	}
	return base.With(fields...)
}
