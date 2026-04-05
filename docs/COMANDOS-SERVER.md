# Guía de Comandos CLI - Prunus API

Esta documentación detalla el uso de la interfaz de línea de comandos (CLI) de Prunus, construida sobre **Cobra** y **Viper**. La CLI permite gestionar el ciclo de vida de la aplicación, desde la ejecución de migraciones hasta el despliegue del servidor.

## 1. Entorno de Desarrollo

Para la ejecución rápida durante el desarrollo sin necesidad de compilar manualmente, utilice `go run`.

### Iniciar Servidor
```bash
# Ejecución estándar (Puerto por defecto: 9090)
go run ./cmd/ serve

# Especificar un puerto dinámico
go run ./cmd/ serve --port 8080

# Cargar un archivo de configuración específico
go run ./cmd/ serve --config .env.dev
```

---

## 2. Producción y Despliegue

En entornos de producción (como Debian/Linux), se recomienda compilar el proyecto para obtener un binario optimizado.

### Compilación
```bash
# Generar el binario 'prunus'
go build -o prunus ./cmd/
```

### Ejecución del Binario
Una vez compilado, el binario puede ejecutarse directamente:
```bash
./prunus serve
```

---

## 3. Gestión de Base de Datos

Prunus utiliza comandos dedicados para la administración del esquema de base de datos. Se recomienda ejecutar las migraciones antes de iniciar el servicio en entornos productivos.

### Ejecutar Migraciones
```bash
./prunus migrate
```

### Flujo de Despliegue Recomendado
```bash
# 1. Actualizar esquema
./prunus migrate

# 2. Iniciar servicio si el paso anterior fue exitoso
./prunus serve
```

---

## 4. Gestión de Caché

Prunus utiliza Redis para el almacenamiento de caché de catálogos y sesiones. Puede gestionar la limpieza de la base de datos de caché de forma manual.

### Limpiar Caché (FlushDB)
```bash
# Ejecución con go run
go run ./cmd/ cache clear

# Con el binario compilado
./prunus cache clear
```

---

## 5. Referencia y Ayuda

La CLI es autodocumentada. Puede consultar la ayuda global o específica de cada comando en cualquier momento.

### Ayuda Global
```bash
./prunus --help
```

### Ayuda por Comando
```bash
./prunus serve --help
./prunus migrate --help
./prunus cache --help
```

---

## Parámetros Comunes (Flags)

| Flag | Descripción | Ejemplo |
|------|-------------|---------|
| `--port` | Define el puerto de escucha del servidor HTTP. | `--port 9090` |
| `--config` | Especifica la ruta del archivo de configuración (.env). | `--config .env.prod` |
| `--help` | Muestra información detallada del comando. | `-h` o `--help` |

---

## 6. Configuración en Debian (systemd)

Para asegurar que Prunus se inicie automáticamente tras un reinicio del servidor y se mantenga en ejecución, utilice el archivo de servicio de `systemd` incluido en `deployment/systemd/prunus.service`.

### Pasos de Instalación

1. **Preparar Directorios y Usuario**:
   ```bash
   # Crear directorio de instalación
   sudo mkdir -p /opt/prunus
   sudo chown www-data:www-data /opt/prunus

   # Copiar binario y archivos necesarios
   sudo cp prunus .env /opt/prunus/
   
   # Crear directorio de logs
   sudo mkdir -p /var/log/prunus
   sudo chown www-data:www-data /var/log/prunus
   ```

2. **Instalar el Servicio**:
   ```bash
   sudo cp deployment/systemd/prunus.service /etc/systemd/system/prunus.service
   sudo systemctl daemon-reload
   ```

3. **Gestionar el Servicio**:
   ```bash
   # Iniciar el servicio y habilitarlo al arranque
   sudo systemctl enable --now prunus

   # Consultar el estado
   sudo systemctl status prunus

   # Ver logs en tiempo real
   sudo journalctl -u prunus -f
   ```