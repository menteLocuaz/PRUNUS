package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	transport "github.com/prunus/pkg/transport/http"
)

func RouterEmpresa(empresaHandler *transport.EmpresaHandler) http.Handler {
	r := chi.NewRouter()

	r.Get("/empresas", empresaHandler.GetAll)
	r.Post("/empresas", empresaHandler.Create)
	r.Get("/empresas/{id}", empresaHandler.GetByID)
	r.Put("/empresas/{id}", empresaHandler.Update)
	r.Delete("/empresas/{id}", empresaHandler.Delete)

	return r
}
