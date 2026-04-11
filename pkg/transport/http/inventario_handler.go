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
	"github.com/prunus/pkg/utils"
	"github.com/prunus/pkg/utils/response"
	"github.com/prunus/pkg/utils/validator"
)

type InventarioHandler struct {
	service *services.ServiceInventario
}

func NewInventarioHandler(s *services.ServiceInventario) *InventarioHandler {
	return &InventarioHandler{service: s}
}

func (h *InventarioHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	params := utils.ParsePaginationParams(r)
	resp, err := h.service.GetAllInventario(r.Context(), params)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.Success(w, "Inventario obtenido correctamente", resp)
}

func (h *InventarioHandler) GetBySucursal(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID de sucursal inválido")
		return
	}

	params := utils.ParsePaginationParams(r)
	resp, err := h.service.GetInventarioBySucursal(r.Context(), id, params)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.Success(w, "Inventario de la sucursal obtenido correctamente", resp)
}

func (h *InventarioHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}

	resp, err := h.service.GetInventarioByID(r.Context(), id)
	if err != nil {
		response.NotFound(w, "Inventario no encontrado")
		return
	}
	response.Success(w, "Inventario obtenido correctamente", resp)
}

func (h *InventarioHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.InventarioCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		response.ValidationError(w, validator.FormatErrors(err))
		return
	}

	inventario := models.Inventario{
		IDProducto:   req.IDProducto,
		IDSucursal:   req.IDSucursal,
		StockActual:  req.StockActual,
		StockMinimo:  req.StockMinimo,
		StockMaximo:  req.StockMaximo,
		PrecioCompra: req.PrecioCompra,
		PrecioVenta:  req.PrecioVenta,
	}

	resp, err := h.service.CreateInventario(r.Context(), inventario)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	response.Created(w, "Inventario creado correctamente", resp)
}

func (h *InventarioHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}

	var req dto.InventarioUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		response.ValidationError(w, validator.FormatErrors(err))
		return
	}

	inventario := models.Inventario{
		StockActual:  req.StockActual,
		StockMinimo:  req.StockMinimo,
		StockMaximo:  req.StockMaximo,
		PrecioCompra: req.PrecioCompra,
		PrecioVenta:  req.PrecioVenta,
	}

	resp, err := h.service.UpdateInventario(r.Context(), id, inventario)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.Success(w, "Inventario actualizado correctamente", resp)
}

func (h *InventarioHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}

	if err := h.service.DeleteInventario(r.Context(), id); err != nil {
		response.NotFound(w, "Inventario no encontrado")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *InventarioHandler) RegistrarMovimiento(w http.ResponseWriter, r *http.Request) {
	var req dto.MovimientoCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		response.ValidationError(w, validator.FormatErrors(err))
		return
	}

	// Obtener ID de usuario del contexto (inyectado por middleware)
	userID, _ := r.Context().Value("user_id").(uuid.UUID)

	movimiento := models.MovimientoInventario{
		IDProducto:     req.IDProducto,
		IDSucursal:     req.IDSucursal,
		TipoMovimiento: req.TipoMovimiento,
		Cantidad:       req.Cantidad,
		Referencia:     req.Referencia,
		IDUsuario:      userID,
	}

	resp, err := h.service.RegistrarMovimiento(r.Context(), movimiento)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	response.Created(w, "Movimiento de inventario registrado correctamente", resp)
}

func (h *InventarioHandler) RegistrarMovimientoMasivo(w http.ResponseWriter, r *http.Request) {
	var req dto.MovimientoMasivoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		response.ValidationError(w, validator.FormatErrors(err))
		return
	}

	userID, _ := r.Context().Value("user_id").(uuid.UUID)

	items := make([]models.MovimientoItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = models.MovimientoItem{
			IDProducto: item.IDProducto,
			Cantidad:   item.Cantidad,
		}
	}

	resp, err := h.service.RegistrarMovimientoMasivo(r.Context(), req.IDSucursal, userID, req.TipoMovimiento, req.Referencia, items)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	response.Created(w, "Movimientos de inventario registrados correctamente", resp)
}

func (h *InventarioHandler) GetMovimientos(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID de producto inválido")
		return
	}
	params := utils.ParsePaginationParams(r)
	resp, err := h.service.GetMovimientos(r.Context(), id, params)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.Success(w, "Movimientos obtenidos correctamente", resp)
}

