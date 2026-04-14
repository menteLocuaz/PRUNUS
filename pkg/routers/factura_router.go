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
		
		// Impuestos
		r.Get("/impuestos", h.GetImpuestos)
		r.Get("/impuestos/{id}", h.GetImpuestoByID)
		r.Post("/impuestos", h.CreateImpuesto)
		r.Put("/impuestos/{id}", h.UpdateImpuesto)
		r.Delete("/impuestos/{id}", h.DeleteImpuesto)

		// Formas de Pago
		r.Get("/formas-pago", h.GetFormasPago)
		r.Post("/formas-pago", h.CreateFormaPago)
		r.Put("/formas-pago/{id}", h.UpdateFormaPago)
		r.Delete("/formas-pago/{id}", h.DeleteFormaPago)
	})

	return r
}
