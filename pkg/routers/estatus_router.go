package routers

import (
	"github.com/go-chi/chi/v5"
	transport "github.com/prunus/pkg/transport/http"
)

func EstatusRouter(h *transport.EstatusHandler) chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		// r.Use(middleware.RequireAuth())
		r.Get("/", h.GetAll)
		r.Get("/catalogo", h.GetMasterCatalog)
		r.Post("/", h.Create)
		r.Get("/{id}", h.GetByID)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
		r.Get("/tipo/{tipo}", h.GetByTipo)
		r.Get("/modulo/{moduloID}", h.GetByModulo)
	})

	return r
}
