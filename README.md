# Prunus

API REST desarrollada en Go para la gestión empresarial. Permite administrar empresas, sucursales, usuarios, roles, productos, proveedores, clientes, categorías, monedas y unidades de medida.

---

## Tecnologías

| Herramienta | Versión | Uso |
|---|---|---|
| Go | 1.25.4 | Lenguaje principal |
| PostgreSQL | 15 | Base de datos relacional |
| Chi Router | v5.2.3 | Enrutamiento HTTP |
| pgx | v5.8.0 | Driver PostgreSQL |
| golang-jwt | v5.3.0 | Autenticación JWT |
| godotenv | v1.5.1 | Variables de entorno |
| bcrypt | golang.org/x/crypto | Hash de contraseñas |

---

## Estructura del Proyecto

```
prunus/
├── cmd/
│   └── main.go                  # Punto de entrada, inyección de dependencias
├── pkg/
│   ├── config/
│   │   └── database/
│   │       ├── database.go      # Conexión a PostgreSQL
│   │       └── migrations/      # Migraciones automáticas de tablas
│   ├── dto/                     # Data Transfer Objects (request/response)
│   ├── helper/                  # Generación y validación de JWT
│   ├── middleware/               # Logger y autenticación JWT
│   ├── models/                  # Modelos de dominio
│   ├── routers/                 # Definición de rutas por recurso
│   ├── services/                # Lógica de negocio
│   ├── store/                   # Capa de acceso a datos (SQL)
│   ├── transport/
│   │   └── http/                # Handlers HTTP
│   └── utils/                   # Utilidades generales
├── docker-compose.yml
├── go.mod
└── go.sum
```

---

## Configuración

### 1. Clonar el repositorio

```bash
git clone https://github.com/tu-usuario/prunus.git
cd prunus
```

### 2. Variables de entorno

Copia el archivo de ejemplo y edítalo con tus valores:

```bash
cp .env.example .env
```

```env
DB_HOST=localhost
DB_USER=admin
DB_PASSWORD=tu_contraseña
DB_NAME=maxpoint
DB_PORT=5432
DB_SSLMODE=disable
JWT_SECRET=tu_clave_secreta_jwt
```

### 3. Instalar dependencias

```bash
go mod download
```

---

## Ejecución

### Con Docker Compose (PostgreSQL)

Levanta solo la base de datos:

```bash
docker-compose up -d
```

### Servidor local

```bash
go run cmd/main.go
```

El servidor corre por defecto en `http://localhost:8080`.

---

## Autenticación

Todos los endpoints **excepto** `POST /api/v1/login` requieren un JWT válido en el header:

```
Authorization: Bearer <token>
```

### Flujo de autenticación

```
1. POST /api/v1/login        → obtener token
2. Usar token en header      → acceder a recursos protegidos
3. POST /api/v1/refresh-token → renovar token antes de que expire
4. POST /api/v1/logout       → cerrar sesión
```

---

## Endpoints

### Autenticación

#### `POST /api/v1/login`

Autentica un usuario y retorna un JWT.

**No requiere token.**

**Body:**
```json
{
  "email": "admin@empresa.com",
  "password": "MiClave123"
}
```

**Respuesta exitosa `200`:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "usuario": {
    "id_usuario": 1,
    "email": "admin@empresa.com",
    "usu_nombre": "Juan Pérez",
    "id_sucursal": 1,
    "id_rol": 1
  },
  "expires_at": 1780000000
}
```

**Casos límite:**
| Caso | Código | Respuesta |
|---|---|---|
| Email o password vacíos | `400` | `Formato de petición inválido` |
| Credenciales incorrectas | `401` | `Credenciales inválidas` |
| Usuario con estado `0` (inactivo) | `401` | `Usuario inactivo` |

---

#### `GET /api/v1/me`

Retorna los datos del usuario autenticado según el JWT.

**Respuesta `200`:**
```json
{
  "id_usuario": 1,
  "email": "admin@empresa.com",
  "id_sucursal": 1,
  "id_rol": 1
}
```

---

#### `POST /api/v1/logout`

Confirma el cierre de sesión. El cliente debe eliminar el token localmente.

**Respuesta `200`:**
```json
{
  "message": "Sesión cerrada exitosamente"
}
```

---

#### `POST /api/v1/refresh-token`

Renueva el token JWT activo.

**Respuesta `200`:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_at": 1780003600
}
```

