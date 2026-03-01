package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	transport "github.com/prunus/pkg/transport/http"
)

func RouterCategoria(categoriaHandler *transport.CategoriaHandler) http.Handler {
	r := chi.NewRouter()

	r.Get("/categoria", categoriaHandler.GetAll)
	r.Post("/categoria", categoriaHandler.Create)
	r.Get("/categoria/{id}", categoriaHandler.GetByID)
	r.Put("/categoria/{id}", categoriaHandler.Update)
	r.Delete("/categoria/{id}", categoriaHandler.Delete)

	return r
}
