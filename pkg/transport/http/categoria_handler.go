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

// CategoriaHandler maneja las solicitudes HTTP relacionadas con categorías
type CategoriaHandler struct {
	service *services.ServiceCategoria
}

// NewCategoriaHandler crea una nueva instancia del handler de categoría
func NewCategoriaHandler(s *services.ServiceCategoria) *CategoriaHandler {
	return &CategoriaHandler{service: s}
}

// GetAll obtiene todas las categorías
func (h *CategoriaHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Llama al servicio para obtener todas las categorías
	resp, err := h.service.GetAllCategorias()
	if err != nil {
		// En caso de error, responde con error 500
		response.InternalServerError(w, err.Error())
		return
	}

	// Responde con éxito y los datos obtenidos
	response.Success(w, "Categorías obtenidas correctamente", resp)
}

// GetByID obtiene una categoría por ID
func (h *CategoriaHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Extrae el parámetro "id" de la URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// Si el ID no es válido, responde con error 400
		response.BadRequest(w, "ID inválido")
		return
	}

	// Llama al servicio para obtener la categoría por ID
	resp, err := h.service.GetCategoriaByID(uint(id))
	if err != nil {
		// Si no se encuentra la categoría, responde con error 404
		response.NotFound(w, "Categoría no encontrada")
		return
	}

	// Responde con éxito y la categoría encontrada
	response.Success(w, "Categoría obtenida correctamente", resp)
}

// Create crea una nueva categoría
func (h *CategoriaHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Decodifica el cuerpo JSON en la estructura de solicitud
	var req dto.CategoriaCreateRequest
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

	// Crea un modelo Categoria con los datos recibidos
	categoria := models.Categoria{
		Nombre:     req.Nombre,
		IDSucursal: req.IDSucursal,
	}

	// Llama al servicio para crear la categoría
	resp, err := h.service.CreateCategoria(categoria)
	if err != nil {
		// Si hay error en la creación, responde con error 400
		response.BadRequest(w, err.Error())
		return
	}

	// Responde con código 201 y la categoría creada
	response.Created(w, "Categoría creada correctamente", resp)
}

// Update actualiza una categoría existente
func (h *CategoriaHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Extrae y valida el ID de la URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// Si el ID no es válido, responde con error 400
		response.BadRequest(w, "ID inválido")
		return
	}

	// Decodifica el JSON de actualización
	var req dto.CategoriaUpdateRequest
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

	// Crea un modelo Categoria con los datos actualizados
	categoria := models.Categoria{
		Nombre:     req.Nombre,
		IDSucursal: req.IDSucursal,
	}

	// Llama al servicio para actualizar la categoría
	resp, err := h.service.UpdateCategoria(uint(id), categoria)
	if err != nil {
		// En caso de error interno, responde con error 500
		response.InternalServerError(w, err.Error())
		return
	}

	// Responde con éxito y la categoría actualizada
	response.Success(w, "Categoría actualizada correctamente", resp)
}

// Delete elimina una categoría
func (h *CategoriaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Extrae y valida el ID de la URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// Si el ID no es válido, responde con error 400
		response.BadRequest(w, "ID inválido")
		return
	}

	// Llama al servicio para eliminar la categoría
	if err := h.service.DeleteCategoria(uint(id)); err != nil {
		// Si no se encuentra la categoría, responde con error 404
		response.NotFound(w, "Categoría no encontrada")
		return
	}

	// Responde con código 204 No Content para indicar eliminación exitosa
	w.WriteHeader(http.StatusNoContent)
}
