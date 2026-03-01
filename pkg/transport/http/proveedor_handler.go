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

type ProveedorHandler struct {
	service *services.ServiceProveedor
}

func NewProveedorHandler(s *services.ServiceProveedor) *ProveedorHandler {
	return &ProveedorHandler{service: s}
}

func (h *ProveedorHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp, err := h.service.GetAllProveedores()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

func (h *ProveedorHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	resp, err := h.service.GetProveedorByID(uint(id))
	if err != nil {
		http.Error(w, "Proveedor no encontrado", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

func (h *ProveedorHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req dto.ProveedorCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	proveedor := models.Proveedor{
		Nombre:     req.Nombre,
		RUC:        req.RUC,
		Telefono:   req.Telefono,
		Direccion:  req.Direccion,
		Email:      req.Email,
		Estado:     req.Estado,
		IDSucursal: req.IDSucursal,
		IDEmpresa:  req.IDEmpresa,
	}

	resp, err := h.service.CreateProveedor(proveedor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *ProveedorHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var req dto.ProveedorUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	proveedor := models.Proveedor{
		Nombre:     req.Nombre,
		RUC:        req.RUC,
		Telefono:   req.Telefono,
		Direccion:  req.Direccion,
		Email:      req.Email,
		Estado:     req.Estado,
		IDSucursal: req.IDSucursal,
		IDEmpresa:  req.IDEmpresa,
	}

	resp, err := h.service.UpdateProveedor(uint(id), proveedor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

func (h *ProveedorHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteProveedor(uint(id)); err != nil {
		http.Error(w, "Proveedor no encontrado", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
