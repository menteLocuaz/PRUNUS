package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	transport "github.com/prunus/pkg/transport/http"
)

func RouterCliente(clienteHandler *transport.ClienteHandler) http.Handler {
	r := chi.NewRouter()

	r.Get("/cliente", clienteHandler.GetAll)
	r.Post("/cliente", clienteHandler.Create)
	r.Get("/cliente/{id}", clienteHandler.GetByID)
	r.Put("/cliente/{id}", clienteHandler.Update)
	r.Delete("/cliente/{id}", clienteHandler.Delete)

	return r
}
