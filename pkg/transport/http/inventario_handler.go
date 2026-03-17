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
	"github.com/prunus/pkg/utils/validator"
)

type InventarioHandler struct {
	service *services.ServiceInventario
}

func NewInventarioHandler(s *services.ServiceInventario) *InventarioHandler {
	return &InventarioHandler{service: s}
}

func (h *InventarioHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.GetAllInventario(r.Context())
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.Success(w, "Inventarios obtenidos correctamente", resp)
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
	var m models.MovimientoInventario
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}
	resp, err := h.service.RegistrarMovimiento(r.Context(), m)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	response.Created(w, "Movimiento de inventario registrado correctamente", resp)
}

func (h *InventarioHandler) GetMovimientos(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID de producto inválido")
		return
	}
	resp, err := h.service.GetMovimientos(r.Context(), id)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.Success(w, "Movimientos obtenidos correctamente", resp)
}
