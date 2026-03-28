package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/prunus/pkg/middleware"
	transport "github.com/prunus/pkg/transport/http"
)

// ConfiguracionPosRouter define las rutas para el handler de configuración de impresión.
func ConfiguracionPosRouter(h *transport.ConfiguracionHandler) chi.Router {
	r := chi.NewRouter()

	// Rutas protegidas por autenticación
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth())

		r.Get("/canales/{chainId}", h.GetCanales)
		r.Get("/impresoras/{restId}", h.GetImpresoras)
		r.Get("/puertos", h.GetPuertos)
	})

	return r
}
