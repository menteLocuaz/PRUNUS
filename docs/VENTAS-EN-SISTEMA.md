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


  ¿Deseas que comencemos a crear los servicios en Go para manejar esta lógica de "Apertura de Caja"?

