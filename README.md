# Prunus - Business Management API (ERP/POS)

**Prunus** es una API RESTful de alto rendimiento diseñada para la gestión integral de empresas, inventarios y operaciones comerciales. Construida bajo los principios de **Clean Architecture**, ofrece una base sólida y escalable para aplicaciones de Punto de Venta (POS) y administración de recursos empresariales.

---

## 🚀 Características Principales

- **Arquitectura Multi-tenant:** Gestión de múltiples empresas y sucursales con segregación estricta de datos.
- **CLI Potente:** Interfaz de línea de comandos basada en **Cobra** para gestión de servicios y migraciones.
- **Control de Acceso (RBAC):** Sistema de roles y permisos dinámicos por sucursal.
- **Gestión de Inventario:** Control de stock, precios, unidades de medida y trazabilidad.
- **Caché Distribuida:** Integración con **Redis** para optimizar consultas de catálogos frecuentes.
- **Observabilidad:** Logging estructurado en JSON compatible con **Grafana Loki**, **CloudWatch** y **Fluentd**.
- **Seguridad:** Autenticación JWT Stateless y hashing de contraseñas con `bcrypt`.

---

## 🛠️ Stack Tecnológico

| Categoría | Herramienta | Uso |
|---|---|---|
| **Lenguaje** | **Go (Golang)** 1.25.4 | Lógica de negocio y concurrencia |
| **Base de Datos** | **PostgreSQL** 15 | Persistencia relacional principal |
| **Caché** | **Redis** 7-alpine | Optimización de lecturas y catálogos |
| **Routing** | **Chi Router** v5.2.3 | Enrutamiento HTTP ligero y rápido |
| **CLI/Config** | **Cobra & Viper** | Interfaz de comandos y gestión de entorno |
| **Persistencia** | **pgx** v5.8.0 | Driver de alto rendimiento para PostgreSQL |
| **Infraestructura** | **Docker & systemd** | Contenerización y gestión de servicios |

---

## 📁 Estructura del Proyecto

```text
prunus/
├── cmd/                # Entry points (serve, migrate, root)
├── deployment/         # Archivos de despliegue (systemd, etc.)
├── docs/               # Documentación técnica y guías de observabilidad
├── pkg/
│   ├── config/         # Conexiones (DB, Redis) y Migraciones
│   ├── dto/            # Data Transfer Objects
│   ├── middleware/     # Auth, Logging (JSON), Rate Limiting
│   ├── models/         # Entidades de dominio
│   ├── services/       # Lógica de negocio (Capa Central)
│   ├── store/          # Repositorios (SQL & Redis)
│   └── transport/http/ # Handlers y controladores API
├── docker-compose.yml  # Orquestación de infraestructura local
└── .env.example        # Plantilla de configuración
```

---

## ⚙️ Instalación y Ejecución

### 1. Requisitos
- Go 1.25+
- Docker & Docker Compose

### 2. Configuración Inicial
```bash
cp .env.example .env
docker-compose up -d  # Levanta Postgres y Redis
```

### 3. Gestión vía CLI
Prunus utiliza una CLI moderna para facilitar la administración y el registro inicial de datos:

```bash
# Ejecutar migraciones de base de datos
go run ./cmd/ migrate

# Iniciar el servidor API
go run ./cmd/ serve --port 9090

# Registro de entidades base (Setup inicial)
# 1. Registrar un estatus
go run ./cmd/ register estatus --desc "Activo" --tipo "GENERAL" --modulo 1

# 2. Registrar una empresa (requiere ID de estatus)
go run ./cmd/ register empresa --nombre "Empresa Central" --rut "12345678-9" --status <UUID_STATUS>

# 3. Registrar una sucursal (requiere ID de empresa y estatus)
go run ./cmd/ register sucursal --empresa <UUID_EMPRESA> --nombre "Sucursal Norte" --status <UUID_STATUS>

# 4. Registrar un rol (requiere ID de sucursal y estatus)
go run ./cmd/ register rol --sucursal <UUID_SUCURSAL> --nombre "Administrador" --status <UUID_STATUS>

# 5. Registrar un usuario (requiere sucursal, rol, email, nombre, dni, password y estatus)
go run ./cmd/ register usuario --sucursal <UUID_SUCURSAL> --rol <UUID_ROL> --email "admin@prunus.com" --nombre "Admin" --dni "12345678" --password "secret123" --status <UUID_STATUS>
```

> Para una guía detallada de comandos y despliegue en producción, consulte: [**Guía de Comandos y Despliegue**](docs/COMANDOS-SERVER.md).

---

## 📊 Observabilidad y Monitoreo

La API genera logs estructurados en formato JSON por defecto, facilitando la integración con stacks modernos de monitoreo.

- **Destinos Soportados**: Grafana Loki, AWS CloudWatch, Fluentd.
- **Métricas**: Latencia por operación, trazabilidad de errores y auditoría.

Consulte la [**Guía de Observabilidad**](docs/grafana.md) para configurar tableros en Grafana.

---

## 🔐 Seguridad e Integridad

- **Soft Deletes**: Implementado en todas las entidades críticas (`deleted_at`).
- **Validación**: Capa de servicios con validación robusta de entrada.
- **Auditoría**: Todas las tablas incluyen `created_at` y `updated_at`.
- **RBAC**: Middleware que valida contexto de sucursal y rol en cada petición.

---

## 📄 Licencia
Este proyecto es propiedad de Prunus y se distribuye bajo la licencia MIT.
