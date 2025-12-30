# Guía Rápida: Configurar y Probar Autenticación

## Pasos para Configurar el Sistema de Autenticación

### 1. Configurar Variables de Entorno

Edita tu archivo `.env` y agrega estas líneas:

```env
# JWT Configuration
JWT_SECRET=cambia_esto_por_una_clave_super_segura_de_al_menos_32_caracteres_random
JWT_EXPIRATION_HOURS=24
```

**IMPORTANTE:** Cambia `JWT_SECRET` por una clave única y segura.

### 2. Instalar Dependencias

Las dependencias ya están instaladas, pero si necesitas reinstalarlas:

```bash
go mod download
```

### 3. Iniciar el Servidor

```bash
cd c:\Users\Lenovo\Music\Go\prunus
go run cmd/main.go
```

Deberías ver:
```
✅ Iniciando servidor
```

El servidor estará corriendo en `http://localhost:9090`

---

## Pruebas Rápidas

### Paso 1: Crear un Usuario de Prueba

Primero necesitas crear un usuario en la base de datos. Puedes hacerlo de dos formas:

#### Opción A: Usando el endpoint de creación (temporalmente sin protección)

Si necesitas crear el primer usuario admin, temporalmente puedes comentar el middleware de autenticación en `main_router.go`:

```go
// Comentar temporalmente esta línea:
// r.Use(middleware.RequireAuth())
```

Luego crear el usuario:

```bash
curl -X POST http://localhost:9090/api/v1/usuario \
  -H "Content-Type: application/json" \
  -d '{
    "id_sucursal": 1,
    "rol": {
      "id_rol": 1
    },
    "usu_email": "admin@prunus.com",
    "usu_nombre": "Administrador",
    "usu_dni": "12345678",
    "usu_telefono": "+51999888777",
    "usu_password": "Admin123",
    "estado": 1
  }'
```

**No olvides volver a descomentar el middleware después.**

#### Opción B: Directamente en la Base de Datos

```sql
-- Primero, asegúrate de tener una sucursal y un rol
INSERT INTO sucursal (nombre_sucursal, direccion, telefono, estado)
VALUES ('Sucursal Central', 'Av. Principal 123', '+51999000111', 1);

INSERT INTO rol (nombre_rol, id_sucursal, estado)
VALUES ('Administrador', 1, 1);

-- Luego crea el usuario (password: Admin123)
-- El hash es para "Admin123"
INSERT INTO usuario (id_sucursal, id_rol, email, usu_nombre, usu_dni, usu_telefono, password, estado)
VALUES (
  1,
  1,
  'admin@prunus.com',
  'Administrador',
  '12345678',
  '+51999888777',
  '$2a$10$rZ5qX8KYqJxY9vxF8F0qLu.YB8G0qH4vMqXVZ4U3B9K3r5nH6K2Gm',
  1
);
```

---

### Paso 2: Hacer Login

```bash
curl -X POST http://localhost:9090/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@prunus.com",
    "password": "Admin123"
  }'
```

**Respuesta esperada:**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "usuario": {
    "id_usuario": 1,
    "usu_email": "admin@prunus.com",
    "usu_nombre": "Administrador",
    ...
  },
  "expires_at": 1735617045
}
```

**Copia el token** de la respuesta.

---

### Paso 3: Probar Endpoint Protegido

Reemplaza `<TU_TOKEN>` con el token que obtuviste:

```bash
curl -X GET http://localhost:9090/api/v1/me \
  -H "Authorization: Bearer <TU_TOKEN>"
```

**Respuesta esperada:**

```json
{
  "id_usuario": 1,
  "usu_email": "admin@prunus.com",
  "usu_nombre": "Administrador",
  ...
}
```

---

### Paso 4: Probar Acceso sin Token (Debe Fallar)

```bash
curl -X GET http://localhost:9090/api/v1/empresas
```

**Respuesta esperada:**

```
Token de autenticación requerido
```

Status: `401 Unauthorized`

---

### Paso 5: Probar Acceso con Token

```bash
curl -X GET http://localhost:9090/api/v1/empresas \
  -H "Authorization: Bearer <TU_TOKEN>"
