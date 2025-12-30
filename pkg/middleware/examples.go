package middleware

import (
	"os"
)

// Este archivo contiene ejemplos de configuración del middleware de logging
// NO es código de producción, solo ejemplos para referencia

/*
===========================================
EJEMPLO 1: Configuración por Defecto (JSON)
===========================================

Logs estructurados en JSON con toda la información básica.

func NewMainRouter(...) http.Handler {
    r := chi.NewRouter()

    // Logging JSON estructurado (recomendado para producción)
    r.Use(middleware.Logger(middleware.DefaultLogConfig()))

    // ... tus rutas
    return r
}

Salida:
{"time":"2025-12-30T10:30:45-05:00","method":"POST","path":"/api/v1/usuario","query":"","status":201,"latency_ms":45,"client_ip":"192.168.1.100","user_agent":"PostmanRuntime/7.26.8","bytes_out":234}


===========================================
EJEMPLO 2: Configuración Simple (Texto con Colores)
===========================================

Logs en texto plano con colores para desarrollo local.

func NewMainRouter(...) http.Handler {
    r := chi.NewRouter()

    // Logging simple con colores (ideal para desarrollo)
    r.Use(middleware.Logger(middleware.SimpleLogConfig()))

    // ... tus rutas
    return r
}

Salida:
[2025-12-30 10:30:45] POST   201 /api/v1/usuario                                    45ms 192.168.1.100


===========================================
EJEMPLO 3: Configuración Personalizada
===========================================

Configuración avanzada con todas las opciones.

func NewMainRouter(...) http.Handler {
    r := chi.NewRouter()

    // Configuración personalizada
    customConfig := &middleware.LogConfig{
        SkipPaths:      []string{"/health", "/metrics"},  // No loguear estas rutas
        LogHeaders:     true,                              // Incluir headers HTTP
        LogQueryParams: true,                              // Incluir query parameters
        LogUserAgent:   true,                              // Incluir User-Agent
        TimeFormat:     time.RFC3339Nano,                 // Formato de tiempo preciso
        Output:         os.Stdout,                         // Salida estándar
        UseJSON:        true,                              // Formato JSON
        ColorOutput:    false,                             // Sin colores
    }

    r.Use(middleware.Logger(customConfig))

    // ... tus rutas
    return r
}


===========================================
EJEMPLO 4: Logging a Archivo
===========================================

Guardar logs en un archivo en lugar de consola.

func NewMainRouter(...) http.Handler {
    r := chi.NewRouter()

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

    // ... tus rutas
    return r
}


===========================================
EJEMPLO 5: Logging Solo para Rutas Específicas
===========================================

Aplicar logging solo a ciertas rutas en lugar de global.

func NewMainRouter(...) http.Handler {
    r := chi.NewRouter()

    r.Route("/api/v1", func(r chi.Router) {
        // Logging solo para /api/v1/usuario
        r.Route("/usuario", func(r chi.Router) {
            r.Use(middleware.Logger(middleware.DefaultLogConfig()))

            r.Get("/", usuarioHandler.GetAll)
            r.Post("/", usuarioHandler.Create)
            // ...
        })

        // Sin logging para /api/v1/empresas
        r.Route("/empresas", func(r chi.Router) {
            r.Get("/", empresaHandler.GetAll)
            r.Post("/", empresaHandler.Create)
            // ...
        })
    })

    return r
}


===========================================
EJEMPLO 6: Ignorar Rutas de Health Check
===========================================

No loguear rutas de monitoreo que generan mucho ruido.

func NewMainRouter(...) http.Handler {
    r := chi.NewRouter()

    // Configurar rutas a ignorar
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

    // ... tus rutas
    return r
}


===========================================
EJEMPLO 7: Múltiples Outputs (Consola + Archivo)
===========================================

Escribir logs tanto a consola como a archivo simultáneamente.

func NewMainRouter(...) http.Handler {
    r := chi.NewRouter()

    // Abrir archivo
    logFile, err := os.OpenFile("logs/api.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatal(err)
    }

    // Crear MultiWriter para escribir a ambos destinos
    multiWriter := io.MultiWriter(os.Stdout, logFile)

    config := &middleware.LogConfig{
        Output:  multiWriter,
        UseJSON: true,
    }

    r.Use(middleware.Logger(config))

    // ... tus rutas
    return r
}


===========================================
EJEMPLO 8: Desactivar Logging Completamente
===========================================

Para desactivar el logging, simplemente comenta o elimina la línea.

func NewMainRouter(...) http.Handler {
    r := chi.NewRouter()

    // Logging DESACTIVADO - simplemente comentar o eliminar
    // r.Use(middleware.Logger(middleware.DefaultLogConfig()))

    // ... tus rutas
    return r
}


===========================================
EJEMPLO 9: Logging Minimalista (Solo Errores)
===========================================

Loguear solo peticiones que resultaron en errores.

// Nota: Este ejemplo requeriría modificar el middleware para agregar
// un filtro de nivel de log. Por ahora, puedes usar SkipPaths para
// reducir el ruido, o implementar un wrapper adicional.

func NewMainRouter(...) http.Handler {
    r := chi.NewRouter()

    // Configuración mínima
    minimalConfig := &middleware.LogConfig{
        LogHeaders:     false,
        LogQueryParams: false,
        LogUserAgent:   false,
        UseJSON:        false,
        ColorOutput:    true,
    }

    r.Use(middleware.Logger(minimalConfig))

    // ... tus rutas
    return r
}


===========================================
CÓMO REMOVER EL MIDDLEWARE
===========================================

Para remover completamente el middleware de logging:

1. Comentar la línea en main_router.go:
   // r.Use(middleware.Logger(middleware.DefaultLogConfig()))

2. O eliminarla completamente:
   (simplemente borra la línea)

3. Opcionalmente, eliminar el import si no se usa:
   // "github.com/prunus/pkg/middleware"

No es necesario tocar ningún otro archivo. El middleware es completamente
modular y no afecta el funcionamiento de los handlers.

*/

// Ejemplos de configuraciones predefinidas adicionales

// VerboseLogConfig retorna una configuración con máximo detalle
func VerboseLogConfig() *LogConfig {
	return &LogConfig{
		SkipPaths:      []string{},
		LogHeaders:     true,
		LogQueryParams: true,
		LogUserAgent:   true,
		TimeFormat:     "2006-01-02 15:04:05.000000",
		Output:         os.Stdout,
		UseJSON:        true,
		ColorOutput:    false,
	}
}

// ProductionLogConfig retorna una configuración optimizada para producción
func ProductionLogConfig() *LogConfig {
	return &LogConfig{
		SkipPaths:      []string{"/health", "/healthz", "/metrics"},
		LogHeaders:     false,
		LogQueryParams: true,
		LogUserAgent:   true,
		TimeFormat:     "2006-01-02T15:04:05Z07:00",
		Output:         os.Stdout,
		UseJSON:        true,
		ColorOutput:    false,
	}
}

// DevelopmentLogConfig retorna una configuración para desarrollo
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
