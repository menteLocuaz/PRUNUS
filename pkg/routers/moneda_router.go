package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	transport "github.com/prunus/pkg/transport/http"
)

func RouterMoneda(monedaHandler *transport.MonedaHandler) http.Handler {
	r := chi.NewRouter()

	r.Get("/moneda", monedaHandler.GetAll)
	r.Post("/moneda", monedaHandler.Create)
	r.Get("/moneda/{id}", monedaHandler.GetByID)
	r.Put("/moneda/{id}", monedaHandler.Update)
	r.Delete("/moneda/{id}", monedaHandler.Delete)

	return r
}
