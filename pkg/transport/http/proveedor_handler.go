package transport

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/prunus/pkg/dto"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/services"
	"github.com/prunus/pkg/utils"
	"github.com/prunus/pkg/utils/response"
	"github.com/prunus/pkg/utils/validator"
)

// ProveedorHandler maneja las solicitudes HTTP relacionadas con proveedores
type ProveedorHandler struct {
	service *services.ServiceProveedor
}

// NewProveedorHandler crea un nuevo handler con el servicio inyectado
func NewProveedorHandler(s *services.ServiceProveedor) *ProveedorHandler {
	return &ProveedorHandler{service: s}
}

// GetAll obtiene una lista paginada de proveedores y responde con JSON
func (h *ProveedorHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	params := utils.ParsePaginationParams(r)
	resp, err := h.service.GetAllProveedores(r.Context(), params)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, "Proveedores obtenidos correctamente", resp)
}

// GetByID obtiene un proveedor por su ID
func (h *ProveedorHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}

	resp, err := h.service.GetProveedorByID(r.Context(), id)
	if err != nil {
		response.NotFound(w, "Proveedor no encontrado")
		return
	}

	response.Success(w, "Proveedor obtenido correctamente", resp)
}

// Create crea un nuevo proveedor a partir del JSON recibido
func (h *ProveedorHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.ProveedorCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		response.ValidationError(w, validator.FormatErrors(err))
		return
	}

	proveedor := models.Proveedor{
		RazonSocial:    req.RazonSocial,
		NitRut:         req.NitRut,
		ContactoNombre: req.ContactoNombre,
		Telefono:       req.Telefono,
		Direccion:      req.Direccion,
		Email:          req.Email,
		IDStatus:       req.IDStatus,
	}

	resp, err := h.service.CreateProveedor(r.Context(), proveedor)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Created(w, "Proveedor creado correctamente", resp)
}

// Update actualiza un proveedor existente identificado por ID
func (h *ProveedorHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}

	var req dto.ProveedorUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		response.ValidationError(w, validator.FormatErrors(err))
		return
	}

	proveedor := models.Proveedor{
		RazonSocial:    req.RazonSocial,
		NitRut:         req.NitRut,
		ContactoNombre: req.ContactoNombre,
		Telefono:       req.Telefono,
		Direccion:      req.Direccion,
		Email:          req.Email,
		IDStatus:       req.IDStatus,
	}

	resp, err := h.service.UpdateProveedor(r.Context(), id, proveedor)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, "Proveedor actualizado correctamente", resp)
}

// Delete elimina un proveedor por su ID
func (h *ProveedorHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}

	if err := h.service.DeleteProveedor(r.Context(), id); err != nil {
		response.NotFound(w, "Proveedor no encontrado")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
