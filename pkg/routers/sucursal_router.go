package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	transport "github.com/prunus/pkg/transport/http"
)

func NewRouterSucursal(sucursalHandler *transport.SucursalHandler) http.Handler {
	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Mount("/", RouterSucursal(sucursalHandler))
		})
	})

	return r
}

func RouterSucursal(sucursalHandler *transport.SucursalHandler) http.Handler {
	r := chi.NewRouter()

	r.Get("/sucursal", sucursalHandler.GetAll)
	r.Post("/sucursal", sucursalHandler.Create)
	r.Get("/sucursal/{id}", sucursalHandler.GetByID)
	r.Put("/sucursal/{id}", sucursalHandler.Update)
	r.Delete("/sucursal/{id}", sucursalHandler.Delete)

	return r
}
