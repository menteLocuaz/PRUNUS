package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/prunus/pkg/middleware"
	transport "github.com/prunus/pkg/transport/http"
)

func FacturaRouter(h *transport.FacturaHandler) chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth())
		r.Get("/", h.GetAll)
		r.Post("/", h.Create)
		r.Post("/completa", h.RegistrarCompleta)
		r.Get("/{id}", h.GetByID)
		r.Get("/impuestos", h.GetImpuestos)
		r.Get("/formas-pago", h.GetFormasPago)
	})

	return r
}