---

### Empresas `/api/v1/empresas`

#### `GET /api/v1/empresas`

Retorna todas las empresas activas (sin `deleted_at`).

```bash
curl -X GET http://localhost:8080/api/v1/empresas \
  -H "Authorization: Bearer <token>"
```

**Respuesta `200`:**
```json
[
  {
    "id": 1,
    "nombre": "Corporación Andina S.A.",
    "rut": "20.123.456-7",
    "estado": 1
  }
]
```

---

#### `POST /api/v1/empresas`

Crea una nueva empresa.

**Body:**
```json
{
  "nombre": "Corporación Andina S.A.",
  "rut": "20.123.456-7",
  "estado": 1
}
```

**Casos límite:**
| Caso | Código |
|---|---|
| `nombre` vacío | `400` |
| `rut` duplicado | `400` |
| `estado` distinto de `0` o `1` | `400` |

---

#### `GET /api/v1/empresas/{id}`

Retorna una empresa por su ID.

```bash
curl -X GET http://localhost:8080/api/v1/empresas/1 \
  -H "Authorization: Bearer <token>"
```

**Casos límite:**
| Caso | Código |
|---|---|
| ID no numérico (ej: `/empresas/abc`) | `400` |
| ID inexistente | `404` |
| Empresa eliminada (soft delete) | `404` |

---

#### `PUT /api/v1/empresas/{id}`

Actualiza todos los campos de una empresa.

**Body:**
```json
{
  "nombre": "Corporación Andina S.A. Actualizada",
  "rut": "20.123.456-7",
  "estado": 1
}
```

---

#### `DELETE /api/v1/empresas/{id}`

Realiza un **soft delete**: actualiza `deleted_at`, el registro no se elimina físicamente.

```bash
curl -X DELETE http://localhost:8080/api/v1/empresas/1 \
  -H "Authorization: Bearer <token>"
```

---

### Sucursales `/api/v1/sucursal`

Cada sucursal pertenece a una empresa (`id_empresa`).

#### `POST /api/v1/sucursal`

```json
{
  "id_empresa": 1,
  "nombre_sucursal": "Sede Central Lima",
  "estado": 1
}
```

**Casos límite:**
| Caso | Código |
|---|---|
| `id_empresa` inexistente | `400` |
| `nombre_sucursal` vacío | `400` |

---

### Roles `/api/v1/rol`

Los roles están vinculados a una sucursal.

#### `POST /api/v1/rol`

```json
{
  "nombre_rol": "Administrador",
  "id_sucursal": 1,
  "estado": 1
}
```

---

### Usuarios `/api/v1/usuario`

#### `POST /api/v1/usuario`

La contraseña se almacena con hash bcrypt. El usuario requiere un rol y una sucursal válidos.

```json
{
  "id_sucursal": 1,
  "id_rol": 1,
  "email": "juan.perez@empresa.com",
  "usu_nombre": "Juan Pérez",
  "usu_dni": "12345678",
  "usu_telefono": "+51 987654321",
  "password": "MiClaveSegura123",
  "estado": 1
}
```

**Respuesta `201`:**
```json
{
  "id_usuario": 3,
  "id_sucursal": 1,
  "id_rol": 1,
  "email": "juan.perez@empresa.com",
  "usu_nombre": "Juan Pérez",
  "usu_dni": "12345678",
  "usu_telefono": "+51 987654321",
  "estado": 1
}
```

**Casos límite:**
| Caso | Código |
|---|---|
| `email` con formato inválido | `400` |
| `email` ya registrado | `400` |
| `password` vacío | `400` |
| `id_rol` o `id_sucursal` inexistentes | `400` |
| `estado` = `0` impide el login | `401` en `/login` |

---

### Categorías `/api/v1/categoria`

Clasifican productos dentro de una sucursal.

#### `POST /api/v1/categoria`

```json
{
  "nombre": "Electrónica",
  "id_sucursal": 1
}
```

---

### Monedas `/api/v1/moneda`

Define las monedas disponibles por sucursal.

#### `POST /api/v1/moneda`

```json
{
  "nombre": "Sol Peruano",
  "id_sucursal": 1,
  "estado": 1
}
```

---

### Unidades de Medida `/api/v1/medida`

#### `POST /api/v1/medida`

