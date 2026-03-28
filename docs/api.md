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

## 3. Administración Contable (Periodos)
Gestiona los ciclos contables necesarios para la operación de las cajas. **Es obligatorio tener un periodo abierto para poder abrir una caja.**

| Método | Ruta | Descripción |
|--------|------|-------------|
| POST | `/abrir` | Iniciar un nuevo periodo contable |
| POST | `/cerrar/{id}` | Finalizar el periodo (Valida que no haya cajas abiertas) |
| GET | `/activo` | Obtener el periodo actual |

**POST /periodos/abrir**
> No requiere body. El sistema registra automáticamente la fecha y el usuario desde el token.

**POST /periodos/cerrar/{id}**
> Si existen estaciones con sesiones abiertas vinculadas a este periodo, el sistema devolverá un error `400` impidiendo el cierre hasta que todas las cajas finalicen su arqueo.

---

## 4. Inventario y Stock

### Inventario (`/inventario`)
Controla el stock físico por sucursal, alertas de existencias y valuación.

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/` | Listar todos los registros de inventario |
| POST | `/` | Crear registro de inventario inicial |
| PUT | `/{id}` | Actualizar stock o precios |
| POST | `/movimientos` | Registrar un movimiento (ENTRADA, SALIDA, AJUSTE, DEVOLUCION, TRASLADO) |
| POST | `/movimientos/masivo` | Registrar movimientos masivos de inventario |
| GET | `/movimientos/{id}` | Listar historial de movimientos (Kardex) de un producto |
| GET | `/alertas` | Listar productos con stock bajo (Query: `id_sucursal`) |
| GET | `/valuacion` | Obtener valor total del inventario (Query: `id_sucursal`) |

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
> El sistema actualiza automáticamente el `stock_actual` en la tabla de inventario y el `stock` total en la tabla de productos.
```json
{
  "id_producto": "uuid",
  "id_sucursal": "uuid",
  "tipo_movimiento": "ENTRADA",
  "cantidad": 50,
  "referencia": "Compra Proveedor #999"
}
```

**POST /inventario/movimientos/masivo**
> Permite registrar movimientos para múltiples productos en una sola petición.
```json
{
  "id_sucursal": "uuid",
  "tipo_movimiento": "ENTRADA",
  "referencia": "Ajuste de inventario anual",
  "items": [
    {
      "id_producto": "uuid",
      "cantidad": 50
    },
    {
      "id_producto": "uuid",
      "cantidad": 30
    }
  ]
}
```

**GET /inventario/valuacion?id_sucursal=uuid**
```json
{
  "status": "success",
  "message": "Valuación de inventario calculada correctamente",
  "data": {
    "id_sucursal": "uuid",
    "total_valor": 15450.75
  }
}
```

---

## 4. Punto de Venta (POS) y Caja

### POS (`/pos`)
| Método | Ruta | Descripción |
|--------|------|-------------|
| POST | `/abrir` | Abrir una caja (Control Estación) |
| GET | `/estado/{id}` | Consultar estado de una estación |

### Estaciones POS (`/estaciones-pos`)
Gestión de estaciones físicas de venta.

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/` | Listar todas las estaciones |
| POST | `/` | Registrar nueva estación |
| GET | `/{id}` | Obtener detalle de estación |
| PUT | `/{id}` | Actualizar estación |
| DELETE | `/{id}` | Eliminar estación (Soft Delete) |

**POST /estaciones-pos**
```json
{
  "codigo": "CAJA-01",
  "nombre": "Caja Principal Planta Baja",
  "ip": "192.168.1.10",
  "id_sucursal": "uuid",
  "id_status": "uuid"
}
```

