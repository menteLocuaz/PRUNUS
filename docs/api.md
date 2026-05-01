# API — Referencia de Endpoints

Base URL: `http://localhost:9090/api/v1`

**Autenticación:** `Authorization: Bearer <token>` en todos los endpoints marcados con 🔒.  
**Content-Type:** `application/json` en todas las peticiones con body.

---

## Respuesta estándar

```json
{ "status": "success", "message": "...", "data": { } }
```
Errores devuelven el código HTTP correspondiente:
```json
{ "status": "error", "message": "Descripción del error" }
```

## Paginación (cursor-based)

Los endpoints de listado aceptan:

| Param       | Tipo    | Descripción                                  |
|-------------|---------|----------------------------------------------|
| `limit`     | int     | Registros por página (default 20)            |
| `last_id`   | UUID    | ID del último ítem recibido                  |
| `last_date` | RFC3339 | Fecha del último ítem (ISO 8601 con timezone)|

---

## 1. Autenticación (`/auth`)

### POST /auth/login
```json
{ "email": "admin@prunus.com", "password": "Admin123" }
```
**200:**
```json
{
  "data": {
    "token": "jwt...",
    "usuario": { },
    "expires_at": 1234567890
  }
}
```

### GET /auth/me 🔒
Devuelve el usuario autenticado.

### POST /auth/refresh-token 🔒
No requiere body. Devuelve `{ "token": "...", "expires_at": 1234567890 }`.

### POST /auth/logout 🔒
No requiere body.

---

## 2. Usuarios (`/usuarios`) 🔒

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/` | Listar (paginado) |
| POST | `/` | Crear usuario |
| GET | `/{id}` | Obtener por ID |
| PUT | `/{id}` | Actualizar |
| DELETE | `/{id}` | Eliminar (204) |
| POST | `/administrar` | Crear con accesos multi-sucursal |
| POST | `/administrar/{id}` | Actualizar con accesos multi-sucursal |

**Body (POST / PUT):**
```json
{
  "id_sucursal": "uuid",
  "id_rol": "uuid",
  "username": "cajero01",
  "email": "cajero01@empresa.com",
  "usu_nombre": "Juan Pérez",
  "usu_dni": "12345678",
  "usu_telefono": "099999999",
  "password": "MiPassword1",
  "usu_tarjeta_nfc": "NFC_CODE",
  "usu_pin_pos": "1234",
  "nombre_ticket": "Juan",
  "id_status": "uuid",
  "sucursales_acceso": ["uuid1", "uuid2"]
}
```
`username` (4–50), `email`, `usu_nombre` (3–100), `usu_dni` (8–15) y `password` (mín 6) son requeridos en POST. En PUT `password` es opcional.

---

## 3. Roles (`/roles`) 🔒

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/` | Listar |
| POST | `/` | Crear |
| GET | `/{id}` | Obtener por ID |
| PUT | `/{id}` | Actualizar |
| DELETE | `/{id}` | Eliminar (204) |

**Body:**
```json
{ "nombre_rol": "Cajero", "id_sucursal": "uuid", "id_status": "uuid" }
```
`nombre_rol` (3–100) requerido.

---

## 4. Empresas (`/empresas`) 🔒

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/` | Listar |
| POST | `/` | Crear |
| GET | `/{id}` | Obtener por ID |
| PUT | `/{id}` | Actualizar |
| DELETE | `/{id}` | Eliminar (204) |

**Body:**
```json
{ "nombre": "Empresa S.A.", "rut": "800123456-7", "id_status": "uuid" }
```

---

## 5. Sucursales (`/sucursales`) 🔒

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/` | Listar |
| POST | `/` | Crear |
| GET | `/{id}` | Obtener por ID |
| PUT | `/{id}` | Actualizar |
| DELETE | `/{id}` | Eliminar (204) |

**Body:**
```json
{ "id_empresa": "uuid", "nombre_sucursal": "Sucursal Central", "id_status": "uuid" }
```

---

