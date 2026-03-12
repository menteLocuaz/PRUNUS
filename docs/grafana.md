Actualmente, el sistema utiliza Structured Logging (slog) de Go, y el almacenamiento de los logs se maneja de la siguiente manera:


  1. Ubicación Actual: Consola (Stdout)
  En cmd/main.go, el logger está configurado para escribir en la salida estándar (os.Stdout) en formato JSON:


   1 // cmd/main.go
   2 logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
   3 slog.SetDefault(logger)


  Esto significa que, por defecto, los logs se ven en la terminal donde corre el servidor. En entornos modernos (como Docker o
  Kubernetes), esto es lo ideal, ya que los recolectores de logs (como CloudWatch, Fluentd o Grafana Loki) capturan automáticamente
  todo lo que sale por la consola.


  2. Formato Estructurado (JSON)
  Al ser JSON, cada log de rendimiento que añadimos se guarda con campos específicos que facilitan las búsquedas:


   1 {
   2   "time": "2026-03-11T...",
   3   "level": "WARN",
   4   "msg": "Operación lenta detectada",
   5   "capa": "store",
   6   "operacion": "GetAllUsuarios",
   7   "latencia_ms": 250,
   8   "threshold": "200ms"
   9 }