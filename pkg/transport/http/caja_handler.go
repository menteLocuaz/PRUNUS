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
)

type CajaHandler struct {
	service *services.ServiceCaja
}

func NewCajaHandler(s *services.ServiceCaja) *CajaHandler {
	return &CajaHandler{service: s}
}

func (h *CajaHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.GetAllCajas(r.Context())
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.Success(w, "Cajas obtenidas correctamente", resp)
}

func (h *CajaHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}
	resp, err := h.service.GetCajaByID(r.Context(), id)
	if err != nil {
		response.NotFound(w, "Caja no encontrada")
		return
	}
	response.Success(w, "Caja obtenida correctamente", resp)
}

func (h *CajaHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.Caja
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}
	resp, err := h.service.CreateCaja(r.Context(), req)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	response.Created(w, "Caja creada correctamente", resp)
}

func (h *CajaHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}
	var req models.Caja
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}
	resp, err := h.service.UpdateCaja(r.Context(), id, req)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.Success(w, "Caja actualizada correctamente", resp)
}

func (h *CajaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}
	if err := h.service.DeleteCaja(r.Context(), id); err != nil {
		response.NotFound(w, "Caja no encontrada")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *CajaHandler) AbrirSesion(w http.ResponseWriter, r *http.Request) {
	var req struct {
		IDCaja        uuid.UUID             `json:"id_caja"`
		MontoApertura float64               `json:"monto_apertura"`
		Desglose      []dto.DenominacionDTO `json:"desglose"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}

	if req.IDCaja == uuid.Nil {
		response.BadRequest(w, "El ID de la caja es requerido")
		return
	}

	// Obtener ID de usuario del contexto (inyectado por middleware)
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Unauthorized(w, "Usuario no autenticado")
		return
	}

	resp, err := h.service.AbrirSesion(r.Context(), req.IDCaja, userID, req.MontoApertura, req.Desglose)
	if err != nil {
		// Si el error es por sesión ya abierta, devolvemos Conflict (409) o BadRequest (400)
		// dependiendo de la convención, aquí usaremos BadRequest para simplificar el manejo en frontend
		response.BadRequest(w, err.Error())
		return
	}
	response.Created(w, "Sesión de caja abierta correctamente", resp)
}

func (h *CajaHandler) CerrarSesion(w http.ResponseWriter, r *http.Request) {
	var req dto.CierreCajaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}

	// Si el ID viene en la URL, priorizarlo
	idStr := chi.URLParam(r, "id")
	if idStr != "" {
		if id, err := uuid.Parse(idStr); err == nil {
			req.IDControlEstacion = id
		}
	}

	// Obtener ID de usuario del contexto (inyectado por middleware)
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Unauthorized(w, "Usuario no autenticado")
		return
	}

	resp, err := h.service.ArqueoYCierre(r.Context(), req, userID)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, "Arqueo y cierre completado exitosamente", resp)
}

func (h *CajaHandler) RegistrarMovimiento(w http.ResponseWriter, r *http.Request) {
	var req models.MovimientoCaja
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}
	resp, err := h.service.RegistrarMovimiento(r.Context(), req)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	response.Created(w, "Movimiento registrado correctamente", resp)
}

func (h *CajaHandler) GetMovimientos(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID de sesión inválido")
		return
	}
	resp, err := h.service.GetMovimientos(r.Context(), id)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.Success(w, "Movimientos obtenidos correctamente", resp)
}