## 6. Productos (`/productos`) 🔒

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/` | Listar (paginado) |
| POST | `/` | Crear producto + inventario inicial |
| GET | `/{id}` | Obtener por ID |
| GET | `/buscar/{codigo}` | Buscar por código de barras o SKU |
| PUT | `/{id}` | Actualizar |
| DELETE | `/{id}` | Eliminar (204) |

**Body POST:**
```json
{
  "nombre": "Aceite Oliva 1L",
  "descripcion": "Prensado en frío",
  "codigo_barras": "7501234567890",
  "sku": "ACE-OLV-001",
  "precio_compra": 85.50,
  "precio_venta": 120.00,
  "stock": 50,
  "fecha_vencimiento": "2027-12-31T00:00:00Z",
  "imagen": "url_o_base64",
  "id_status": "uuid",
  "id_sucursal": "uuid",
  "id_categoria": "uuid",
  "id_moneda": "uuid",
  "id_unidad": "uuid"
}
```
`nombre`, `precio_compra`, `precio_venta`, `stock`, `id_sucursal`, `id_categoria`, `id_moneda`, `id_unidad` requeridos. Crear producto crea automáticamente su registro en inventario.

**Body PUT:** igual que POST pero sin `stock` ni `id_sucursal`.

---

## 7. Categorías (`/categorias`) 🔒

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/` | Listar |
| POST | `/` | Crear |
| GET | `/{id}` | Obtener por ID |
| PUT | `/{id}` | Actualizar |
| DELETE | `/{id}` | Eliminar (204) |

**Body:** `{ "nombre": "Lácteos", "id_sucursal": "uuid" }`

---

## 8. Clientes (`/clientes`) 🔒

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/` | Listar |
| POST | `/` | Crear |
| GET | `/{id}` | Obtener por ID |
| PUT | `/{id}` | Actualizar |
| DELETE | `/{id}` | Eliminar (204) |

**Body:**
```json
{
  "empresa_cliente": "Corporación XYZ",
  "nombre": "María López",
  "ruc": "1234567890001",
  "direccion": "Av. Principal 123",
  "telefono": "0991234567",
  "email": "maria@xyz.com",
  "id_status": "uuid"
}
```
Todos los campos son requeridos.

---

## 9. Estatus (`/estatus`) 🔒

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/` | Listar todos |
| GET | `/catalogo` | Catálogo completo agrupado por módulo — **usar en carga inicial del frontend** |
| GET | `/tipo/{tipo}` | Filtrar por tipo (ej. `STOCK`, `GENERAL`) |
| GET | `/modulo/{moduloID}` | Filtrar por ID de módulo |
| POST | `/` | Crear |
| GET | `/{id}` | Obtener por ID |
| PUT | `/{id}` | Actualizar |
| DELETE | `/{id}` | Eliminar (204) |

**GET /estatus/catalogo — respuesta:**
```json
{
  "data": {
    "1": { "modulo": "Empresa", "items": [{ "id": "uuid", "descripcion": "Activa", "tipo": "GENERAL" }] },
    "4": { "modulo": "Producto", "items": [{ "id": "uuid", "descripcion": "Disponible", "tipo": "STOCK" }] }
  }
}
```

---

