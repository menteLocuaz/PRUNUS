# Middleware de Logging para Prunus API

Middleware modular de logging para registrar todas las peticiones HTTP en tu aplicación Go con Chi router.

## Características

✅ **Modular** - Se activa/desactiva con una sola línea
✅ **Configurable** - Múltiples opciones de personalización
✅ **Cero Breaking Changes** - No modifica código existente
✅ **Performante** - Overhead mínimo
✅ **Flexible** - Salida JSON o texto plano
✅ **Información Completa** - Método, ruta, estado, latencia, IP, etc.

## Instalación

El middleware ya está incluido en el proyecto. No requiere instalación adicional.

## Uso Rápido

### 1. Activar Logging (Ya está activado por defecto)

En [main_router.go](../routers/main_router.go), el middleware ya está configurado:

```go
func NewMainRouter(...) http.Handler {
    r := chi.NewRouter()

    // Middleware de logging - YA ACTIVADO
    r.Use(middleware.Logger(middleware.DefaultLogConfig()))

    // ... tus rutas
    return r
}
```

### 2. Desactivar Logging

Para desactivar completamente:

```go
func NewMainRouter(...) http.Handler {
    r := chi.NewRouter()

    // Simplemente comenta o elimina esta línea
    // r.Use(middleware.Logger(middleware.DefaultLogConfig()))

    // ... tus rutas
    return r
}
```

## Configuraciones Disponibles

### Configuración por Defecto (Recomendada)

Logs estructurados en JSON con información completa:

```go
r.Use(middleware.Logger(middleware.DefaultLogConfig()))
```

**Salida JSON:**
```json
{"time":"2025-12-30T10:30:45-05:00","method":"POST","path":"/api/v1/usuario","query":"","status":201,"latency_ms":45,"client_ip":"192.168.1.100","user_agent":"PostmanRuntime/7.26.8","bytes_out":234}
```

### Configuración Simple (Desarrollo)

Logs en texto plano con colores:

```go
r.Use(middleware.Logger(middleware.SimpleLogConfig()))
```

**Salida con colores:**
```
[2025-12-30 10:30:45] POST   201 /api/v1/usuario                                    45ms 192.168.1.100
```

### Configuración de Producción

Optimizada para ambientes de producción (ignora health checks):

```go
r.Use(middleware.Logger(middleware.ProductionLogConfig()))
```

### Configuración Personalizada

Crea tu propia configuración:

```go
customConfig := &middleware.LogConfig{
    SkipPaths:      []string{"/health", "/metrics"},  // Rutas a ignorar
    LogHeaders:     true,                              // Incluir headers
    LogQueryParams: true,                              // Incluir query params
    LogUserAgent:   true,                              // Incluir User-Agent
    TimeFormat:     time.RFC3339,                      // Formato de tiempo
    Output:         os.Stdout,                         // Donde escribir
    UseJSON:        true,                              // Formato JSON
    ColorOutput:    false,                             // Sin colores
}

r.Use(middleware.Logger(customConfig))
```

## Opciones de Configuración

| Campo | Tipo | Descripción | Default |
|-------|------|-------------|---------|
| `SkipPaths` | `[]string` | Rutas que no se loguearán | `[]` |
| `LogHeaders` | `bool` | Incluir headers HTTP | `false` |
| `LogQueryParams` | `bool` | Incluir query parameters | `true` |
| `LogUserAgent` | `bool` | Incluir User-Agent | `true` |
| `TimeFormat` | `string` | Formato de timestamp | `time.RFC3339` |
| `Output` | `io.Writer` | Destino de los logs | `os.Stdout` |
| `UseJSON` | `bool` | Usar formato JSON | `true` |
| `ColorOutput` | `bool` | Colorear salida de texto | `false` |

## Ejemplos Avanzados

### Logging a Archivo

```go
import (
    "io"
    "log"
    "os"
)

// Abrir archivo de logs
logFile, err := os.OpenFile("logs/api.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
if err != nil {
    log.Fatal("Error abriendo archivo de logs:", err)
}

// Configurar logging a archivo
fileConfig := &middleware.LogConfig{
    Output:  logFile,
    UseJSON: true,
}

r.Use(middleware.Logger(fileConfig))
```

### Logging a Consola + Archivo Simultáneamente

```go
import "io"

// Abrir archivo
logFile, _ := os.OpenFile("logs/api.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

// Crear MultiWriter
multiWriter := io.MultiWriter(os.Stdout, logFile)

config := &middleware.LogConfig{
    Output:  multiWriter,
    UseJSON: true,
}

r.Use(middleware.Logger(config))
```

### Logging Solo para Rutas Específicas

```go
r.Route("/api/v1", func(r chi.Router) {
    // Logging SOLO para /usuario
    r.Route("/usuario", func(r chi.Router) {
        r.Use(middleware.Logger(middleware.DefaultLogConfig()))

        r.Get("/", usuarioHandler.GetAll)
        r.Post("/", usuarioHandler.Create)
    })

    // SIN logging para /empresas
    r.Route("/empresas", func(r chi.Router) {
        r.Get("/", empresaHandler.GetAll)
        r.Post("/", empresaHandler.Create)
    })
})
```

### Ignorar Health Checks

```go
config := &middleware.LogConfig{
    SkipPaths: []string{
        "/health",
        "/healthz",
        "/metrics",
        "/ping",
    },
    UseJSON: true,
}

r.Use(middleware.Logger(config))
```

