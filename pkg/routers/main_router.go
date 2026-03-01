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

				r.Get("/sucursal", sucursalHandler.GetAll)
				r.Post("/sucursal", sucursalHandler.Create)
				r.Get("/sucursal/{id}", sucursalHandler.GetByID)
				r.Put("/sucursal/{id}", sucursalHandler.Update)
				r.Delete("/sucursal/{id}", sucursalHandler.Delete)

				r.Get("/rol", rolHandler.GetAll)
				r.Post("/rol", rolHandler.Create)
				r.Get("/rol/{id}", rolHandler.GetByID)
				r.Put("/rol/{id}", rolHandler.Update)
				r.Delete("/rol/{id}", rolHandler.Delete)

				r.Get("/usuario", usuarioHandler.GetAll)
				r.Post("/usuario", usuarioHandler.Create)
				r.Get("/usuario/{id}", usuarioHandler.GetByID)
				r.Put("/usuario/{id}", usuarioHandler.Update)
				r.Delete("/usuario/{id}", usuarioHandler.Delete)

				r.Get("/categoria", categoriaHandler.GetAll)
				r.Post("/categoria", categoriaHandler.Create)
				r.Get("/categoria/{id}", categoriaHandler.GetByID)
				r.Put("/categoria/{id}", categoriaHandler.Update)
				r.Delete("/categoria/{id}", categoriaHandler.Delete)

				r.Get("/cliente", clienteHandler.GetAll)
				r.Post("/cliente", clienteHandler.Create)
				r.Get("/cliente/{id}", clienteHandler.GetByID)
				r.Put("/cliente/{id}", clienteHandler.Update)
				r.Delete("/cliente/{id}", clienteHandler.Delete)

				r.Get("/medida", medidaHandler.GetAll)
				r.Post("/medida", medidaHandler.Create)
				r.Get("/medida/{id}", medidaHandler.GetByID)
				r.Put("/medida/{id}", medidaHandler.Update)
				r.Delete("/medida/{id}", medidaHandler.Delete)

				r.Get("/moneda", monedaHandler.GetAll)
				r.Post("/moneda", monedaHandler.Create)
				r.Get("/moneda/{id}", monedaHandler.GetByID)
				r.Put("/moneda/{id}", monedaHandler.Update)
				r.Delete("/moneda/{id}", monedaHandler.Delete)

				r.Get("/producto", productoHandler.GetAll)
				r.Post("/producto", productoHandler.Create)
				r.Get("/producto/{id}", productoHandler.GetByID)
				r.Put("/producto/{id}", productoHandler.Update)
				r.Delete("/producto/{id}", productoHandler.Delete)

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