### Dispositivos POS (`/dispositivos-pos`)
Gestión de periféricos (Impresoras, Kioskos, Datáfonos, Scanners).

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/` | Listar todos los dispositivos |
| POST | `/` | Registrar nuevo dispositivo |
| GET | `/{id}` | Obtener detalle de dispositivo |
| PUT | `/{id}` | Actualizar dispositivo |
| DELETE | `/{id}` | Eliminar dispositivo (Soft Delete) |

**POST /dispositivos-pos**
```json
{
  "nombre": "Escáner de mano Almacén",
  "tipo": "SCANNER", // Opciones: IMPRESORA, DATAFONO, KIOSKO, MONITOR, SCANNER
  "ip": "192.168.1.55",
  "id_estacion": "uuid"
}
```
 ### Configuración de Impresión (`/configuracion-pos`)
Endpoints auxiliares para obtener catálogos de hardware y canales de comunicación necesarios para la configuración de
estaciones de impresión.
    
    | Método | Ruta | Descripción |
    |--------|------|-------------|
    | GET | `/canales/{chainId}` | Listar canales de impresión activos por ID de cadena |
    | GET | `/impresoras/{restId}` | Listar impresoras configuradas por ID de restaurante |
    | GET | `/puertos` | Listar todos los puertos de comunicación (COM/LPT) activos |
   
   **GET /configuracion-pos/canales/1**
    - **Respuesta Exitosa (200 OK):**
  {
    "success": true,
    "message": "Canales obtenidos",
    "data": [
      { "id": 1, "descripcion": "Canal Local" },
      { "id": 2, "descripcion": "Canal Remoto" }
    ]
  }

   
    **GET /configuracion-pos/puertos**
    - **Respuesta Exitosa (200 OK):**
  {
    "success": true,
    "message": "Puertos obtenidos",
    "data": [
      { "id": 1, "descripcion": "COM1" },
      { "id": 2, "descripcion": "LPT1" }
    ]
  }

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
| POST | `/` | Generar una nueva factura básica |
| POST | `/completa` | Registrar factura con detalle y pagos |
| GET | `/{id}` | Obtener detalle de una factura específica |
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

**POST /facturas/completa**
> Registra una factura con toda su información: cabecera, detalles (items) y formas de pago en una sola transacción atómica.
```json
{
  "cabecera": {
    "fac_numero": "FAC-002",
    "subtotal": 100.0,
    "iva": 19.0,
    "total": 119.0,
    "observacion": "Venta de mostrador",
    "id_estacion": "uuid",
    "id_cliente": "uuid",
    "id_periodo": "uuid",
    "id_control_estacion": "uuid"
  },
  "detalles": [
    {
      "id_producto": "uuid",
      "cantidad": 1,
      "precio": 100.0,
      "subtotal": 100.0,
      "impuesto": 19.0,
      "total": 119.0
    }
  ],
  "pagos": [
    {
      "id_forma_pago": "uuid",
      "valor_billete": 120.0,
      "total_pagar": 119.0
    }
  ]
}
```

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

## 7. Proveedores y Compras

### Proveedores (`/proveedores`)
Gestión de entidades que suministran productos al negocio.

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/` | Listar todos los proveedores |
| POST | `/` | Registrar un nuevo proveedor |
| GET | `/{id}` | Obtener detalle de un proveedor |
| PUT | `/{id}` | Actualizar información del proveedor |
| DELETE | `/{id}` | Eliminación lógica del proveedor |

**POST /proveedores**
```json
{
  "nombre": "Distribuidora Global S.A.",
  "ruc": "1234567890123",
  "telefono": "+507 6655-4433",
  "direccion": "Parque Industrial, Galera 4",
  "email": "ventas@distglobal.com",
  "id_sucursal": "uuid",
  "id_empresa": "uuid",
  "id_status": "uuid"
}
```

---

### Compras (`/compras`)
Módulo encargado del abastecimiento de inventario mediante órdenes de compra.

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/` | Listar historial de órdenes de compra |
| POST | `/` | Crear una nueva Orden de Compra (OC) |
| POST | `/recepcion` | Recibir mercancía y abastecer inventario |

**POST /compras** (Crear Orden)
> Registra la intención de compra. El stock NO se ve afectado en este paso.
```json
{
  "numero_orden": "OC-2024-001",
  "id_proveedor": "uuid",
  "id_sucursal": "uuid",
  "id_moneda": "uuid",
  "id_status": "uuid", // Estado inicial (ej: SOLICITADO)
  "observaciones": "Pedido mensual de bebidas",
  "detalles": [
    {
      "id_producto": "uuid",
      "cantidad_pedida": 100,
      "precio_unitario": 0.50,
      "impuesto": 7.0
    }
  ]
}
```

**POST /compras/recepcion** (Entrada de Mercancía)
> Al ejecutar este endpoint, el sistema genera automáticamente movimientos de `ENTRADA` en el inventario para cada producto recibido.
```json
{
  "id_orden_compra": "uuid",
  "id_status": "uuid", // Nuevo estado (ej: RECIBIDO)
  "items": [
    {
      "id_detalle_compra": "uuid", // ID de la línea de la OC
      "id_producto": "uuid",
      "cantidad_recibida": 100
    }
  ]
}
```

---

## Casos Límite y Errores
1. **Duplicados en Inventario:** No se permite crear dos registros de inventario para el mismo producto en la misma sucursal.
2. **Caja con Sesión Activa:** Intentar abrir una caja que ya tiene una sesión `ABIERTA` devolverá un error `400`.
3. **Stock Negativo:** Los movimientos que resulten en stock negativo generarán una advertencia en los logs del sistema, aunque la operación se complete si no hay validación estricta habilitada.
...
4. **Token Expirado:** Cualquier petición después de 24h (por defecto) devolverá `401 Unauthorized`.
5. **Recepciones Parciales:** El sistema permite recibir cantidades menores a las pedidas en una Orden de Compra. El stock se actualizará únicamente por el valor de `cantidad_recibida`.
6. **Inexistencia en Inventario:** Si intenta recibir un producto en una sucursal donde no ha sido inicializado el registro de inventario, el sistema devolverá un error `400`. Primero debe usar `POST /inventario` para ese producto/sucursal.

---

## Notas Técnicas
- **UUIDs:** Todos los campos `id_*` deben ser UUID v4 válidos.
- **Fechas:** Formato ISO-8601 (`YYYY-MM-DDTHH:MM:SSZ`).
- **Soft Delete:** Los registros eliminados mediante `DELETE` no se borran físicamente, solo se ocultan de las consultas estándar.
