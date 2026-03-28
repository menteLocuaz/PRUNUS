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

type CompraHandler struct {
	service *services.ServiceCompra
}

func NewCompraHandler(s *services.ServiceCompra) *CompraHandler {
	return &CompraHandler{service: s}
}

func (h *CompraHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.GetAllOrdenes(r.Context())
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.Success(w, "Órdenes de compra obtenidas correctamente", resp)
}

func (h *CompraHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID de orden de compra inválido")
		return
	}

	resp, err := h.service.GetOrdenByID(r.Context(), id)
	if err != nil {
		response.NotFound(w, err.Error())
		return
	}
	response.Success(w, "Orden de compra obtenida correctamente", resp)
}

func (h *CompraHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.OrdenCompraCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		response.ValidationError(w, validator.FormatErrors(err))
		return
	}

	userID, _ := r.Context().Value("user_id").(uuid.UUID)

	orden := models.OrdenCompra{
		NumeroOrden:   req.NumeroOrden,
		IDProveedor:   req.IDProveedor,
		IDSucursal:    req.IDSucursal,
		IDUsuario:     userID,
		IDMoneda:      req.IDMoneda,
		IDStatus:      req.IDStatus,
		Observaciones: req.Observaciones,
	}

	var subtotal, impuesto float64
	for _, d := range req.Detalles {
		detTotal := d.CantidadPedida * d.PrecioUnitario
		detImp := detTotal * (d.Impuesto / 100)

		orden.Detalles = append(orden.Detalles, &models.DetalleOrdenCompra{
			IDProducto:     d.IDProducto,
			CantidadPedida: d.CantidadPedida,
			PrecioUnitario: d.PrecioUnitario,
			Impuesto:       detImp,
			Total:          detTotal + detImp,
		})

		subtotal += detTotal
		impuesto += detImp
	}

	orden.Subtotal = subtotal
	orden.Impuesto = impuesto
	orden.Total = subtotal + impuesto

	resp, err := h.service.CreateOrden(r.Context(), &orden)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	response.Created(w, "Orden de compra creada correctamente", resp)
}

func (h *CompraHandler) Recepcion(w http.ResponseWriter, r *http.Request) {
	var req dto.RecepcionCompraRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		response.ValidationError(w, validator.FormatErrors(err))
		return
	}

	userID, _ := r.Context().Value("user_id").(uuid.UUID)

	var items []*models.DetalleOrdenCompra
	for _, it := range req.Items {
		items = append(items, &models.DetalleOrdenCompra{
			IDDetalleCompra:  it.IDDetalleCompra,
			IDProducto:       it.IDProducto,
			CantidadRecibida: it.CantidadRecibida,
		})
	}

	err := h.service.ProcesarRecepcion(r.Context(), req.IDOrdenCompra, req.IDStatus, items, userID)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, "Mercancía recibida e inventario abastecido correctamente", nil)
}
