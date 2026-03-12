Monitorización de Rendimiento Proactiva, implementando un sistema de trazado de latencia tanto para la
  base de datos como para servicios externos (Redis e integraciones futuras).


  Cambios Realizados:


   1. Utilidad de Rendimiento Centralizada:
       * Se creó pkg/utils/performance/performance.go que contiene la función Trace. Esta función mide el tiempo de ejecución y
         registra un log de nivel WARN (usando slog.WarnContext) si se superan los umbrales predefinidos.
       * Umbrales establecidos:
           * Base de Datos (SQL): 200ms
           * Redis: 50ms
           * Servicios Externos: 500ms


   2. Instrumentación en la Capa Store (SQL):
       * Se actualizaron todos los archivos en pkg/store/ (Usuario, Categoría, Cliente, Empresa, Estatus, Medida, Moneda, POS,
         Producto, Proveedor, Rol, Sucursal).
       * Cada método que realiza operaciones de base de datos ahora incluye un defer performance.Trace(...) al inicio para capturar
         consultas lentas de forma proactiva.


   3. Instrumentación de Redis:
       * Se actualizó pkg/store/redis_store.go para monitorear las operaciones de caché, tratando a Redis como un servicio de
         infraestructura crítica con su propio umbral de rendimiento.


   4. Preparación para Servicios Externos:
       * La infraestructura de trazado está lista para ser utilizada en futuras integraciones con proveedores de pagos, correos o
         agregadores (UberEats, Rappi), cumpliendo con el objetivo de capturar cuellos de botella antes de que causen tiempos de
         espera (timeouts).

Estos cambios están diseñados para ser transparentes y automáticos en las operaciones existentes, pero aquí te explico cómo
  aprovecharlos y aplicarlos a futuro:


  1. Monitoreo Automático (Ya activo)
  No necesitas hacer nada para las funciones actuales del store. Si una consulta a la base de datos tarda más de 200ms o una
  operación en Redis más de 50ms, verás un log de advertencia en la consola automáticamente:


   1 # Ejemplo de lo que verás en los logs si algo va lento:
   2 WARN Operación lenta detectada capa=store operacion=GetAllProductos latencia_ms=245 threshold=200ms

  2. Cómo usarlo en nuevas funciones
  Si creas un nuevo método en el store o un servicio externo (ej: una API de pagos), solo debes añadir una línea al inicio del
  método:


   1 func (s *tuStore) MiNuevaConsulta(ctx context.Context) error {
   2     // Esta línea mide el tiempo desde 'ahora' hasta que la función termine (defer)
   3     defer performance.Trace(ctx, "store", "MiNuevaConsulta", performance.DbThreshold, time.Now())
   4
   5     // ... resto del código (query SQL, etc)
   6 }

  3. Uso para Servicios Externos (Integraciones)
  Si integras un proveedor externo (como Stripe o un servicio de correos), usa el threshold de servicios externos:


   1 func (s *EmailService) EnviarCorreo(ctx context.Context) {
   2     // Los servicios externos suelen ser más lentos, por eso el threshold es de 500ms
   3     defer performance.Trace(ctx, "external", "EnviarCorreo", performance.ExternalThreshold, time.Now())
   4
   5     // Llamada a la API externa...
   6 }

  4. Dónde ajustar los límites
  Si consideras que 200ms es muy poco (o mucho) para la base de datos, puedes centralizar los cambios en un solo lugar:

  Archivo: pkg/utils/performance/performance.go


   1 const (
   2     DbThreshold       = 200 * time.Millisecond // Cambia este valor aquí
   3     RedisThreshold    = 50 * time.Millisecond
   4     ExternalThreshold = 500 * time.Millisecond
   5 )


  Resumen: Los cambios te permiten identificar cuellos de botella sin necesidad de herramientas externas complejas, simplemente
  revisando los logs de nivel WARN en tu sistema de logging actual.