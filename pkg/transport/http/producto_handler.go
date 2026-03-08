package transport

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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
	resp, err := h.service.GetAllProductos()
	if err != nil {
		// En caso de error, responde con error 500
		response.InternalServerError(w, err.Error())
		return
	}

	// Responde con éxito y los datos obtenidos
	response.Success(w, "Productos obtenidos correctamente", resp)
}

// GetByID obtiene un producto por su ID
func (h *ProductoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Extrae el parámetro "id" de la URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// Si el ID no es válido, responde con error 400
		response.BadRequest(w, "ID inválido")
		return
	}

	// Llama al servicio para obtener el producto por ID
	resp, err := h.service.GetProductoByID(uint(id))
	if err != nil {
		// Si no se encuentra el producto, responde con error 404
		response.NotFound(w, "Producto no encontrado")
		return
	}

	// Responde con éxito y el producto encontrado
	response.Success(w, "Producto obtenido correctamente", resp)
}

// Create crea un nuevo producto a partir del JSON recibido
func (h *ProductoHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Decodifica el cuerpo JSON en la estructura de solicitud
	var req dto.ProductoCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Si el JSON es inválido, responde con error 400
		response.BadRequest(w, "JSON inválido")
		return
	}

	// Validar la estructura
	if err := validator.Validate.Struct(req); err != nil {
		response.ValidationError(w, validator.FormatErrors(err))
		return
	}

	// Crea un modelo Producto con los datos recibidos
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

	// Llama al servicio para crear el producto
	resp, err := h.service.CreateProducto(producto)
	if err != nil {
		// Si hay error en la creación, responde con error 400
		response.BadRequest(w, err.Error())
		return
	}

	// Responde con código 201 y el producto creado
	response.Created(w, "Producto creado correctamente", resp)
}

// Update actualiza un producto existente identificado por ID
func (h *ProductoHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Extrae y valida el ID de la URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}

	// Decodifica el JSON de actualización
	var req dto.ProductoUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}

	// Validar la estructura
	if err := validator.Validate.Struct(req); err != nil {
		response.ValidationError(w, validator.FormatErrors(err))
		return
	}

	// Crea un modelo Producto con los datos actualizados
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

	// Llama al servicio para actualizar el producto
	resp, err := h.service.UpdateProducto(uint(id), producto)
	if err != nil {
		// En caso de error interno, responde con error 500
		response.InternalServerError(w, err.Error())
		return
	}

	// Responde con éxito y el producto actualizado
	response.Success(w, "Producto actualizado correctamente", resp)
}

// Delete elimina un producto por su ID
func (h *ProductoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Extrae y valida el ID de la URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}

	// Llama al servicio para eliminar el producto
	if err := h.service.DeleteProducto(uint(id)); err != nil {
		// Si no se encuentra el producto, responde con error 404
		response.NotFound(w, "Producto no encontrado")
		return
	}

	// Responde con código 204 No Content para indicar eliminación exitosa
	w.WriteHeader(http.StatusNoContent)
}
