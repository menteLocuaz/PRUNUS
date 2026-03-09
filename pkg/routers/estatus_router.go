package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	transport "github.com/prunus/pkg/transport/http"
)

func RouterEstatus(estatusHandler *transport.EstatusHandler) http.Handler {
	r := chi.NewRouter()

	r.Get("/estatus", estatusHandler.GetAll)
	r.Post("/estatus", estatusHandler.Create)
	r.Get("/estatus/{id}", estatusHandler.GetByID)
	r.Put("/estatus/{id}", estatusHandler.Update)
	r.Delete("/estatus/{id}", estatusHandler.Delete)

	// Rutas adicionales específicas
	r.Get("/estatus/tipo/{tipo}", estatusHandler.GetByTipo)
	r.Get("/estatus/modulo/{moduloID}", estatusHandler.GetByModulo)

	return r
}
