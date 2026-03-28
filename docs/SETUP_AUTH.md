# Guía Rápida: Configurar y Probar Autenticación

Este documento describe los pasos para levantar el servicio y configurar el acceso inicial, considerando que el sistema utiliza **UUID v4** para todos sus identificadores y una base de datos PostgreSQL 15.

## Pasos para Configurar el Sistema

### 1. Configurar Variables de Entorno

Asegúrate de tener un archivo `.env` en la raíz del proyecto con la configuración de la base de datos y JWT:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=prunus_db
DB_SSLMODE=disable

# JWT Configuration
JWT_SECRET=cambia_esto_por_una_clave_super_segura_de_al_menos_32_caracteres_random
JWT_EXPIRATION_HOURS=24
```

### 2. Iniciar el Servidor

El servidor realiza las migraciones automáticas al iniciar.

```bash
go run cmd/main.go
```

Deberías ver en los logs:
```
✅ Conexión a la base de datos establecida
✅ Migraciones completadas
✅ Iniciando servidor en puerto :9090
```

---

## Pruebas Rápidas e Inicialización de Datos

Debido a que el sistema utiliza UUIDs y llaves foráneas, es necesario insertar los datos maestros iniciales (Empresa, Sucursal, Rol) antes de crear el primer usuario.

### Paso 1: Inicializar Datos en la Base de Datos (SQL)

Ejecuta el siguiente script en tu cliente de base de datos (psql, DBeaver, etc.):

```sql
-- 1. Insertar Empresa (ID: 7f7b...)
INSERT INTO empresa (id_empresa, nombre, rut, id_status)
VALUES ('7f7b0e11-1234-4a21-9591-316279f06742', 'Empresa Matriz', '800123456-7', '59039503-85CF-E511-80C1-000C29C9E0E0');

-- 2. Insertar Sucursal (ID: a3b4...) vinculado a la Empresa
INSERT INTO sucursal (id_sucursal, id_empresa, nombre_sucursal, id_status)
VALUES ('a3b4c5d6-e7f8-4a1b-9c2d-3e4f5a6b7c8d', '7f7b0e11-1234-4a21-9591-316279f06742', 'Sucursal Central', '59039503-85CF-E511-80C1-000C29C9E0E0');

-- 3. Insertar Rol (ID: b4c5...) vinculado a la Sucursal
INSERT INTO rol (id_rol, nombre_rol, id_sucursal, id_status)
VALUES ('b4c5d6e7-f8a1-4b2c-9d3e-4f5a6b7c8d9e', 'Administrador', 'a3b4c5d6-e7f8-4a1b-9c2d-3e4f5a6b7c8d', '59039503-85CF-E511-80C1-000C29C9E0E0');

-- 4. Crear el usuario administrador (password: Admin123)
-- El hash corresponde a "Admin123" generado con bcrypt
INSERT INTO usuario (id_usuario, id_sucursal, id_rol, email, usu_nombre, usu_dni, usu_telefono, password, id_status)
VALUES (
  'c5d6e7f8-a1b2-4c3d-9e4f-5a6b7c8d9e0f',
  'a3b4c5d6-e7f8-4a1b-9c2d-3e4f5a6b7c8d',
  'b4c5d6e7-f8a1-4b2c-9d3e-4f5a6b7c8d9e',
  'admin@prunus.com',
  'Administrador Sistema',
  '12345678',
  '+51999888777',
  '$2a$10$U.sUS/qwAXlDPrJZ9wAaLe78DmRtcnWVY39wFp85YLiL0iIVPVkkK',
  '59039503-85CF-E511-80C1-000C29C9E0E0'
);
```

---

### Paso 2: Hacer Login para Obtener Token

```bash
curl -X POST http://localhost:9090/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@prunus.com",
    "password": "Admin123"
  }'
```

**Respuesta esperada:**

```json
{
  "status": "success",
  "message": "Inicio de sesión exitoso",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "usuario": {
      "id_usuario": "c5d6e7f8-a1b2-4c3d-9e4f-5a6b7c8d9e0f",
      "email": "admin@prunus.com",
      "usu_nombre": "Administrador Sistema"
    },
    "expires_at": "2026-03-26T15:00:00Z"
  }
}
```

**Copia el valor del token.**

---

### Paso 3: Probar Acceso a Endpoints Protegidos

Reemplaza `<TU_TOKEN>` con el token obtenido:

```bash
# Listar inventario (requiere token)
curl -X GET http://localhost:9090/api/v1/inventario \
  -H "Authorization: Bearer <TU_TOKEN>"
```

---

## Verificación de Configuración

### ✅ Checklist de Solución de Problemas

1. **UUIDs en las peticiones:** Asegúrate de enviar los IDs en formato UUID (ej: `550e8400-e29b-41d4-a716-446655440000`) y no números enteros.
2. **Estatus ACTIVO:** Todos los registros (Empresa, Sucursal, Rol, Usuario) deben tener el `id_status` correspondiente a `EstatusActivo` (`59039503-85CF-E511-80C1-000C29C9E0E0`).
3. **JWT_SECRET:** Si cambias esta variable en el `.env`, todos los tokens anteriores dejarán de ser válidos.
4. **Contexto de Usuario:** El sistema inyecta automáticamente `user_id`, `user_sucursal` y `user_rol` en el contexto después de validar el token.

---

## Documentación Adicional
- [API.md](./api.md) - Referencia completa de endpoints.
- [DATABASE.md](./DATABASE.md) - Esquema y diseño de tablas.