## 10. Inventario (`/inventario`) 🔒

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/` | Listar (paginado) |
| POST | `/` | Crear registro manual de inventario |
| GET | `/{id}` | Obtener por ID |
| GET | `/sucursal/{id}` | Inventario de una sucursal (paginado) |
| PUT | `/{id}` | Actualizar stock / precios |
| DELETE | `/{id}` | Eliminar (204) |
| POST | `/movimientos` | Registrar movimiento individual |
| POST | `/movimientos/masivo` | Registrar movimientos para múltiples productos |
| GET | `/movimientos/{id_producto}` | Historial (Kardex) de un producto (paginado) |
| GET | `/alertas?id_sucursal={uuid}` | Productos con stock ≤ stock_mínimo |
| GET | `/alertas/detalle?id_sucursal={uuid}` | Alertas con detalle completo del producto |
| GET | `/valuacion?id_sucursal={uuid}&metodo={peps\|ueps\|promedio}` | Valor contable del inventario |
| GET | `/rotacion?id_sucursal={uuid}` | Análisis ABC de rotación |
| GET | `/rotacion/detalle?id_sucursal={uuid}` | Detalle de rotación por producto |
| GET | `/composicion-categoria?id_sucursal={uuid}` | Composición del inventario por categoría |
| POST | `/historico/snapshot` | Capturar snapshot del valor del inventario |
| GET | `/historico?id_sucursal={uuid}` | Histórico de valor del inventario |
| GET | `/perdidas?id_sucursal={uuid}` | Análisis de pérdidas / mermas |
| GET | `/margen?id_sucursal={uuid}` | Análisis de margen de ganancia |

> Si `id_sucursal` se omite en alertas/valuacion/rotacion/composicion/historico/perdidas/margen, se toma del token JWT.

**Body POST /inventario:**
```json
{
  "id_producto": "uuid", "id_sucursal": "uuid",
  "stock_actual": 50.0, "stock_minimo": 5.0, "stock_maximo": 200.0,
  "precio_compra": 85.50, "precio_venta": 120.00
}
```

**Body POST /inventario/movimientos:**
```json
{
  "id_producto": "uuid",
  "id_sucursal": "uuid",
  "tipo_movimiento": "AJUSTE",
  "cantidad": 5.0,
  "referencia": "Producto dañado"
}
```
`tipo_movimiento` válidos: `ENTRADA` | `SALIDA` | `AJUSTE` | `DEVOLUCION` | `TRASLADO`

**Body POST /inventario/movimientos/masivo:**
```json
{
  "id_sucursal": "uuid",
  "tipo_movimiento": "ENTRADA",
  "referencia": "Recepción OC-001",
  "items": [
    { "id_producto": "uuid", "cantidad": 25.0 },
    { "id_producto": "uuid", "cantidad": 10.0 }
  ]
}
```

---

## 11. Proveedores (`/proveedores`) 🔒

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/` | Listar |
| POST | `/` | Crear |
| GET | `/{id}` | Obtener por ID |
| PUT | `/{id}` | Actualizar |
| DELETE | `/{id}` | Eliminar (204) |

**Body:**
```json
{
  "nombre": "Distribuidora Norte",
  "ruc": "1234567890001",
  "telefono": "0991234567",
  "direccion": "Calle 1 y Av. 2",
  "email": "contacto@norte.com",
  "id_status": "uuid",
  "id_sucursal": "uuid",
  "id_empresa": "uuid"
}
```

---

## 12. Compras (`/compras`) 🔒

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/` | Listar órdenes de compra |
| POST | `/` | Crear orden de compra |
| GET | `/{id}` | Detalle de orden de compra |
| POST | `/recepcion` | Recibir mercancía (genera Lotes y actualiza stock) |

**Body POST /compras:**
```json
{
  "numero_orden": "OC-2024-001",
  "id_proveedor": "uuid",
  "id_sucursal": "uuid",
  "id_moneda": "uuid",
  "id_status": "uuid",
  "observaciones": "Urgente",
  "detalles": [
    { "id_producto": "uuid", "cantidad_pedida": 50.0, "precio_unitario": 85.50, "impuesto": 0.0 }
  ]
}
```

**Body POST /compras/recepcion:**
```json
{
  "id_orden_compra": "uuid",
  "id_status": "uuid",
  "items": [
    { "id_detalle_compra": "uuid", "id_producto": "uuid", "cantidad_recibida": 25.0 }
  ]
}
```
La recepción actualiza stock en inventario y genera un `Lote` por producto con el código `[NumeroOrden]-[FragmentoIDProducto]`.

---

## 13. Períodos (`/periodos`) 🔒

**Obligatorio tener un período abierto para operar facturas y POS.**

| Método | Ruta | Descripción |
|--------|------|-------------|
| POST | `/abrir` | Abrir nuevo período contable |
| POST | `/cerrar/{id}` | Cerrar período |
| GET | `/activo` | Obtener período activo actual |

No requieren body.

---

## 14. POS (`/pos`) 🔒

| Método | Ruta | Descripción |
|--------|------|-------------|
| POST | `/abrir` | Apertura de sesión de cajero con fondo base |
| POST | `/desmontar` | Cierre de sesión de cajero |
| POST | `/actualizar-valores` | Declaración de arqueo (efectivo / tarjetas) |
| GET | `/estado/{id}` | Estado actual de una estación |

**Body POST /pos/abrir:**
```json
{ "id_estacion": "uuid", "fondo_base": 100.00, "id_user_pos": "uuid" }
```

**Body POST /pos/desmontar:**
```json
{ "id_control_estacion": "uuid", "id_restaurante": "uuid", "motivo_descuadre": "..." }
```

**Body POST /pos/actualizar-valores:**
```json
{ "id_control_estacion": "uuid", "id_forma_pago": "uuid", "valor": 500.00 }
```

---

## 15. Cajas (`/cajas`) 🔒

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/` | Listar cajas |
| POST | `/` | Crear caja |
| GET | `/{id}` | Obtener por ID |
| POST | `/sesion/abrir` | Abrir sesión de caja con fondo base |
| POST | `/sesion/cerrar/{id}` | Cerrar sesión de caja (arqueo y cierre) |

