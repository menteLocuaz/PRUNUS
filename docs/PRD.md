# Documento de Requisitos del Producto (PRD) - Prunus

| Proyecto | Prunus |
| :--- | :--- |
| **Versión** | 1.1 (Actualizado) |
| **Fecha** | 09 de Marzo, 2026 |
| **Estado** | En Desarrollo Avanzado |
| **Tecnología Principal** | Go (Golang), PostgreSQL, Redis, Docker |

---

## 1. Introducción
Prunus es una API RESTful robusta y escalable diseñada para la gestión integral de operaciones comerciales. El sistema utiliza una arquitectura multi-tenencia (multi-empresa y multi-sucursal), permitiendo una administración segregada de inventarios, ventas, finanzas y personal.

El objetivo es proveer un núcleo transaccional seguro y de alto rendimiento que sirva para aplicaciones de Punto de Venta (POS), control de almacenes y administración de recursos empresariales (ERP).

## 2. Objetivos del Producto
1.  **Multi-Tenencia:** Segregación total de datos por empresa y sucursal.
2.  **Ciclo de Venta Completo:** Desde la apertura de caja y gestión de pedidos hasta la facturación y cierre de caja.
3.  **Control de Inventario Preciso:** Seguimiento de stock en tiempo real y registro histórico de movimientos.
4.  **Optimización de Rendimiento:** Uso de Redis para caché de catálogos y minimización de latencia en el POS.
5.  **Auditabilidad y Seguridad:** Registro de auditoría en caja, logs de sistema y control de acceso basado en roles (RBAC).

## 3. Perfiles de Usuario (Roles)
-   **Super Admin:** Acceso total a la configuración global y todas las empresas.
-   **Administrador de Empresa:** Gestiona múltiples sucursales de su propia organización.
-   **Administrador de Sucursal:** Controla el inventario, personal y cierres de caja de su sede.
-   **Cajero / Vendedor:** Operaciones de venta, apertura/cierre de caja y gestión de pedidos.
-   **Bodeguero:** Gestión de entradas de mercancía, proveedores y ajustes de inventario.

## 4. Requisitos Funcionales

### 4.1. Gestión Organizacional y Usuarios
-   **Multi-empresa:** Registro de múltiples razones sociales.
-   **Sucursales:** Configuración de puntos de venta físicos independientes.
-   **RBAC Dinámico:** Creación de roles personalizados con permisos específicos.
-   **Auth:** Login seguro, gestión de perfiles (`/me`) y refresco de tokens.

### 4.2. Catálogo e Inventario
-   **Productos:** Ficha técnica con precios de compra/venta, imágenes y categorías.
-   **Inventario:** Seguimiento de existencias y registro de movimientos (entradas/salidas).
-   **Atributos:** Gestión de unidades de medida (kg, unidad), monedas y categorías.

### 4.3. Operaciones de Punto de Venta (POS)
-   **Estaciones POS:** Configuración y control de estaciones físicas de venta.
-   **Apertura y Cierre de Caja:** Control de efectivo inicial, retiros y arqueo final.
-   **Gestión de Pedidos:** Registro de órdenes de pedido (incluyendo integraciones con agregadores externos).
-   **Facturación:** Emisión de facturas con detalle de ítems, impuestos y múltiples formas de pago.

### 4.4. Finanzas y Auditoría
-   **Impuestos:** Configuración de tasas impositivas por producto/región.
-   **Formas de Pago:** Soporte para efectivo, tarjetas, transferencias, etc.
-   **Anulaciones:** Registro de motivos de anulación para facturas y pedidos.
-   **Logs y Auditoría:** Registro detallado de eventos críticos del sistema y movimientos de caja.

## 5. Requisitos No Funcionales
-   **Arquitectura:** Clean Architecture (Capas: Models, Store, Service, Transport).
-   **Escalabilidad:** Contenerización con Docker y optimización con Redis.
-   **Disponibilidad:** Base de datos PostgreSQL con migraciones automáticas.
-   **Seguridad:** Encriptación de datos sensibles y validación estricta de FKs.

## 6. Modelo de Datos (Entidades Clave)

### Núcleo Organizacional
-   `empresa`, `sucursal`, `usuario`, `rol`.

### Catálogo
-   `producto`, `categoria`, `medida`, `moneda`, `impuesto`.

### Transaccional (Ventas/POS)
-   `estaciones_pos`, `control_estacion` (Apertura/Cierre).
-   `factura`, `detalle_factura`, `forma_pago_factura`.
-   `orden_pedido`, `agregadores` (UberEats, Rappi, etc.).
-   `caja`, `retiros`, `auditoria_caja`.

### Logística
-   `inventario`, `movimientos_inventario`.
-   `proveedores`, `clientes`.

---

## 7. Roadmap y Próximos Pasos
-   [x] **Implementación de Redis:** Caché para roles y categorías.
-   [ ] **Módulo de Reportes:** Generación de reportes de ventas diarias y márgenes de utilidad.
-   [ ] **Alertas de Stock:** Notificaciones automáticas cuando un producto alcance su stock mínimo.
-   [ ] **Integración de Facturación Electrónica:** Adaptación a normativas fiscales locales.
-   [ ] **Dashboard Analítico:** Visualización de KPIs en tiempo real para administradores.