```

**Respuesta esperada:**

```json
[
  {
    "id_empresa": 1,
    "nombre_empresa": "...",
    ...
  }
]
```

---

## Verificación de Configuración

### ✅ Checklist

- [ ] Archivo `.env` tiene `JWT_SECRET` configurado
- [ ] Archivo `.env` tiene `JWT_EXPIRATION_HOURS` configurado
- [ ] Servidor inicia sin errores
- [ ] Existe al menos un usuario en la BD
- [ ] Usuario tiene un rol asignado
- [ ] Rol está activo (`estado = 1`)
- [ ] Usuario está activo (`estado = 1`)
- [ ] Login retorna un token válido
- [ ] Endpoint `/api/v1/me` funciona con token
- [ ] Endpoints protegidos rechazan requests sin token

---

## Endpoints Disponibles

### Públicos (Sin Token)

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| POST | `/api/v1/login` | Iniciar sesión |

### Protegidos (Requieren Token)

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| GET | `/api/v1/me` | Información del usuario actual |
| POST | `/api/v1/logout` | Cerrar sesión |
| POST | `/api/v1/refresh-token` | Renovar token |
| GET | `/api/v1/empresas` | Listar empresas |
| POST | `/api/v1/empresas` | Crear empresa |
| GET | `/api/v1/empresas/{id}` | Obtener empresa |
| PUT | `/api/v1/empresas/{id}` | Actualizar empresa |
| DELETE | `/api/v1/empresas/{id}` | Eliminar empresa |
| ... | `/api/v1/sucursal/*` | CRUD sucursales |
| ... | `/api/v1/rol/*` | CRUD roles |
| ... | `/api/v1/usuario/*` | CRUD usuarios |

---

## Ejemplo Completo: Flujo de Trabajo

```bash
# 1. Login
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:9090/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@prunus.com","password":"Admin123"}')

# 2. Extraer token (requiere jq)
TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.token')

echo "Token: $TOKEN"

# 3. Obtener información del usuario actual
curl -X GET http://localhost:9090/api/v1/me \
  -H "Authorization: Bearer $TOKEN"

# 4. Listar empresas
curl -X GET http://localhost:9090/api/v1/empresas \
  -H "Authorization: Bearer $TOKEN"

# 5. Crear nueva empresa
curl -X POST http://localhost:9090/api/v1/empresas \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "nombre_empresa": "Mi Empresa",
    "ruc": "20123456789",
    "direccion": "Av. Test 123",
    "telefono": "+51999888777",
    "estado": 1
  }'

# 6. Logout
curl -X POST http://localhost:9090/api/v1/logout \
  -H "Authorization: Bearer $TOKEN"
```

---

## Solución de Problemas Comunes

### "JWT_SECRET no configurado"

**Solución:**
1. Abre `.env`
2. Agrega: `JWT_SECRET=tu_clave_secreta_aqui`
3. Reinicia el servidor

### "Credenciales inválidas"

**Causas:**
- Email o password incorrectos
- Usuario no existe
- Usuario inactivo
- Rol inactivo

**Solución:**
```sql
-- Verificar usuario
SELECT u.*, r.*
FROM usuario u
LEFT JOIN rol r ON u.id_rol = r.id_rol
WHERE u.email = 'admin@prunus.com';

-- Debe retornar el usuario con estado=1 y rol con estado=1
```

### "Token de autenticación requerido"

**Causa:** No se envió el token o está mal formado

**Solución:**
- Asegúrate de incluir el header: `Authorization: Bearer <token>`
- Verifica que no haya espacios extras
- El formato correcto es: `Bearer ` + token (con un espacio)

### "Token inválido"

**Causas:**
- Token manipulado
- JWT_SECRET diferente al usado para generarlo
- Token mal formado

**Solución:**
- Hacer login de nuevo para obtener un token válido
- Verificar que JWT_SECRET sea el mismo

### "Token expirado"

**Causa:** El token excedió las 24 horas

**Solución:**
- Hacer login de nuevo
- O usar `/refresh-token` antes de que expire

---

## Siguiente Paso

Una vez que todo funciona correctamente:

1. **Crear más usuarios** con diferentes roles
2. **Implementar autorización por roles** usando `middleware.RequireRole()`
3. **Configurar CORS** si vas a usar desde un frontend
4. **Configurar HTTPS** en producción
5. **Implementar rate limiting** en el endpoint de login

---

## Documentación Completa

Para más detalles, consulta:
- [AUTENTICACION.md](./AUTENTICACION.md) - Documentación completa del sistema

---

## Características Implementadas

✅ Login con email y password
✅ Generación de JWT tokens
✅ Validación de tokens
✅ Rutas protegidas con middleware
✅ Endpoint de usuario actual (/me)
✅ Logout
✅ Refresh token
✅ Middleware opcional de roles
✅ Verificación de estado de usuario y rol
✅ Passwords hasheadas con bcrypt
✅ Claims personalizados en JWT
✅ Extracción de IP del cliente
✅ Manejo de errores de autenticación

---

**¡Sistema de autenticación completamente funcional! 🎉**
