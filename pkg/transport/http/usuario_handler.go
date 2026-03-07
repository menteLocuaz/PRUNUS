package transport

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/prunus/pkg/dto"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/services"
	"github.com/prunus/pkg/utils/response" // Importa el paquete response para respuestas estandarizadas
)

// UsuarioHandler maneja las peticiones HTTP relacionadas con usuarios
type UsuarioHandler struct {
	service *services.ServiceUsuario
}

// NewUsuarioHandler crea una nueva instancia del handler de usuario
func NewUsuarioHandler(s *services.ServiceUsuario) *UsuarioHandler {
	return &UsuarioHandler{service: s}
}

// GetAll maneja la petición GET para obtener todos los usuarios
func (h *UsuarioHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Llama al servicio para obtener todos los usuarios
	usuarios, err := h.service.GetAllUsuarios()
	if err != nil {
		// En caso de error, responde con error 500
		response.InternalServerError(w, err.Error())
		return
	}

	// Responde con éxito y los datos obtenidos
	response.Success(w, "Usuarios obtenidos correctamente", usuarios)
}

// GetByID maneja la petición GET para obtener un usuario por ID
func (h *UsuarioHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID del parámetro de la URL
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		// Si el ID no es válido, responde con error 400
		response.BadRequest(w, "ID inválido")
		return
	}

	// Llama al servicio para obtener el usuario por ID
	usuario, err := h.service.GetUsuarioByID(uint(id))
	if err != nil {
		// Si no se encuentra el usuario, responde con error 404
		response.NotFound(w, err.Error())
		return
	}

	// Responde con éxito y el usuario encontrado
	response.Success(w, "Usuario obtenido correctamente", usuario)
}

// Create maneja la petición POST para crear un nuevo usuario
func (h *UsuarioHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Decodificar el cuerpo de la petición
	var req dto.UsuarioCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Si el JSON es inválido, responde con error 400
		response.BadRequest(w, "JSON inválido")
		return
	}

	// Convertir DTO a modelo
	usuario := models.Usuario{
		IDSucursal:  req.IDSucursal,
		UsuEmail:    req.UsuEmail,
		UsuNombre:   req.UsuNombre,
		UsuDni:      req.UsuDni,
		UsuTelefono: req.UsuTelefono,
		UsuPassword: req.UsuPassword,
		Estado:      req.Estado,
		Rol: &models.Rol{
			IDRol: req.IDRol,
		},
	}

	// Crear el usuario usando el servicio
	resp, err := h.service.CreateUsuario(usuario)
	if err != nil {
		// Si hay error en la creación, responde con error 400
		response.BadRequest(w, err.Error())
		return
	}

	// Responde con código 201 y el usuario creado
	response.Created(w, "Usuario creado correctamente", resp)
}

// Update maneja la petición PUT para actualizar un usuario existente
func (h *UsuarioHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID del parámetro de la URL
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		// Si el ID no es válido, responde con error 400
		response.BadRequest(w, "ID inválido")
		return
	}

	// Decodificar el cuerpo de la petición
	var req dto.UsuarioUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Si el JSON es inválido, responde con error 400
		response.BadRequest(w, "JSON inválido")
		return
	}

	// Convertir DTO a modelo
	usuario := models.Usuario{
		IDSucursal:  req.IDSucursal,
		UsuEmail:    req.UsuEmail,
		UsuNombre:   req.UsuNombre,
		UsuDni:      req.UsuDni,
		UsuTelefono: req.UsuTelefono,
		UsuPassword: req.UsuPassword,
		Estado:      req.Estado,
		Rol: &models.Rol{
			IDRol: req.IDRol,
		},
	}

	// Actualizar el usuario usando el servicio
	resp, err := h.service.UpdateUsuario(uint(id), usuario)
	if err != nil {
		// Si hay error en la actualización, responde con error 400
		response.BadRequest(w, err.Error())
		return
	}

	// Responde con éxito y el usuario actualizado
	response.Success(w, "Usuario actualizado correctamente", resp)
}

// Delete maneja la petición DELETE para eliminar un usuario (soft delete)
func (h *UsuarioHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID del parámetro de la URL
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		// Si el ID no es válido, responde con error 400
		response.BadRequest(w, "ID inválido")
		return
	}

	// Eliminar el usuario usando el servicio
	if err := h.service.DeleteUsuario(uint(id)); err != nil {
		// Si hay error en la eliminación, responde con error 400
		response.BadRequest(w, err.Error())
		return
	}

	// Responde con código 204 No Content para indicar eliminación exitosa
	w.WriteHeader(http.StatusNoContent)
}
