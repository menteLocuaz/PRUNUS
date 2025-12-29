package transport

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/prunus/pkg/dto"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/services"
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
	w.Header().Set("Content-Type", "application/json")

	usuarios, err := h.service.GetAllUsuarios()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(usuarios)
}

// GetByID maneja la petición GET para obtener un usuario por ID
func (h *UsuarioHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Obtener el ID del parámetro de la URL
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	usuario, err := h.service.GetUsuarioByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(usuario)
}

// Create maneja la petición POST para crear un nuevo usuario
func (h *UsuarioHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Decodificar el cuerpo de la petición
	var req dto.UsuarioCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// Update maneja la petición PUT para actualizar un usuario existente
func (h *UsuarioHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Obtener el ID del parámetro de la URL
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Decodificar el cuerpo de la petición
	var req dto.UsuarioUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// Delete maneja la petición DELETE para eliminar un usuario (soft delete)
func (h *UsuarioHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Obtener el ID del parámetro de la URL
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Eliminar el usuario usando el servicio
	if err := h.service.DeleteUsuario(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Usuario eliminado exitosamente",
	})
}
