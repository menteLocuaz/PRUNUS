# Prunus - Sistema de Gestión Empresarial (ERP/POS API)

**Prunus** es una API RESTful robusta y escalable diseñada para la gestión integral de empresas, sucursales, inventarios y operaciones comerciales. Construida bajo los principios de **Clean Architecture**, ofrece una base sólida para aplicaciones de Punto de Venta (POS) y administración de recursos empresariales.

---

## 🚀 Características Principales

- **Arquitectura Multi-tenant:** Gestión de múltiples empresas y sucursales con segregación estricta de datos.
- **Control de Acceso (RBAC):** Sistema de roles y permisos dinámicos por sucursal.
- **Gestión de Inventario:** Control de stock, precios de compra/venta, unidades de medida y vencimientos.
- **Caché de Alto Rendimiento:** Integración con **Redis** para optimizar la consulta de catálogos frecuentes.
- **Seguridad:** Autenticación basada en JWT (Stateless) y hashing de contraseñas con bcrypt.
- **Integridad de Datos:** Implementación de **Soft Deletes** en todas las entidades críticas.
- **Migraciones Automáticas:** Sistema de versionado de base de datos que se ejecuta al iniciar la aplicación.

---

## 🛠️ Tecnologías

| Herramienta | Versión | Uso |
|---|---|---|
| **Go (Golang)** | 1.25.4 | Lenguaje principal y lógica de negocio |
| **PostgreSQL** | 15 | Base de datos relacional principal |
| **Redis** | 7-alpine | Sistema de caché para optimización de lectura |
| **Chi Router** | v5.2.3 | Enrutamiento HTTP ligero y rápido |
| **pgx** | v5.8.0 | Driver de alto rendimiento para PostgreSQL |
| **JWT (golang-jwt)** | v5.3.0 | Gestión de tokens de autenticación |
| **Docker & Compose** | - | Contenerización y orquestación de servicios |

---

## 📁 Estructura del Proyecto

```text
prunus/
├── cmd/
│   └── main.go                  # Punto de entrada e Inyección de Dependencias
├── docs/                        # Documentación técnica, PRD y guías de implementación
├── pkg/
│   ├── config/
│   │   └── database/
│   │       ├── connection.go    # Conexión a PostgreSQL
│   │       ├── redis.go         # Conexión a Redis
│   │       └── migrations/      # Scripts de migración de base de datos (001-031)
│   ├── dto/                     # Data Transfer Objects (Request/Response)
│   ├── helper/                  # Utilidades para JWT y hashing
│   ├── middleware/              # Auth, Logging, CORS y Rate Limiting
│   ├── models/                  # Entidades de dominio e interfaces de persistencia
│   ├── routers/                 # Definición de rutas por módulo
│   ├── services/                # Lógica de negocio y coordinación
│   ├── store/                   # Implementación de persistencia (SQL & Redis)
│   ├── transport/http/          # Handlers de la capa de transporte
│   └── utils/                   # Validadores y respuestas estandarizadas
├── docker-compose.yml           # Orquestación de Postgres y Redis
└── .env.example                 # Plantilla de variables de entorno
```

---

## ⚙️ Configuración y Despliegue

### 1. Requisitos Previos
- Go 1.25+
- Docker y Docker Compose

### 2. Preparar el Entorno
Copia el archivo de ejemplo y configura tus credenciales:
```bash
cp .env.example .env
```

### 3. Levantar Infraestructura (Postgres & Redis)
```bash
docker-compose up -d
```

### 4. Ejecutar la Aplicación
Las migraciones se ejecutarán automáticamente al iniciar:
```bash
go run cmd/main.go
```
El servidor estará disponible en `http://localhost:9090` (por defecto).

---

## 🔐 Autenticación

El sistema utiliza **JWT (JSON Web Tokens)**. Excepto por el endpoint de login, todos los demás requieren el header:
`Authorization: Bearer <tu_token>`

### Flujo Principal:
1. `POST /api/v1/login` -> Obtiene Access Token.
2. `GET /api/v1/me` -> Valida identidad y contexto (Sucursal/Rol).
3. `POST /api/v1/refresh-token` -> Renueva la sesión.

---

## 📦 Módulos y Endpoints

El sistema cuenta con una amplia gama de módulos para la gestión empresarial:

### Organización
- **Empresas:** `/api/v1/empresas`
- **Sucursales:** `/api/v1/sucursal`
- **Roles:** `/api/v1/rol`
- **Usuarios:** `/api/v1/usuario`

### Catálogos e Inventario
- **Categorías:** `/api/v1/categoria` (Con caché en Redis)
- **Productos:** `/api/v1/producto`
- **Proveedores:** `/api/v1/proveedor`
- **Clientes:** `/api/v1/cliente`
- **Monedas:** `/api/v1/moneda`
- **Unidades:** `/api/v1/medida`

### Operaciones POS (En Desarrollo/Implementados)
- **Estatus:** `/api/v1/estatus`
- **Estaciones POS:** Gestión de puntos de venta físicos.
- **Caja y Retiros:** Control de flujo de efectivo.
- **Facturación:** Emisión de comprobantes y detalles de venta.
- **Inventario:** Movimientos y ajustes de stock.

---

## ⚡ Optimización con Redis

Se ha implementado una capa de caché para los catálogos de alta frecuencia (Roles y Categorías) siguiendo el patrón de **Invalidación por Escritura**:
- **Lectura:** Se consulta primero Redis; ante un *miss*, se lee de DB y se puebla la caché (TTL 1h).
- **Escritura (Create/Update/Delete):** Se invalida automáticamente la entrada en Redis para garantizar consistencia.

---

## 🛡️ Estándares de Desarrollo

- **Idiomática:** Mensajes de error y comentarios en **Español**.
- **Seguridad:** Uso estricto de consultas parametrizadas para evitar Inyección SQL.
- **Calidad:** Validaciones robustas en la capa de servicios antes de la persistencia.
- **Auditoría:** Todas las tablas incluyen `created_at`, `updated_at` y `deleted_at`.

---

## 📄 Licencia
Este proyecto es de código abierto bajo la licencia MIT.
