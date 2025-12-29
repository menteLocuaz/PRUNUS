package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	transport "github.com/prunus/pkg/transport/http"
)

func RouterSucursal(sucursalHandler *transport.SucursalHandler) http.Handler {
	r := chi.NewRouter()

	r.Get("/sucursal", sucursalHandler.GetAll)
	r.Post("/sucursal", sucursalHandler.Create)
	r.Get("/sucursal/{id}", sucursalHandler.GetByID)
	r.Put("/sucursal/{id}", sucursalHandler.Update)
	r.Delete("/sucursal/{id}", sucursalHandler.Delete)

	return r
}
