# Middleware de Logging para Prunus API

Middleware modular de logging basado en **log/slog** (estándar de Go 1.21+) para registrar todas las peticiones HTTP en tu aplicación Go con Chi router.

## Características

✅ **Estándar** - Utiliza la librería oficial `log/slog` de Go
✅ **Estructurado** - Logs consistentes y fáciles de parsear
✅ **Modular** - Se activa/desactiva con una sola línea
✅ **Configurable** - Múltiples opciones de personalización
✅ **Cero Breaking Changes** - Mantiene compatibilidad con la API anterior
✅ **Flexible** - Salida JSON o texto plano nativa de slog
✅ **Información Completa** - Método, ruta, estado, latencia, IP, etc.

## Instalación

El middleware ya está incluido en el proyecto. Utiliza las capacidades nativas de Go.

## Uso Rápido

### 1. Activar Logging (Ya está activado por defecto)

En [main_router.go](../routers/main_router.go), el middleware ya está configurado:

```go
func NewMainRouter(...) http.Handler {
    r := chi.NewRouter()

    // Middleware de logging - YA ACTIVADO
    r.Use(middleware.Logger(middleware.ProductionLogConfig()))

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
    // r.Use(middleware.Logger(middleware.ProductionLogConfig()))

    // ... tus rutas
    return r
}
```

## Configuraciones Disponibles

### Configuración por Defecto (Recomendada)

Logs estructurados en JSON con información completa usando `slog.JSONHandler`:

```go
r.Use(middleware.Logger(middleware.DefaultLogConfig()))
```

**Salida JSON (slog):**
```json
{
  "time": "2026-03-11T10:30:45Z",
  "level": "INFO",
  "msg": "Petición HTTP completada",
  "method": "POST",
  "path": "/api/v1/usuario",
  "status": 201,
  "latency": 45000000,
  "latency_ms": 45,
  "client_ip": "192.168.1.100",
  "user_agent": "PostmanRuntime/7.26.8",
  "bytes_out": 234
}
```

### Configuración Simple (Desarrollo)

Logs en texto plano usando `slog.TextHandler`:

```go
r.Use(middleware.Logger(middleware.SimpleLogConfig()))
```

**Salida Text (slog):**
```
time=2026-03-11T10:30:45.000Z level=INFO msg="Petición HTTP completada" method=POST path=/api/v1/usuario status=201 latency=45ms latency_ms=45 client_ip=192.168.1.100
```

### Configuración de Producción

Optimizada para ambientes de producción (ignora health checks y usa JSON):

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
    Output:         os.Stdout,                         // Donde escribir
    UseJSON:        true,                              // Formato JSON
}

r.Use(middleware.Logger(customConfig))
```

## Opciones de Configuración

| Campo | Tipo | Descripción | Default |
|-------|------|-------------|---------|
| `SkipPaths` | `[]string` | Rutas que no se loguearán | `[]` |
| `LogHeaders` | `bool` | Incluir headers HTTP (agrupados) | `false` |
| `LogQueryParams` | `bool` | Incluir query parameters | `true` |
| `LogUserAgent` | `bool` | Incluir User-Agent | `true` |
| `Output` | `io.Writer` | Destino de los logs | `os.Stdout` |
| `UseJSON` | `bool` | Usar `JSONHandler` vs `TextHandler` | `true` |

> **Nota sobre Compatibilidad:** Los campos `TimeFormat` y `ColorOutput` se mantienen en la estructura para evitar errores de compilación, pero son ignorados ya que `slog` gestiona el formato de tiempo automáticamente y no soporta colores nativamente en su `TextHandler` estándar.

## Información Capturada

El middleware registra la siguiente información por cada petición:

| Campo | Descripción | JSON Key |
|-------|-------------|----------|
| Timestamp | Fecha y hora (ISO8601) | `time` |
| Nivel | INFO, WARN (4xx), ERROR (5xx) | `level` |
| Mensaje | "Petición HTTP completada" | `msg` |
| Método HTTP | GET, POST, PUT, DELETE, etc. | `method` |
| Ruta | Path de la URL | `path` |
| Query Parameters | Parámetros de consulta (opcional) | `query` |
| Código de Estado | 200, 404, 500, etc. | `status` |
| Latencia | Duración legible y en ms | `latency`, `latency_ms` |
| IP del Cliente | Dirección IP real | `client_ip` |
| User-Agent | Navegador/cliente (opcional) | `user_agent` |
| Bytes Enviados | Tamaño de la respuesta | `bytes_out` |
| Headers HTTP | Headers agrupados (opcional) | `headers` |

## Niveles de Log Automáticos

El middleware ajusta el nivel de log según el resultado de la petición:
- **INFO**: Códigos 1xx, 2xx, 3xx
- **WARN**: Códigos 4xx (Errores del cliente)
- **ERROR**: Códigos 5xx (Errores del servidor)

## Ventajas de log/slog

1. **Eficiencia**: `slog` está altamente optimizado y es parte del core de Go.
2. **Estandarización**: Todos tus logs seguirán el mismo formato estructurado.
3. **Interoperabilidad**: Fácil integración con cualquier sistema de agregación de logs (ELK, Datadog, CloudWatch).
4. **Agrupación**: Los headers se guardan como un grupo de atributos, facilitando su visualización en herramientas como Kibana.

## Mantenimiento

### Actualizar a log/slog

Si tienes código antiguo que usaba `LogEntry`, simplemente elíminalo. `slog` maneja los atributos de forma dinámica, lo que hace el código mucho más limpio y fácil de mantener.

Consulta [examples.go](examples.go) para ver ejemplos actualizados.
