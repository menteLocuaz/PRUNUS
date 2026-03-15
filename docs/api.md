# API - Referencia Detallada de Endpoints

Base URL: `http://localhost:9090/api/v1`

> **Nota:** Todos los endpoints que modifican datos requieren `Content-Type: application/json`.
> La mayoría de los endpoints requieren autenticación mediante un token JWT en el header `Authorization: Bearer <token>`.

---

## 1. Autenticación

### POST /auth/login
Inicia sesión en el sistema.
- **Body:**
```json
{
  "email": "admin@prunus.com",
  "password": "password123"
}
```
- **Respuesta Exitosa (200 OK):**
```json
{
  "status": "success",
  "message": "Inicio de sesión exitoso",
  "data": {
    "token": "jwt_token_here",
    "usuario": { ... },
    "expires_at": "2026-03-15T15:00:00Z"
  }
}
```
- **Casos Límite / Errores:**
  - `400 Bad Request`: Formato de JSON inválido o faltan campos obligatorios.
  - `401 Unauthorized`: Credenciales incorrectas.

### GET /auth/me
Obtiene información del usuario autenticado actual.
- **Header:** `Authorization: Bearer <token>`
- **Respuesta Exitosa (200 OK):** Retorna el objeto usuario (sin password).

### POST /auth/logout
Cierra la sesión (Informa al servidor, el cliente debe borrar el token).
- **Header:** `Authorization: Bearer <token>`

### POST /auth/refresh-token
Genera un nuevo token basado en el token actual antes de que expire.
- **Header:** `Authorization: Bearer <token>`
- **Respuesta:** Nuevo token y nueva fecha de expiración.

---

## 2. Administración de Usuarios y Roles

### Usuarios
Base URL: `/usuarios`

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/` | Listar todos los usuarios |
| POST | `/` | Crear un nuevo usuario |
| GET | `/{id}` | Obtener usuario por ID (UUID) |
| PUT | `/{id}` | Actualizar usuario |
| DELETE | `/{id}` | Eliminar usuario (Soft Delete) |

**POST /usuarios**
- **Validaciones:**
  - `usu_nombre`: min 3, max 100.
  - `email`: formato email válido.
  - `password`: min 6 caracteres.
  - `usu_dni`: min 8, max 15.
- **Body:**
```json
{
  "id_sucursal": "uuid",
  "id_rol": "uuid",
  "email": "usuario@ejemplo.com",
  "usu_nombre": "Nombre Apellido",
  "usu_dni": "12345678",
  "usu_telefono": "987654321",
  "password": "securepassword",
  "id_status": "uuid"
}
```

### Roles
Base URL: `/roles`

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/` | Listar todos los roles |
| POST | `/` | Crear un nuevo rol |
| GET | `/{id}` | Obtener rol por ID |
| PUT | `/{id}` | Actualizar rol |
| DELETE | `/{id}` | Eliminar rol |

**POST /roles**
- **Body:**
```json
{
  "nombre_rol": "Administrador",
  "id_sucursal": "uuid",
  "id_status": "uuid"
}
```

---

## 3. Estructura Organizacional

### Empresas
Base URL: `/empresas`

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/` | Listar todas las empresas |
| POST | `/` | Crear una empresa |
| GET | `/{id}` | Obtener empresa por ID |
| PUT | `/{id}` | Actualizar empresa |
| DELETE | `/{id}` | Eliminar empresa |

**POST /empresas**
```json
{
  "nombre": "Prunus Corp",
  "rut": "12345678-9",
  "id_status": "uuid"
}
```

### Sucursales
Base URL: `/sucursales`

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/` | Listar sucursales |
| POST | `/` | Crear sucursal |
| GET | `/{id}` | Obtener por ID |
| PUT | `/{id}` | Actualizar sucursal |
| DELETE | `/{id}` | Eliminar sucursal |

