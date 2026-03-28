package transport

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/prunus/pkg/dto"
	"github.com/prunus/pkg/services"
	"github.com/prunus/pkg/utils"
	"github.com/prunus/pkg/utils/response"
	"github.com/prunus/pkg/utils/validator"
)

// ProductoHandler maneja las solicitudes HTTP relacionadas con productos
type ProductoHandler struct {
	service *services.ServiceProducto
}

// NewProductoHandler crea un nuevo handler con el servicio inyectado
func NewProductoHandler(s *services.ServiceProducto) *ProductoHandler {
	return &ProductoHandler{service: s}
}

// GetAll obtiene todos los productos y responde con JSON
func (h *ProductoHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	params := utils.ParsePaginationParams(r)
	resp, err := h.service.GetAllProductos(r.Context(), params)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.Success(w, "Productos obtenidos correctamente", resp)
}

// GetByID obtiene un producto por su ID
func (h *ProductoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}

	resp, err := h.service.GetProductoByID(r.Context(), id)
	if err != nil {
		response.NotFound(w, "Producto no encontrado")
		return
	}
	response.Success(w, "Producto obtenido correctamente", resp)
}

// Create crea un nuevo producto a partir del JSON recibido
func (h *ProductoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.ProductoCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		response.ValidationError(w, validator.FormatErrors(err))
		return
	}

	resp, err := h.service.CreateProducto(r.Context(), req)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	response.Created(w, "Producto creado correctamente", resp)
}

// Update actualiza un producto existente identificado por ID
func (h *ProductoHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}

	var req dto.ProductoUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		response.ValidationError(w, validator.FormatErrors(err))
		return
	}

	resp, err := h.service.UpdateProducto(r.Context(), id, req)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.Success(w, "Producto actualizado correctamente", resp)
}

// Delete elimina un producto por su ID
func (h *ProductoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}

	if err := h.service.DeleteProducto(r.Context(), id); err != nil {
		response.NotFound(w, "Producto no encontrado")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
