package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/prunus/pkg/middleware"
	transport "github.com/prunus/pkg/transport/http"
)

func UsuarioRouter(h *transport.UsuarioHandler) chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth())
		r.Get("/", h.GetAll)
		r.Post("/", h.Create)
		r.Post("/administrar", h.Administrar)
		r.Post("/administrar/{id}", h.Administrar)
		r.Get("/{id}", h.GetByID)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})

	return r
}
