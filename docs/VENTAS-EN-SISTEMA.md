Para comenzar a realizar ventas en el sistema, la estructura que hemos diseñado sigue un proceso riguroso de control de efectivo y
  auditoría. No se puede simplemente "facturar"; el sistema requiere que se cumpla un ciclo de apertura de "Caja" (denominado
  técnicamente como Control de Estación).


  Aquí te explico el flujo lógico basado en las tablas que creamos:


  1. El Concepto de "Caja" en el Sistema
  La "Caja" no es solo una tabla, es la combinación de tres elementos:
   * Estación POS: El equipo físico o terminal (computador, tablet).
   * Periodo: El marco de tiempo (ej: el turno de la mañana o el día de hoy).
   * Control de Estación: El registro que dice "El Usuario X abrió la Estación Y en el Periodo Z con un saldo inicial".

  ---


  2. Pasos para comenzar las ventas (Flujo Operativo)


  Paso 1: Apertura del Periodo (Tabla periodo)
  Antes de cualquier venta, un administrador debe abrir un Periodo.
   * Por qué: Para agrupar todas las ventas de un día o un turno específico y facilitar el cuadre contable al final.


  Paso 2: Apertura de la Estación (Tabla control_estacion)
  El cajero ingresa al sistema y realiza la "Apertura de Caja":
   * Se registra el fondo_base (el dinero en efectivo que hay físicamente para dar cambio).
   * El estado del control_estacion cambia a "Abierto".
   * Sin este paso, el sistema bloqueará cualquier intento de generar una factura.


  Paso 3: Registro de Pedidos (Tabla orden_pedido)
  Una vez la caja está abierta, se comienzan a tomar los pedidos:
   * Se selecciona el canal (Salón, Rappi, Para llevar).
   * Se agregan productos al pedido.


  Paso 4: Facturación y Pago (Tablas factura y forma_pago_factura)
  Cuando el cliente paga:
   * Se genera la Factura.
   * Se registra el pago en forma_pago_factura. Si el cliente paga una parte en efectivo y otra con tarjeta, el sistema lo permite
     (Pagos Mixtos) gracias a esta tabla.

  ---


  3. Durante el día (Seguridad)
  Si hay mucho efectivo en la caja, el cajero puede realizar Retiros (Tabla retiros).
   * Esto disminuye el saldo que el cajero "debe" tener al final del turno, mejorando la seguridad.


  4. Cierre de Ventas
  Al finalizar el turno:
   1. Se cierra el Control de Estación: Se cuenta el dinero físico y se compara con lo que el sistema dice que debería haber
      (pos_calculado vs fondo_retirado).
   2. Si hay diferencias, se registran en la tabla auditoria_caja.
   3. Finalmente, se cierra el Periodo global.


  Resumen Técnico
  Para implementar esto en el código (Go), las APIs de ventas siempre deberán validar:
   1. ¿Existe un Periodo activo?
   2. ¿Tiene el usuario un Control de Estación abierto para la estación desde la que está intentando vender?



## Estados
Para que el sistema de estatus sea eficaz y rápido, tanto en el backend como para el consumo del frontend, propongo un plan basado
  en un Catálogo Maestro Estructurado.


  1. Definición de Identificadores de Módulos (mdl_id)
  Primero, estandarizamos los IDs de los módulos en el código (vía constantes) para evitar "números mágicos":



  ┌────┬──────────┬─────────────────────────────────────────┐
  │ ID │ Módulo   │ Descripción                             │
  ├────┼──────────┼─────────────────────────────────────────┤
  │ 1  │ EMPRESA  │ Configuración global de la empresa      │
  │ 2  │ SUCURSAL │ Gestión de sucursales                   │
  │ 3  │ USUARIO  │ Gestión de accesos y perfiles           │
  │ 4  │ PRODUCTO │ Catálogo de productos e inventario      │
  │ 5  │ VENTA    │ Facturación y pedidos                   │
  │ 6  │ COMPRA   │ Órdenes de compra y proveedores         │
  │ 7  │ FINANZAS │ Tesorería, monedas y pagos              │
  │ 8  │ CAJA_POS │ Control de estaciones, turnos y arqueos │
  └────┴──────────┴─────────────────────────────────────────┘



  2. Estructura de Respuesta JSON "Caché-Friendly"
  En lugar de que el frontend pida estatus uno por uno, diseñaremos un endpoint de Catálogo Maestro que devuelva un objeto indexado
  por el ID del módulo. Esto permite al frontend cargar una sola vez y acceder en $O(1)$.

  Endpoint Propuesto: GET /api/v1/estatus/catalogo


    1 {
    2   "status": "success",
    3   "data": {
    4     "1": {
    5       "modulo": "Empresa",
    6       "items": [
    7         { "id": "uuid-1", "descripcion": "Activo", "tipo": "1" },
    8         { "id": "uuid-2", "descripcion": "Suspendido", "tipo": "1" }
    9       ]
   10     },
   11     "8": {
   12       "modulo": "Caja/POS",
   13       "items": [
   14         { "id": "59039503...", "descripcion": "Activo", "tipo": "1" },
   15         { "id": "99039503...", "descripcion": "Fondo Asignado", "tipo": "1" }
   16       ]
   17     }
   18   }
   19 }


  3. Implementación de "Slugs" o Códigos Rápidos
  Para que el backend sea rápido al validar lógica de negocio (ej: "solo permitir venta si el estatus es ACTIVO"), añadiremos una
  columna opcional std_codigo (ej: ACT, INA, ARC) o usaremos la descripción normalizada.


  4. Estrategia de Carga y Rendimiento
   1. Warm-up de Caché: Al iniciar la aplicación, el ServiceEstatus precargará todos los estatus en Redis agrupados por módulo.
   2. Single Source of Truth: El frontend solicita este JSON al iniciar sesión y lo guarda en su estado global (Redux/Pinia/Context).
   3. Validación en Base de Datos: Las tablas transaccionales (como factura o inventario) solo guardarán el id_status (UUID),
      garantizando integridad referencial.


  5. Acción Inmediata: Endpoint de Catálogo
  Implementaré un método en el servicio que transforme la lista plana de la base de datos en este mapa estructurado por módulo para
  maximizar la velocidad de lectura del cliente.