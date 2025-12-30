# Sistema de Autenticación - Prunus API

Documentación completa del sistema de login y autenticación con JWT para la API de Prunus.

## Tabla de Contenidos

1. [Resumen](#resumen)
2. [Arquitectura](#arquitectura)
3. [Configuración](#configuración)
4. [Endpoints de Autenticación](#endpoints-de-autenticación)
5. [Uso de Tokens](#uso-de-tokens)
6. [Ejemplos de Uso](#ejemplos-de-uso)
7. [Seguridad](#seguridad)
8. [Troubleshooting](#troubleshooting)

---

## Resumen

El sistema de autenticación implementa:

- ✅ **Login con email y password**
- ✅ **Generación de JWT tokens**
- ✅ **Validación de tokens en rutas protegidas**
- ✅ **Logout (client-side)**
- ✅ **Refresh token**
- ✅ **Información del usuario autenticado**
- ✅ **Middleware de autorización por roles**

**Tecnologías utilizadas:**
- JWT (JSON Web Tokens) - `github.com/golang-jwt/jwt/v5`
- Bcrypt para passwords - `golang.org/x/crypto/bcrypt`
- Chi Router Middleware

---

## Arquitectura

### Flujo de Autenticación

```
┌─────────────┐
│   Cliente   │
│             │
│ 1. POST     │
│  /login     │
│  {email,    │
│   password} │
└──────┬──────┘
       │
       ▼
┌─────────────────────────────┐
│   AuthHandler.Login()       │
│                             │
│ - Decodifica request        │
│ - Valida formato            │
└──────┬──────────────────────┘
       │
       ▼
┌─────────────────────────────┐
│ ServiceUsuario.Authenticate │
│                             │
│ - Busca usuario por email   │
│ - Verifica estado activo    │
│ - Verifica password bcrypt  │
│ - Verifica rol activo       │
└──────┬──────────────────────┘
       │
       ▼
┌─────────────────────────────┐
│   helper.GenerateToken()    │
│                             │
│ - Crea JWT con claims       │
│ - Firma con JWT_SECRET      │
│ - Expira en 24h (default)   │
└──────┬──────────────────────┘
       │
       ▼
┌─────────────┐
│   Cliente   │
│             │
│ Recibe:     │
│ - token     │
│ - usuario   │
│ - expires_at│
└─────────────┘
```

### Estructura de Archivos

```
pkg/
├── models/
│   └── auth.go              # LoginRequest, LoginResponse, JWTClaims
├── helper/
│   ├── hashear.go           # Bcrypt (ya existía)
│   └── jwt.go               # GenerateToken, ValidateToken, RefreshToken
├── store/
│   └── usuario_store.go     # GetUsuarioByEmail()
├── services/
│   └── usuario_services.go  # AuthenticateUsuario()
├── transport/http/
│   └── auth_handler.go      # Login, Logout, GetMe, RefreshToken
├── middleware/
│   └── auth.go              # RequireAuth, RequireRole, OptionalAuth
└── routers/
    └── main_router.go       # Rutas públicas y protegidas
```

---

## Configuración

### 1. Variables de Entorno

Agrega estas variables a tu archivo `.env`:

```env
# JWT Configuration
JWT_SECRET=tu_clave_secreta_muy_larga_y_segura_minimo_32_caracteres
JWT_EXPIRATION_HOURS=24
```

**IMPORTANTE:**
- `JWT_SECRET` debe ser una cadena aleatoria y segura
- En producción, usa al menos 32 caracteres
- NUNCA subas el `.env` a control de versiones
- Cambia el secret en cada ambiente (dev, staging, prod)

### 2. Generar JWT_SECRET Seguro

```bash
# Opción 1: OpenSSL
openssl rand -base64 32

# Opción 2: En Go
import "crypto/rand"
import "encoding/base64"

b := make([]byte, 32)
rand.Read(b)
secret := base64.StdEncoding.EncodeToString(b)
```

### 3. Verificar Configuración

Asegúrate de que tu `.env` tenga todos estos campos:

```env
# Base de Datos
DB_HOST=localhost
DB_USER=tu_usuario
DB_PASSWORD=tu_password
DB_NAME=prunus_db
DB_PORT=5432
DB_SSLMODE=disable

# JWT
JWT_SECRET=abc123xyz789...
JWT_EXPIRATION_HOURS=24
```

---

## Endpoints de Autenticación

### 1. Login (Iniciar Sesión)

**Endpoint:** `POST /api/v1/login`
**Autenticación:** No requerida
**Descripción:** Autentica un usuario y retorna un JWT token

#### Request

```http
POST /api/v1/login HTTP/1.1
Host: localhost:9090
Content-Type: application/json

{
  "email": "usuario@ejemplo.com",
  "password": "MiPassword123"
}
```

#### Response Exitosa (200 OK)

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZF91c3VhcmlvIjoxLCJlbWFpbCI6InVzdWFyaW9AZWplbXBsby5jb20iLCJpZF9yb2wiOjIsInJvbF9ub21icmUiOiJBZG1pbmlzdHJhZG9yIiwiaWRfc3VjdXJzYWwiOjUsImV4cCI6MTczNTYxNzA0NSwiaWF0IjoxNzM1NTMwNjQ1LCJuYmYiOjE3MzU1MzA2NDUsImlzcyI6InBydW51cy1hcGkiLCJzdWIiOiIxIn0.abc123xyz...",
  "usuario": {
    "id_usuario": 1,
    "id_sucursal": 5,
    "usu_email": "usuario@ejemplo.com",
    "usu_nombre": "Juan Pérez",
    "usu_dni": "12345678",
    "usu_telefono": "+51999888777",
    "estado": 1,
    "rol": {
      "id_rol": 2,
      "rol_nombre": "Administrador",
      "estado": 1
    },
    "sucursal": {
      "id_sucursal": 5,
      "nombre_sucursal": "Sucursal Central",
      "estado": 1
    },
    "created_at": "2025-12-01T10:00:00Z",
    "updated_at": "2025-12-30T10:00:00Z"
  },
  "expires_at": 1735617045
}
```

#### Errores Posibles

| Código | Descripción |
|--------|-------------|
| `400` | Formato de petición inválido |
| `401` | Credenciales inválidas (email o password incorrectos) |
| `401` | Usuario inactivo |
| `401` | Rol del usuario inactivo |

**Ejemplo de Error (401):**

```json
{
  "error": "credenciales inválidas"
}
```

---

### 2. Get Me (Usuario Actual)

**Endpoint:** `GET /api/v1/me`
**Autenticación:** Requerida
**Descripción:** Obtiene información del usuario autenticado actual

#### Request

```http
GET /api/v1/me HTTP/1.1
Host: localhost:9090
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

#### Response Exitosa (200 OK)

```json
{
  "id_usuario": 1,
  "id_sucursal": 5,
  "usu_email": "usuario@ejemplo.com",
  "usu_nombre": "Juan Pérez",
  "usu_dni": "12345678",
  "usu_telefono": "+51999888777",
  "estado": 1,
  "rol": {
    "id_rol": 2,
    "rol_nombre": "Administrador",
    "estado": 1
  },
  "sucursal": {
    "id_sucursal": 5,
    "nombre_sucursal": "Sucursal Central",
    "estado": 1
  }
}
```

---

### 3. Logout (Cerrar Sesión)

**Endpoint:** `POST /api/v1/logout`
**Autenticación:** Requerida
**Descripción:** Cierra la sesión del usuario (client-side)

#### Request

```http
POST /api/v1/logout HTTP/1.1
Host: localhost:9090
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

#### Response Exitosa (200 OK)

```json
{
  "message": "Sesión cerrada exitosamente"
}
```

**Nota:** En un sistema JWT stateless, el logout se maneja principalmente en el cliente eliminando el token. El servidor confirma el logout pero no invalida el token (a menos que implementes una blacklist).

---

### 4. Refresh Token (Renovar Token)

**Endpoint:** `POST /api/v1/refresh-token`
**Autenticación:** Requerida
**Descripción:** Genera un nuevo token basado en el token actual

#### Request

```http
POST /api/v1/refresh-token HTTP/1.1
Host: localhost:9090
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

#### Response Exitosa (200 OK)

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.nuevo_token...",
  "expires_at": 1735703445
}
```

---

## Uso de Tokens

### Estructura del Token JWT

El token contiene los siguientes claims:

```json
{
  "id_usuario": 1,
  "email": "usuario@ejemplo.com",
  "id_rol": 2,
  "rol_nombre": "Administrador",
  "id_sucursal": 5,
  "exp": 1735617045,
  "iat": 1735530645,
  "nbf": 1735530645,
  "iss": "prunus-api",
  "sub": "1"
}
```

**Claims estándar:**
- `exp` (Expiration): Cuándo expira el token
- `iat` (Issued At): Cuándo se generó el token
- `nbf` (Not Before): Desde cuándo es válido
- `iss` (Issuer): Quién emitió el token
- `sub` (Subject): ID del usuario

**Claims personalizados:**
- `id_usuario`: ID del usuario
- `email`: Email del usuario
- `id_rol`: ID del rol
- `rol_nombre`: Nombre del rol
- `id_sucursal`: ID de la sucursal

### Cómo Usar el Token

#### 1. Guardar el Token (Cliente)

Después del login exitoso:

```javascript
// En el navegador (localStorage)
localStorage.setItem('token', response.token);
localStorage.setItem('user', JSON.stringify(response.usuario));

// O en sessionStorage para sesión temporal
sessionStorage.setItem('token', response.token);
```

#### 2. Enviar el Token en Requests

Incluye el token en el header `Authorization`:

```http
GET /api/v1/empresas HTTP/1.1
Host: localhost:9090
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Formato:** `Bearer <token>`

#### 3. Verificar Expiración (Cliente)

```javascript
const token = localStorage.getItem('token');
const expiresAt = localStorage.getItem('expires_at');

if (Date.now() / 1000 > expiresAt) {
  // Token expirado, hacer login de nuevo
  redirectToLogin();
}
```

---

## Ejemplos de Uso

### Ejemplo con cURL

#### Login

```bash
curl -X POST http://localhost:9090/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@prunus.com",
    "password": "MiPassword123"
  }'
```

#### Acceder a Recurso Protegido

```bash
# Primero guarda el token
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# Luego úsalo en requests
curl -X GET http://localhost:9090/api/v1/empresas \
  -H "Authorization: Bearer $TOKEN"
```

#### Get Me

```bash
curl -X GET http://localhost:9090/api/v1/me \
  -H "Authorization: Bearer $TOKEN"
```

---

### Ejemplo con JavaScript (Fetch API)

```javascript
// 1. Login
async function login(email, password) {
  const response = await fetch('http://localhost:9090/api/v1/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ email, password }),
  });

  if (!response.ok) {
    throw new Error('Login fallido');
  }

  const data = await response.json();

  // Guardar token
  localStorage.setItem('token', data.token);
  localStorage.setItem('expires_at', data.expires_at);

  return data;
}

// 2. Hacer request autenticado
async function getEmpresas() {
  const token = localStorage.getItem('token');

  const response = await fetch('http://localhost:9090/api/v1/empresas', {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  });

  if (response.status === 401) {
    // Token expirado o inválido
    redirectToLogin();
    return;
  }

  return await response.json();
}

// 3. Logout
async function logout() {
  const token = localStorage.getItem('token');

  await fetch('http://localhost:9090/api/v1/logout', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  });

  // Limpiar storage
  localStorage.removeItem('token');
  localStorage.removeItem('expires_at');

  redirectToLogin();
}

// 4. Refresh token
async function refreshToken() {
  const token = localStorage.getItem('token');

  const response = await fetch('http://localhost:9090/api/v1/refresh-token', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  });

  if (!response.ok) {
    redirectToLogin();
    return;
  }

  const data = await response.json();
  localStorage.setItem('token', data.token);
  localStorage.setItem('expires_at', data.expires_at);
}
```

---

### Ejemplo con Postman

#### 1. Login

1. Crear nueva request POST
2. URL: `http://localhost:9090/api/v1/login`
3. Headers: `Content-Type: application/json`
4. Body (raw, JSON):
   ```json
   {
     "email": "admin@prunus.com",
     "password": "MiPassword123"
   }
   ```
5. Enviar
6. Copiar el `token` de la respuesta

#### 2. Configurar Token en Postman

**Opción A: Manual**
1. En cada request, ir a Headers
2. Agregar: `Authorization: Bearer <tu_token>`

**Opción B: Variables de Colección**
1. Crear variable `auth_token` en la colección
2. En la tab "Tests" del request de login:
   ```javascript
   pm.collectionVariables.set("auth_token", pm.response.json().token);
   ```
3. En otros requests, usar: `Authorization: Bearer {{auth_token}}`

#### 3. Probar Endpoints Protegidos

1. GET `http://localhost:9090/api/v1/me`
2. Headers: `Authorization: Bearer {{auth_token}}`
3. Enviar

---

## Seguridad

### Medidas Implementadas

1. **Passwords Hasheadas**
   - Bcrypt con costo 10
   - Nunca se almacenan passwords en texto plano
   - No se retornan passwords en las respuestas

2. **Tokens Firmados**
   - JWT firmados con HMAC-SHA256
   - Secret key almacenado en variables de entorno
   - Imposible falsificar sin el secret

3. **Expiración de Tokens**
   - Tokens expiran después de 24 horas (configurable)
   - Se verifica la expiración en cada request

4. **Validación de Estado**
   - Usuario debe estar activo (`estado = 1`)
   - Rol del usuario debe estar activo
   - Sucursal debe existir y estar activa

5. **Mensajes de Error Genéricos**
   - No revela si un email existe en la BD
   - Errores genéricos: "credenciales inválidas"
   - Previene enumeración de usuarios

6. **HTTPS Recomendado**
   - En producción, usa HTTPS
   - Previene interceptación de tokens
   - Configura CORS apropiadamente

### Mejores Prácticas

#### En el Servidor

- ✅ Usar JWT_SECRET largo y aleatorio (32+ caracteres)
- ✅ Cambiar JWT_SECRET en cada ambiente
- ✅ No loguear tokens en los logs
- ✅ Usar HTTPS en producción
- ✅ Configurar CORS restrictivamente
- ✅ Implementar rate limiting en /login

#### En el Cliente

- ✅ Almacenar tokens en localStorage o sessionStorage
- ✅ No almacenar en cookies sin httpOnly
- ✅ Verificar expiración antes de hacer requests
- ✅ Limpiar tokens al hacer logout
- ✅ Manejar respuestas 401 redirigiendo a login
- ✅ Usar HTTPS para todas las peticiones

### Vulnerabilidades Conocidas

❌ **NO hacer:**

1. **NO** usar `JWT_SECRET` débil como "secret" o "123456"
2. **NO** enviar el token en la URL (query params)
3. **NO** almacenar tokens en cookies sin `httpOnly`
4. **NO** loguear tokens en archivos de log
5. **NO** compartir el mismo secret entre ambientes
6. **NO** ignorar errores de validación de tokens

### Mejoras Futuras Opcionales

- 🔄 Implementar refresh tokens de larga duración
- 🚫 Implementar blacklist de tokens (para logout real)
- 🔐 Agregar autenticación de dos factores (2FA)
- 🔑 Implementar recuperación de contraseña
- 📊 Registrar intentos de login fallidos
- ⏱️ Implementar rate limiting
- 🔄 Rotación automática de JWT_SECRET

---

## Troubleshooting

### Error: "JWT_SECRET no configurado"

**Causa:** No existe la variable de entorno `JWT_SECRET`

**Solución:**
1. Agregar `JWT_SECRET` a tu archivo `.env`
2. Reiniciar el servidor
3. Verificar que godotenv carga el `.env` correctamente

```env
JWT_SECRET=tu_clave_secreta_aqui
```

---

### Error: "Token inválido"

**Causas posibles:**
1. Token mal formado
2. Token firmado con otro secret
3. Token manipulado

**Solución:**
1. Verificar que el token se envía en el formato correcto: `Bearer <token>`
2. Asegurarse de que el JWT_SECRET sea el mismo que se usó para generar el token
3. Hacer login nuevamente para obtener un token válido

---

### Error: "Token expirado"

**Causa:** El token ha excedido su tiempo de vida (24 horas por defecto)

**Solución:**
1. Hacer login nuevamente
2. O usar el endpoint `/refresh-token` antes de que expire

---

### Error: "Credenciales inválidas"

**Causas posibles:**
1. Email incorrecto
2. Password incorrecta
3. Usuario no existe
4. Usuario inactivo
5. Rol inactivo

**Solución:**
1. Verificar email y password
2. Verificar que el usuario esté activo en la BD: `SELECT * FROM usuario WHERE email = 'tu@email.com'`
3. Verificar que el rol esté activo: `SELECT * FROM rol WHERE id_rol = X`

---

### Las rutas protegidas no requieren autenticación

**Causa:** El middleware `RequireAuth()` no está aplicado

**Solución:**
1. Verificar que las rutas estén dentro del bloque `r.Group()`
2. Verificar que `r.Use(middleware.RequireAuth())` esté presente

```go
r.Group(func(r chi.Router) {
    r.Use(middleware.RequireAuth())  // ← Debe estar aquí

    r.Get("/empresas", empresaHandler.GetAll)
    // ...
})
```

---

### Error: "No se pudo obtener información del usuario"

**Causa:** Los claims no están en el contexto

**Solución:**
1. Asegurarse de que el middleware `RequireAuth()` se ejecuta antes del handler
2. Verificar el orden de los middlewares

---

## Información Adicional

### Archivos Creados/Modificados

```
✅ pkg/models/auth.go               (NUEVO)
✅ pkg/helper/jwt.go                (NUEVO)
✅ pkg/store/usuario_store.go       (MODIFICADO - agregado GetUsuarioByEmail)
✅ pkg/services/usuario_services.go (MODIFICADO - agregado AuthenticateUsuario)
✅ pkg/transport/http/auth_handler.go (NUEVO)
✅ pkg/middleware/auth.go           (NUEVO)
✅ pkg/routers/main_router.go       (MODIFICADO - rutas públicas/protegidas)
✅ cmd/main.go                      (MODIFICADO - inyección de AuthHandler)
✅ .env.example                     (MODIFICADO - agregado JWT_SECRET)
```

### Rutas Públicas (Sin Autenticación)

| Método | Ruta | Descripción |
|--------|------|-------------|
| POST | `/api/v1/login` | Iniciar sesión |

### Rutas Protegidas (Requieren Autenticación)

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/api/v1/me` | Usuario actual |
| POST | `/api/v1/logout` | Cerrar sesión |
| POST | `/api/v1/refresh-token` | Renovar token |
| GET/POST/PUT/DELETE | `/api/v1/empresas/*` | CRUD empresas |
| GET/POST/PUT/DELETE | `/api/v1/sucursal/*` | CRUD sucursales |
| GET/POST/PUT/DELETE | `/api/v1/rol/*` | CRUD roles |
| GET/POST/PUT/DELETE | `/api/v1/usuario/*` | CRUD usuarios |

---

## Contacto y Soporte

Para dudas o problemas con el sistema de autenticación, contacta al equipo de desarrollo.

**Recursos útiles:**
- [JWT.io](https://jwt.io/) - Decodificar y verificar tokens JWT
- [Bcrypt Calculator](https://bcrypt-generator.com/) - Generar/verificar hashes bcrypt
- Documentación del proyecto en `/docs`
