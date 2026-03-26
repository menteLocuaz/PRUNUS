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

type FacturaHandler struct {
	service *services.ServiceFactura
}

func NewFacturaHandler(s *services.ServiceFactura) *FacturaHandler {
	return &FacturaHandler{service: s}
}

func (h *FacturaHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Factura models.Factura           `json:"factura"`
		Items   []*models.DetalleFactura `json:"items"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}
	resp, err := h.service.CreateFactura(r.Context(), req.Factura, req.Items)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	response.Created(w, "Factura creada correctamente", resp)
}

func (h *FacturaHandler) RegistrarCompleta(w http.ResponseWriter, r *http.Request) {
	var req dto.FacturaCompletaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}

	// Obtener ID de usuario del contexto (inyectado por middleware de auth)
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Unauthorized(w, "Usuario no autenticado")
		return
	}

	resp, err := h.service.RegistrarFacturaCompleta(r.Context(), req, userID)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Created(w, "Factura procesada correctamente", resp)
}

func (h *FacturaHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}
	f, items, err := h.service.GetFactura(r.Context(), id)
	if err != nil {
		response.NotFound(w, "Factura no encontrada")
		return
	}
	response.Success(w, "Factura obtenida correctamente", map[string]interface{}{
		"factura": f,
		"items":   items,
	})
}

func (h *FacturaHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.GetAllFacturas(r.Context())
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.Success(w, "Facturas obtenidas correctamente", resp)
}

func (h *FacturaHandler) GetImpuestos(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.GetImpuestos(r.Context())
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.Success(w, "Impuestos obtenidos correctamente", resp)
}

func (h *FacturaHandler) GetFormasPago(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.GetFormasPago(r.Context())
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.Success(w, "Formas de pago obtenidas correctamente", resp)
}
