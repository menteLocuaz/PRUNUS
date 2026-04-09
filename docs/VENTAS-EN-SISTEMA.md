# Ciclo Completo de Ventas: Guía Técnica de Implementación

Esta guía describe el flujo operativo y técnico necesario para realizar ventas en el sistema Prunus. El proceso sigue un ciclo riguroso de auditoría y control de efectivo, estructurado en capas de validación.

---

## 1. Prerrequisitos de Venta
Para que el sistema permita cualquier transacción, se deben cumplir tres condiciones jerárquicas:
1.  **Periodo Activo:** Un marco de tiempo contable abierto por administración.
2.  **Estación Identificada:** El dispositivo físico (PC/Tablet) debe estar registrado.
3.  **Control de Estación Abierto:** El cajero debe haber iniciado su turno con un fondo base.

---

## 2. Gestión de Periodos (Nivel Administrativo)

El **Periodo** (`periodo`) es la entidad de mayor jerarquía. Agrupa todas las transacciones de un día o turno global.

### Apertura de Periodo
*   **Endpoint:** `POST /api/v1/periodo/abrir`
*   **Lógica de Negocio:**
    *   El sistema valida que no exista ya un periodo con estatus "ABIERTO".
    *   Se registra la fecha de apertura y el usuario administrador responsable.
    *   **Importante:** Sin un periodo activo, el servicio de POS bloqueará cualquier intento de `AbrirCaja`.

### Cierre de Periodo
*   **Endpoint:** `POST /api/v1/periodo/cerrar/{id}`
*   **Validación Crítica:** No se puede cerrar un periodo si aún existen **Controles de Estación** (cajas) abiertos. Todas las estaciones deben estar "Desmontadas" o "Cerradas" antes del cierre global.

---

## 3. Identificación y Apertura de Estación (Cajero)

### Detección de la Estación (Matching por IP)
El sistema identifica automáticamente desde qué equipo se está operando mediante la tabla `dispositivo_pos`:
*   El frontend o middleware detecta la IP del cliente.
*   Se busca en `dispositivo_pos` el registro que coincida con esa `ip`.
*   Dicho registro provee el `id_estacion` necesario para todas las operaciones subsiguientes.

### Apertura de Caja (Control de Estación)
Una vez identificada la estación, el cajero debe realizar la apertura formal:
*   **Endpoint:** `POST /api/v1/pos/abrir`
*   **Payload:** `{ "id_estacion": "uuid", "fondo_base": 50.00, "id_user_pos": "uuid" }`
*   **Resultado:** Se crea un registro en `control_estacion` con el estatus `FONDO_ASIGNADO`. Este `id_control_estacion` es obligatorio para generar facturas.

---

## 4. Registro de Pedidos (Orden de Pedido)

La **Orden de Pedido** (`orden_pedido`) representa la intención de compra y el canal de venta.

*   **Endpoint:** `POST /api/v1/orden-pedido`
*   **Atributos Clave:**
    *   **Canal:** Determina el origen (Salón, Rappi, Para llevar, Delivery).
    *   **Observación:** Notas especiales para cocina o despacho.
    *   **Total:** Monto calculado de la orden.
*   **Estado:** Una orden nace en estado `PENDIENTE` o `EN PREPARACIÓN` según la configuración del módulo.

---

## 5. Facturación y Pago (Transacción Atómica)

Para garantizar la integridad de los datos, Prunus utiliza un proceso de **Facturación Completa**. No se permite crear la factura sin sus detalles o sus pagos en pasos separados.

### Registro Integral
*   **Endpoint:** `POST /api/v1/factura/registrar-completa`
*   **Estructura del JSON (`FacturaCompletaRequest`):**
    1.  **Cabecera:** Datos generales (`id_cliente`, `id_orden_pedido`, `id_control_estacion`, totales e impuestos).
    2.  **Detalles:** Listado de productos (`id_producto`, `cantidad`, `precio`, `impuesto`).
    3.  **Pagos:** Listado de formas de pago (`id_forma_pago`, `total_pagar`). Soporta **Pagos Mixtos** (ej: parte en efectivo y parte con tarjeta).

**Validación:** El sistema verifica que la suma de los `Pagos` coincida exactamente con el `Total` de la cabecera antes de persistir la transacción.

---

## 6. Cierre de Turno y Auditoría

Al finalizar el turno, se debe conciliar el dinero físico con el sistema.

### Arqueo (Actualizar Valores Declarados)
El cajero cuenta el dinero y lo registra:
*   **Endpoint:** `POST /api/v1/pos/actualizar-valores`
*   **Proceso:** El sistema compara el `pos_calculado` (ventas registradas) vs el `valor_declarado` (lo que el cajero dice tener).
*   Cualquier diferencia genera automáticamente un registro en `auditoria_caja` con el motivo del descuadre.

### Desmontado del Cajero
*   **Endpoint:** `POST /api/v1/pos/desmontar`
*   Cambia el estatus de la estación a `DISPONIBLE` y cierra el ciclo del usuario actual, permitiendo que otro cajero pueda abrir la misma estación en el siguiente turno.

---

## Resumen de Validaciones de Seguridad
| Error | Causa Probable | Solución |
| :--- | :--- | :--- |
| **403 Forbidden** | No hay un periodo activo. | El administrador debe abrir un periodo. |
| **400 Bad Request** | La estación ya tiene una sesión activa. | Cerrar/Desmontar la sesión previa antes de abrir una nueva. |
| **422 Unprocessable** | El total de pagos no coincide con el total factura. | Validar cálculos en el frontend antes de enviar el JSON. |
