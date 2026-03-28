package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/prunus/pkg/middleware"
	transport "github.com/prunus/pkg/transport/http"
)

func POSRouter(handler *transport.POSHandler) chi.Router {
	r := chi.NewRouter()

	// Rutas protegidas por autenticación
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth())

		r.Post("/abrir", handler.AbrirCajaHandler)
		r.Post("/desmontar", handler.DesmontarCajeroHandler)
		r.Get("/estado/{id}", handler.GetEstadoCajaHandler)
	})

	return r
}
