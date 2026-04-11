# Documento de Requisitos del Producto (PRD) - Prunus

| Proyecto | Prunus |
| :--- | :--- |
| **Versión** | 1.2 (Optimizado) |
| **Fecha** | 09 de Abril, 2026 |
| **Estado** | En Desarrollo Avanzado / Fase de Optimización |
| **Tecnología Principal** | Go (1.25+), PostgreSQL 15, Redis, Docker |

---

## 1. Introducción
Prunus es una API RESTful de grado empresarial diseñada para la gestión integral de operaciones de supermercados y retail. Utiliza una arquitectura multi-tenencia (multi-empresa y multi-sucursal) que garantiza la segregación total de datos y un alto rendimiento transaccional bajo cargas intensas de Punto de Venta (POS).

## 2. Objetivos del Producto
1.  **Alta Disponibilidad y Rendimiento:** Minimizar latencia en el POS mediante estrategias de caching agresivas y optimización de consultas SQL.
2.  **Integridad Transaccional:** Garantizar la consistencia de inventarios y finanzas mediante operaciones atómicas (ej: Facturación Completa).
3.  **Control Operativo Estricto:** Gestión de periodos contables, arqueos de caja y auditoría detallada por transacción.
4.  **Escalabilidad Multi-Sucursal:** Facilitar la administración centralizada de múltiples puntos de venta con inventarios independientes.
5.  **Optimización de Bulk Operations:** Procesamiento eficiente de movimientos masivos de inventario y asignaciones de seguridad.

## 3. Perfiles de Usuario (Roles)
-   **Super Admin:** Configuración global, gestión de empresas y monitoreo de sistema.
-   **Administrador de Empresa:** Gestión de catálogos maestros y visualización consolidada de sucursales.
-   **Administrador de Sucursal:** Control local de inventario, personal y validación de arqueos de caja.
-   **Cajero / Vendedor:** Operaciones de venta, manejo de efectivo y atención al cliente.
-   **Bodeguero:** Recepción de mercancía por lotes, gestión de proveedores y ajustes de stock.

## 4. Requisitos Funcionales

### 4.1. Seguridad y Acceso (RBAC)
-   **RBAC con Caching:** Control de acceso basado en roles con permisos cacheados en Redis para latencia mínima.
-   **Multi-Sucursal:** Usuarios con acceso habilitado a una o varias sedes específicas.
-   **Detección de Estación:** Identificación automática de estaciones POS basada en IP del dispositivo.

### 4.2. Gestión de Catálogo y Productos
-   **Maestro de Productos:** Soporte para códigos de barras, SKU, imágenes y fechas de vencimiento.
-   **Unidades y Medidas:** Conversión y gestión flexible de unidades (kg, lt, unidad).
-   **Categorización:** Estructura jerárquica para análisis de ventas y rotación.

### 4.3. Inventario y Logística
-   **Movimientos Masivos:** Registro de entradas/salidas por lote mediante optimización de base de datos (`UNNEST`).
-   **Trazabilidad por Lotes:** Control de stock por fecha de recepción y vencimiento.
-   **Valuación Avanzada:** Soporte para métodos de costo Promedio, PEPS (FIFO) y UEPS (LIFO).
-   **Alertas de Stock:** Monitoreo en tiempo real de niveles críticos y stock mínimo.

### 4.4. Ciclo de Venta POS
-   **Control de Periodos:** Marco contable obligatorio para habilitar operaciones de venta.
-   **Apertura/Cierre de Caja:** Fondo base, retiros parciales y arqueo ciego con detección de descuadres.
-   **Facturación Atómica:** Registro integral de cabecera, detalle y pagos en una única transacción de base de datos.
-   **Pagos Mixtos:** Soporte para transacciones pagadas con múltiples medios simultáneamente (Efectivo + Tarjeta).

## 5. Requisitos No Funcionales
-   **Arquitectura:** Clean Architecture con inyección de dependencias modular.
-   **Optimización de Store:** Centralización de campos SQL y helpers de escaneo (Pattern DRY) para evitar redundancias.
-   **Caching Layer:** Redis implementado para Roles, Permisos, Categorías y Monedas con TTL de 24h e invalidación reactiva.
-   **Auditoría Nativa:** Uso de triggers en PostgreSQL y `SET LOCAL app.current_user_id` para trazabilidad a nivel de fila.

## 6. Modelo de Datos (Entidades Clave)

### Gestión y Seguridad
-   `empresa`, `sucursal`, `usuario`, `rol`, `permiso_rol`.

### POS y Finanzas
-   `periodo`, `estaciones_pos`, `control_estacion`, `auditoria_caja`.
-   `factura`, `detalle_factura`, `forma_pago_factura`.
-   `orden_pedido` (Canales: Salón, Delivery, Agregadores).

### Inventario y Catálogo
-   `producto`, `inventario`, `lotes`, `movimientos_inventario`.
-   `categoria`, `unidad`, `moneda`.

---

## 7. Roadmap y Progreso Técnico
-   [x] **Optimización de Capa de Datos:** Implementación de constantes de campos y helpers de escaneo (DRY).
-   [x] **Inserciones Masivas:** Uso de `UNNEST` para mejorar rendimiento en un 80% en operaciones de lote.
-   [x] **Caching de Permisos:** Reducción de carga en BD mediante caching de RBAC en Redis.
-   [x] **Valuación de Inventario:** Implementación de algoritmos PEPS/UEPS.
-   [ ] **Reportes Analíticos Avanzados:** Análisis ABC de rotación y márgenes de utilidad reales.
-   [ ] **Dashboard Gerencial:** Visualización de KPIs (Ventas por hora, ticket promedio).
-   [ ] **Integración de Facturación Electrónica:** Módulo de comunicación con entes fiscales.
-   [ ] **Alertas Push:** Notificaciones a móviles para administradores sobre stock crítico.