**Body POST /cajas/sesion/abrir:**
```json
{ "id_caja": "uuid", "fondo_base": 100.00 }
```

**Body POST /cajas/sesion/cerrar/{id}:**
```json
{
  "efectivo_declarado": 650.00,
  "motivo_descuadre": "Diferencia en cambio"
}
```

---

## 16. Facturas (`/facturas`) 🔒

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/` | Listar (paginado) |
| POST | `/` | Crear factura simple |
| GET | `/{id}` | Detalle de factura con ítems |
| POST | `/completa` | Registro atómico: factura + detalles + pagos |
| GET | `/impuestos` | Catálogo de impuestos disponibles |
| GET | `/impuestos/{id}` | Obtener impuesto por ID |
| POST | `/impuestos` | Crear impuesto |
| PUT | `/impuestos/{id}` | Actualizar impuesto |
| DELETE | `/impuestos/{id}` | Eliminar impuesto (204) |
| GET | `/formas-pago` | Catálogo de formas de pago |
| POST | `/formas-pago` | Crear forma de pago |
| PUT | `/formas-pago/{id}` | Actualizar forma de pago |
| DELETE | `/formas-pago/{id}` | Eliminar forma de pago (204) |

**Body POST /facturas/completa:**
```json
{
  "cabecera": {
    "fac_numero": "FAC-0001",
    "subtotal": 100.00,
    "iva": 15.00,
    "total": 115.00,
    "observacion": "...",
    "id_estacion": "uuid",
    "id_orden_pedido": "uuid",
    "id_cliente": "uuid",
    "id_periodo": "uuid",
    "id_control_estacion": "uuid",
    "base_impuesto": 100.00,
    "impuesto": 0.15,
    "valor_impuesto": 15.00,
    "metadata": { }
  },
  "detalles": [
    { "id_producto": "uuid", "cantidad": 2.0, "precio": 50.00, "subtotal": 100.00, "impuesto": 15.00, "total": 115.00 }
  ],
  "pagos": [
    { "id_forma_pago": "uuid", "valor_billete": 120.00, "total_pagar": 115.00 }
  ]
}
```

---

## 17. Órdenes de Pedido (`/ordenes`) 🔒

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/` | Listar |
| POST | `/` | Crear orden |
| GET | `/{id}` | Detalle |
| PUT | `/{id}/status` | Cambiar estado |

**Body PUT /ordenes/{id}/status:**
```json
{ "id_status": "uuid" }
```

---

