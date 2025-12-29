package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	transport "github.com/prunus/pkg/transport/http"
)

// RouterRol configura las rutas específicas para el recurso rol
func RouterRol(rolHandler *transport.RolHandler) http.Handler {
	r := chi.NewRouter()

	// Rutas CRUD para rol
	r.Get("/rol", rolHandler.GetAll)         // GET /api/v1/rol - Obtener todos los roles
	r.Post("/rol", rolHandler.Create)        // POST /api/v1/rol - Crear un nuevo rol
	r.Get("/rol/{id}", rolHandler.GetByID)   // GET /api/v1/rol/{id} - Obtener rol por ID
	r.Put("/rol/{id}", rolHandler.Update)    // PUT /api/v1/rol/{id} - Actualizar rol
	r.Delete("/rol/{id}", rolHandler.Delete) // DELETE /api/v1/rol/{id} - Eliminar rol (soft delete)

	return r
}
