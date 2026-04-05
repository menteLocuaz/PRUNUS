# Documentación del Sistema de Inventario - Prunus

Esta guía detalla cómo utilizar las nuevas funcionalidades del sistema de inventario, incluyendo la gestión de productos por código de barras, trazabilidad por lotes, valuación avanzada y análisis de rotación.

## 1. Gestión de Productos (Catálogo)

### Crear Producto con Inventario Inicial
**Endpoint:** `POST /api/productos`
**Descripción:** Registra un producto en el catálogo maestro y crea automáticamente su registro de inventario inicial en la sucursal especificada.

**JSON Body:**
```json
{
  "nombre": "Aceite de Oliva Extra Virgen 1L",
  "descripcion": "Aceite de prensado en frío",
  "codigo_barras": "7501234567890",
  "sku": "ACE-OLV-001",
  "precio_compra": 85.50,
  "precio_venta": 120.00,
  "stock": 50,
  "id_sucursal": "uuid-sucursal",
  "id_categoria": "uuid-categoria",
  "id_moneda": "uuid-moneda",
  "id_unidad": "uuid-unidad",
  "fecha_vencimiento": "2027-12-31T00:00:00Z"
}
```

### Buscar Producto por Código de Barras o SKU
**Endpoint:** `GET /api/productos/buscar/{codigo}`
**Descripción:** Busca un producto utilizando su código de barras de 13 dígitos o su SKU interno. Ideal para integración con escáneres.

### Obtener Inventario por Sucursal
**Endpoint:** `GET /api/v1/inventario/sucursal/{id_sucursal}`
**Descripción:** Obtiene todos los registros de inventario (productos y sus stocks) asociados a una sucursal específica. Soporta paginación.

---

## 2. Trazabilidad por Lotes

El sistema ahora genera automáticamente un **Lote** cada vez que se recibe mercancía a través de una Orden de Compra. Esto permite rastrear costos específicos y fechas de vencimiento por cada entrada.

### Recepción de Mercancía (Generación de Lotes)
**Endpoint:** `POST /api/compras/recepcion`
**Descripción:** Registra la entrada física de productos solicitados en una Orden de Compra (OC). Este proceso es crítico para la trazabilidad, ya que vincula la mercancía recibida con un costo específico y una fecha de entrada.

**Comportamiento del Sistema:**
1.  **Actualización de OC**: Incrementa la `cantidad_recibida` en los detalles de la orden de compra.
2.  **Aumento de Stock**: Suma la cantidad recibida al `stock_actual` en la tabla de inventario de la sucursal.
3.  **Movimiento Histórico**: Registra un movimiento de tipo `ENTRADA` con la referencia de la OC.
4.  **Generación de Lote**: Crea un nuevo registro en la tabla `lotes`. El código del lote se genera automáticamente como `[NumeroOrden]-[FragmentoIDProducto]`.

**JSON Body:**
```json
{
  "id_orden_compra": "uuid-de-la-orden",
  "id_status": "uuid-estatus-recibido",
  "items": [
    {
      "id_detalle_compra": "uuid-detalle-oc-item1",
      "id_producto": "uuid-producto-1",
      "cantidad_recibida": 25.0
    },
    {
      "id_detalle_compra": "uuid-detalle-oc-item2",
      "id_producto": "uuid-producto-2",
      "cantidad_recibida": 10.0
    }
  ]
}
```

**Notas importantes:**
*   Se permite recepción parcial (recibir menos de lo pedido).
*   El costo capturado en el Lote será el `precio_unitario` definido en la Orden de Compra original.
*   Si el producto tiene fecha de vencimiento global, esta se heredará al lote (o puede ser extendida en futuras versiones).

---

## 3. Valuación de Inventario

Permite conocer el valor contable de las existencias en una sucursal utilizando diferentes metodologías.

**Endpoint:** `GET /api/inventario/valuacion?id_sucursal={uuid}&metodo={peps|ueps|promedio}`

### Parámetros:
*   `id_sucursal`: Identificador de la sucursal a valuar.
*   `metodo`:
    *   `promedio`: (Por defecto) Calcula `stock_actual * precio_compra_promedio`.
    *   `peps`: (FIFO) Valúa basándose en el costo de los lotes más antiguos que aún tienen stock.
    *   `ueps`: (LIFO) Valúa basándose en el costo de los lotes más recientes.

**Respuesta Exitosa:**
```json
{
  "status": "success",
  "message": "Valuación de inventario calculada correctamente",
  "data": {
    "id_sucursal": "uuid-sucursal",
    "metodo": "peps",
    "total_valor": 15420.50
  }
}
```

---

## 4. Análisis de Rotación ABC

Clasifica los productos de una sucursal según su importancia económica (Principio de Pareto).

**Endpoint:** `GET /api/inventario/rotacion?id_sucursal={uuid}`

### Categorías de Salida:
*   **Clase A**: Productos que representan el ~80% del valor total (Alta inversión, requieren control estricto).
*   **Clase B**: Productos que representan el siguiente ~15% del valor.
*   **Clase C**: Productos de baja inversión (~5% del valor), representan el mayor volumen de ítems pero menor impacto financiero.

**Respuesta:**
```json
{
  "status": "success",
  "data": {
    "A": ["uuid-prod-1", "uuid-prod-2"],
    "B": ["uuid-prod-3", "uuid-prod-4", "uuid-prod-5"],
    "C": ["uuid-prod-6", "...", "uuid-prod-100"]
  }
}
```

---

## 5. Movimientos de Inventario

### Registrar Movimiento Manual (Ajuste/Dañado)
**Endpoint:** `POST /api/inventario/movimientos`
**JSON Body:**
```json
{
  "id_producto": "uuid-producto",
  "id_sucursal": "uuid-sucursal",
  "tipo_movimiento": "AJUSTE_NEGATIVO",
  "cantidad": 5,
  "referencia": "Producto dañado en exhibición"
}
```

---

## Casos Límite y Validaciones

1.  **Stock Negativo**: El sistema permite registros negativos en ajustes, pero las ventas validan disponibilidad según la configuración de la estación.
2.  **Lotes Agotados**: Los métodos PEPS/UEPS ignoran automáticamente lotes con `cantidad_actual = 0`.
3.  **Códigos Duplicados**: El sistema validará que el `codigo_barras` sea único por empresa para evitar colisiones en el escaneo.
4.  **Fechas de Vencimiento**: Si un producto no tiene fecha de vencimiento, el campo debe enviarse como `null` o emitirse en el JSON.
