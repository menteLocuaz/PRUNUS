package routers

import (
	"net/http"

	"github.com/prunus/pkg/middleware"
	transport "github.com/prunus/pkg/transport/http"

	"github.com/go-chi/chi/v5"
)

func DashboardRouter(h *transport.DashboardHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequireAuth())

	r.Get("/resumen", h.GetResumen)
	r.Get("/antiguedad-deuda", h.GetAntiguedadDeuda)

	return r
}
