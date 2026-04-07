# Implementación del Dashboard de Inventario y Finanzas

Este documento sirve como referencia técnica y estratégica para la implementación del tablero de control (dashboard) de Prunus. El objetivo central es transformar los datos transaccionales en información accionable para responder tres preguntas críticas: **¿Qué comprar?**, **¿Qué liquidar?** y **¿Cuánto dinero hay en riesgo?**

---

## 1. Filosofía del Dashboard: Orientación a la Acción

A diferencia de un reporte estático, este dashboard está diseñado bajo el principio de **Gestión por Excepción**. Las gráficas deben resaltar anomalías y riesgos financieros de forma inmediata mediante el uso de colores semafóricos (Rojo para riesgo, Amarillo para advertencia, Verde para salud).

---

## 2. Gestión de Inventario (Operativo)

Ayuda a evitar quiebres de stock y excesos de mercancía que inmovilizan capital.

### Nivel de Stock por Almacén
*   **Objetivo:** Visualizar la distribución física de la mercancía.
*   **Tipo Recomendado:** Gráfico de barras apiladas.
*   **Fuente de Datos:** `GET /api/v1/inventario` (agrupado por sucursal).
*   **Acción Sugerida:** Realizar traslados internos entre sucursales si una tiene exceso y otra tiene quiebre.

### Artículos con Stock Bajo (Alertas)
*   **Objetivo:** Identificar productos por debajo del punto de reorden.
*   **Tipo Recomendado:** Gráfico de barras horizontales (ordenado de menor a mayor).
*   **Fuente de Datos:** `GET /api/v1/dashboard/resumen` -> `productos_bajo_stock`.
*   **Acción Sugerida:** Generar órdenes de compra inmediatas.

### Composición del Inventario por Categoría
*   **Objetivo:** Identificar la concentración de capital por tipo de producto.
*   **Tipo Recomendado:** Gráfico de dona (Donut chart).
*   **Fuente de Datos:** `GET /api/v1/dashboard/composicion-categoria`.
*   **Lógica:** Muestra el `%` del valor total del inventario que representa cada categoría.

---

## 3. Estado Financiero del Inventario

Traduce unidades físicas en valor monetario real.

### Valor del Inventario Actual
*   **Objetivo:** Conocer el costo total de la mercancía en estantería.
*   **Tipo Recomendado:** KPI Card (Número grande) + Barras.
*   **Fuente de Datos:** `GET /api/v1/dashboard/resumen` -> `valor_inventario_total`.

### Pérdidas por Merma o Caducidad
*   **Objetivo:** Visualizar el dinero perdido por daños o productos vencidos.
*   **Tipo Recomendado:** Gráfico de cascada (Waterfall) o barras rojas.
*   **Fuente de Datos:** `GET /api/v1/dashboard/mermas`.
*   **Lógica:** Suma de movimientos tipo `MERMA` y `CADUCADO`.

### Valor del Inventario en el Tiempo
*   **Objetivo:** Observar si el capital inmovilizado está creciendo peligrosamente.
*   **Tipo Recomendado:** Gráfico de áreas o líneas.
*   **Fuente de Datos:** Requiere consulta a la tabla `inventario_historico` ( snapshots diarios).

---

## 4. Rendimiento de Ventas y Compras

Entendimiento del flujo de caja (Cash Flow).

### Ventas vs. Compras Mensuales
*   **Objetivo:** Saber si el negocio gasta más de lo que ingresa.
*   **Tipo Recomendado:** Gráfico de líneas doble eje.
*   **Fuente de Datos:** `GET /api/v1/dashboard/resumen` -> `ventas_vs_compras`.

### Antigüedad de Cuentas por Cobrar
*   **Objetivo:** Estado de las deudas de clientes.
*   **Tipo Recomendado:** Gráfico de barras por antigüedad (0-30, 31-60, 61-90, 90+ días).
*   **Fuente de Datos:** `GET /api/v1/dashboard/antiguedad-deuda`.

### Top 10 Productos Más Rentables (Pareto)
*   **Objetivo:** Identificar productos que generan el mayor margen neto, no solo volumen.
*   **Tipo Recomendado:** Gráfico de barras tipo "Pareto".
*   **Fuente de Datos:** `GET /api/v1/dashboard/resumen` -> `top_productos`.
*   **Lógica:** `(Precio Venta - Costo Unitario) * Cantidad`.

---

## 5. Salud del Negocio (KPIs Críticos)

### Punto de Equilibrio (Break-even)
*   **Objetivo:** Visualizar el volumen de ventas necesario para cubrir costos fijos y variables.
*   **Tipo Recomendado:** Gráfico de líneas cruzadas.
*   **Fuente de Datos:** `GET /api/v1/dashboard/resumen` -> `punto_equilibrio` y `gastos_mensuales`.

### Ciclo de Conversión de Efectivo (CCC)
*   **Objetivo:** Medir la eficiencia en días para recuperar la inversión.
*   **Tipo Recomendado:** Indicador de reloj (Gauge chart).
*   **Fuente de Datos:** `GET /api/v1/dashboard/resumen` -> `ciclo_conversion_efectivo`, `dio`, `dso`, `dpo`.
*   **Interpretación:** 
    *   **DIO (Inventory):** Días que la mercancía está en almacén.
    *   **DSO (Receivables):** Días que tardamos en cobrar a clientes.
    *   **DPO (Payables):** Días que tardamos en pagar a proveedores.

---

## Apéndice: Mapeo Técnico de Endpoints

| Endpoint | DTO / Campo Clave | Uso en Frontend |
| :--- | :--- | :--- |
| `/dashboard/resumen` | `DashboardResumen` | KPIs generales, Pareto y CCC. |
| `/dashboard/antiguedad-deuda` | `[]AntiguedadDeudaDTO` | Gráfica de barras de cobranza. |
| `/dashboard/composicion-categoria` | `[]InventarioCategoriaDTO` | Gráfica de dona de categorías. |
| `/dashboard/mermas` | `[]TopProductoDTO` | Top 10 pérdidas monetarias. |

> **Nota para Desarrolladores:** Todos los montos monetarios se entregan en la moneda base del sistema (definida por `id_moneda` en la sucursal). Los cálculos de promedios (DIO, DSO, DPO) utilizan una ventana deslizante de **90 días**.
