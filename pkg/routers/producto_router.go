package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	transport "github.com/prunus/pkg/transport/http"
)

func RouterProducto(productoHandler *transport.ProductoHandler) http.Handler {
	r := chi.NewRouter()

	r.Get("/producto", productoHandler.GetAll)
	r.Post("/producto", productoHandler.Create)
	r.Get("/producto/{id}", productoHandler.GetByID)
	r.Put("/producto/{id}", productoHandler.Update)
	r.Delete("/producto/{id}", productoHandler.Delete)

	return r
}
