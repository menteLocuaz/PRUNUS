# API de Inventario — Referencia para Frontend

Base URL: `/api/v1/inventario`

Todos los endpoints requieren autenticación:
```
Authorization: Bearer <token>
```

---

## Paginación (cursor-based)

Los endpoints que devuelven listas aceptan estos query params:

| Parámetro  | Tipo   | Descripción                                      |
|------------|--------|--------------------------------------------------|
| `limit`    | int    | Cantidad de registros a retornar (default: 20)   |
| `last_id`  | UUID   | ID del último registro recibido (para siguiente página) |
| `last_date`| RFC3339| Fecha del último registro (ISO 8601, con timezone) |

---

## Respuesta estándar

```json
{
  "status": "success",
  "message": "Descripción del resultado",
  "data": { }
}
```

Errores devuelven el código HTTP correspondiente con:
```json
{
  "status": "error",
  "message": "Descripción del error"
}
```

---

## Endpoints

### 1. Listar todo el inventario
`GET /api/v1/inventario/`

Soporta paginación.

**Respuesta `200`:**
```json
{
  "status": "success",
  "message": "Inventario obtenido correctamente",
  "data": [
    {
      "id_inventario": "uuid",
      "id_producto": "uuid",
      "id_sucursal": "uuid",
      "stock_actual": 50.0,
      "stock_minimo": 5.0,
      "stock_maximo": 200.0,
      "precio_compra": 85.50,
      "precio_venta": 120.00,
      "ubicacion": "Pasillo A",
      "created_at": "2024-01-15T10:00:00Z",
      "updated_at": "2024-01-15T10:00:00Z"
    }
  ]
}
```

---

### 2. Obtener inventario por sucursal
`GET /api/v1/inventario/sucursal/{id_sucursal}`

Soporta paginación.

**Respuesta `200`:** misma forma que listar todo.

**Errores:**
- `400` — ID de sucursal inválido

---

### 3. Obtener registro por ID
`GET /api/v1/inventario/{id}`

**Respuesta `200`:** objeto `Inventario` individual.

**Errores:**
- `400` — ID inválido
- `404` — Inventario no encontrado

---

### 4. Crear registro de inventario
`POST /api/v1/inventario/`

Crea un registro de inventario para un producto en una sucursal. Solo puede existir un registro por combinación `(id_producto, id_sucursal)`.

**Body:**
```json
{
  "id_producto": "uuid",
  "id_sucursal": "uuid",
  "stock_actual": 50.0,
  "stock_minimo": 5.0,
  "stock_maximo": 200.0,
  "precio_compra": 85.50,
  "precio_venta": 120.00
}
```

Todos los campos son requeridos. Los valores numéricos deben ser `>= 0`.

**Respuesta `201`:** objeto `Inventario` creado.

**Errores:**
- `400` — JSON inválido, validación fallida, o ya existe un registro para ese producto en esa sucursal

---

### 5. Actualizar registro de inventario
`PUT /api/v1/inventario/{id}`

**Body:**
```json
{
  "stock_actual": 45.0,
  "stock_minimo": 5.0,
  "stock_maximo": 200.0,
  "precio_compra": 85.50,
  "precio_venta": 120.00
}
```

Todos los campos son requeridos. No se puede cambiar `id_producto` ni `id_sucursal`.

**Respuesta `200`:** objeto `Inventario` actualizado.

---

### 6. Eliminar registro de inventario
`DELETE /api/v1/inventario/{id}`

**Respuesta `204`:** sin cuerpo.

**Errores:**
- `404` — Inventario no encontrado

---

### 7. Registrar movimiento individual
`POST /api/v1/inventario/movimientos`

**Body:**
```json
{
  "id_producto": "uuid",
  "id_sucursal": "uuid",
  "tipo_movimiento": "AJUSTE",
  "cantidad": 5.0,
  "referencia": "Producto dañado en exhibición"
}
```

**`tipo_movimiento` — valores válidos:**

| Valor       | Descripción                        |
|-------------|------------------------------------|
| `ENTRADA`   | Ingreso de mercancía               |
| `SALIDA`    | Egreso de mercancía                |
| `AJUSTE`    | Corrección de inventario (±)       |
| `DEVOLUCION`| Devolución de cliente o proveedor  |
| `TRASLADO`  | Transferencia entre sucursales     |

- `cantidad` debe ser `> 0`
- `referencia` es opcional
- `id_usuario` se toma automáticamente del token JWT

