package transport

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/prunus/pkg/services"
	"github.com/prunus/pkg/utils/response"
)

type PeriodoHandler struct {
	service *services.ServicePeriodo
}

func NewPeriodoHandler(s *services.ServicePeriodo) *PeriodoHandler {
	return &PeriodoHandler{service: s}
}

// AbrirPeriodoHandler inicia un nuevo periodo contable
func (h *PeriodoHandler) AbrirPeriodoHandler(w http.ResponseWriter, r *http.Request) {
	idUsuario, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Unauthorized(w, "Sesión de usuario no válida")
		return
	}

	result, err := h.service.AbrirNuevoPeriodo(r.Context(), idUsuario)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Created(w, "Periodo contable abierto exitosamente", result)
}

// FinalizarPeriodoHandler realiza el cierre seguro (valida estaciones cerradas)
func (h *PeriodoHandler) FinalizarPeriodoHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	idPeriodo, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID de periodo inválido")
		return
	}

	idUsuario, _ := r.Context().Value("user_id").(uuid.UUID)

	err = h.service.FinalizarPeriodo(r.Context(), idPeriodo, idUsuario)
	if err != nil {
		// Aquí capturamos el error de "estaciones abiertas"
		response.BadRequest(w, err.Error())
		return
	}

	response.Success(w, "Periodo finalizado y cerrado correctamente", nil)
}

// GetPeriodoActivoHandler obtiene la información del periodo actual
func (h *PeriodoHandler) GetPeriodoActivoHandler(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.GetActivePeriodo(r.Context())
	if err != nil {
		response.InternalServerError(w, "Error al consultar el periodo")
		return
	}

	if result == nil {
		response.NotFound(w, "No existe un periodo activo actualmente")
		return
	}

	response.Success(w, "Periodo activo obtenido", result)
}
