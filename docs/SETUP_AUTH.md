# Cómo Configurar el Sistema de Autenticación

Esta guía describe los pasos para inicializar el sistema de seguridad de Prunus POS: ejecutar migraciones, cargar datos iniciales y validar los métodos de acceso disponibles.

---

## Prerrequisitos

Asegúrate de que tu archivo `.env` contenga las variables requeridas antes de ejecutar cualquier comando:

```env
# Base de datos
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=tu_password
DB_NAME=prunus

# JWT
JWT_SECRET=tu_clave_secreta_de_minimo_32_caracteres
JWT_EXPIRATION_HOURS=24
```

> `JWT_SECRET` debe tener **mínimo 32 caracteres**. El servidor no arrancará si no se cumple este requisito.

---

## Opción A — Inicialización Automática (Recomendado)

Ejecuta los tres comandos en orden. Cada uno depende del anterior.

### 1. Ejecutar migraciones

Crea el esquema completo de la base de datos (migraciones 000001–000038).

```bash
go run ./cmd/ migrate
```

Esto incluye todas las tablas del sistema: entidades maestras, POS, inventario, caja, compras, auditoría y configuración de impresoras.

### 2. Cargar datos seed

Inserta los datos base necesarios para operar: estatus del sistema, empresa demo, sucursal central, 8 módulos de navegación y el rol **Administrador Global**.

```bash
go run ./cmd/ seed
```

Tras este paso, los siguientes registros estarán disponibles:

| Entidad            | Nombre                  | UUID                                   |
|--------------------|-------------------------|----------------------------------------|
| Empresa            | Prunus Business Demo    | `11111111-1111-4111-a111-111111111111` |
| Sucursal           | Sucursal Central        | `22222222-2222-4222-a222-222222222222` |
| Rol                | Administrador Global    | `7d7b0e11-1234-4a21-9591-316279f06742` |
| Estatus (Activo)   | Activo                  | `7f7b0e11-1234-4a21-9591-316279f06742` |

### 3. Crear el usuario administrador

```bash
go run ./cmd/ register usuario \
  --sucursal "22222222-2222-4222-a222-222222222222" \
  --rol      "7d7b0e11-1234-4a21-9591-316279f06742" \
  --username "admin" \
  --email    "admin@prunus.com" \
  --nombre   "Admin Master" \
  --dni      "12345678" \
  --password "Admin123" \
  --status   "7f7b0e11-1234-4a21-9591-316279f06742"
```

Una salida exitosa muestra:

```
✅ Usuario registrado con ID: <uuid-generado>
```

---

## Opción B — Inicialización Manual vía SQL

Si prefieres usar un cliente como DBeaver o `psql`, ejecuta este script en orden.

```sql
-- 1. DATOS BASE
INSERT INTO empresa (id_empresa, nombre, rut, id_status)
VALUES ('11111111-1111-4111-a111-111111111111', 'Prunus Business Demo', '12345678-9', '7f7b0e11-1234-4a21-9591-316279f06742')
ON CONFLICT (id_empresa) DO NOTHING;

INSERT INTO sucursal (id_sucursal, id_empresa, nombre_sucursal, id_status)
VALUES ('22222222-2222-4222-a222-222222222222', '11111111-1111-4111-a111-111111111111', 'Sucursal Central', '7f7b0e11-1234-4a21-9591-316279f06742')
ON CONFLICT (id_sucursal) DO NOTHING;

INSERT INTO rol (id_rol, nombre_rol, id_sucursal, id_status)
VALUES ('7d7b0e11-1234-4a21-9591-316279f06742', 'Administrador Global', '22222222-2222-4222-a222-222222222222', '7f7b0e11-1234-4a21-9591-316279f06742')
ON CONFLICT (id_rol) DO NOTHING;

-- 2. MÓDULOS DE NAVEGACIÓN
-- Nota: ON CONFLICT requiere la cláusula WHERE porque mdl_id tiene índice parcial.
INSERT INTO modulo (mdl_id, mdl_descripcion, abreviatura, ruta, icono, id_status, is_active, orden)
VALUES
    (1, 'Configuración Empresa', 'EMP',  '/config/empresa',    'settings',      '7f7b0e11-1234-4a21-9591-316279f06742', true, 1),
    (2, 'Gestión Sucursales',    'SUC',  '/config/sucursales', 'store',         '7f7b0e11-1234-4a21-9591-316279f06742', true, 2),
    (3, 'Seguridad y Usuarios',  'USR',  '/config/usuarios',   'users',         '7f7b0e11-1234-4a21-9591-316279f06742', true, 3),
    (4, 'Catálogo Productos',    'PROD', '/productos',         'package',       '7f7b0e11-1234-4a21-9591-316279f06742', true, 4),
    (5, 'Ventas y POS',          'VENT', '/ventas',            'shopping-cart', '7f7b0e11-1234-4a21-9591-316279f06742', true, 5),
    (6, 'Control de Inventario', 'INV',  '/inventario',        'archive',       '7f7b0e11-1234-4a21-9591-316279f06742', true, 6),
    (7, 'Reportes y Dashboards', 'DASH', '/dashboard',         'pie-chart',     '7f7b0e11-1234-4a21-9591-316279f06742', true, 7),
    (8, 'Gastos Operativos',     'GAST', '/gastos',            'dollar-sign',   '7f7b0e11-1234-4a21-9591-316279f06742', true, 8)
ON CONFLICT (mdl_id) WHERE deleted_at IS NULL DO UPDATE SET
    mdl_descripcion = EXCLUDED.mdl_descripcion,
    ruta            = EXCLUDED.ruta,
    icono           = EXCLUDED.icono,
    orden           = EXCLUDED.orden;

-- 3. PERMISOS TOTALES AL ADMINISTRADOR GLOBAL
INSERT INTO permiso_rol (id_rol, id_modulo, can_read, can_write, can_update, can_delete)
SELECT '7d7b0e11-1234-4a21-9591-316279f06742', id_modulo, true, true, true, true
FROM modulo
WHERE deleted_at IS NULL
ON CONFLICT (id_rol, id_modulo) DO NOTHING;

-- 4. USUARIO ADMINISTRADOR
-- Se recomienda usar el CLI (Opción A) en lugar de SQL para que el hash del password
-- sea generado correctamente por el sistema (bcrypt costo 12).
-- Si usas SQL directamente, genera el hash con:
--   go run ./cmd/ register usuario --username admin --email admin@prunus.com ...
INSERT INTO usuario (id_sucursal, id_rol, username, email, usu_nombre, usu_dni, password, id_status)
VALUES (
    '22222222-2222-4222-a222-222222222222',
    '7d7b0e11-1234-4a21-9591-316279f06742',
    'admin',
    'admin@prunus.com',
    'Admin Master',
    '12345678',
    '$2a$10$U.sUS/qwAXlDPrJZ9wAaLe78DmRtcnWVY39wFp85YLiL0iIVPVkkK',
    '7f7b0e11-1234-4a21-9591-316279f06742'
) ON CONFLICT DO NOTHING;
```

