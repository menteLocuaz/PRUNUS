package transport

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/prunus/pkg/dto"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/services"
	"github.com/prunus/pkg/utils"
	"github.com/prunus/pkg/utils/response"
	"github.com/prunus/pkg/utils/validator"
	"go.uber.org/zap"
)

type FacturaHandler struct {
	service *services.ServiceFactura
	logger  *zap.Logger
}

func NewFacturaHandler(s *services.ServiceFactura, l *zap.Logger) *FacturaHandler {
	return &FacturaHandler{service: s, logger: l}
}

func (h *FacturaHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Factura models.Factura           `json:"factura"`
		Items   []*models.DetalleFactura `json:"items"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Error decodificando factura simple", zap.Error(err))
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

	// Leer body para auditoría interna en caso de fallo
	bodyBytes, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Fallo crítico decodificación JSON Factura",
			zap.Error(err),
			zap.String("body", string(bodyBytes)))
		response.BadRequest(w, "JSON malformado o tipos de datos incompatibles (UUIDs inválidos)")
		return
	}

	// Validar
	if err := validator.Validate.Struct(req); err != nil {
		errs := validator.FormatErrors(err)
		h.logger.Warn("Error validación Factura Completa", zap.Any("errors", errs))
		response.ValidationError(w, errs)
		return
	}

	// Obtener ID de usuario
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Unauthorized(w, "Sesión de usuario inválida")
		return
	}

	resp, err := h.service.RegistrarFacturaCompleta(r.Context(), req, userID)
	if err != nil {
		h.logger.Error("Error al registrar factura completa", zap.Error(err))
		response.InternalServerError(w, "Error interno al procesar la factura")
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
	params := utils.ParsePaginationParams(r)
	resp, err := h.service.GetAllFacturas(r.Context(), params)
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

func (h *FacturaHandler) GetImpuestoByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}
	resp, err := h.service.GetImpuestoByID(r.Context(), id)
	if err != nil {
		response.NotFound(w, "Impuesto no encontrado")
		return
	}
	response.Success(w, "Impuesto obtenido correctamente", resp)
}

func (h *FacturaHandler) CreateImpuesto(w http.ResponseWriter, r *http.Request) {
	var req models.Impuesto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}

	resp, err := h.service.CreateImpuesto(r.Context(), req)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	response.Created(w, "Impuesto creado correctamente", resp)
}

func (h *FacturaHandler) UpdateImpuesto(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}

	var req models.Impuesto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}

	resp, err := h.service.UpdateImpuesto(r.Context(), id, req)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	response.Success(w, "Impuesto actualizado correctamente", resp)
}

func (h *FacturaHandler) DeleteImpuesto(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}

	if err := h.service.DeleteImpuesto(r.Context(), id); err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *FacturaHandler) GetFormasPago(ctx http.ResponseWriter, r *http.Request) {
	resp, err := h.service.GetFormasPago(r.Context())
	if err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}
	response.Success(ctx, "Formas de pago obtenidas correctamente", resp)
}

func (h *FacturaHandler) CreateFormaPago(w http.ResponseWriter, r *http.Request) {
	var req models.FormaPago
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}

	resp, err := h.service.CreateFormaPago(r.Context(), req)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	response.Created(w, "Forma de pago creada correctamente", resp)
}

func (h *FacturaHandler) UpdateFormaPago(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}

	var req models.FormaPago
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}

	resp, err := h.service.UpdateFormaPago(r.Context(), id, req)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	response.Success(w, "Forma de pago actualizada correctamente", resp)
}

func (h *FacturaHandler) DeleteFormaPago(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}

	if err := h.service.DeleteFormaPago(r.Context(), id); err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
