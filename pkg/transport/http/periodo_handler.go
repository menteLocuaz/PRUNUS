package transport

import (
	"encoding/json"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/prunus/pkg/services"
	"github.com/prunus/pkg/utils/response"
	"go.uber.org/zap"
)

type PeriodoHandler struct {
	service *services.ServicePeriodo
}

func NewPeriodoHandler(s *services.ServicePeriodo) *PeriodoHandler {
	return &PeriodoHandler{service: s}
}

// AbrirPeriodoHandler inicia un nuevo periodo con auditoría completa
func (h *PeriodoHandler) AbrirPeriodoHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Obtener datos del contexto (inyectados por Auth Middleware)
	idUsuario, _ := r.Context().Value("user_id").(uuid.UUID)
	idSucursal, _ := r.Context().Value("user_sucursal").(uuid.UUID)

	// 2. Capturar IP Real
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ip = xff
	}

	// 3. Decodificar motivo opcional
	var body struct {
		Motivo string `json:"motivo"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)

	// 4. Procesar apertura
	result, err := h.service.AbrirNuevoPeriodo(r.Context(), idUsuario, idSucursal, ip, body.Motivo)
	if err != nil {
		h.service.Logger().Error("Fallo crítico al intentar abrir periodo",
			zap.Error(err),
			zap.String("sucursal", idSucursal.String()))
		response.InternalServerError(w, "Error al procesar la apertura del periodo")
		return
	}

	response.Created(w, "Periodo contable procesado exitosamente", result)
}

// FinalizarPeriodoHandler realiza el cierre seguro
func (h *PeriodoHandler) FinalizarPeriodoHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	idPeriodo, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID de periodo inválido")
		return
	}

	idUsuario, _ := r.Context().Value("user_id").(uuid.UUID)
	
	// Capturar IP de cierre
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ip = xff
	}

	err = h.service.FinalizarPeriodo(r.Context(), idPeriodo, idUsuario, ip)
	if err != nil {
		response.InternalServerError(w, "Error al finalizar el periodo: "+err.Error())
		return
	}

	response.Success(w, "Periodo finalizado y cerrado correctamente", nil)
}

// GetPeriodoActivoHandler obtiene el periodo de la sucursal actual
func (h *PeriodoHandler) GetPeriodoActivoHandler(w http.ResponseWriter, r *http.Request) {
	idSucursal, _ := r.Context().Value("user_sucursal").(uuid.UUID)
	
	result, err := h.service.GetActivePeriodo(r.Context(), idSucursal)
	if err != nil {
		response.InternalServerError(w, "Error al consultar el periodo")
		return
	}

	if result == nil {
		response.NotFound(w, "No existe un periodo activo para su sucursal")
		return
	}

	response.Success(w, "Periodo activo obtenido", result)
}
