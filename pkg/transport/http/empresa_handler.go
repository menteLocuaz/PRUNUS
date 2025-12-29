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

type EmpresaHandler struct {
	service *services.ServiceEmpresa
}

func NewEmpresaHandler(s *services.ServiceEmpresa) *EmpresaHandler {
	return &EmpresaHandler{
		service: s,
	}

}

// GetAll obtiene todas las empresas
func (h *EmpresaHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp, err := h.service.GetAllEmpresa()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

// GetByID obtiene una empresa por ID
func (h *EmpresaHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	resp, err := h.service.GetByIDEmpresa(uint(id))
	if err != nil {
		http.Error(w, "Empresa no encontrada", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

// Create crea una nueva empresa
func (h *EmpresaHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req dto.EmpresaCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	empresa := models.Empresa{
		Nombre: req.Nombre,
		RUT:    req.RUT,
		Estado: req.Estado,
	}

	resp, err := h.service.CrearEmpresa(empresa)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// Update actualiza una empresa existente
func (h *EmpresaHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var req dto.EmpresaUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	empresa := models.Empresa{
		Nombre: req.Nombre,
		RUT:    req.RUT,
		Estado: req.Estado,
	}

	resp, err := h.service.UpdateEmpresa(uint(id), empresa)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

// Delete elimina una empresa
func (h *EmpresaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.service.ElimminarEmpresa(uint(id)); err != nil {
		http.Error(w, "Empresa no encontrada", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
