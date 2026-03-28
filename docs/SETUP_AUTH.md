# Guía Rápida: Configurar y Probar Autenticación

Este documento describe los pasos para levantar el servicio y configurar el acceso inicial, considerando que el sistema utiliza **UUID v4** para todos sus identificadores y una arquitectura de seguridad basada en JWT.

## Pasos para Configurar el Sistema

### 1. Configurar Variables de Entorno

Asegúrate de tener un archivo `.env` en la raíz del proyecto:

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

El servidor realiza las migraciones automáticas al iniciar, incluyendo la creación de los catálogos maestros de estatus.

```bash
go run cmd/main.go serve
```

---

## Pruebas Rápidas e Inicialización de Datos

Debido a que el sistema utiliza UUIDs y llaves foráneas, es necesario insertar los datos maestros iniciales (Empresa, Sucursal, Rol) antes de crear el primer usuario.

### Opción A: Inicialización mediante CLI (Recomendado)

Prunus incluye una herramienta CLI integrada para registrar las entidades base de forma segura.

#### 1. Registrar Empresa
```bash
go run cmd/main.go register empresa --nombre "Empresa Matriz" --rut "800123456-7" --status "fc273a6a-ab7b-4453-a560-ac62fa64348b"
```

#### 2. Registrar Sucursal (Usa el ID devuelto arriba)
```bash
go run cmd/main.go register sucursal --empresa "<ID_EMPRESA>" --nombre "Sucursal Central" --status "6cf06fbe-b21c-46e3-a34b-b24f5167cd9a"
```

#### 3. Registrar Rol
```bash
go run cmd/main.go register rol --sucursal "<ID_SUCURSAL>" --nombre "Administrador" --status "fc273a6a-ab7b-4453-a560-ac62fa64348b"
```

#### 4. Registrar Usuario Administrador (Password: Admin123)
```bash
go run cmd/main.go register usuario --sucursal "<ID_SUCURSAL>" --rol "<ID_ROL>" --email "admin@prunus.com" --nombre "Admin Master" --dni "12345678" --password "Admin123" --status "3a99d245-b34f-48a5-ac08-a5a010c5822f"
```

---

### Opción B: Inicialización Manual vía SQL (Directo en DB)

Si prefieres usar un cliente como DBeaver o psql, usa este script sincronizado con los UUIDs de `012_estatus.go`:

```sql
-- 1. Insertar Empresa (Estatus 'Activa': fc273a6a...)
INSERT INTO empresa (id_empresa, nombre, rut, id_status)
VALUES (
  '7f7b0e11-1234-4a21-9591-316279f06742', 
  'Empresa Matriz', 
  '800123456-7', 
  'fc273a6a-ab7b-4453-a560-ac62fa64348b'
);

-- 2. Insertar Sucursal (Estatus 'Abierta': 6cf06fbe...)
INSERT INTO sucursal (id_sucursal, id_empresa, nombre_sucursal, id_status)
VALUES (
  'a3b4c5d6-e7f8-4a1b-9c2d-3e4f5a6b7c8d', 
  '7f7b0e11-1234-4a21-9591-316279f06742', 
  'Sucursal Central', 
  '6cf06fbe-b21c-46e3-a34b-b24f5167cd9a'
);

-- 3. Insertar Rol (Estatus 'Activa': fc273a6a...)
INSERT INTO rol (id_rol, nombre_rol, id_sucursal, id_status)
VALUES (
  'b4c5d6e7-f8a1-4b2c-9d3e-4f5a6b7c8d9e', 
  'Administrador', 
  'a3b4c5d6-e7f8-4a1b-9c2d-3e4f5a6b7c8d', 
  'fc273a6a-ab7b-4453-a560-ac62fa64348b'
);

-- 4. Crear el usuario administrador (password: Admin123)
-- El hash corresponde a "Admin123" (Estatus 'Activo': 3a99d245...)
INSERT INTO usuario (id_usuario, id_sucursal, id_rol, email, usu_nombre, usu_dni, password, id_status)
VALUES (
  'c5d6e7f8-a1b2-4c3d-9e4f-5a6b7c8d9e0f',
  'a3b4c5d6-e7f8-4a1b-9c2d-3e4f5a6b7c8d',
  'b4c5d6e7-f8a1-4b2c-9d3e-4f5a6b7c8d9e',
  'admin@prunus.com',
  'Admin Master',
  '12345678',
  '$2a$10$U.sUS/qwAXlDPrJZ9wAaLe78DmRtcnWVY39wFp85YLiL0iIVPVkkK',
  '3a99d245-b34f-48a5-ac08-a5a010c5822f'
);
```

---

### Paso 3: Hacer Login para Obtener Token

```bash
curl -X POST http://localhost:9090/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@prunus.com", "password": "Admin123"}'
```

---

## Verificación de Configuración

### ✅ Checklist de Seguridad e Integridad

1. **Contexto de Usuario:** El sistema inyecta automáticamente `user_id`, `user_sucursal` y `user_rol` tras validar el token.
2. **Estatus Coherente:** El sistema utiliza disparadores (Triggers) para asegurar que solo se asignen estatus válidos según el módulo.
3. **Auditoría:** Los cambios críticos se registran en tablas dedicadas (`historial_precios`, `factura_audit`).
4. **JWT_SECRET:** Si cambias esta variable, todos los tokens anteriores invalidarán.

---

## Documentación Adicional
- [API.md](./api.md) - Referencia completa de endpoints.
- [DATABASE.md](./DATABASE.md) - Esquema y diseño de tablas.