> **Importante:** No uses un hash escrito a mano. El proyecto usa bcrypt con costo 12. Usa siempre el CLI (`Opción A`) para que el password sea hasheado correctamente.

---

## Validar el Login

El endpoint único de autenticación es:

```
POST /api/v1/login
Content-Type: application/json
```

Prunus admite tres métodos de acceso:

### Por Email

```json
{ "email": "admin@prunus.com", "password": "Admin123" }
```

### Por Username

```json
{ "username": "admin", "password": "Admin123" }
```

### Por PIN (acceso rápido POS)

```json
{ "pin": "1234" }
```

Una respuesta exitosa devuelve un JWT con la siguiente estructura relevante:

```json
{
  "token": "<jwt>",
  "usuario": { "id_usuario": "...", "email": "admin@prunus.com" },
  "permisos": ["/config/empresa", "/config/sucursales", "/ventas", "..."]
}
```

---

## Validar Permisos en el Frontend

El campo `permisos` contiene las **rutas habilitadas** para el rol del usuario. Úsalo para renderizar dinámicamente el sidebar o guardar rutas.

```ts
// Ejemplo en TypeScript
const rutasHabilitadas: string[] = response.permisos;
const sidebar = rutasHabilitadas.map(ruta => menuConfig[ruta]).filter(Boolean);
```

Los contextos JWT inyectados por `RequireAuth()` en cada request son:

| Clave           | Tipo        | Descripción                    |
|-----------------|-------------|--------------------------------|
| `user_id`       | `uuid.UUID` | ID del usuario autenticado     |
| `user_email`    | `string`    | Email del usuario              |
| `user_rol`      | `string`    | Nombre del rol                 |
| `user_sucursal` | `uuid.UUID` | Sucursal asignada al usuario   |

---

## Referencia rápida — Tablas del sistema

A partir de la migración **000038**, el esquema incluye las siguientes tablas operacionales:

### Módulo de Caja

| Tabla            | Descripción                                      |
|------------------|--------------------------------------------------|
| `caja`           | Punto de venta físico (terminal de caja)         |
| `sesion_caja`    | Turno de un cajero: apertura, cierre y montos    |
| `movimiento_caja`| Ingresos y egresos de efectivo dentro de un turno|

### Auditoría y Logs

| Tabla          | Descripción                                        |
|----------------|----------------------------------------------------|
| `log_sistema`  | Registro de acciones: usuario, tabla, IP y fecha   |

### Compras

| Tabla                  | Descripción                                    |
|------------------------|------------------------------------------------|
| `detalle_orden_compra` | Líneas de producto de una orden de compra      |

Columnas agregadas a `orden_compra`: `numero_orden`, `id_moneda`, `subtotal`, `impuesto`, `observaciones`, `fecha_vencimiento`.

### Configuración POS / Impresión

| Tabla             | Descripción                                      |
|-------------------|--------------------------------------------------|
| `canal_impresion` | Canales de impresión configurados por cadena     |
| `impresora`       | Impresoras disponibles por sucursal              |
| `puertos`         | Puertos COM/USB disponibles para dispositivos    |

---

## Solución de problemas comunes

### Error: `dirty state` en migraciones

Si una migración falló a mitad y la BD quedó en estado `dirty`, fuérzala a la versión anterior:

```bash
# Reemplaza N con el número de la última migración exitosa
go run ./cmd/ migrate force N
go run ./cmd/ migrate
```

### Error: `no unique or exclusion constraint matching ON CONFLICT`

Ocurre al ejecutar SQL manual con `ON CONFLICT (mdl_id)` en la tabla `modulo`. El índice de `mdl_id` es parcial. Asegúrate de usar:

```sql
ON CONFLICT (mdl_id) WHERE deleted_at IS NULL DO UPDATE SET ...
```

### Error: `JWT_SECRET must be at least 32 characters`

Aumenta la longitud del valor de `JWT_SECRET` en el archivo `.env`.
