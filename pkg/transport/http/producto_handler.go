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

type ProductoHandler struct {
	service *services.ServiceProducto
}

func NewProductoHandler(s *services.ServiceProducto) *ProductoHandler {
	return &ProductoHandler{service: s}
}

func (h *ProductoHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp, err := h.service.GetAllProductos()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

func (h *ProductoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	resp, err := h.service.GetProductoByID(uint(id))
	if err != nil {
		http.Error(w, "Producto no encontrado", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

func (h *ProductoHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req dto.ProductoCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	producto := models.Producto{
		Nombre:           req.Nombre,
		Descripcion:      req.Descripcion,
		PrecioCompra:     req.PrecioCompra,
		PrecioVenta:      req.PrecioVenta,
		Stock:            req.Stock,
		FechaVencimiento: req.FechaVencimiento,
		Imagen:           req.Imagen,
		Estado:           req.Estado,
		IDSucursal:       req.IDSucursal,
		IDCategoria:      req.IDCategoria,
		IDMoneda:         req.IDMoneda,
		IDUnidad:         req.IDUnidad,
	}

	resp, err := h.service.CreateProducto(producto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *ProductoHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var req dto.ProductoUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	producto := models.Producto{
		Nombre:           req.Nombre,
		Descripcion:      req.Descripcion,
		PrecioCompra:     req.PrecioCompra,
		PrecioVenta:      req.PrecioVenta,
		Stock:            req.Stock,
		FechaVencimiento: req.FechaVencimiento,
		Imagen:           req.Imagen,
		Estado:           req.Estado,
		IDSucursal:       req.IDSucursal,
		IDCategoria:      req.IDCategoria,
		IDMoneda:         req.IDMoneda,
		IDUnidad:         req.IDUnidad,
	}

	resp, err := h.service.UpdateProducto(uint(id), producto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

func (h *ProductoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteProducto(uint(id)); err != nil {
		http.Error(w, "Producto no encontrado", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