func (h *InventarioHandler) GetAllAlertas(w http.ResponseWriter, r *http.Request) {
	sucursalIDStr := r.URL.Query().Get("id_sucursal")
	if sucursalIDStr == "" {
		// Si no viene sucursal, intentar obtenerla del token/contexto
		ctxSucursal, ok := r.Context().Value("user_sucursal").(uuid.UUID)
		if !ok {
			response.BadRequest(w, "Debe proporcionar id_sucursal")
			return
		}
		sucursalIDStr = ctxSucursal.String()
	}

	sucursalID, err := uuid.Parse(sucursalIDStr)
	if err != nil {
		response.BadRequest(w, "id_sucursal inválido")
		return
	}

	resp, err := h.service.GetAlertasStock(r.Context(), sucursalID)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.Success(w, "Alertas de stock obtenidas correctamente", resp)
}

func (h *InventarioHandler) GetValuacion(w http.ResponseWriter, r *http.Request) {
	sucursalIDStr := r.URL.Query().Get("id_sucursal")
	metodo := r.URL.Query().Get("metodo") // peps, ueps, promedio
	if metodo == "" {
		metodo = "promedio"
	}

	if sucursalIDStr == "" {
		ctxSucursal, ok := r.Context().Value("user_sucursal").(uuid.UUID)
		if !ok {
			response.BadRequest(w, "Debe proporcionar id_sucursal")
			return
		}
		sucursalIDStr = ctxSucursal.String()
	}

	sucursalID, err := uuid.Parse(sucursalIDStr)
	if err != nil {
		response.BadRequest(w, "id_sucursal inválido")
		return
	}

	total, err := h.service.GetValuacion(r.Context(), sucursalID, metodo)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	data := map[string]interface{}{
		"id_sucursal": sucursalID,
		"metodo":      metodo,
		"total_valor": total,
	}
	response.Success(w, "Valuación de inventario calculada correctamente", data)
}

func (h *InventarioHandler) GetRotacion(w http.ResponseWriter, r *http.Request) {
	sucursalIDStr := r.URL.Query().Get("id_sucursal")
	if sucursalIDStr == "" {
		ctxSucursal, ok := r.Context().Value("user_sucursal").(uuid.UUID)
		if !ok {
			response.BadRequest(w, "Debe proporcionar id_sucursal")
			return
		}
		sucursalIDStr = ctxSucursal.String()
	}

	sucursalID, err := uuid.Parse(sucursalIDStr)
	if err != nil {
		response.BadRequest(w, "id_sucursal inválido")
		return
	}

	abc, err := h.service.GetAnalisisRotacion(r.Context(), sucursalID)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, "Análisis de rotación ABC obtenido correctamente", abc)
}

func (h *InventarioHandler) GetRotacionDetalle(w http.ResponseWriter, r *http.Request) {
	sucursalID, ok := parseSucursalID(r)
	if !ok {
		response.BadRequest(w, "Debe proporcionar id_sucursal")
		return
	}

	fechaInicioStr := r.URL.Query().Get("fecha_inicio")
	fechaFinStr := r.URL.Query().Get("fecha_fin")

	var fechaInicio, fechaFin time.Time
	var err error

	// Si no se proporcionan fechas, usamos los últimos 30 días por defecto
	if fechaInicioStr == "" {
		fechaInicio = time.Now().AddDate(0, 0, -30)
	} else {
		fechaInicio, err = time.Parse(time.RFC3339, fechaInicioStr)
		if err != nil {
			// Intentar formato corto si falla el largo
			fechaInicio, err = time.Parse("2006-01-02", fechaInicioStr)
			if err != nil {
				response.BadRequest(w, "fecha_inicio inválida, use formato YYYY-MM-DD o RFC3339")
				return
			}
		}
	}

	if fechaFinStr == "" {
		fechaFin = time.Now()
	} else {
		fechaFin, err = time.Parse(time.RFC3339, fechaFinStr)
		if err != nil {
			fechaFin, err = time.Parse("2006-01-02", fechaFinStr)
			if err != nil {
				response.BadRequest(w, "fecha_fin inválida, use formato YYYY-MM-DD o RFC3339")
				return
			}
		}
	}

	params := dto.RotacionFiltroParams{FechaInicio: fechaInicio, FechaFin: fechaFin}
	data, err := h.service.GetRotacionDetalle(r.Context(), sucursalID, params)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	response.Success(w, "Rotación de inventario calculada correctamente", data)
}

