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

type DispositivoPosHandler struct {
	service *services.ServiceDispositivoPos
}

func NewDispositivoPosHandler(s *services.ServiceDispositivoPos) *DispositivoPosHandler {
	return &DispositivoPosHandler{service: s}
}

func (h *DispositivoPosHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.GetAll(r.Context())
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.Success(w, "Dispositivos obtenidos correctamente", resp)
}

func (h *DispositivoPosHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}
	resp, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		response.NotFound(w, "Dispositivo no encontrado")
		return
	}
	response.Success(w, "Dispositivo obtenido correctamente", resp)
}

func (h *DispositivoPosHandler) Create(w http.ResponseWriter, r *http.Request) {
	var d models.DispositivoPos
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}
	resp, err := h.service.Create(r.Context(), d)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	response.Created(w, "Dispositivo creado correctamente", resp)
}

func (h *DispositivoPosHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}
	var d models.DispositivoPos
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}
	resp, err := h.service.Update(r.Context(), id, d)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	response.Success(w, "Dispositivo actualizado correctamente", resp)
}

func (h *DispositivoPosHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
	response.Success(w, "Dispositivo eliminado correctamente", nil)
}
