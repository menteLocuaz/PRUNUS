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

// RolHandler maneja las peticiones HTTP relacionadas con roles
type RolHandler struct {
	service *services.ServiceRol
}

// NewRolHandler crea una nueva instancia del handler de rol
func NewRolHandler(s *services.ServiceRol) *RolHandler {
	return &RolHandler{service: s}
}

// GetAll maneja la petición GET para obtener todos los roles
func (h *RolHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	roles, err := h.service.GetAllRoles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(roles)
}

// GetByID maneja la petición GET para obtener un rol por ID
func (h *RolHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Obtener el ID del parámetro de la URL
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	rol, err := h.service.GetRolByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rol)
}

// Create maneja la petición POST para crear un nuevo rol
func (h *RolHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Decodificar el cuerpo de la petición
	var req dto.RolCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Convertir DTO a modelo
	rol := models.Rol{
		RolNombre:  req.RolNombre,
		IDSucursal: req.IDSucursal,
		Estado:     req.Estado,
	}

	// Crear el rol usando el servicio
	resp, err := h.service.CreateRol(rol)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// Update maneja la petición PUT para actualizar un rol existente
func (h *RolHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Obtener el ID del parámetro de la URL
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Decodificar el cuerpo de la petición
	var req dto.RolUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Convertir DTO a modelo
	rol := models.Rol{
		RolNombre:  req.RolNombre,
		IDSucursal: req.IDSucursal,
		Estado:     req.Estado,
	}

	// Actualizar el rol usando el servicio
	resp, err := h.service.UpdateRol(uint(id), rol)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// Delete maneja la petición DELETE para eliminar un rol (soft delete)
func (h *RolHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Obtener el ID del parámetro de la URL
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Eliminar el rol usando el servicio
	if err := h.service.DeleteRol(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Rol eliminado exitosamente",
	})
}
