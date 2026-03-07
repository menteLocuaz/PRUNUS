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
	categoriaHandler *transport.CategoriaHandler,
	clienteHandler *transport.ClienteHandler,
	medidaHandler *transport.MedidaHandler,
	monedaHandler *transport.MonedaHandler,
	productoHandler *transport.ProductoHandler,
	proveedorHandler *transport.ProveedorHandler,
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger(middleware.ProductionLogConfig()))

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Post("/login", authHandler.Login)

			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireAuth())

				r.Get("/me", authHandler.GetMe)
				r.Post("/logout", authHandler.Logout)
				r.Post("/refresh-token", authHandler.RefreshToken)

				r.Get("/empresas", empresaHandler.GetAll)
				r.Post("/empresas", empresaHandler.Create)
				r.Get("/empresas/{id}", empresaHandler.GetByID)
				r.Put("/empresas/{id}", empresaHandler.Update)
				r.Delete("/empresas/{id}", empresaHandler.Delete)

				r.Get("/sucursales", sucursalHandler.GetAll)
				r.Post("/sucursales", sucursalHandler.Create)
				r.Get("/sucursales/{id}", sucursalHandler.GetByID)
				r.Put("/sucursales/{id}", sucursalHandler.Update)
				r.Delete("/sucursales/{id}", sucursalHandler.Delete)

				r.Get("/roles", rolHandler.GetAll)
				r.Post("/roles", rolHandler.Create)
				r.Get("/roles/{id}", rolHandler.GetByID)
				r.Put("/roles/{id}", rolHandler.Update)
				r.Delete("/roles/{id}", rolHandler.Delete)

				r.Get("/usuarios", usuarioHandler.GetAll)
				r.Post("/usuarios", usuarioHandler.Create)
				r.Get("/usuarios/{id}", usuarioHandler.GetByID)
				r.Put("/usuarios/{id}", usuarioHandler.Update)
				r.Delete("/usuarios/{id}", usuarioHandler.Delete)

				r.Get("/categorias", categoriaHandler.GetAll)
				r.Post("/categorias", categoriaHandler.Create)
				r.Get("/categorias/{id}", categoriaHandler.GetByID)
				r.Put("/categorias/{id}", categoriaHandler.Update)
				r.Delete("/categorias/{id}", categoriaHandler.Delete)

				r.Get("/clientes", clienteHandler.GetAll)
				r.Post("/clientes", clienteHandler.Create)
				r.Get("/clientes/{id}", clienteHandler.GetByID)
				r.Put("/clientes/{id}", clienteHandler.Update)
				r.Delete("/clientes/{id}", clienteHandler.Delete)

				r.Get("/medidas", medidaHandler.GetAll)
				r.Post("/medidas", medidaHandler.Create)
				r.Get("/medidas/{id}", medidaHandler.GetByID)
				r.Put("/medidas/{id}", medidaHandler.Update)
				r.Delete("/medidas/{id}", medidaHandler.Delete)

				r.Get("/monedas", monedaHandler.GetAll)
				r.Post("/monedas", monedaHandler.Create)
				r.Get("/monedas/{id}", monedaHandler.GetByID)
				r.Put("/monedas/{id}", monedaHandler.Update)
				r.Delete("/monedas/{id}", monedaHandler.Delete)

				r.Get("/productos", productoHandler.GetAll)
				r.Post("/productos", productoHandler.Create)
				r.Get("/productos/{id}", productoHandler.GetByID)
				r.Put("/productos/{id}", productoHandler.Update)
				r.Delete("/productos/{id}", productoHandler.Delete)

				r.Get("/proveedores", proveedorHandler.GetAll)
				r.Post("/proveedores", proveedorHandler.Create)
				r.Get("/proveedores/{id}", proveedorHandler.GetByID)
				r.Put("/proveedores/{id}", proveedorHandler.Update)
				r.Delete("/proveedores/{id}", proveedorHandler.Delete)
			})
		})
	})

	return r
}
