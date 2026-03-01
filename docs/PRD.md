 Documento de Requisitos del Producto (PRD) - Prunus



  ┌──────────────────────┬─────────────────────────────────┐
  │ Proyecto             │ Prunus                          │
  ├──────────────────────┼─────────────────────────────────┤
  │ Versión              │ 1.0 (Borrador)                  │
  │ Fecha                │ 21 de Febrero, 2026             │
  │ Estado               │ En Desarrollo                   │
  │ Tecnología Principal │ Go (Golang), PostgreSQL, Docker │
  └──────────────────────┴─────────────────────────────────┘

  ---


  1. Introducción
  Prunus es una API RESTful robusta y escalable diseñada para la gestión integral de compras y ventas de empresas. El
  sistema está construido bajo una arquitectura multi-tenencia (multi-empresa y multi-sucursal), permitiendo que una
  organización administre sus inventarios, usuarios, proveedores y clientes de manera centralizada pero segregada por
  sucursales.


  El objetivo es proveer un backend sólido, seguro y de alto rendimiento que sirva como núcleo para aplicaciones web o
  móviles de punto de venta (POS) y administración de recursos empresariales (ERP simplificado).


  2. Objetivos del Producto
   1. Centralización: Permitir a una empresa gestionar múltiples sucursales, usuarios y roles desde una única
      plataforma.
   2. Gestión de Inventario: Controlar el stock, precios de compra/venta y vencimientos de productos en tiempo real.
   3. Seguridad: Implementar un control de acceso basado en roles (RBAC) estricto y autenticación segura vía JWT.
   4. Auditabilidad: Mantener la integridad de los datos mediante "Soft Deletes" (eliminación lógica) y registros de
      fecha de creación/actualización.
   5. Escalabilidad: Arquitectura modular (Clean Architecture) que facilite la adición de nuevos módulos sin deuda
      técnica.


  3. Perfiles de Usuario (Roles)
  El sistema soporta roles dinámicos, pero se identifican los siguientes arquetipos base:


   * Super Admin / Dueño: Tiene acceso total a todas las sucursales de la empresa. Puede crear sucursales y asignar
     gerentes.
   * Administrador de Sucursal: Gestiona el inventario, usuarios y reportes de una sucursal específica.
   * Vendedor / Cajero: Acceso limitado a realizar ventas y consultar stock.
   * Bodeguero: Encargado de la entrada de mercancía y gestión de proveedores.


  4. Requisitos Funcionales


  4.1. Autenticación y Seguridad
   * Login: Autenticación mediante correo electrónico y contraseña.
   * JWT: Generación de Tokens de Acceso (Access Token) y Tokens de Refresco (Refresh Token).
   * Seguridad: Encriptación de contraseñas con bcrypt.
   * Sesión: Endpoint /me para obtener datos del usuario actual y su contexto (sucursal/rol).


  4.2. Gestión Organizacional (Multi-tenant)
   * Empresas:
       * Registro y edición de información fiscal (RUT, Razón Social).
       * Gestión de estados (Activo/Inactivo).
   * Sucursales:
       * Creación de múltiples sucursales vinculadas a una empresa.
       * Segregación de datos por id_sucursal.
   * Roles y Permisos:
       * Creación de roles personalizados por sucursal.
       * Asignación de roles a usuarios.


  4.3. Gestión de Usuarios
   * CRUD de Usuarios: Alta, baja (lógica) y modificación de personal.
   * Asignación: Vinculación estricta de un usuario a una Sucursal y un Rol.
   * Datos Personales: Gestión de DNI, teléfono, email y estado.


  4.4. Gestión de Catálogo e Inventario (Compra y Venta)
   * Productos:
       * Ficha completa: Nombre, Descripción, Imagen.
       * Precios: Manejo de PrecioCompra y PrecioVenta para cálculo de márgenes.
       * Stock: Control de cantidad actual (Stock) y alertas (planificado).
       * Vencimiento: Registro de FechaVencimiento para perecederos.
   * Categorización: Clasificación de productos por Categoría.
   * Unidades y Monedas: Soporte para diferentes unidades de medida (kg, unidad, litro) y tipos de moneda.


  4.5. Gestión de Terceros
   * Proveedores:
       * Registro de proveedores con RUC, contacto y dirección.
       * Vinculación a empresa y sucursal.
   * Clientes:
       * Base de datos de clientes para facturación o fidelización.
       * Datos fiscales (RUC/DNI) y de contacto.

  5. Requisitos No Funcionales


   * Rendimiento: Respuestas de la API en < 200ms para operaciones estándar.
   * Disponibilidad: Base de datos PostgreSQL robusta para integridad transaccional.
   * Mantenibilidad: Código estructurado en Clean Architecture (Capas: Transport, Service, Store, Models).
   * Portabilidad: Despliegue contenerizado mediante Docker y Docker Compose.
   * Integridad de Datos: Uso estricto de Claves Foráneas (Foreign Keys) y Soft Deletes (deleted_at) para nunca perder
     historial.

  6. Stack Tecnológico



  ┌─────────────────┬───────────────────┬───────────────────────────────────────────────┐
  │ Componente      │ Tecnología        │ Descripción                                   │
  ├─────────────────┼───────────────────┼───────────────────────────────────────────────┤
  │ Lenguaje        │ Go (Golang) 1.25+ │ Alto rendimiento y concurrencia.              │
  │ Base de Datos   │ PostgreSQL        │ Relacional, robusta y ACID.                   │
  │ Router HTTP     │ Chi Router        │ Ligero y compatible con net/http.             │
  │ Driver DB       │ pgx               │ Driver de alto rendimiento para Postgres.     │
  │ Autenticación   │ JWT (Go-JWT)      │ Stateless authentication.                     │
  │ Infraestructura │ Docker            │ Contenerización para desarrollo y despliegue. │
  └─────────────────┴───────────────────┴───────────────────────────────────────────────┘



  7. Modelo de Datos (Entidades Clave)

  Breve descripción de las tablas principales identificadas:


   * `empresa`: Entidad raíz.
   * `sucursal`: Pertenencia física, hija de Empresa.
   * `usuario`: Operador del sistema, pertenece a una Sucursal y tiene un Rol.
   * `rol`: Define permisos dentro de una Sucursal.
   * `producto`: Ítem comercializable con precio, costo y stock.
   * `proveedor`: Entidad que suministra productos.
   * `cliente`: Entidad que compra productos.

  ---


  Notas para Futuras Iteraciones (Roadmap)
   * Módulo de Transacciones: Implementar tablas de Ventas y DetalleVenta para registrar las salidas de stock.
   * Módulo de Compras: Implementar registro de Compras para aumentar el stock automáticamente.
   * Reportes: Endpoints analíticos para "Ventas por día", "Productos más vendidos", etc.