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

### Rol
GET /api/v1/rol - Obtener todos los roles
POST /api/v1/rol - Crear un nuevo rol
GET /api/v1/rol/{id} - Obtener rol por ID
PUT /api/v1/rol/{id} - Actualizar rol
DELETE /api/v1/rol/{id} - Eliminar rol (soft delete)
### Usuario
GET /api/v1/usuario - Obtener todos los usuarios (con información del rol)
POST /api/v1/usuario - Crear un nuevo usuario
GET /api/v1/usuario/{id} - Obtener usuario por ID
PUT /api/v1/usuario/{id} - Actualizar usuario
DELETE /api/v1/usuario/{id} - Eliminar usuario (soft delete)
### 🔧 Características implementadas:
Soft Delete: Todos los deletes actualizan deleted_at en lugar de eliminar físicamente
Validaciones: Email, campos requeridos, formato de email
Relaciones: Usuario incluye información del Rol mediante LEFT JOIN
Comentarios descriptivos: Todo el código está documentado
Manejo de errores: Errores claros y específicos en español
Patrón arquitectónico: Clean Architecture con separación de capas

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

Crear un usuario
```json
{
  "id_sucursal": 1,
  "id_rol": 1,
  "email": "usuario@example.com",
  "usu_nombre": "Juan Pérez",
  "usu_dni": "12345678",
  "usu_telefono": "+57 3001234567",
  "password": "MiClaveSegura123",
  "estado": 1
}

```

## Contribuir

Las contribuciones son bienvenidas. Por favor, abre un issue o un pull request para sugerencias o mejoras.

## Licencia

Este proyecto es de código abierto y está disponible bajo la licencia MIT.
