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

// SucursalHandler maneja las solicitudes HTTP relacionadas con sucursales
type SucursalHandler struct {
	service *services.ServiceSucursal
}

// NewSucursalHandler crea una nueva instancia del handler de sucursal
func NewSucursalHandler(s *services.ServiceSucursal) *SucursalHandler {
	return &SucursalHandler{
		service: s,
	}
}

// GetAll obtiene todas las sucursales
func (h *SucursalHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Llama al servicio para obtener todas las sucursales
	resp, err := h.service.GetAllSucursales()
	if err != nil {
		// En caso de error, responde con error 500
		response.InternalServerError(w, err.Error())
		return
	}

	// Responde con éxito y los datos obtenidos
	response.Success(w, "Sucursales obtenidas correctamente", resp)
}

// GetByID obtiene una sucursal por ID
func (h *SucursalHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Extrae el parámetro "id" de la URL
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		// Si el ID no es válido, responde con error 400
		response.BadRequest(w, "ID inválido")
		return
	}

	// Llama al servicio para obtener la sucursal por ID
	resp, err := h.service.GetSucursalByID(id)
	if err != nil {
		// Si no se encuentra la sucursal, responde con error 404
		response.NotFound(w, "Sucursal no encontrada")
		return
	}

	// Responde con éxito y la sucursal encontrada
	response.Success(w, "Sucursal obtenida correctamente", resp)
}

// Create crea una nueva sucursal
func (h *SucursalHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Decodifica el cuerpo JSON en la estructura de solicitud
	var req dto.SucursalCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Si el JSON es inválido, responde con error 400
		response.BadRequest(w, "JSON inválido")
		return
	}

	// Validar la estructura
	if err := validator.Validate.Struct(req); err != nil {
		response.ValidationError(w, validator.FormatErrors(err))
		return
	}

	// Crea un modelo Sucursal con los datos recibidos
	sucursal := models.Sucursal{
		IDEmpresa:      req.IDEmpresa,
		NombreSucursal: req.NombreSucursal,
		IDStatus:       req.IDStatus,
	}

	// Llama al servicio para crear la sucursal
	resp, err := h.service.CreateSucursal(sucursal)
	if err != nil {
		// Si hay error en la creación, responde con error 400
		response.BadRequest(w, err.Error())
		return
	}

	// Responde con código 201 y la sucursal creada
	response.Created(w, "Sucursal creada correctamente", resp)
}

// Update actualiza una sucursal existente
func (h *SucursalHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Extrae y valida el ID de la URL
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		// Si el ID no es válido, responde con error 400
		response.BadRequest(w, "ID inválido")
		return
	}

	// Decodifica el JSON de actualización
	var req dto.SucursalUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Si el JSON es inválido, responde con error 400
		response.BadRequest(w, "JSON inválido")
		return
	}

	// Validar la estructura
	if err := validator.Validate.Struct(req); err != nil {
		response.ValidationError(w, validator.FormatErrors(err))
		return
	}

	// Crea un modelo Sucursal con los datos actualizados
	sucursal := models.Sucursal{
		IDEmpresa:      req.IDEmpresa,
		NombreSucursal: req.NombreSucursal,
		IDStatus:       req.IDStatus,
	}

	// Llama al servicio para actualizar la sucursal
	resp, err := h.service.UpdateSucursal(id, sucursal)
	if err != nil {
		// En caso de error interno, responde con error 500
		response.InternalServerError(w, err.Error())
		return
	}

	// Responde con éxito y la sucursal actualizada
	response.Success(w, "Sucursal actualizada correctamente", resp)
}

// Delete elimina una sucursal
func (h *SucursalHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Extrae y valida el ID de la URL
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		// Si el ID no es válido, responde con error 400
		response.BadRequest(w, "ID inválido")
		return
	}

	// Llama al servicio para eliminar la sucursal
	if err := h.service.DeleteSucursal(id); err != nil {
		// Si no se encuentra la sucursal, responde con error 404
		response.NotFound(w, "Sucursal no encontrada")
		return
	}

	// Responde con código 204 No Content para indicar eliminación exitosa
	w.WriteHeader(http.StatusNoContent)
}
