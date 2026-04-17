# Ciclo Completo de Ventas: Guía Técnica de Implementación

Esta guía describe el flujo operativo y técnico necesario para realizar ventas en el sistema Prunus. El proceso sigue un ciclo riguroso de auditoría y control de efectivo, estructurado en capas de validación.

---

## 0. Carga Inicial (Frontend)
Antes de iniciar cualquier proceso, el frontend debe obtener los catálogos maestros para mapear estados y configuraciones.

*   **Endpoint:** `GET /api/catalogo` 🔒
*   **Propósito:** Obtener todos los estados (`PENDIENTE`, `PAGADO`, `ACTIVO`, etc.) agrupados por módulo.
*   **Uso:** Almacenar en caché local para evitar consultas repetitivas durante el ciclo de venta.

---

## 1. Prerrequisitos de Venta
Para que el sistema permita cualquier transacción, se deben cumplir tres condiciones jerárquicas:
1.  **Periodo Activo:** Un marco de tiempo contable abierto por administración.
2.  **Estación Identificada:** El dispositivo físico (PC/Tablet) debe estar registrado.
3.  **Control de Estación Abierto:** El cajero debe haber iniciado su turno con un fondo base.

---

## 2. Gestión de Periodos (Nivel Administrativo)

El **Periodo** (`periodos`) es la entidad de mayor jerarquía. Agrupa todas las transacciones de un día o turno global.

### Apertura de Periodo
*   **Endpoint:** `POST /api/v1/periodos/abrir` 🔒
*   **Lógica de Negocio:**
    *   Valida que no exista un periodo con estatus "ABIERTO".
    *   Registra la fecha de apertura y el usuario administrador responsable.
    *   **Importante:** Sin un periodo activo, el servicio de POS bloqueará cualquier intento de `AbrirCaja`.

### Cierre de Periodo
*   **Endpoint:** `POST /api/v1/periodos/cerrar/{id}` 🔒
*   **Validación Crítica:** No se puede cerrar un periodo si existen **Controles de Estación** (cajas) abiertos. Todas las estaciones deben estar "Cerradas" o "Desmontadas".

---

## 3. Identificación y Apertura de Estación (Cajero)

### Detección de la Estación (Matching por IP)
El sistema identifica automáticamente el equipo mediante la tabla `dispositivo_pos`:
*   Se busca el registro en `dispositivo_pos` que coincida con la `ip` del cliente.
*   Este registro provee el `id_estacion` necesario para operar.

### Apertura de Caja (Control de Estación)
*   **Endpoint:** `POST /api/v1/pos/abrir` 🔒
*   **Payload:** `{ "id_estacion": "uuid", "fondo_base": 50.00, "id_user_pos": "uuid" }`
*   **Resultado:** Crea un registro en `control_estacion` con estatus `FONDO_ASIGNADO`. El `id_control_estacion` es obligatorio para facturar.

---

## 4. Registro de Pedidos (Orden de Pedido)

La **Orden de Pedido** (`ordenes`) representa la intención de compra y el canal de venta.

*   **Endpoint:** `POST /api/v1/ordenes` 🔒
*   **Atributos Clave:**
    *   **Canal:** Origen (Salón, Rappi, Para llevar, Delivery).
    *   **Estado:** Inicia generalmente como `PENDIENTE` (consultar `/api/catalogo` para el ID exacto del módulo 5).

---

## 5. Facturación y Pago (Transacción Atómica)

Prunus utiliza un proceso de **Facturación Completa** para garantizar la integridad.

### Registro Integral
*   **Endpoint:** `POST /api/v1/facturas/completa` 🔒
*   **Estructura del JSON (`FacturaCompletaRequest`):**
    1.  **Cabecera:** Datos generales, `id_cliente`, `id_orden_pedido`, `id_control_estacion`.
    2.  **Detalles:** Listado de productos (`id_producto`, `cantidad`, `precio`, `impuesto`).
    3.  **Pagos:** Listado de formas de pago (`id_forma_pago`, `total_pagar`). Soporta **Pagos Mixtos**.

**Validación:** La suma de los `Pagos` debe coincidir exactamente con el `Total` de la factura.

---

## 6. Cierre de Turno y Auditoría

Al finalizar el turno, se debe conciliar el dinero físico con el sistema.

### Arqueo (Actualizar Valores Declarados)
*   **Endpoint:** `POST /api/v1/pos/actualizar-valores` 🔒
*   **Proceso:** Compara `pos_calculado` (ventas) vs `valor_declarado` (físico). Las diferencias se registran en `auditoria_caja`.

### Desmontado del Cajero
*   **Endpoint:** `POST /api/v1/pos/desmontar` 🔒
*   Cambia el estatus de la estación a `DISPONIBLE` y cierra el ciclo del usuario, permitiendo un nuevo turno.

---

## Resumen de Validaciones de Seguridad
| Error | Causa Probable | Solución |
| :--- | :--- | :--- |
| **401 Unauthorized** | Token ausente o expirado. | Reautenticar el usuario. |
| **403 Forbidden** | No hay un periodo activo. | El administrador debe abrir un periodo. |
| **400 Bad Request** | La estación ya tiene una sesión activa. | Cerrar/Desmontar la sesión previa. |
| **422 Unprocessable** | El total de pagos no coincide con el total factura. | Validar cálculos en el frontend. |
