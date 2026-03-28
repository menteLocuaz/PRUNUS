# API - Referencia Detallada de Endpoints

Base URL: `http://localhost:9090/api/v1`

> **Nota:** Todos los endpoints que modifican datos requieren `Content-Type: application/json`.
> La mayoría de los endpoints requieren autenticación mediante un token JWT en el header `Authorization: Bearer <token>`.

---

## 1. Autenticación y Sistema

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

### GET /estatus/catalogo
Obtiene el catálogo maestro de estados agrupado por módulo. **Ideal para carga inicial del frontend.**
- **Respuesta Exitosa (200 OK):**
```json
{
  "status": "success",
  "data": {
    "1": {
      "modulo": "Empresa",
      "items": [
        { "id": "uuid", "descripcion": "Activa", "tipo": "GENERAL" }
      ]
    },
    "4": {
      "modulo": "Producto",
      "items": [
        { "id": "uuid", "descripcion": "Disponible", "tipo": "STOCK" }
      ]
    }
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
| POST | `/administrar` | Gestión integral (NFC, PIN, Multi-Sucursal) |
| GET | `/{id}` | Obtener usuario por ID |

**POST /usuarios/administrar**
> Diseñado para entornos de alta rotación. Permite configurar accesos rápidos y movilidad de personal.
- **Body:**
```json
{
  "email": "cajero01@super.com",
  "usu_nombre": "Juan Pérez",
  "usu_dni": "12345678",
  "usu_pin_pos": "1234",
  "id_sucursal": "uuid_base",
  "id_rol": "uuid_rol",
  "id_status": "uuid_status",
  "sucursales_acceso": ["uuid_1", "uuid_2"]
}
```
### Roles (`/roles`)
| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/` | Listar todos los roles |
| POST | `/` | Crear un nuevo rol |
---

## 3. Administración Contable (Periodos)
Gestiona los ciclos contables. **Es obligatorio tener un periodo abierto para operar.**

| Método | Ruta | Descripción |
|--------|------|-------------|
| POST | `/periodos/abrir` | Iniciar un nuevo periodo contable |
| POST | `/periodos/cerrar/{id}` | Finalizar el periodo |
| GET | `/periodos/activo` | Obtener el periodo actual |

---

## 4. Inventario y Stock

### Inventario (`/inventario`)
Controla el stock físico y precios por sucursal.

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/` | Listar registros de inventario |
| POST | `/` | Crear registro de inventario inicial |
| PUT | `/{id}` | Actualizar stock o precios |
| POST | `/movimientos` | Registrar un movimiento (VENTA, COMPRA, AJUSTE, DEVOLUCION) |
| GET | `/movimientos/{prod_id}` | Historial (Kardex) de un producto |

**POST /inventario**
```json
{
  "id_producto": "uuid",
  "id_sucursal": "uuid",
  "stock_actual": 100,
  "precio_compra": 50.0,
  "precio_venta": 85.0
}
```
> **Nota:** El sistema ahora registra automáticamente cada cambio de precio en la tabla de auditoría `historial_precios`.

---

## 5. Ventas y Pedidos

### Órdenes (`/ordenes`)
| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/` | Listar todas las órdenes |
| POST | `/` | Crear una nueva orden de pedido |
| GET | `/{id}` | Detalle de la orden |
| PUT | `/{id}/status` | Cambiar estatus de la orden |

**POST /ordenes**
```json
{
  "odp_observacion": "Sin cebolla",
  "id_user_pos": "uuid",
  "id_periodo": "uuid",
  "id_estacion": "uuid",
  "id_status": "uuid",
  "canal": "DELIVERY",
  "odp_total": 45.50
}
```

### Agregadores (`/agregadores`)
Integración con UberEats, Rappi, etc.

| Método | Ruta | Descripción |
|--------|------|-------------|
| POST | `/` | Registrar nueva plataforma |
| POST | `/orden` | Vincular orden con ID externo y comisión |

**POST /agregadores/orden**
```json
{
  "id_orden_pedido": "uuid",
  "id_agregador": "uuid",
  "codigo_externo": "UBER-12345",
  "comision_agregador": 5.20,
  "datos_agregador": {
    "repartidor": "Carlos P.",
    "tiempo_estimado": "20 min"
  }
}
```

### Facturación (`/facturas`)
| Método | Ruta | Descripción |
|--------|------|-------------|
| POST | `/completa` | Registro atómico de factura, detalles y pagos |
| GET | `/{id}` | Detalle de factura |

---

## 6. Punto de Venta (POS)

### POS (`/pos`)
| Método | Ruta | Descripción |
|--------|------|-------------|
| POST | `/abrir` | Apertura de caja (Fondo Base) |
| POST | `/desmontar` | Cierre de sesión de cajero |
| POST | `/actualizar-valores` | Declaración de arqueo (Efectivo/Tarjetas) |

---

## Casos de Uso del Frontend
1. **Carga Inicial:** Llamar a `GET /estatus/catalogo` para mapear IDs de estados a nombres legibles.
2. **Venta Delivery:**
    - Crear orden: `POST /ordenes`.
    - Vincular plataforma: `POST /agregadores/orden`.
    - Facturar: `POST /facturas/completa`.
3. **Auditoría:** Los cambios de precios y estados de facturas se registran en tablas de auditoría dedicadas para mayor seguridad.
