package transport

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/prunus/pkg/dto"
	"github.com/prunus/pkg/services"
	"github.com/prunus/pkg/utils/response"
	"github.com/prunus/pkg/utils/validator"
)

// UsuarioHandler maneja las solicitudes HTTP relacionadas con usuarios
type UsuarioHandler struct {
	service *services.ServiceUsuario
}

// NewUsuarioHandler crea una nueva instancia del handler de usuario
func NewUsuarioHandler(s *services.ServiceUsuario) *UsuarioHandler {
	return &UsuarioHandler{service: s}
}

// GetAll obtiene todos los usuarios
func (h *UsuarioHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.GetAllUsuarios(r.Context())
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.Success(w, "Usuarios obtenidos correctamente", resp)
}

// GetByID obtiene un usuario por ID
func (h *UsuarioHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		response.ValidationError(w, validator.FormatErrors(err))
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
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}

	var req dto.UsuarioUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		response.ValidationError(w, validator.FormatErrors(err))
		return
	}

	resp, err := h.service.UpdateUsuario(r.Context(), id, req.ToModel())
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.Success(w, "Usuario actualizado correctamente", resp)
}

// Administrar maneja la gestión integral de un usuario (Supermercado)
func (h *UsuarioHandler) Administrar(w http.ResponseWriter, r *http.Request) {
	var req dto.UsuarioCreateRequest // Reutilizamos el create para administración integral
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		response.ValidationError(w, validator.FormatErrors(err))
		return
	}

	// Obtener ID del admin desde el token
	adminID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Unauthorized(w, "Admin no autenticado")
		return
	}

	idStr := chi.URLParam(r, "id")
	var userID uuid.UUID
	if idStr != "" {
		userID, _ = uuid.Parse(idStr)
	}

	usuario := req.ToModel()
	usuario.IDUsuario = userID

	resp, err := h.service.AdministrarUsuario(r.Context(), usuario, adminID)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.Success(w, "Usuario gestionado correctamente", resp)
}

// Delete elimina un usuario
func (h *UsuarioHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}

	if err := h.service.DeleteUsuario(r.Context(), id); err != nil {
		response.NotFound(w, "Usuario no encontrado")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
