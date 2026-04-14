# Referencia de Comandos CLI - Prunus API

Esta guía proporciona una referencia técnica detallada de la interfaz de línea de comandos (CLI) de Prunus. La CLI está diseñada para gestionar el ciclo de vida de la aplicación, la base de datos y los servicios auxiliares.

## 1. Comando `serve`

Inicia el servidor API REST y levanta los servicios en segundo plano (workers).

### Uso
```bash
./prunus serve [flags]
```

### Flags Específicos
| Flag | Tipo | Descripción | Por defecto |
|------|------|-------------|-------------|
| `--port`, `-p` | string | Puerto de escucha del servidor HTTP. | `9090` |
| `--auto-migrate` | bool | Ejecuta automáticamente las migraciones pendientes (`migrate up`) antes de iniciar el servidor. | `false` |

### Comportamiento al inicio
- Valida las variables de entorno críticas (`DB_*`, `JWT_SECRET`).
- Establece conexiones a PostgreSQL y Redis.
- **Trazabilidad de Versión:** Si `--auto-migrate` no está activo, el servidor consulta y muestra en los logs la versión actual de la base de datos y su estado (limpio/dirty).
- Inicia el worker de snapshots de inventario (cada 24 horas).

---

## 2. Comando `migrate`

Gestión avanzada del esquema de la base de datos mediante archivos SQL en `database/migrations`.

### Subcomandos

#### `migrate up`
Aplica todas las migraciones pendientes para actualizar el esquema a la última versión disponible.
```bash
./prunus migrate up
```

#### `migrate down [n]`
Revierte las migraciones aplicadas. Por defecto revierte la última.
```bash
# Revertir solo la última
./prunus migrate down

# Revertir las últimas 3
./prunus migrate down 3
```

#### `migrate version`
Muestra la versión actual aplicada en la base de datos y detecta si el estado es **DIRTY** (indicando una migración fallida que requiere intervención).
```bash
./prunus migrate version
```

#### `migrate force <version>`
Fuerza la base de datos a una versión específica. Se utiliza principalmente para limpiar el estado **DIRTY** después de corregir manualmente un error de migración.
```bash
./prunus migrate force 41
```

---

## 3. Comando `seed`

Puebla la base de datos con datos maestros esenciales para el funcionamiento inicial del sistema.

### Uso
```bash
./prunus seed
```

### Datos Incluidos
- **Módulos Base:** Inserta o actualiza los módulos de Configuración, Sucursales, Usuarios, Productos, Ventas y Caja.
- **Permisos Administrativos:** Asigna permisos completos (Lectura, Escritura, Actualización, Borrado) sobre todos los módulos a cualquier rol con el nombre "Administrador".

---

## 4. Comando `cache`

Gestión del almacenamiento volátil en Redis.

### Subcomandos
- `clear`: Ejecuta un `FlushDB` en la instancia de Redis configurada para limpiar todos los datos cacheados (catálogos, sesiones, etc.).

---

## Flags Globales

Estos flags están disponibles para todos los comandos de la CLI.

| Flag | Descripción |
|------|-------------|
| `--config` | Ruta al archivo de configuración de entorno (ej. `.env.production`). Por defecto busca `.env` en el directorio raíz. |
| `--help`, `-h` | Muestra la ayuda detallada del comando o subcomando. |

---

## Flujo de Despliegue Recomendado (Receta)

Para un despliegue seguro en entornos de producción, se recomienda el siguiente orden de ejecución:

1. **Migrar:** `./prunus migrate up` (o usar el flag `--auto-migrate` en el paso 3).
2. **Seed (Opcional):** `./prunus seed` (solo si es una instalación nueva o se añadieron módulos).
3. **Servir:** `./prunus serve --port 80`
