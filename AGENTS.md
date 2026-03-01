# AGENTS.md - Guía para Agentes de Código

Este archivo contiene las directrices y comandos esenciales para que los agentes de código trabajen eficientemente en este proyecto Go (Prunus).

### Arquitectura del Proyecto

Este es un proyecto REST API en Go que sigue el patrón de Clean Architecture con las siguientes capas:

- `cmd/` - Punto de entrada de la aplicación
- `pkg/models/` - Modelos de datos y entidades del dominio
- `pkg/dto/` - Data Transfer Objects para requests/responses
- `pkg/store/` - Capa de acceso a datos (repository pattern)
- `pkg/services/` - Lógica de negocio y validaciones
- `pkg/transport/http/` - Handlers HTTP y controladores
- `pkg/routers/` - Configuración de rutas y middleware
- `pkg/middleware/` - Middleware de autenticación, logging, etc.
- `pkg/helper/` - Utilidades (JWT, hashing, etc.)
- `pkg/config/database/` - Configuración de base de datos

## 🚀 Comandos Esenciales

### Ejecución
```bash
# Ejecutar la aplicación
go run cmd/main.go

# Construir binario
go build -o prunus cmd/main.go
```

### Dependencias
```bash
# Descargar dependencias
go mod download

# Actualizar dependencias
go mod tidy

# Verificar dependencias
go mod verify
```

### Pruebas
```bash
# Ejecutar todas las pruebas
go test ./...

# Ejecutar pruebas con cobertura
go test -cover ./...

# Ejecutar una prueba específica
go test ./pkg/services -run TestCreateUsuario

# Ejecutar pruebas en un paquete específico
go test ./pkg/services/
```

### Formateo y Linting
```bash
# Formatear código
go fmt ./...

# Verificar formato
go vet ./...

# Ejecutar golint (si está instalado)
golint ./...
```

## 📋 Convenciones de Código

### Imports
- Agrupar imports en tres bloques: estándar, terceros, locales
- Usar alias para paquetes largos cuando sea apropiado
```go
import (
    "encoding/json"
    "net/http"
    
    "github.com/go-chi/chi/v5"
    "github.com/golang-jwt/jwt/v5"
    
    "github.com/prunus/pkg/models"
    "github.com/prunus/pkg/services"
)
```

### Nomenclatura
- **Structs**: PascalCase (ej: `Usuario`, `EmpresaHandler`)
- **Funciones**: PascalCase para exportadas, camelCase para privadas
- **Variables**: camelCase (ej: `idUsuario`, `emailUsuario`)
- **Constantes**: UPPER_SNAKE_CASE (ej: `MAX_CONNECTIONS`)
- **Interfaces**: Prefijo "Store" para repositories (ej: `StoreUsuario`)
- **Nombres de paquetes**: lowercase, una palabra (ej: `models`, `services`)

### Estructura de Archivos
1. Package declaration
2. Imports (agrupados)
3. Types y constants
4. Constructor functions (`New...`)
5. Methods (agrupados por funcionalidad)
6. Funciones privadas de soporte

### JSON Tags
- Usar snake_case para campos JSON
- Incluir `omitempty` para campos opcionales
```go
type Usuario struct {
    IDUsuario  uint   `json:"id_usuario"`
    UsuEmail   string `json:"email"`
    CreatedAt  time.Time `json:"created_at"`
    DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}
```

## 🔧 Patrones y Convenciones

### Dependency Injection
- Usar constructor pattern con `New...` functions
- Inyectar dependencias a través de interfaces
```go
func NewUsuarioHandler(service *services.ServiceUsuario) *UsuarioHandler {
    return &UsuarioHandler{service: service}
}
```

### Manejo de Errores
- Retornar errores explícitos en español
- Usar `errors.New()` para errores simples
- Incluir contexto en mensajes de error
```go
if usuario.UsuEmail == "" {
    return nil, errors.New("el email del usuario es requerido")
}
```

### Validaciones
- Implementar validaciones en la capa de services
- Usar regex para validaciones de formato
- Validar campos obligatorios antes de operaciones CRUD
```go
emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
if !emailRegex.MatchString(usuario.UsuEmail) {
    return nil, errors.New("el formato del email es inválido")
}
```

### Base de Datos
- Usar pgx/v5 como driver de PostgreSQL
- Implementar soft delete con campo `deleted_at`
- Usar transacciones cuando sea necesario
- SQL queries como strings con formato consistente

### Autenticación
- Usar JWT tokens para autenticación
- Implementar middleware `RequireAuth()`
- Hashear passwords con bcrypt
- Limpiar passwords de objetos antes de retornarlos

### HTTP Handlers
- Setear `Content-Type: application/json` en todos los handlers
- Usar códigos de estado HTTP apropiados
- Decodificar JSON con validación de errores
```go
w.Header().Set("Content-Type", "application/json")
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    http.Error(w, "JSON inválido", http.StatusBadRequest)
    return
}
```

### Rutas
- Usar Chi router v5
- Versionar API con `/api/v1/`
- Agrupar rutas protegidas con middleware de autenticación
- Usar nombres de recursos en plural (ej: `/usuarios`, `/empresas`)

## 🎯 Buenas Prácticas

### Logging
- Usar el middleware de logging configurado
- No hacer log de información sensible (passwords, tokens)
- Incluir contexto relevante en los logs

### Seguridad
- Nunca exponer passwords en respuestas JSON
- Validar todos los inputs del usuario
- Usar parámetros en lugar de concatenar SQL queries
- Implementar rate limiting si es necesario

### Testing
- Escribir tests unitarios para la lógica de negocio
- Mockear dependencias externas (base de datos)
- Usar table-driven tests para múltiples casos
- Mantener cobertura de código > 80%

### Performance
- Usar connection pooling con pgx
- Implementar paginación para listas grandes
- Cerrar conexiones a base de datos con `defer`
- Evitar N+1 queries en consultas relacionales

## 🚨 Consideraciones Especiales

- **Soft Delete**: Todos los deletes deben actualizar `deleted_at` en lugar de eliminar físicamente
- **Campos de Auditoría**: Incluir `created_at`, `updated_at`, `deleted_at` en todas las tablas
- **Idioma**: Todos los mensajes de error y comentarios deben estar en español
- **Estado**: Usar campo `estado` (1=activo, 0=inactivo) además de soft delete
- **Relaciones**: Incluir objetos relacionados (ej: `Rol` en `Usuario`) cuando sea relevante

## 📦 Dependencias Principales

- `github.com/go-chi/chi/v5` - Router HTTP
- `github.com/jackc/pgx/v5` - Driver PostgreSQL
- `github.com/golang-jwt/jwt/v5` - Tokens JWT
- `github.com/joho/godotenv` - Variables de entorno
- `golang.org/x/crypto` - Funciones criptográficas (bcrypt)

## 🔍 Depuración

- Usar `fmt.Printf()` o `log.Printf()` para debugging temporal
- Revisar logs del middleware para trace de peticiones
- Usar `go mod why` para entender dependencias
- Verificar conexión a base de datos con logs de error

---

**Nota**: Esta guía está diseñada para ser usada por agentes de código. Mantener este archivo actualizado cuando se añadan nuevas convenciones o patrones al proyecto.