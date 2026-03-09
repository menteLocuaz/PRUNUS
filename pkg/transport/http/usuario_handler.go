package transport

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/prunus/pkg/dto"
	"github.com/prunus/pkg/models"
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
	resp, err := h.service.GetAllUsuarios()
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

	resp, err := h.service.GetUsuarioByID(id)
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

	usuario := models.Usuario{
		IDSucursal:  req.IDSucursal,
		IDRol:       req.IDRol,
		Email:       req.UsuEmail,
		UsuNombre:   req.UsuNombre,
		UsuDNI:      req.UsuDni,
		UsuTelefono: req.UsuTelefono,
		Password:    req.UsuPassword,
		IDStatus:    req.IDStatus,
	}

	resp, err := h.service.CreateUsuario(usuario)
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

	usuario := models.Usuario{
		IDSucursal:  req.IDSucursal,
		IDRol:       req.IDRol,
		Email:       req.UsuEmail,
		UsuNombre:   req.UsuNombre,
		UsuDNI:      req.UsuDni,
		UsuTelefono: req.UsuTelefono,
		Password:    req.UsuPassword,
		IDStatus:    req.IDStatus,
	}

	resp, err := h.service.UpdateUsuario(id, usuario)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.Success(w, "Usuario actualizado correctamente", resp)
}

// Delete elimina un usuario
func (h *UsuarioHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}

	if err := h.service.DeleteUsuario(id); err != nil {
		response.NotFound(w, "Usuario no encontrado")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
