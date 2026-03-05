// Package routers define y configura el enrutador principal de la aplicación.
// Centraliza el registro de todas las rutas HTTP agrupadas por versión de API
// y aplica los middlewares globales correspondientes.
package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prunus/pkg/middleware"
	transport "github.com/prunus/pkg/transport/http"
)

// NewMainRouter construye y retorna el enrutador HTTP principal de la aplicación.
// Recibe los handlers de cada dominio como dependencias y los registra bajo
// el prefijo /api/v1, aplicando autenticación JWT a todos los endpoints
// protegidos mediante el middleware RequireAuth.
//
// Rutas públicas (sin autenticación):
//   - POST /api/v1/login
//
// Rutas protegidas (requieren token JWT válido):
//   - Sesión:    GET /me, POST /logout, POST /refresh-token
//   - Empresas, Sucursales, Roles, Usuarios, Categorías,
//     Clientes, Medidas, Monedas, Productos, Proveedores
//     → operaciones CRUD estándar (GET, POST, PUT, DELETE)
func NewMainRouter(
	empresaHandler *transport.EmpresaHandler,
	sucursalHandler *transport.SucursalHandler,
	rolHandler *transport.RolHandler,
	usuarioHandler *transport.UsuarioHandler,
	authHandler *transport.AuthHandler,
	categoriaHandler *transport.CategoriaHandler,
	clienteHandler *transport.ClienteHandler,
	medidaHandler *transport.MedidaHandler,
	monedaHandler *transport.MonedaHandler,
	productoHandler *transport.ProductoHandler,
	proveedorHandler *transport.ProveedorHandler,
) http.Handler {
	r := chi.NewRouter()

	// Aplica el middleware de logging en modo producción a todas las peticiones entrantes.
	r.Use(middleware.Logger(middleware.ProductionLogConfig()))

	// Prefijo base de la API. Permite versionar los endpoints de forma limpia.
	r.Route("/api", func(r chi.Router) {

		// Versión 1 de la API. Agrupa todos los recursos bajo /api/v1.
		r.Route("/v1", func(r chi.Router) {

			// --- Ruta pública ---
			// No requiere autenticación. Valida credenciales y retorna un token JWT.
			r.Post("/login", authHandler.Login)

			// --- Rutas protegidas ---
			// El middleware RequireAuth valida el token JWT en cada petición.
			// Si el token es inválido o está ausente, retorna 401 Unauthorized.
			r.Group(func(r chi.Router) {
				// Se quita para el regitro
				r.Use(middleware.RequireAuth())

				// Gestión de sesión del usuario autenticado.
				r.Get("/me", authHandler.GetMe)                    // Retorna el perfil del usuario en sesión.
				r.Post("/logout", authHandler.Logout)              // Invalida el token activo.
				r.Post("/refresh-token", authHandler.RefreshToken) // Renueva el token JWT antes de su expiración.

				// CRUD: Empresa — gestión de las empresas registradas en el sistema.
				r.Get("/empresas", empresaHandler.GetAll)
				r.Post("/empresas", empresaHandler.Create)
				r.Get("/empresas/{id}", empresaHandler.GetByID)
				r.Put("/empresas/{id}", empresaHandler.Update)
				r.Delete("/empresas/{id}", empresaHandler.Delete)

				// CRUD: Sucursal — administración de sucursales por empresa.
				r.Get("/sucursal", sucursalHandler.GetAll)
				r.Post("/sucursal", sucursalHandler.Create)
				r.Get("/sucursal/{id}", sucursalHandler.GetByID)
				r.Put("/sucursal/{id}", sucursalHandler.Update)
				r.Delete("/sucursal/{id}", sucursalHandler.Delete)

				// CRUD: Rol — definición de roles para el control de acceso.
				r.Get("/rol", rolHandler.GetAll)
				r.Post("/rol", rolHandler.Create)
				r.Get("/rol/{id}", rolHandler.GetByID)
				r.Put("/rol/{id}", rolHandler.Update)
				r.Delete("/rol/{id}", rolHandler.Delete)

				// CRUD: Usuario — administración de cuentas de usuario del sistema.
				r.Get("/usuario", usuarioHandler.GetAll)
				r.Post("/usuario", usuarioHandler.Create)
				r.Get("/usuario/{id}", usuarioHandler.GetByID)
				r.Put("/usuario/{id}", usuarioHandler.Update)
				r.Delete("/usuario/{id}", usuarioHandler.Delete)

				// CRUD: Categoría — clasificación de productos o servicios.
				r.Get("/categoria", categoriaHandler.GetAll)
				r.Post("/categoria", categoriaHandler.Create)
				r.Get("/categoria/{id}", categoriaHandler.GetByID)
				r.Put("/categoria/{id}", categoriaHandler.Update)
				r.Delete("/categoria/{id}", categoriaHandler.Delete)

				// CRUD: Cliente — catálogo de clientes del negocio.
				r.Get("/cliente", clienteHandler.GetAll)
				r.Post("/cliente", clienteHandler.Create)
				r.Get("/cliente/{id}", clienteHandler.GetByID)
				r.Put("/cliente/{id}", clienteHandler.Update)
				r.Delete("/cliente/{id}", clienteHandler.Delete)

				// CRUD: Medida — unidades de medida utilizadas en productos (kg, lt, pz, etc.).
				r.Get("/medida", medidaHandler.GetAll)
				r.Post("/medida", medidaHandler.Create)
				r.Get("/medida/{id}", medidaHandler.GetByID)
				r.Put("/medida/{id}", medidaHandler.Update)
				r.Delete("/medida/{id}", medidaHandler.Delete)

				// CRUD: Moneda — divisas disponibles para operaciones y transacciones.
				r.Get("/moneda", monedaHandler.GetAll)
				r.Post("/moneda", monedaHandler.Create)
				r.Get("/moneda/{id}", monedaHandler.GetByID)
				r.Put("/moneda/{id}", monedaHandler.Update)
				r.Delete("/moneda/{id}", monedaHandler.Delete)

				// CRUD: Producto — catálogo principal de productos del inventario.
				r.Get("/producto", productoHandler.GetAll)
				r.Post("/producto", productoHandler.Create)
				r.Get("/producto/{id}", productoHandler.GetByID)
				r.Put("/producto/{id}", productoHandler.Update)
				r.Delete("/producto/{id}", productoHandler.Delete)

				// CRUD: Proveedor — catálogo de proveedores de productos o servicios.
				r.Get("/proveedor", proveedorHandler.GetAll)
				r.Post("/proveedor", proveedorHandler.Create)
				r.Get("/proveedor/{id}", proveedorHandler.GetByID)
				r.Put("/proveedor/{id}", proveedorHandler.Update)
				r.Delete("/proveedor/{id}", proveedorHandler.Delete)
			})
		})
	})

	return r
}
