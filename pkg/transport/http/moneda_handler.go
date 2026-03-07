package transport

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/prunus/pkg/dto"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/services"
	"github.com/prunus/pkg/utils/response" // Importa el paquete response para respuestas estandarizadas
)

// MonedaHandler maneja las solicitudes HTTP relacionadas con las monedas
type MonedaHandler struct {
	service *services.ServiceMoneda
}

// NewMonedaHandler crea un nuevo handler con el servicio inyectado
func NewMonedaHandler(s *services.ServiceMoneda) *MonedaHandler {
	return &MonedaHandler{service: s}
}

// GetAll obtiene todas las monedas y responde con JSON
func (h *MonedaHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Llama al servicio para obtener todas las monedas
	resp, err := h.service.GetAllMonedas()
	if err != nil {
		// En caso de error, responde con error 500
		response.InternalServerError(w, err.Error())
		return
	}

	// Responde con éxito y los datos obtenidos
	response.Success(w, "Monedas obtenidas correctamente", resp)
}

// GetByID obtiene una moneda por su ID
func (h *MonedaHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Extrae el parámetro "id" de la URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// Si el ID no es válido, responde con error 400
		response.BadRequest(w, "ID inválido")
		return
	}

	// Llama al servicio para obtener la moneda por ID
	resp, err := h.service.GetMonedaByID(uint(id))
	if err != nil {
		// Si no se encuentra la moneda, responde con error 404
		response.NotFound(w, "Moneda no encontrada")
		return
	}

	// Responde con éxito y la moneda encontrada
	response.Success(w, "Moneda obtenida correctamente", resp)
}

// Create crea una nueva moneda a partir del JSON recibido
func (h *MonedaHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Decodifica el cuerpo JSON en la estructura de solicitud
	var req dto.MonedaCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Si el JSON es inválido, responde con error 400
		response.BadRequest(w, "JSON inválido")
		return
	}

	// Crea un modelo Moneda con los datos recibidos
	moneda := models.Moneda{
		Nombre:     req.Nombre,
		IDSucursal: req.IDSucursal,
		Estado:     req.Estado,
	}

	// Llama al servicio para crear la moneda
	resp, err := h.service.CreateMoneda(moneda)
	if err != nil {
		// Si hay error en la creación, responde con error 400
		response.BadRequest(w, err.Error())
		return
	}

	// Responde con código 201 y la moneda creada
	response.Created(w, "Moneda creada correctamente", resp)
}

// Update actualiza una moneda existente identificada por ID
func (h *MonedaHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Extrae y valida el ID de la URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}

	// Decodifica el JSON de actualización
	var req dto.MonedaUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}

	// Crea un modelo Moneda con los datos actualizados
	moneda := models.Moneda{
		Nombre:     req.Nombre,
		IDSucursal: req.IDSucursal,
		Estado:     req.Estado,
	}

	// Llama al servicio para actualizar la moneda
	resp, err := h.service.UpdateMoneda(uint(id), moneda)
	if err != nil {
		// En caso de error interno, responde con error 500
		response.InternalServerError(w, err.Error())
		return
	}

	// Responde con éxito y la moneda actualizada
	response.Success(w, "Moneda actualizada correctamente", resp)
}

// Delete elimina una moneda por su ID
func (h *MonedaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Extrae y valida el ID de la URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}

	// Llama al servicio para eliminar la moneda
	if err := h.service.DeleteMoneda(uint(id)); err != nil {
		// Si no se encuentra la moneda, responde con error 404
		response.NotFound(w, "Moneda no encontrada")
		return
	}

	// Responde con código 204 No Content para indicar eliminación exitosa
	w.WriteHeader(http.StatusNoContent)
}
