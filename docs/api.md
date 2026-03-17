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

---

## 2. Administración de Usuarios y Roles

### Usuarios (`/usuarios`)
| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/` | Listar todos los usuarios |
| POST | `/` | Crear un nuevo usuario |
| GET | `/{id}` | Obtener usuario por ID |

### Roles (`/roles`)
| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/` | Listar todos los roles |
| POST | `/` | Crear un nuevo rol |

---

## 3. Inventario y Stock

### Inventario (`/inventario`)
Controla el stock físico por sucursal.

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/` | Listar todos los registros de inventario |
| POST | `/` | Crear registro de inventario inicial |
| PUT | `/{id}` | Actualizar stock o precios |
| POST | `/movimientos` | Registrar un movimiento (VENTA, COMPRA, AJUSTE) |
| GET | `/movimientos/{id}` | Listar movimientos de un producto |

**POST /inventario**
```json
{
  "id_producto": "uuid",
  "id_sucursal": "uuid",
  "stock_actual": 100,
  "stock_minimo": 10,
  "stock_maximo": 500,
  "precio_compra": 50.0,
  "precio_venta": 85.0
}
```

**POST /inventario/movimientos**
```json
{
  "id_producto": "uuid",
  "tipo_movimiento": "VENTA",
  "cantidad": 5,
  "id_usuario": "uuid",
  "referencia": "Factura #123"
}
```

---

## 4. Punto de Venta (POS) y Caja

### POS (`/pos`)
| Método | Ruta | Descripción |
|--------|------|-------------|
| POST | `/abrir` | Abrir una caja (Control Estación) |
| GET | `/estado/{id}` | Consultar estado de una estación |
| GET | `/dispositivos` | Listar impresoras/periféricos |

### Caja (`/caja`)
| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/` | Listar cajas físicas |
| POST | `/abrir` | Iniciar una sesión de cajero |
| POST | `/cerrar/{id}`| Cerrar sesión con arqueo |

**POST /caja/abrir**
```json
{
  "id_caja": "uuid",
  "id_usuario": "uuid",
  "monto_apertura": 1000.00
}
```

---

## 5. Facturación

### Facturas (`/facturas`)
| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/` | Listar historial de facturas |
| POST | `/` | Generar una nueva factura con detalles |
| GET | `/impuestos` | Listar catálogo de impuestos |
| GET | `/formas-pago`| Listar catálogo de formas de pago |

**POST /facturas**
```json
{
  "factura": {
    "fac_numero": "FAC-001",
    "cfac_total": 119.0,
    "id_cliente": "uuid",
    "id_periodo": "uuid",
    "id_status": "uuid"
  },
  "items": [
    {
      "id_producto": "uuid",
      "cantidad": 1,
      "precio": 100.0,
      "impuesto": 19.0,
      "total": 119.0
    }
  ]
}
```

---

## 6. Pedidos y Agregadores

### Ordenes (`/ordenes`)
| Método | Ruta | Descripción |
|--------|------|-------------|
| POST | `/` | Crear orden de pedido |
| PUT | `/{id}/status` | Cambiar estado (EJ: De "Pendiente" a "Cocinando") |

### Agregadores (`/agregadores`)
Para integración con plataformas externas (UberEats, Rappi, etc).

| Método | Ruta | Descripción |
|--------|------|-------------|
| POST | `/` | Registrar nuevo agregador |
| POST | `/orden` | Vincular orden de pedido con ID externo |

---

## Casos Límite y Errores

1. **Duplicados en Inventario:** No se permite crear dos registros de inventario para el mismo producto en la misma sucursal.
2. **Caja con Sesión Activa:** Intentar abrir una caja que ya tiene una sesión `ABIERTA` devolverá un error `400`.
3. **Stock Negativo:** Los movimientos que resulten en stock negativo generarán una advertencia en los logs del sistema, aunque la operación se complete si no hay validación estricta habilitada.
4. **Token Expirado:** Cualquier petición después de 24h (por defecto) devolverá `401 Unauthorized`.

---

## Notas Técnicas
- **UUIDs:** Todos los campos `id_*` deben ser UUID v4 válidos.
- **Fechas:** Formato ISO-8601 (`YYYY-MM-DDTHH:MM:SSZ`).
- **Soft Delete:** Los registros eliminados mediante `DELETE` no se borran físicamente, solo se ocultan de las consultas estándar.