**POST /sucursales**
```json
{
  "id_empresa": "uuid",
  "nombre_sucursal": "Sucursal Norte",
  "id_status": "uuid"
}
```

---

## 4. Inventario y Catálogos

### Productos
Base URL: `/productos`

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/` | Listar productos |
| POST | `/` | Crear producto |
| GET | `/{id}` | Obtener por ID |
| PUT | `/{id}` | Actualizar producto |
| DELETE | `/{id}` | Eliminar producto |

**POST /productos**
- **Body:**
```json
{
  "nombre": "Laptop HP Pro",
  "descripcion": "16GB RAM, 512GB SSD",
  "precio_compra": 850.00,
  "precio_venta": 1200.00,
  "stock": 50,
  "fecha_vencimiento": "2028-12-31T23:59:59Z",
  "imagen": "https://cdn.example.com/img.jpg",
  "id_status": "uuid",
  "id_sucursal": "uuid",
  "id_categoria": "uuid",
  "id_moneda": "uuid",
  "id_unidad": "uuid"
}
```
- **Casos Límite:**
  - El sistema genera una alerta (slog) si el producto se crea con stock < 5.
  - `404 Not Found` si se intenta actualizar o eliminar un ID inexistente.

### Categorías
Base URL: `/categorias`
```json
{
  "nombre": "Electrónica",
  "id_sucursal": "uuid"
}
```

### Medidas (Unidades)
Base URL: `/medidas`
```json
{
  "nombre": "Unidad",
  "id_sucursal": "uuid"
}
```

### Monedas
Base URL: `/monedas`
```json
{
  "nombre": "Dólar Americano",
  "id_sucursal": "uuid",
  "id_status": "uuid"
}
```

---

## 5. Estatus (Estados del Sistema)
Base URL: `/estatus`

Este módulo maneja los estados para diferentes entidades (Productos, Facturas, Usuarios, etc.).

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/` | Listar todos los estados definidos |
| GET | `/catalogo` | Obtener catálogo maestro agrupado por módulo |
| GET | `/tipo/{tipo}` | Obtener estados por tipo (ej: `ACTIVO`, `INACTIVO`) |
| GET | `/modulo/{id}` | Obtener estados de un módulo específico |
| POST | `/` | Crear un nuevo estado |

**Parámetros:**
- `moduloID` (int): ID del módulo (1: Empresa, 2: Sucursal, 4: Usuario, 10: Producto, etc.)

---

## 6. Punto de Venta (POS)
Base URL: `/pos`

### POST /pos/abrir
Abre una caja para iniciar un turno.
- **Body:**
```json
{
  "id_estacion": "uuid",
  "fondo_base": 500.00,
  "id_user_pos": "uuid"
}
```

### GET /pos/estado/{id}
Consulta el resumen actual de una estación de POS específica.
- **ID:** UUID de la **estación**.

---

## 7. Clientes y Proveedores

### Clientes
Base URL: `/clientes`
- **Campos Obligatorios:** `empresa_cliente`, `nombre`, `ruc`, `direccion`, `telefono`, `email`, `id_status`.

### Proveedores
Base URL: `/proveedores`
- **Campos Obligatorios:** `nombre`, `ruc`, `id_status`, `id_sucursal`, `id_empresa`.

---

## Notas Técnicas y Errores Comunes

- **Identificadores:** Todos los `id_*` son de tipo `UUID` v4.
- **Soft Delete:** Los endpoints `DELETE` marcan el registro con `deleted_at`. Los registros marcados no aparecen en los `GET` pero siguen en la DB.
- **Formatos de Respuesta:**
  - **Éxito (200/201):** `{"status": "success", "message": "...", "data": { ... }}`
  - **Error de Validación (400):** `{"status": "error", "message": "Fallo de validación", "errors": { "campo": "mensaje" }}`
  - **No Encontrado (404):** `{"status": "error", "message": "Recurso no encontrado"}`
- **Rate Limit:** Máximo 100 peticiones por minuto por dirección IP.
