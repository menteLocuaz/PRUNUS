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

type EmpresaHandler struct {
	service *services.ServiceEmpresa
}

func NewEmpresaHandler(s *services.ServiceEmpresa) *EmpresaHandler {
	return &EmpresaHandler{
		service: s,
	}

}

// GetAll obtiene todas las empresas
func (h *EmpresaHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.GetAllEmpresa(r.Context())
	if err != nil {
		response.InternalServerError(w, "Error al obtener las empresas")
		return
	}

	response.Success(w, "Empresas obtenidas correctamente", resp)
}

// GetByID obtiene una empresa por ID
func (h *EmpresaHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}

	resp, err := h.service.GetByIDEmpresa(r.Context(), id)
	if err != nil {
		response.NotFound(w, "Empresa no encontrada")
		return
	}

	response.Success(w, "Empresa obtenida correctamente", resp)
}

// Create crea una nueva empresa
func (h *EmpresaHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.EmpresaCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}

	// Validar la estructura
	if err := validator.Validate.Struct(req); err != nil {
		response.ValidationError(w, validator.FormatErrors(err))
		return
	}

	empresa := models.Empresa{
		Nombre:   req.Nombre,
		RUT:      req.RUT,
		IDStatus: req.IDStatus,
	}

	resp, err := h.service.CrearEmpresa(r.Context(), empresa)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Created(w, "Empresa creada exitosamente", resp)
}

// Update actualiza una empresa existente
func (h *EmpresaHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}

	var req dto.EmpresaUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}

	// Validar la estructura
	if err := validator.Validate.Struct(req); err != nil {
		response.ValidationError(w, validator.FormatErrors(err))
		return
	}

	empresa := models.Empresa{
		Nombre:   req.Nombre,
		RUT:      req.RUT,
		IDStatus: req.IDStatus,
	}

	resp, err := h.service.UpdateEmpresa(r.Context(), id, empresa)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, "Empresa actualizada correctamente", resp)
}

// Delete elimina una empresa
func (h *EmpresaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}

	if err := h.service.ElimminarEmpresa(r.Context(), id); err != nil {
		response.NotFound(w, "Empresa no encontrada")
		return
	}

	// Responde con código 204 No Content para indicar eliminación exitosa
	w.WriteHeader(http.StatusNoContent)
}
