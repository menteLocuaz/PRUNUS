package transport

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/prunus/pkg/dto"
	"github.com/prunus/pkg/services"
	"github.com/prunus/pkg/utils"
	"github.com/prunus/pkg/utils/request"
	"github.com/prunus/pkg/utils/response"
)

// UsuarioHandler maneja las solicitudes HTTP relacionadas con usuarios
type UsuarioHandler struct {
	service *services.ServiceUsuario
}

// NewUsuarioHandler crea una nueva instancia del handler de usuario
func NewUsuarioHandler(s *services.ServiceUsuario) *UsuarioHandler {
	return &UsuarioHandler{service: s}
}

// GetAll obtiene una lista paginada de usuarios
func (h *UsuarioHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	params := utils.ParsePaginationParams(r)
	resp, err := h.service.GetAllUsuarios(r.Context(), params)
	if err != nil {
		response.InternalServerError(w, "Error al obtener usuarios")
		return
	}
	response.Success(w, "Usuarios obtenidos correctamente", resp)
}

// GetByID obtiene un usuario por ID
func (h *UsuarioHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := request.GetID(r, "id")
	if err != nil {
		response.HandleError(w, err)
		return
	}

	resp, err := h.service.GetUsuarioByID(r.Context(), id)
	if err != nil {
		response.NotFound(w, "Usuario no encontrado")
		return
	}
	response.Success(w, "Usuario obtenido correctamente", resp)
}

// Create crea un nuevo usuario
func (h *UsuarioHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.UsuarioCreateRequest
	if err := request.DecodeAndValidate(r, &req); err != nil {
		response.HandleError(w, err)
		return
	}

	resp, err := h.service.CreateUsuario(r.Context(), req.ToModel())
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	response.Created(w, "Usuario creado correctamente", resp)
}

// Update actualiza un usuario existente
func (h *UsuarioHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := request.GetID(r, "id")
	if err != nil {
		response.HandleError(w, err)
		return
	}

	var req dto.UsuarioUpdateRequest
	if err := request.DecodeAndValidate(r, &req); err != nil {
		response.HandleError(w, err)
		return
	}

	resp, err := h.service.UpdateUsuario(r.Context(), id, req.ToModel())
	if err != nil {
		response.InternalServerError(w, "Error al actualizar usuario")
		return
	}
	response.Success(w, "Usuario actualizado correctamente", resp)
}

// Administrar maneja la gestión integral de un usuario (Supermercado)
func (h *UsuarioHandler) Administrar(w http.ResponseWriter, r *http.Request) {
	var req dto.UsuarioCreateRequest // Reutilizamos el create para administración integral
	if err := request.DecodeAndValidate(r, &req); err != nil {
		response.HandleError(w, err)
		return
	}

	// Obtener ID del admin desde el token
	adminID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Unauthorized(w, "Admin no autenticado")
		return
	}

	id, _ := request.GetID(r, "id") // Puede ser uuid.Nil si no está presente

	usuario := req.ToModel()
	usuario.IDUsuario = id

	resp, err := h.service.AdministrarUsuario(r.Context(), usuario, adminID)
	if err != nil {
		response.InternalServerError(w, "Error al gestionar usuario")
		return
	}
	response.Success(w, "Usuario gestionado correctamente", resp)
}

// Delete elimina un usuario
func (h *UsuarioHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := request.GetID(r, "id")
	if err != nil {
		response.HandleError(w, err)
		return
	}

	if err := h.service.DeleteUsuario(r.Context(), id); err != nil {
		response.NotFound(w, "Usuario no encontrado")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
