package transport

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/services"
	"github.com/prunus/pkg/utils/response"
)

type EstacionPosHandler struct {
	service *services.ServiceEstacionPos
}

func NewEstacionPosHandler(s *services.ServiceEstacionPos) *EstacionPosHandler {
	return &EstacionPosHandler{service: s}
}

func (h *EstacionPosHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.GetAll(r.Context())
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.Success(w, "Estaciones obtenidas correctamente", resp)
}

func (h *EstacionPosHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}
	resp, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		response.NotFound(w, "Estación no encontrada")
		return
	}
	response.Success(w, "Estación obtenida correctamente", resp)
}

func (h *EstacionPosHandler) Create(w http.ResponseWriter, r *http.Request) {
	var e models.EstacionPos
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}
	resp, err := h.service.Create(r.Context(), e)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	response.Created(w, "Estación creada correctamente", resp)
}

func (h *EstacionPosHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}
	var e models.EstacionPos
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}
	resp, err := h.service.Update(r.Context(), id, e)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	response.Success(w, "Estación actualizada correctamente", resp)
}

func (h *EstacionPosHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}
	if err := h.service.Delete(r.Context(), id); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.Success(w, "Estación eliminada correctamente", nil)
}
