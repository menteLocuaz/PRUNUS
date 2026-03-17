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

type AgregadoresHandler struct {
	service *services.ServiceAgregadores
}

func NewAgregadoresHandler(s *services.ServiceAgregadores) *AgregadoresHandler {
	return &AgregadoresHandler{service: s}
}

func (h *AgregadoresHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.GetAllAgregadores(r.Context())
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.Success(w, "Agregadores obtenidos correctamente", resp)
}

func (h *AgregadoresHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}
	resp, err := h.service.GetAgregadorByID(r.Context(), id)
	if err != nil {
		response.NotFound(w, "Agregador no encontrado")
		return
	}
	response.Success(w, "Agregador obtenido correctamente", resp)
}

func (h *AgregadoresHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.AgregadorCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		response.ValidationError(w, validator.FormatErrors(err))
		return
	}
	a := models.Agregador{Nombre: req.Nombre, Descripcion: req.Descripcion}
	resp, err := h.service.CreateAgregador(r.Context(), a)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	response.Created(w, "Agregador creado correctamente", resp)
}

func (h *AgregadoresHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}
	var req dto.AgregadorUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		response.ValidationError(w, validator.FormatErrors(err))
		return
	}
	a := models.Agregador{Nombre: req.Nombre, Descripcion: req.Descripcion}
	resp, err := h.service.UpdateAgregador(r.Context(), id, a)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.Success(w, "Agregador actualizado correctamente", resp)
}

func (h *AgregadoresHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}
	if err := h.service.DeleteAgregador(r.Context(), id); err != nil {
		response.NotFound(w, "Agregador no encontrado")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AgregadoresHandler) CreateOrden(w http.ResponseWriter, r *http.Request) {
	var req dto.OrdenAgregadorCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		response.ValidationError(w, validator.FormatErrors(err))
		return
	}
	o := models.OrdenAgregador{
		IDOrdenPedido:  req.IDOrdenPedido,
		IDAgregador:    req.IDAgregador,
		CodigoExterno:  req.CodigoExterno,
		DatosAgregador: req.DatosAgregador,
		Fecha:          time.Now(),
	}
	resp, err := h.service.CreateOrdenAgregador(r.Context(), o)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	response.Created(w, "Orden de agregador creada correctamente", resp)
}
