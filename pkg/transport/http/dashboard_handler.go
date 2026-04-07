package transport

import (
	"net/http"

	"github.com/prunus/pkg/services"
	"github.com/prunus/pkg/utils/response"

	"github.com/google/uuid"
)

type DashboardHandler struct {
	service services.ServiceDashboard
}

func NewDashboardHandler(s services.ServiceDashboard) *DashboardHandler {
	return &DashboardHandler{service: s}
}

func (h *DashboardHandler) GetResumen(w http.ResponseWriter, r *http.Request) {
	sucursalID, ok := r.Context().Value("user_sucursal").(uuid.UUID)
	if !ok {
		response.BadRequest(w, "Sucursal no identificada en el contexto")
		return
	}

	data, err := h.service.GetDashboardData(r.Context(), sucursalID)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, "Resumen de dashboard obtenido correctamente", data)
}

func (h *DashboardHandler) GetAntiguedadDeuda(w http.ResponseWriter, r *http.Request) {
	sucursalID, ok := r.Context().Value("user_sucursal").(uuid.UUID)
	if !ok {
		response.BadRequest(w, "Sucursal no identificada en el contexto")
		return
	}

	data, err := h.service.GetAntiguedadDeuda(r.Context(), sucursalID)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, "Reporte de antigüedad de deuda obtenido correctamente", data)
}

func (h *DashboardHandler) GetComposicionCategoria(w http.ResponseWriter, r *http.Request) {
	sucursalID, ok := r.Context().Value("user_sucursal").(uuid.UUID)
	if !ok {
		response.BadRequest(w, "Sucursal no identificada en el contexto")
		return
	}

	data, err := h.service.GetComposicionCategoria(r.Context(), sucursalID)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, "Composición por categoría obtenida correctamente", data)
}

func (h *DashboardHandler) GetMermas(w http.ResponseWriter, r *http.Request) {
	sucursalID, ok := r.Context().Value("user_sucursal").(uuid.UUID)
	if !ok {
		response.BadRequest(w, "Sucursal no identificada en el contexto")
		return
	}

	data, err := h.service.GetMermas(r.Context(), sucursalID)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, "Reporte de mermas obtenido correctamente", data)
}
