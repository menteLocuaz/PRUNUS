package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/prunus/pkg/middleware"
	transport "github.com/prunus/pkg/transport/http"
)

func PeriodoRouter(handler *transport.PeriodoHandler) chi.Router {
	r := chi.NewRouter()

	// Rutas protegidas (Generalmente solo para administradores)
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth())

		r.Post("/abrir", handler.AbrirPeriodoHandler)
		r.Post("/cerrar/{id}", handler.FinalizarPeriodoHandler)
		r.Get("/activo", handler.GetPeriodoActivoHandler)
	})

	return r
}
