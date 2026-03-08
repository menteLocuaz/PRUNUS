la implementación de Redis para el almacenamiento en caché de catálogos con alta frecuencia de acceso (categorías y
  roles), siguiendo los principios de Clean Architecture.

  Resumen de la Implementación:


   1. Infraestructura: Se añadió Redis al docker-compose.yml y las variables correspondientes a .env.example.
   2. Capa de Datos: Se creó un RedisStore que implementa una interfaz genérica CacheStore definida en la capa de modelos. Esto
      permite que los servicios no dependan directamente de la librería de Redis.
   3. Lógica de Negocio (Servicios):
       * Lectura: GetAll y GetByID en los servicios de Categoría y Rol ahora consultan primero Redis. Si hay un cache miss, consultan
         la base de datos y guardan el resultado en Redis por 1 hora.
       * Escritura/Invalidación: Los métodos Create, Update y Delete ahora invalidan automáticamente las llaves de caché afectadas
         para asegurar la consistencia de los datos.
   4. Inyección de Dependencias: El punto de entrada cmd/main.go ahora inicializa el cliente de Redis y lo inyecta en los servicios
      correspondientes.


  Próximos Pasos Sugeridos:
   * Monitoreo: Podría añadir un middleware de logging para ver los Cache Hits y Misses en la consola.
   * Otros Catálogos: Este mismo patrón puede extenderse fácilmente a otros servicios (ej. Moneda, Unidad de Medida) simplemente
     inyectando el cacheStore en sus constructores.
   * Pruebas: Iniciar los servicios con docker-compose up -d y realizar peticiones GET para verificar que la velocidad de respuesta
     mejora tras la primera consulta.


  El proyecto compila correctamente y está listo para ser probado en un entorno con Redis.