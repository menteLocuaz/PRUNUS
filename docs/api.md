# API - Referencia de endpoints

Base URL: `http://localhost:9090/api/v1`

> Todos los endpoints que modifican datos requieren `Content-Type: application/json`.
> La mayoría de los endpoints requieren autenticación mediante un token JWT en el header `Authorization: Bearer <token>`.

---

## Autenticación

| Método | Ruta | Descripción |
|--------|------|-------------|
| POST | `/login` | Iniciar sesión |
| GET | `/auth/me` | Obtener información del usuario actual |
| POST | `/auth/logout` | Cerrar sesión |
| POST | `/auth/refresh-token` | Refrescar el token de acceso |

**POST /login**
```json
{
  "email": "admin@prunus.com",
  "password": "password123"
}
```

---

## Administración de Usuarios y Roles

### Usuarios
| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/usuarios` | Listar todos los usuarios |
| POST | `/usuarios` | Crear un nuevo usuario |
| GET | `/usuarios/{id}` | Obtener usuario por ID |
| PUT | `/usuarios/{id}` | Actualizar usuario |
| DELETE | `/usuarios/{id}` | Eliminar usuario |

**POST /usuarios**
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
| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/roles` | Listar todos los roles |
| POST | `/roles` | Crear un nuevo rol |
| GET | `/roles/{id}` | Obtener rol por ID |
| PUT | `/roles/{id}` | Actualizar rol |
| DELETE | `/roles/{id}` | Eliminar rol |

---

## Estructura Organizacional

### Empresas
| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/empresas` | Listar todas las empresas |
| POST | `/empresas` | Crear una nueva empresa |
| GET | `/empresas/{id}` | Obtener empresa por ID |
| PUT | `/empresas/{id}` | Actualizar empresa |
| DELETE | `/empresas/{id}` | Eliminar empresa |

### Sucursales
| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/sucursales` | Listar todas las sucursales |
| POST | `/sucursales` | Crear una nueva sucursal |
| GET | `/sucursales/{id}` | Obtener sucursal por ID |
| PUT | `/sucursales/{id}` | Actualizar sucursal |
| DELETE | `/sucursales/{id}` | Eliminar sucursal |

---

## Catálogos y Productos

### Categorías
| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/categorias` | Listar todas |
| POST | `/categorias` | Crear |
| GET | `/categorias/{id}` | Obtener por ID |
| PUT | `/categorias/{id}` | Actualizar |
| DELETE | `/categorias/{id}` | Eliminar |

### Productos
| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/productos` | Listar todos |
| POST | `/productos` | Crear |
| GET | `/productos/{id}` | Obtener por ID |
| PUT | `/productos/{id}` | Actualizar |
| DELETE | `/productos/{id}` | Eliminar |

**POST /productos**
```json
{
  "nombre": "Laptop HP",
  "descripcion": "Laptop HP 15 pulgadas",
  "precio_compra": 2500.00,
  "precio_venta": 3200.00,
  "stock": 10,
  "fecha_vencimiento": "2027-01-01T00:00:00Z",
  "imagen": "https://ejemplo.com/imagen.jpg",
  "id_sucursal": "uuid",
  "id_categoria": "uuid",
  "id_moneda": "uuid",
  "id_unidad": "uuid",
  "estado": 1
}
```

### Medidas (Unidades)
Ruta: `/medidas`

### Monedas
Ruta: `/monedas`

### Estatus (Estados del Sistema)
| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/estatus` | Listar todos los estados |
| GET | `/estatus/catalogo` | Obtener catálogo maestro agrupado por módulo |
| GET | `/estatus/tipo/{tipo}` | Obtener estados por tipo (PRODUCTO, FACTURA, etc.) |
| GET | `/estatus/modulo/{id}` | Obtener estados por ID de módulo |

---

## Punto de Venta (POS)

| Método | Ruta | Descripción |
|--------|------|-------------|
| POST | `/pos/abrir` | Abrir caja / Iniciar turno |
| GET | `/pos/estado/{id}` | Consultar estado actual de una caja |

**POST /pos/abrir**
```json
{
  "id_estacion": "uuid",
  "fondo_base": 150.00,
  "id_user_pos": "uuid"
}
```

---

## Clientes y Proveedores

### Clientes
Ruta: `/clientes`

### Proveedores
Ruta: `/proveedores`

---

## Notas Técnicas

- **Formatos de ID:** El sistema utiliza `UUID` para todos los identificadores únicos.
- **Soft Deletes:** Los registros no se eliminan físicamente; se marcan como eliminados en la base de datos.
- **Fechas:** Todas las fechas se manejan en formato ISO 8601 UTC (`YYYY-MM-DDTHH:MM:SSZ`).
- **Estados:** Utilice el endpoint `/estatus/catalogo` para conocer los IDs de estado válidos para cada módulo.
