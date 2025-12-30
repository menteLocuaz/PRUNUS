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
) http.Handler {
	r := chi.NewRouter()

	// Middleware de logging - Registra todas las peticiones HTTP
	// Para desactivar el logging, simplemente comenta la siguiente línea
	r.Use(middleware.Logger(middleware.SimpleLogConfig()))

	// Configurar rutas de API versionada
	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
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

	return r
}
