package transport

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/services"
	"github.com/prunus/pkg/utils"
	"github.com/prunus/pkg/utils/response"
)

type OrdenPedidoHandler struct {
	service *services.ServiceOrdenPedido
}

func NewOrdenPedidoHandler(s *services.ServiceOrdenPedido) *OrdenPedidoHandler {
	return &OrdenPedidoHandler{service: s}
}

func (h *OrdenPedidoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var o models.OrdenPedido
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}
	resp, err := h.service.CreateOrden(r.Context(), o)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	response.Created(w, "Orden de pedido creada correctamente", resp)
}

func (h *OrdenPedidoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}
	resp, err := h.service.GetOrdenByID(r.Context(), id)
	if err != nil {
		response.NotFound(w, "Orden no encontrada")
		return
	}
	response.Success(w, "Orden obtenida correctamente", resp)
}

func (h *OrdenPedidoHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	params := utils.ParsePaginationParams(r)
	resp, err := h.service.GetAllOrdenes(r.Context(), params)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.Success(w, "Ordenes obtenidas correctamente", resp)
}

func (h *OrdenPedidoHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}
	var req struct {
		IDStatus uuid.UUID `json:"id_status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}
	if err := h.service.UpdateStatus(r.Context(), id, req.IDStatus); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.Success(w, "Estatus de la orden actualizado correctamente", nil)
}
