package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prunus/pkg/middleware"
	transport "github.com/prunus/pkg/transport/http"
)

// NewMainRouter crea el router principal que combina todos los recursos
func NewMainRouter(
	empresaHandler *transport.EmpresaHandler,
	sucursalHandler *transport.SucursalHandler,
	rolHandler *transport.RolHandler,
	usuarioHandler *transport.UsuarioHandler,
	authHandler *transport.AuthHandler,
) http.Handler {
	r := chi.NewRouter()

	// Middleware de logging - Registra todas las peticiones HTTP
	// Para desactivar el logging, simplemente comenta la siguiente línea
	r.Use(middleware.Logger(middleware.ProductionLogConfig()))

	// Configurar rutas de API versionada
	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			// ==========================================
			// RUTAS PÚBLICAS (sin autenticación)
			// ==========================================
			r.Post("/login", authHandler.Login)

			// ==========================================
			// RUTAS PROTEGIDAS (requieren autenticación)
			// ==========================================
			r.Group(func(r chi.Router) {
				// Aplicar middleware de autenticación a todas las rutas de este grupo
				r.Use(middleware.RequireAuth())

				// Rutas de autenticación
				r.Get("/me", authHandler.GetMe)
				r.Post("/logout", authHandler.Logout)
				r.Post("/refresh-token", authHandler.RefreshToken)

				// Rutas para empresas
				r.Get("/empresas", empresaHandler.GetAll)
				r.Post("/empresas", empresaHandler.Create)
				r.Get("/empresas/{id}", empresaHandler.GetByID)
				r.Put("/empresas/{id}", empresaHandler.Update)
				r.Delete("/empresas/{id}", empresaHandler.Delete)

				// Rutas para sucursales
				r.Get("/sucursal", sucursalHandler.GetAll)
				r.Post("/sucursal", sucursalHandler.Create)
				r.Get("/sucursal/{id}", sucursalHandler.GetByID)
				r.Put("/sucursal/{id}", sucursalHandler.Update)
				r.Delete("/sucursal/{id}", sucursalHandler.Delete)

				// Rutas para roles
				r.Get("/rol", rolHandler.GetAll)
				r.Post("/rol", rolHandler.Create)
				r.Get("/rol/{id}", rolHandler.GetByID)
				r.Put("/rol/{id}", rolHandler.Update)
				r.Delete("/rol/{id}", rolHandler.Delete)

				// Rutas para usuarios
				r.Get("/usuario", usuarioHandler.GetAll)
				r.Post("/usuario", usuarioHandler.Create)
				r.Get("/usuario/{id}", usuarioHandler.GetByID)
				r.Put("/usuario/{id}", usuarioHandler.Update)
				r.Delete("/usuario/{id}", usuarioHandler.Delete)
			})
		})
	})

	return r
}
