package transport

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/prunus/pkg/dto"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/services"
	"github.com/prunus/pkg/utils/response"
	"github.com/prunus/pkg/utils/validator"
)

type EstatusHandler struct {
	service *services.ServiceEstatus
}

func NewEstatusHandler(s *services.ServiceEstatus) *EstatusHandler {
	return &EstatusHandler{service: s}
}

func (h *EstatusHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.GetAllEstatus()
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.Success(w, "Estatus obtenidos correctamente", resp)
}

func (h *EstatusHandler) GetMasterCatalog(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.GetMasterCatalog()
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.Success(w, "Catálogo maestro de estatus obtenido correctamente", resp)
}

func (h *EstatusHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}

	resp, err := h.service.GetEstatusByID(id)
	if err != nil {
		response.NotFound(w, "Estatus no encontrado")
		return
	}

	response.Success(w, "Estatus obtenido correctamente", resp)
}

func (h *EstatusHandler) GetByTipo(w http.ResponseWriter, r *http.Request) {
	tipo := chi.URLParam(r, "tipo")
	if tipo == "" {
		response.BadRequest(w, "Tipo de estado es obligatorio")
		return
	}

	resp, err := h.service.GetEstatusByTipo(tipo)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, "Estatus por tipo obtenidos correctamente", resp)
}

func (h *EstatusHandler) GetByModulo(w http.ResponseWriter, r *http.Request) {
	moduloIDStr := chi.URLParam(r, "moduloID")
	moduloID, err := strconv.Atoi(moduloIDStr)
	if err != nil {
		response.BadRequest(w, "Modulo ID inválido")
		return
	}

	resp, err := h.service.GetEstatusByModulo(moduloID)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, "Estatus por modulo obtenidos correctamente", resp)
}

func (h *EstatusHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateEstatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		response.ValidationError(w, validator.FormatErrors(err))
		return
	}

	estatus := models.Estatus{
		StdDescripcion: req.StdDescripcion,
		StpTipoEstado:  req.StpTipoEstado,
		MdlID:          req.MdlID,
	}

	resp, err := h.service.CreateEstatus(estatus)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Created(w, "Estatus creado correctamente", resp)
}

func (h *EstatusHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}

	var req dto.UpdateEstatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		response.ValidationError(w, validator.FormatErrors(err))
		return
	}

	estatus := models.Estatus{
		StdDescripcion: req.StdDescripcion,
		StpTipoEstado:  req.StpTipoEstado,
		MdlID:          req.MdlID,
	}

	resp, err := h.service.UpdateEstatus(id, estatus)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, "Estatus actualizado correctamente", resp)
}

func (h *EstatusHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}

	if err := h.service.DeleteEstatus(id); err != nil {
		response.NotFound(w, "Estatus no encontrado")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
