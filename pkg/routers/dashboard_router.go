package routers

import (
	"net/http"
	"time"

	"github.com/prunus/pkg/middleware"
	transport "github.com/prunus/pkg/transport/http"

	"github.com/go-chi/chi/v5"
)

func DashboardRouter(h *transport.DashboardHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequireAuth())
	// Las rutas de dashboard ejecutan queries analíticas complejas; 60s supera el timeout global.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/resumen", h.GetResumen)
	r.Get("/antiguedad-deuda", h.GetAntiguedadDeuda)
	r.Get("/composicion-categoria", h.GetComposicionCategoria)
	r.Get("/mermas", h.GetMermas)

	return r
}
