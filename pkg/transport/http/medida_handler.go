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

type MedidaHandler struct {
	service *services.ServiceUnidad
}

func NewMedidaHandler(s *services.ServiceUnidad) *MedidaHandler {
	return &MedidaHandler{service: s}
}

func (h *MedidaHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.GetAllUnidades()
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, "Unidades obtenidas correctamente", resp)
}

func (h *MedidaHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}

	resp, err := h.service.GetUnidadByID(uint(id))
	if err != nil {
		response.NotFound(w, "Medida no encontrada")
		return
	}

	response.Success(w, "Medida obtenida correctamente", resp)
}

func (h *MedidaHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.MedidaCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}

	// Validar la estructura
	if err := validator.Validate.Struct(req); err != nil {
		response.ValidationError(w, validator.FormatErrors(err))
		return
	}

	unidad := models.Unidad{
		Nombre:     req.Nombre,
		IDSucursal: req.IDSucursal,
	}

	resp, err := h.service.CreateUnidad(unidad)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Created(w, "Unidad creada correctamente", resp)
}

func (h *MedidaHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}

	var req dto.MedidaUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}

	// Validar la estructura
	if err := validator.Validate.Struct(req); err != nil {
		response.ValidationError(w, validator.FormatErrors(err))
		return
	}

	unidad := models.Unidad{
		Nombre:     req.Nombre,
		IDSucursal: req.IDSucursal,
	}

	resp, err := h.service.UpdateUnidad(uint(id), unidad)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, "Unidad actualizada correctamente", resp)
}

func (h *MedidaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}

	if err := h.service.DeleteUnidad(uint(id)); err != nil {
		response.NotFound(w, "Medida no encontrada")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
