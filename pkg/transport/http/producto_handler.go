package transport

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/prunus/pkg/dto"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/services"
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
	resp, err := h.service.GetAllProductos(r.Context())
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

	var fechaVenc *time.Time
	if !req.FechaVencimiento.IsZero() {
		fechaVenc = &req.FechaVencimiento
	}

	producto := models.Producto{
		Nombre:           req.Nombre,
		Descripcion:      req.Descripcion,
		PrecioCompra:     req.PrecioCompra,
		PrecioVenta:      req.PrecioVenta,
		Stock:            req.Stock,
		FechaVencimiento: fechaVenc,
		Imagen:           req.Imagen,
		IDStatus:         req.IDStatus,
		IDSucursal:       req.IDSucursal,
		IDCategoria:      req.IDCategoria,
		IDMoneda:         req.IDMoneda,
		IDUnidad:         req.IDUnidad,
	}

	resp, err := h.service.CreateProducto(r.Context(), producto)
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

	var fechaVenc *time.Time
	if !req.FechaVencimiento.IsZero() {
		fechaVenc = &req.FechaVencimiento
	}

	producto := models.Producto{
		Nombre:           req.Nombre,
		Descripcion:      req.Descripcion,
		PrecioCompra:     req.PrecioCompra,
		PrecioVenta:      req.PrecioVenta,
		Stock:            req.Stock,
		FechaVencimiento: fechaVenc,
		Imagen:           req.Imagen,
		IDStatus:         req.IDStatus,
		IDSucursal:       req.IDSucursal,
		IDCategoria:      req.IDCategoria,
		IDMoneda:         req.IDMoneda,
		IDUnidad:         req.IDUnidad,
	}

	resp, err := h.service.UpdateProducto(r.Context(), id, producto)
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
