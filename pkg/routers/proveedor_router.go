package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	transport "github.com/prunus/pkg/transport/http"
)

func RouterProveedor(proveedorHandler *transport.ProveedorHandler) http.Handler {
	r := chi.NewRouter()

	r.Get("/proveedor", proveedorHandler.GetAll)
	r.Post("/proveedor", proveedorHandler.Create)
	r.Get("/proveedor/{id}", proveedorHandler.GetByID)
	r.Put("/proveedor/{id}", proveedorHandler.Update)
	r.Delete("/proveedor/{id}", proveedorHandler.Delete)

	return r
}
