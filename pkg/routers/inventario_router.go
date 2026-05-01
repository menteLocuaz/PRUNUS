package routers

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prunus/pkg/middleware"
	transport "github.com/prunus/pkg/transport/http"
)

func InventarioRouter(h *transport.InventarioHandler) chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth())

		// CRUD y movimientos — timeout heredado del global (10s).
		r.Get("/", h.GetAll)
		r.Post("/", h.Create)
		r.Get("/{id}", h.GetByID)
		r.Get("/sucursal/{id}", h.GetBySucursal)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
		r.Post("/movimientos", h.RegistrarMovimiento)
		r.Post("/movimientos/masivo", h.RegistrarMovimientoMasivo)
		r.Get("/movimientos/{id}", h.GetMovimientos)
		r.Get("/alertas", h.GetAllAlertas)
		r.Get("/alertas/detalle", h.GetAlertasDetalle)

		// Rutas analíticas — pueden tardar más en datasets grandes.
		r.Group(func(r chi.Router) {
			r.Use(middleware.Timeout(60 * time.Second))
			r.Get("/valuacion", h.GetValuacion)
			r.Get("/rotacion", h.GetRotacion)
			r.Get("/rotacion/detalle", h.GetRotacionDetalle)
			r.Get("/composicion-categoria", h.GetComposicionCategoria)
			r.Post("/historico/snapshot", h.CapturarSnapshot)
			r.Get("/historico", h.GetValorHistorico)
			r.Get("/perdidas", h.GetPerdidas)
			r.Get("/margen", h.GetMargenGanancia)
		})
	})

	return r
}
