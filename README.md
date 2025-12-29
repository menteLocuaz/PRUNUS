# Prunus

API REST desarrollada en Go para la gestión de empresas.

## Tecnologías

- **Go** 1.25.4
- **PostgreSQL** - Base de datos
- **Chi Router** - Framework de enrutamiento HTTP
- **pgx** - Driver de PostgreSQL para Go
- **godotenv** - Manejo de variables de entorno

## Estructura del Proyecto

```
prunus/
├── cmd/
│   └── main.go           # Punto de entrada de la aplicación
├── pkg/
│   ├── config/
│   │   └── database/     # Configuración de base de datos
│   ├── dto/              # Data Transfer Objects
│   ├── models/           # Modelos de datos
│   ├── routers/          # Definición de rutas
│   ├── services/         # Lógica de negocio
│   ├── store/            # Capa de acceso a datos
│   └── transport/
│       └── http/         # Handlers HTTP
├── .env.example          # Ejemplo de variables de entorno
├── docker-compose.yml    # Configuración de Docker
└── go.mod                # Dependencias del proyecto
```

## Configuración

1. Clona el repositorio:
```bash
git clone https://github.com/tu-usuario/prunus.git
cd prunus
```

2. Copia el archivo de ejemplo de variables de entorno:
```bash
cp .env.example .env
```

3. Configura las variables de entorno en `.env`:
```env
DB_HOST=localhost
DB_USER=tu_usuario
DB_PASSWORD=tu_contraseña
DB_NAME=nombre_base_datos
DB_PORT=5432
DB_SSLMODE=disable
```

4. Instala las dependencias:
```bash
go mod download
```

## Ejecución

### Con Docker Compose
```bash
docker-compose up -d
```

### Local
```bash
go run cmd/main.go
```

## API Endpoints

### Empresas

- `GET /empresas` - Obtener todas las empresas
- `GET /empresas/{id}` - Obtener una empresa por ID
- `POST /empresas` - Crear una nueva empresa
- `PUT /empresas/{id}` - Actualizar una empresa
- `DELETE /empresas/{id}` - Eliminar una empresa (soft delete)

### Sucursal

- GET /api/v1/sucursal - Obtener todas las sucursales
- POST /api/v1/sucursal - Crear sucursal
- GET /api/v1/sucursal/{id} - Obtener sucursal por ID
- PUT /api/v1/sucursal/{id} - Actualizar sucursal
- DELETE /api/v1/sucursal/{id} - Eliminar sucursa

### Ejemplo de uso

Crear una empresa:
```bash
curl -X POST http://localhost:8080/empresas \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Mi Empresa S.A.",
    "rut": "12.345.678-9",
    "estado": 1
  }'
```

## Contribuir

Las contribuciones son bienvenidas. Por favor, abre un issue o un pull request para sugerencias o mejoras.

## Licencia

Este proyecto es de código abierto y está disponible bajo la licencia MIT.
