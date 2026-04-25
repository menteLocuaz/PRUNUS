package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/prunus/pkg/middleware"
	transport "github.com/prunus/pkg/transport/http"
)

func CajaRouter(h *transport.CajaHandler) chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth())
		r.Get("/", h.GetAll)
		r.Post("/", h.Create)
		r.Get("/{id}", h.GetByID)
		r.Post("/sesion/abrir", h.AbrirSesion)
		r.Post("/sesion/cerrar/{id}", h.CerrarSesion)
	})

	return r
}
