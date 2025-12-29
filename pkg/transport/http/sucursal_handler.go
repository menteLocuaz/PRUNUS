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

type SucursalHandler struct {
	service *services.ServiceSucursal
}

func NewSucursalHandler(s *services.ServiceSucursal) *SucursalHandler {
	return &SucursalHandler{
		service: s,
	}

}

// GetAll obtiene todas las sucursales
func (h *SucursalHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp, err := h.service.GetAllSucursales()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

// GetByID obtiene una sucursal por ID
func (h *SucursalHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	resp, err := h.service.GetSucursalByID(uint(id))
	if err != nil {
		http.Error(w, "Sucursal no encontrada", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

// Create crea una nueva sucursal
func (h *SucursalHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req dto.SucursalCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	sucursal := models.Sucursal{
		IDEmpresa:      req.IDEmpresa,
		NombreSucursal: req.NombreSucursal,
		Estado:         req.Estado,
	}

	resp, err := h.service.CreateSucursal(sucursal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// Update actualiza una sucursal existente
func (h *SucursalHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var req dto.SucursalUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	sucursal := models.Sucursal{
		IDEmpresa:      req.IDEmpresa,
		NombreSucursal: req.NombreSucursal,
		Estado:         req.Estado,
	}

	resp, err := h.service.UpdateSucursal(uint(id), sucursal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

// Delete elimina una sucursal
func (h *SucursalHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteSucursal(uint(id)); err != nil {
		http.Error(w, "Sucursal no encontrada", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