**Respuesta `201`:**
```json
{
  "status": "success",
  "message": "Movimiento de inventario registrado correctamente",
  "data": {
    "id_movimiento": "uuid",
    "id_producto": "uuid",
    "id_sucursal": "uuid",
    "tipo_movimiento": "AJUSTE",
    "cantidad": 5.0,
    "costo_unitario": 85.50,
    "precio_unitario": 120.00,
    "stock_anterior": 50.0,
    "stock_posterior": 45.0,
    "fecha": "2024-01-15T10:00:00Z",
    "id_usuario": "uuid",
    "referencia": "Producto dañado en exhibición",
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:00:00Z"
  }
}
```

---

### 8. Registrar movimientos masivos
`POST /api/v1/inventario/movimientos/masivo`

Registra el mismo tipo de movimiento para múltiples productos en una sola operación.

**Body:**
```json
{
  "id_sucursal": "uuid",
  "tipo_movimiento": "ENTRADA",
  "referencia": "Recepción OC-2024-001",
  "items": [
    { "id_producto": "uuid-producto-1", "cantidad": 25.0 },
    { "id_producto": "uuid-producto-2", "cantidad": 10.0 }
  ]
}
```

- `items` requiere al menos 1 elemento
- Los mismos valores válidos de `tipo_movimiento` aplican aquí

**Respuesta `201`:** array de objetos `MovimientoInventario`.

---

### 9. Historial de movimientos de un producto
`GET /api/v1/inventario/movimientos/{id_producto}`

Retorna el historial de movimientos para un producto específico. Soporta paginación.

**Respuesta `200`:** array de objetos `MovimientoInventario`.

**Errores:**
- `400` — ID de producto inválido

---

### 10. Alertas de stock bajo
`GET /api/v1/inventario/alertas?id_sucursal={uuid}`

Retorna los productos cuyo `stock_actual <= stock_minimo`.

- Si no se envía `id_sucursal`, se usa la sucursal del token JWT.

**Respuesta `200`:** array de objetos `Inventario`.

**Errores:**
- `400` — `id_sucursal` no proporcionado ni disponible en token

---

### 11. Valuación de inventario
`GET /api/v1/inventario/valuacion?id_sucursal={uuid}&metodo={peps|ueps|promedio}`

Calcula el valor contable del inventario de una sucursal.

**Query params:**

| Parámetro    | Requerido | Default    | Descripción                         |
|--------------|-----------|------------|-------------------------------------|
| `id_sucursal`| No*       | —          | UUID de la sucursal                 |
| `metodo`     | No        | `promedio` | Método de valuación                 |

*Si no se envía, se toma del token JWT.

**Métodos disponibles:**

| Valor      | Descripción                                                         |
|------------|---------------------------------------------------------------------|
| `promedio` | `stock_actual × precio_compra` (costo promedio ponderado)           |
| `peps`     | FIFO — valúa con el costo de los lotes más antiguos con stock       |
| `ueps`     | LIFO — valúa con el costo de los lotes más recientes con stock      |

**Respuesta `200`:**
```json
{
  "status": "success",
  "message": "Valuación de inventario calculada correctamente",
  "data": {
    "id_sucursal": "uuid",
    "metodo": "peps",
    "total_valor": 15420.50
  }
}
```

---

### 12. Análisis de rotación ABC
`GET /api/v1/inventario/rotacion?id_sucursal={uuid}`

Clasifica los productos por importancia económica (Principio de Pareto).

- Si no se envía `id_sucursal`, se usa la sucursal del token JWT.

**Categorías:**

| Clase | Representa                | Acción recomendada              |
|-------|---------------------------|---------------------------------|
| A     | ~80% del valor total      | Control estricto, reorden rápido|
| B     | ~15% del valor total      | Control moderado                |
| C     | ~5% del valor total       | Control básico                  |

**Respuesta `200`:**
```json
{
  "status": "success",
  "message": "Análisis de rotación ABC obtenido correctamente",
  "data": {
    "A": ["uuid-prod-1", "uuid-prod-2"],
    "B": ["uuid-prod-3", "uuid-prod-4"],
    "C": ["uuid-prod-5", "uuid-prod-6"]
  }
}
```

---

## Notas para el frontend

- **Lotes y trazabilidad**: Se generan automáticamente al recibir mercancía vía `POST /api/v1/compras/recepcion`. Los lotes con `cantidad_actual = 0` son ignorados en valuación PEPS/UEPS.
- **Búsqueda por código de barras/SKU**: `GET /api/v1/productos/buscar/{codigo}` — devuelve el producto independientemente de la sucursal.
- **Stock negativo**: El sistema lo permite en ajustes manuales; las ventas validan disponibilidad según configuración de estación POS.
- **`deleted_at`**: Solo aparece en la respuesta si tiene valor (soft delete). Filtra registros activos verificando `deleted_at == null` o ausente.