func (h *InventarioHandler) GetComposicionCategoria(w http.ResponseWriter, r *http.Request) {
	sucursalIDStr := r.URL.Query().Get("id_sucursal")
	if sucursalIDStr == "" {
		ctxSucursal, ok := r.Context().Value("user_sucursal").(uuid.UUID)
		if !ok {
			response.BadRequest(w, "Debe proporcionar id_sucursal")
			return
		}
		sucursalIDStr = ctxSucursal.String()
	}

	sucursalID, err := uuid.Parse(sucursalIDStr)
	if err != nil {
		response.BadRequest(w, "id_sucursal inválido")
		return
	}

	data, err := h.service.GetComposicionCategoria(r.Context(), sucursalID)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.Success(w, "Composición por categoría obtenida correctamente", data)
}

// parseSucursalID obtiene id_sucursal del query param o del contexto JWT.
func parseSucursalID(r *http.Request) (uuid.UUID, bool) {
	s := r.URL.Query().Get("id_sucursal")
	if s == "" {
		id, ok := r.Context().Value("user_sucursal").(uuid.UUID)
		return id, ok
	}
	id, err := uuid.Parse(s)
	return id, err == nil
}

// parseFechaParams lee fecha_inicio y fecha_fin (RFC3339) del query string.
func parseFechaParams(r *http.Request) (dto.RotacionFiltroParams, bool) {
	fi, errI := time.Parse(time.RFC3339, r.URL.Query().Get("fecha_inicio"))
	ff, errF := time.Parse(time.RFC3339, r.URL.Query().Get("fecha_fin"))
	return dto.RotacionFiltroParams{FechaInicio: fi, FechaFin: ff}, errI == nil && errF == nil
}

func (h *InventarioHandler) CapturarSnapshot(w http.ResponseWriter, r *http.Request) {
	sucursalID, ok := parseSucursalID(r)
	if !ok {
		response.BadRequest(w, "Debe proporcionar id_sucursal válido")
		return
	}
	if err := h.service.CapturarSnapshotInventario(r.Context(), sucursalID); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.Success(w, "Snapshot de inventario capturado correctamente", nil)
}

func (h *InventarioHandler) GetValorHistorico(w http.ResponseWriter, r *http.Request) {
	sucursalID, ok := parseSucursalID(r)
	if !ok {
		response.BadRequest(w, "Debe proporcionar id_sucursal válido")
		return
	}
	params, ok := parseFechaParams(r)
	if !ok {
		response.BadRequest(w, "Debe proporcionar fecha_inicio y fecha_fin en formato RFC3339")
		return
	}
	data, err := h.service.GetValorHistorico(r.Context(), sucursalID, params)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	response.Success(w, "Historial de valor del inventario obtenido correctamente", data)
}

func (h *InventarioHandler) GetPerdidas(w http.ResponseWriter, r *http.Request) {
	sucursalID, ok := parseSucursalID(r)
	if !ok {
		response.BadRequest(w, "Debe proporcionar id_sucursal válido")
		return
	}
	params, ok := parseFechaParams(r)
	if !ok {
		response.BadRequest(w, "Debe proporcionar fecha_inicio y fecha_fin en formato RFC3339")
		return
	}
	data, err := h.service.GetPerdidas(r.Context(), sucursalID, params)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	response.Success(w, "Pérdidas por merma y caducidad obtenidas correctamente", data)
}

func (h *InventarioHandler) GetMargenGanancia(w http.ResponseWriter, r *http.Request) {
	sucursalID, ok := parseSucursalID(r)
	if !ok {
		response.BadRequest(w, "Debe proporcionar id_sucursal válido")
		return
	}
	params, ok := parseFechaParams(r)
	if !ok {
		response.BadRequest(w, "Debe proporcionar fecha_inicio y fecha_fin en formato RFC3339")
		return
	}
	data, err := h.service.GetMargenGanancia(r.Context(), sucursalID, params)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	response.Success(w, "Margen de ganancia real obtenido correctamente", data)
}

func (h *InventarioHandler) GetAlertasDetalle(w http.ResponseWriter, r *http.Request) {
	sucursalIDStr := r.URL.Query().Get("id_sucursal")
	if sucursalIDStr == "" {
		ctxSucursal, ok := r.Context().Value("user_sucursal").(uuid.UUID)
		if !ok {
			response.BadRequest(w, "Debe proporcionar id_sucursal")
			return
		}
		sucursalIDStr = ctxSucursal.String()
	}

	sucursalID, err := uuid.Parse(sucursalIDStr)
	if err != nil {
		response.BadRequest(w, "id_sucursal inválido")
		return
	}

	data, err := h.service.GetAlertasStockDetalle(r.Context(), sucursalID)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.Success(w, "Alertas de stock crítico obtenidas correctamente", data)
}