```json
{
  "nombre": "Kilogramo",
  "id_sucursal": 1
}
```

---

### Productos `/api/v1/producto`

#### `POST /api/v1/producto`

Requiere que `id_sucursal`, `id_categoria`, `id_moneda` e `id_unidad` existan previamente.

```json
{
  "nombre": "Laptop HP 15",
  "descripcion": "Procesador Intel i5, 8GB RAM, 512GB SSD",
  "precio_compra": 1800.00,
  "precio_venta": 2500.00,
  "stock": 10,
  "fecha_vencimiento": "2027-12-31T00:00:00Z",
  "imagen": "https://cdn.empresa.com/productos/laptop-hp.jpg",
  "estado": 1,
  "id_sucursal": 1,
  "id_categoria": 2,
  "id_moneda": 1,
  "id_unidad": 3
}
```

**Casos límite:**
| Caso | Código |
|---|---|
| `precio_venta` o `precio_compra` negativos | `400` |
| `stock` negativo | `400` |
| `fecha_vencimiento` con formato incorrecto | `400` |
| Alguna FK inexistente | `400` |

---

### Proveedores `/api/v1/proveedor`

#### `POST /api/v1/proveedor`

Vinculado a una empresa y una sucursal.

```json
{
  "nombre": "Distribuidora Tech S.A.C.",
  "ruc": "20456789012",
  "telefono": "+51 012345678",
  "direccion": "Av. Industrial 456, Lima",
  "email": "ventas@disttech.com",
  "estado": 1,
  "id_sucursal": 1,
  "id_empresa": 1
}
```

**Casos límite:**
| Caso | Código |
|---|---|
| `email` con formato inválido | `400` |
| `ruc` duplicado | `400` |
| `id_empresa` o `id_sucursal` inexistentes | `400` |

---

### Clientes `/api/v1/cliente`

#### `POST /api/v1/cliente`

```json
{
  "empresa_cliente": "Retail Norte S.A.",
  "nombre": "Carlos Mendoza",
  "ruc": "20987654321",
  "direccion": "Jr. Comercio 789, Trujillo",
  "telefono": "+51 987123456",
  "email": "carlos@retailnorte.com",
  "estado": 1
}
```

---

## Resumen de Endpoints

| Recurso | Ruta base | Auth requerida |
|---|---|---|
| Login | `POST /api/v1/login` | No |
| Perfil | `GET /api/v1/me` | Sí |
| Logout | `POST /api/v1/logout` | Sí |
| Refresh token | `POST /api/v1/refresh-token` | Sí |
| Empresas | `/api/v1/empresas` | Sí |
| Sucursales | `/api/v1/sucursal` | Sí |
| Roles | `/api/v1/rol` | Sí |
| Usuarios | `/api/v1/usuario` | Sí |
| Categorías | `/api/v1/categoria` | Sí |
| Monedas | `/api/v1/moneda` | Sí |
| Medidas | `/api/v1/medida` | Sí |
| Productos | `/api/v1/producto` | Sí |
| Proveedores | `/api/v1/proveedor` | Sí |
| Clientes | `/api/v1/cliente` | Sí |

Todos los recursos protegidos soportan: `GET /`, `POST /`, `GET /{id}`, `PUT /{id}`, `DELETE /{id}`.

---

## Características Técnicas

- **Soft Delete** — Los `DELETE` actualizan `deleted_at` en lugar de eliminar el registro físicamente. Las consultas `GET` filtran automáticamente los registros eliminados.
- **Migraciones automáticas** — Al iniciar la aplicación se ejecutan las migraciones de la base de datos definidas en `pkg/config/database/migrations/`.
- **JWT stateless** — Los tokens incluyen `id_usuario`, `email`, `id_rol`, `rol_nombre` e `id_sucursal` en sus claims.
- **Hash de contraseñas** — Se usa `bcrypt` para almacenar contraseñas. Nunca se guarda texto plano.
- **Logger HTTP** — Middleware configurable que registra método, ruta, estado y duración de cada petición.
- **Clean Architecture** — Separación en capas: `store` (datos) → `services` (lógica) → `transport` (HTTP).

---

## Contribuir

Las contribuciones son bienvenidas. Por favor, abre un issue o un pull request para sugerencias o mejoras.

## Licencia

Este proyecto es de código abierto y está disponible bajo la licencia MIT.