## 18. Agregadores (`/agregadores`) 🔒

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/` | Listar plataformas |
| POST | `/` | Registrar plataforma |
| GET | `/{id}` | Obtener por ID |
| PUT | `/{id}` | Actualizar |
| DELETE | `/{id}` | Eliminar (204) |
| POST | `/orden` | Vincular orden con plataforma externa |

**Body POST /agregadores/orden:**
```json
{
  "id_orden_pedido": "uuid",
  "id_agregador": "uuid",
  "codigo_externo": "UBER-12345",
  "comision_agregador": 5.20,
  "datos_agregador": { "repartidor": "Carlos P.", "tiempo_estimado": "20 min" }
}
```

---

## 19. Medidas / Unidades (`/medidas`) 🔒

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/` | Listar |
| POST | `/` | Crear |
| GET | `/{id}` | Obtener por ID |
| PUT | `/{id}` | Actualizar |
| DELETE | `/{id}` | Eliminar (204) |

**Body:** `{ "nombre": "Kilogramo", "id_sucursal": "uuid" }`

---

## 20. Monedas (`/monedas`) 🔒

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/` | Listar |
| POST | `/` | Crear |
| GET | `/{id}` | Obtener por ID |
| PUT | `/{id}` | Actualizar |
| DELETE | `/{id}` | Eliminar (204) |

**Body:** `{ "nombre": "Dólar", "id_sucursal": "uuid", "id_status": "uuid" }`

---

## 21. Dispositivos POS (`/dispositivos-pos`) 🔒

CRUD estándar: `GET /`, `POST /`, `GET /{id}`, `PUT /{id}`, `DELETE /{id}`.

---

## 22. Estaciones POS (`/estaciones-pos`) 🔒

CRUD estándar: `GET /`, `POST /`, `GET /{id}`, `PUT /{id}`, `DELETE /{id}`.

---

## 23. Configuración POS (`/configuracion-pos`) 🔒

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/canales/{chainId}` | Canales disponibles para una cadena |
| GET | `/impresoras/{restId}` | Impresoras configuradas |
| GET | `/puertos` | Puertos disponibles |

---

## 24. Dashboard (`/dashboard`) 🔒

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/resumen` | Resumen general: ventas, stock, alertas |
| GET | `/antiguedad-deuda` | Reporte de antigüedad de deuda |
| GET | `/composicion-categoria` | Composición del inventario por categoría |
| GET | `/mermas` | Reporte de mermas y pérdidas |

> Todos los endpoints aceptan `?id_sucursal={uuid}`; si se omite, se toma del token JWT.

---

## Flujos típicos del frontend

### Carga inicial
1. `POST /auth/login` → guardar token
2. `GET /estatus/catalogo` → mapear IDs de estado a etiquetas legibles
3. `GET /periodos/activo` → validar que hay un período contable abierto

### Venta con agregador (delivery)
1. `POST /ordenes` → crear orden
2. `POST /agregadores/orden` → vincular con plataforma externa
3. `POST /facturas/completa` → facturar

### Recepción de mercancía
1. `POST /compras` → crear orden de compra
2. `POST /compras/recepcion` → recibir ítems → actualiza stock e inventario automáticamente

### Ajuste de inventario
- `POST /inventario/movimientos` con `tipo_movimiento: "AJUSTE"` para correcciones individuales
- `POST /inventario/movimientos/masivo` para múltiples productos a la vez

### Apertura de caja (cajero)
1. `POST /pos/abrir` → abrir sesión POS con fondo base
2. `POST /cajas/sesion/abrir` → registrar apertura de caja física

### Cierre / arqueo de caja
1. `POST /pos/actualizar-valores` → declarar efectivo y tarjetas
2. `POST /cajas/sesion/cerrar/{id}` → cerrar sesión con arqueo
3. `POST /pos/desmontar` → desmontar cajero de la estación

### Consultas analíticas (dashboard / reportes)
- `GET /dashboard/resumen` → KPIs del negocio
- `GET /inventario/valuacion?metodo=promedio` → valor contable
- `GET /inventario/rotacion` → análisis ABC
- `GET /inventario/margen` → márgenes por producto
- `GET /inventario/perdidas` → mermas registradas