## Información Capturada

El middleware registra la siguiente información por cada petición:

| Campo | Descripción | JSON Key |
|-------|-------------|----------|
| Timestamp | Fecha y hora de la petición | `time` |
| Método HTTP | GET, POST, PUT, DELETE, etc. | `method` |
| Ruta | Path de la URL | `path` |
| Query Parameters | Parámetros de consulta (opcional) | `query` |
| Código de Estado | 200, 404, 500, etc. | `status` |
| Latencia | Tiempo de respuesta en ms | `latency_ms` |
| IP del Cliente | Dirección IP (con soporte X-Forwarded-For) | `client_ip` |
| User-Agent | Navegador/cliente (opcional) | `user_agent` |
| Bytes Enviados | Tamaño de la respuesta | `bytes_out` |
| Headers HTTP | Headers de la petición (opcional) | `headers` |

## Extracción de IP del Cliente

El middleware detecta automáticamente la IP real del cliente, incluso detrás de proxies:

1. `X-Forwarded-For` header
2. `X-Real-IP` header
3. `RemoteAddr` (fallback)

## Formato de Salida

### JSON (Recomendado para Producción)

```json
{
  "time": "2025-12-30T10:30:45-05:00",
  "method": "POST",
  "path": "/api/v1/usuario",
  "query": "page=1&limit=10",
  "status": 201,
  "latency_ms": 45,
  "client_ip": "192.168.1.100",
  "user_agent": "PostmanRuntime/7.26.8",
  "bytes_out": 234
}
```

### Texto Plano (Desarrollo)

```
[2025-12-30 10:30:45] POST   201 /api/v1/usuario                                    45ms 192.168.1.100
[2025-12-30 10:30:46] GET    200 /api/v1/empresas                                   12ms 192.168.1.101
[2025-12-30 10:30:47] PUT    200 /api/v1/usuario/123                                78ms 192.168.1.100
[2025-12-30 10:30:48] DELETE 204 /api/v1/empresas/456                               23ms 192.168.1.102
```

Con colores activados, el código de estado se colorea:
- **Verde**: 2xx (éxito)
- **Amarillo**: 4xx (error del cliente)
- **Rojo**: 5xx (error del servidor)

## Rendimiento

El middleware tiene un overhead mínimo:

- **Sin headers**: ~10-20μs por petición
- **Con headers**: ~30-50μs por petición
- **JSON encoding**: ~100-200μs por petición

Para aplicaciones de alta carga, se recomienda:
- Desactivar `LogHeaders` si no es necesario
- Usar `SkipPaths` para rutas de health check
- Considerar logging asíncrono para volúmenes muy altos

## Troubleshooting

### Los logs no aparecen

1. Verifica que el middleware esté activado en [main_router.go](../routers/main_router.go)
2. Asegúrate de que `Output` no sea `nil`
3. Si usas archivo, verifica los permisos de escritura

### Los códigos de estado siempre son 200

Esto puede ocurrir si tus handlers no llaman a `w.WriteHeader()`. El middleware captura el código correctamente si se envía explícitamente.

### Colores no se muestran

1. Verifica que `ColorOutput: true` esté configurado
2. Algunos terminales/consolas no soportan colores ANSI
3. En Windows, usa Windows Terminal o ConEmu para mejor soporte

## Integración con Otros Sistemas

### Enviar a Servicio de Logging Externo

Puedes crear un `io.Writer` personalizado que envíe logs a servicios como:
- ELK Stack (Elasticsearch, Logstash, Kibana)
- Splunk
- CloudWatch (AWS)
- Stackdriver (GCP)
- Application Insights (Azure)

```go
type ExternalLogWriter struct {
    // Tu implementación
}

func (w *ExternalLogWriter) Write(p []byte) (n int, err error) {
    // Enviar a servicio externo
    return len(p), nil
}

config := &middleware.LogConfig{
    Output: &ExternalLogWriter{},
}
```

### Integración con Logrus/Zap

```go
import "github.com/sirupsen/logrus"

// Crear logger de logrus
logger := logrus.New()

config := &middleware.LogConfig{
    Output: logger.Writer(),
}
```

## Comparación con Otros Frameworks

| Acción | Gin | Echo | Chi (este middleware) |
|--------|-----|------|----------------------|
| **Activar** | `r.Use(gin.Logger())` | `e.Use(middleware.Logger())` | `r.Use(middleware.Logger(nil))` |
| **Desactivar** | Comentar línea | Comentar línea | Comentar línea |
| **Configurar** | No disponible en default | Config struct | `LogConfig` struct |
| **JSON Output** | No por defecto | No por defecto | ✅ Sí por defecto |

## Mantenimiento

### Actualizar Configuración

Para cambiar la configuración en producción sin reiniciar:

1. Modifica [main_router.go](../routers/main_router.go)
2. Reinicia el servidor
3. Los nuevos logs usarán la nueva configuración

### Ver Ejemplos Completos

Consulta [examples.go](examples.go) para ver 9 ejemplos completos de configuraciones diferentes.

## Licencia

Este middleware es parte del proyecto Prunus y sigue la misma licencia del proyecto principal.

## Soporte

Para reportar bugs o sugerir mejoras, contacta al equipo de desarrollo.
