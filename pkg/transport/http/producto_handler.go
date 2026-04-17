package transport

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/prunus/pkg/dto"
	"github.com/prunus/pkg/services"
	"github.com/prunus/pkg/utils"
	"github.com/prunus/pkg/utils/response"
	"github.com/prunus/pkg/utils/validator"
)

// ProductoHandler maneja las solicitudes HTTP relacionadas con productos.
// Sigue el patrón de inyección de dependencias para el servicio de productos.
type ProductoHandler struct {
	service *services.ServiceProducto
}

// NewProductoHandler crea una nueva instancia de ProductoHandler.
func NewProductoHandler(s *services.ServiceProducto) *ProductoHandler {
	return &ProductoHandler{service: s}
}

// GetAll obtiene una lista paginada de productos.
func (h *ProductoHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	params := utils.ParsePaginationParams(r)

	resp, err := h.service.GetAllProductos(r.Context(), params)
	if err != nil {
		response.InternalServerError(w, "Error al obtener productos: "+err.Error())
		return
	}

	response.Success(w, "Productos obtenidos correctamente", resp)
}

// GetByID busca un producto por su identificador único (UUID).
func (h *ProductoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "El ID proporcionado no es un UUID válido")
		return
	}

	resp, err := h.service.GetProductoByID(r.Context(), id)
	if err != nil {
		response.NotFound(w, "Producto no encontrado")
		return
	}

	response.Success(w, "Producto obtenido correctamente", resp)
}

// GetByCodigo busca un producto por su código de barras o SKU.
func (h *ProductoHandler) GetByCodigo(w http.ResponseWriter, r *http.Request) {
	codigo := chi.URLParam(r, "codigo")
	if codigo == "" {
		response.BadRequest(w, "Se requiere un código de barras o SKU")
		return
	}

	resp, err := h.service.GetProductoByCodigo(r.Context(), codigo)
	if err != nil {
		response.NotFound(w, "Producto no encontrado con el código proporcionado")
		return
	}

	response.Success(w, "Producto obtenido correctamente", resp)
}

// Create procesa la creación de un nuevo producto.
// Ahora soporta múltiples formatos de fecha gracias a JSONDate.
func (h *ProductoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.ProductoCreateRequest

	// Decodificación con manejo de errores descriptivo
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, fmt.Sprintf("Error en el formato del JSON: %v", err))
		return
	}

	// Validación de reglas de negocio en el DTO
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

// Update actualiza los datos de un producto existente.
func (h *ProductoHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "ID de producto inválido")
		return
	}

	var req dto.ProductoUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, fmt.Sprintf("Error en el formato del JSON: %v", err))
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		response.ValidationError(w, validator.FormatErrors(err))
		return
	}

	resp, err := h.service.UpdateProducto(r.Context(), id, req)
	if err != nil {
		// Diferenciamos errores de negocio de errores internos
		response.InternalServerError(w, "Error al actualizar producto: "+err.Error())
		return
	}

	response.Success(w, "Producto actualizado correctamente", resp)
}

// Delete realiza un borrado lógico del producto.
func (h *ProductoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "ID de producto inválido")
		return
	}

	if err := h.service.DeleteProducto(r.Context(), id); err != nil {
		response.NotFound(w, "No se pudo eliminar: producto no encontrado")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
