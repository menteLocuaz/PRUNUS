package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	transport "github.com/prunus/pkg/transport/http"
)

func RouterMedida(medidaHandler *transport.MedidaHandler) http.Handler {
	r := chi.NewRouter()

	r.Get("/medida", medidaHandler.GetAll)
	r.Post("/medida", medidaHandler.Create)
	r.Get("/medida/{id}", medidaHandler.GetByID)
	r.Put("/medida/{id}", medidaHandler.Update)
	r.Delete("/medida/{id}", medidaHandler.Delete)

	return r
}
