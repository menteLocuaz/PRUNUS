package transport

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/prunus/pkg/dto"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/services"
	"github.com/prunus/pkg/utils/response"
	"github.com/prunus/pkg/utils/validator"
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
	resp, err := h.service.GetAllMonedas()
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.Success(w, "Monedas obtenidas correctamente", resp)
}

// GetByID obtiene una moneda por su ID
func (h *MonedaHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}

	resp, err := h.service.GetMonedaByID(id)
	if err != nil {
		response.NotFound(w, "Moneda no encontrada")
		return
	}
	response.Success(w, "Moneda obtenida correctamente", resp)
}

// Create crea una nueva moneda a partir del JSON recibido
func (h *MonedaHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.MonedaCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		response.ValidationError(w, validator.FormatErrors(err))
		return
	}

	moneda := models.Moneda{
		Nombre:     req.Nombre,
		IDSucursal: req.IDSucursal,
		IDStatus:   req.IDStatus,
	}

	resp, err := h.service.CreateMoneda(moneda)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	response.Created(w, "Moneda creada correctamente", resp)
}

// Update actualiza una moneda existente identificada por ID
func (h *MonedaHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}

	var req dto.MonedaUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		response.ValidationError(w, validator.FormatErrors(err))
		return
	}

	moneda := models.Moneda{
		Nombre:     req.Nombre,
		IDSucursal: req.IDSucursal,
		IDStatus:   req.IDStatus,
	}

	resp, err := h.service.UpdateMoneda(id, moneda)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.Success(w, "Moneda actualizada correctamente", resp)
}

// Delete elimina una moneda por su ID
func (h *MonedaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}

	if err := h.service.DeleteMoneda(id); err != nil {
		response.NotFound(w, "Moneda no encontrada")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
