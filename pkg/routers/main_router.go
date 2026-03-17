package routers

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prunus/pkg/middleware"
	transport "github.com/prunus/pkg/transport/http"
)

// Handlers agrupa todos los handlers de la aplicación
type Handlers struct {
	Empresa    *transport.EmpresaHandler
	Sucursal   *transport.SucursalHandler
	Rol        *transport.RolHandler
	Usuario    *transport.UsuarioHandler
	Auth       *transport.AuthHandler
	Categoria  *transport.CategoriaHandler
	Cliente    *transport.ClienteHandler
	Medida     *transport.MedidaHandler
	Moneda     *transport.MonedaHandler
	Producto   *transport.ProductoHandler
	Proveedor  *transport.ProveedorHandler
	Estatus    *transport.EstatusHandler
	POS        *transport.POSHandler
	Inventario *transport.InventarioHandler
}

// NewMainRouter crea el router principal que combina todos los recursos
func NewMainRouter(h *Handlers) http.Handler {
	r := chi.NewRouter()

	// Middleware Global
	r.Use(middleware.CORS())
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger(middleware.ProductionLogConfig()))

	r.Route("/api", func(r chi.Router) {
		// Rate limiting para toda la API (100 peticiones por minuto por IP)
		r.Use(middleware.RateLimit(100, 1*time.Minute))

		r.Route("/v1", func(r chi.Router) {
			// Auth Routes
			r.Mount("/auth", AuthRouter(h.Auth))

			// Mantener login en /v1/login por compatibilidad
			r.Post("/login", h.Auth.Login)

			// Resource Routes
			r.Mount("/empresas", EmpresaRouter(h.Empresa))
			r.Mount("/sucursales", SucursalRouter(h.Sucursal))
			r.Mount("/roles", RolRouter(h.Rol))
			r.Mount("/usuarios", UsuarioRouter(h.Usuario))
			r.Mount("/categorias", CategoriaRouter(h.Categoria))
			r.Mount("/clientes", ClienteRouter(h.Cliente))
			r.Mount("/medidas", MedidaRouter(h.Medida))
			r.Mount("/monedas", MonedaRouter(h.Moneda))
			r.Mount("/productos", ProductoRouter(h.Producto))
			r.Mount("/proveedores", ProveedorRouter(h.Proveedor))
			r.Mount("/estatus", EstatusRouter(h.Estatus))
			r.Mount("/pos", POSRouter(h.POS))
			r.Mount("/inventario", InventarioRouter(h.Inventario))
		})
	})

	return r
}
