package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	transport "github.com/prunus/pkg/transport/http"
)

// RouterUsuario configura las rutas específicas para el recurso usuario
func RouterUsuario(usuarioHandler *transport.UsuarioHandler) http.Handler {
	r := chi.NewRouter()

	// Rutas CRUD para usuario
	r.Get("/usuario", usuarioHandler.GetAll)         // GET /api/v1/usuario - Obtener todos los usuarios
	r.Post("/usuario", usuarioHandler.Create)        // POST /api/v1/usuario - Crear un nuevo usuario
	r.Get("/usuario/{id}", usuarioHandler.GetByID)   // GET /api/v1/usuario/{id} - Obtener usuario por ID
	r.Put("/usuario/{id}", usuarioHandler.Update)    // PUT /api/v1/usuario/{id} - Actualizar usuario
	r.Delete("/usuario/{id}", usuarioHandler.Delete) // DELETE /api/v1/usuario/{id} - Eliminar usuario (soft delete)

	return r
}
