package transport

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/prunus/pkg/dto"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/services"
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

// GetAll obtiene todos los proveedores y responde con JSON
func (h *ProveedorHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Llama al servicio para obtener todos los proveedores
	resp, err := h.service.GetAllProveedores()
	if err != nil {
		// En caso de error, responde con error 500
		response.InternalServerError(w, err.Error())
		return
	}

	// Responde con éxito y los datos obtenidos
	response.Success(w, "Proveedores obtenidos correctamente", resp)
}

// GetByID obtiene un proveedor por su ID
func (h *ProveedorHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Extrae el parámetro "id" de la URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// Si el ID no es válido, responde con error 400
		response.BadRequest(w, "ID inválido")
		return
	}

	// Llama al servicio para obtener el proveedor por ID
	resp, err := h.service.GetProveedorByID(uint(id))
	if err != nil {
		// Si no se encuentra el proveedor, responde con error 404
		response.NotFound(w, "Proveedor no encontrado")
		return
	}

	// Responde con éxito y el proveedor encontrado
	response.Success(w, "Proveedor obtenido correctamente", resp)
}

// Create crea un nuevo proveedor a partir del JSON recibido
func (h *ProveedorHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Decodifica el cuerpo JSON en la estructura de solicitud
	var req dto.ProveedorCreateRequest
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

	// Crea un modelo Proveedor con los datos recibidos
	proveedor := models.Proveedor{
		Nombre:     req.Nombre,
		RUC:        req.RUC,
		Telefono:   req.Telefono,
		Direccion:  req.Direccion,
		Email:      req.Email,
		Estado:     req.Estado,
		IDSucursal: req.IDSucursal,
		IDEmpresa:  req.IDEmpresa,
	}

	// Llama al servicio para crear el proveedor
	resp, err := h.service.CreateProveedor(proveedor)
	if err != nil {
		// Si hay error en la creación, responde con error 400
		response.BadRequest(w, err.Error())
		return
	}

	// Responde con código 201 y el proveedor creado
	response.Created(w, "Proveedor creado correctamente", resp)
}

// Update actualiza un proveedor existente identificado por ID
func (h *ProveedorHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Extrae y valida el ID de la URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}

	// Decodifica el JSON de actualización
	var req dto.ProveedorUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON inválido")
		return
	}

	// Validar la estructura
	if err := validator.Validate.Struct(req); err != nil {
		response.ValidationError(w, validator.FormatErrors(err))
		return
	}

	// Crea un modelo Proveedor con los datos actualizados
	proveedor := models.Proveedor{
		Nombre:     req.Nombre,
		RUC:        req.RUC,
		Telefono:   req.Telefono,
		Direccion:  req.Direccion,
		Email:      req.Email,
		Estado:     req.Estado,
		IDSucursal: req.IDSucursal,
		IDEmpresa:  req.IDEmpresa,
	}

	// Llama al servicio para actualizar the proveedor
	resp, err := h.service.UpdateProveedor(uint(id), proveedor)
	if err != nil {
		// En caso de error interno, responde con error 500
		response.InternalServerError(w, err.Error())
		return
	}

	// Responde con éxito y el proveedor actualizado
	response.Success(w, "Proveedor actualizado correctamente", resp)
}

// Delete elimina un proveedor por su ID
func (h *ProveedorHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Extrae y valida el ID de la URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, "ID inválido")
		return
	}

	// Llama al servicio para eliminar el proveedor
	if err := h.service.DeleteProveedor(uint(id)); err != nil {
		// Si no se encuentra el proveedor, responde con error 404
		response.NotFound(w, "Proveedor no encontrado")
		return
	}

	// Responde con código 204 No Content para indicar eliminación exitosa
	w.WriteHeader(http.StatusNoContent)
}
