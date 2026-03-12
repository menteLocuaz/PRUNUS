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

// ClienteHandler maneja las solicitudes HTTP relacionadas con clientes
type ClienteHandler struct {
	service *services.ServiceCliente
}

// NewClienteHandler crea una nueva instancia del handler de cliente
func NewClienteHandler(s *services.ServiceCliente) *ClienteHandler {
	return &ClienteHandler{service: s}
}

// GetAll obtiene todos los clientes
func (h *ClienteHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Llama al servicio para obtener todos los clientes
	resp, err := h.service.GetAllClientes(r.Context())
	if err != nil {
		// En caso de error, responde con error 500
		response.InternalServerError(w, err.Error())
		return
	}

	// Responde con éxito y los datos obtenidos
	response.Success(w, "Clientes obtenidos correctamente", resp)
}

// GetByID obtiene un cliente por ID
func (h *ClienteHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Extrae el parámetro "id" de la URL
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		// Si el ID no es válido, responde con error 400
		response.BadRequest(w, "ID inválido")
		return
	}

	// Llama al servicio para obtener el cliente por ID
	resp, err := h.service.GetClienteByID(r.Context(), id)
	if err != nil {
		// Si no se encuentra el cliente, responde con error 404
		response.NotFound(w, "Cliente no encontrado")
		return
	}

	// Responde con éxito y el cliente encontrado
	response.Success(w, "Cliente obtenido correctamente", resp)
}

// Create crea un nuevo cliente
func (h *ClienteHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Decodifica el cuerpo JSON en la estructura de solicitud
	var req dto.ClienteCreateRequest
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

	// Crea un modelo Cliente con los datos recibidos
	cliente := models.Cliente{
		EmpresaCliente: req.EmpresaCliente,
		Nombre:         req.Nombre,
		RUC:            req.RUC,
		Direccion:      req.Direccion,
		Telefono:       req.Telefono,
		Email:          req.Email,
		IDStatus:       req.IDStatus,
	}

	// Llama al servicio para crear the client
	resp, err := h.service.CreateCliente(r.Context(), cliente)
	if err != nil {
		// Si hay error en la creación, responde con error 400
		response.BadRequest(w, err.Error())
		return
	}

	// Responde con código 201 y el cliente creado
	response.Created(w, "Cliente creado correctamente", resp)
}

// Update actualiza un cliente existente
func (h *ClienteHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Extrae y valida el ID de la URL
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		// Si el ID no es válido, responde con error 400
		response.BadRequest(w, "ID inválido")
		return
	}

	// Decodifica el JSON de actualización
	var req dto.ClienteUpdateRequest
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

	// Crea un modelo Cliente con los datos actualizados
	cliente := models.Cliente{
		EmpresaCliente: req.EmpresaCliente,
		Nombre:         req.Nombre,
		RUC:            req.RUC,
		Direccion:      req.Direccion,
		Telefono:       req.Telefono,
		Email:          req.Email,
		IDStatus:       req.IDStatus,
	}

	// Llama al servicio para actualizar el cliente
	resp, err := h.service.UpdateCliente(r.Context(), id, cliente)
	if err != nil {
		// En caso de error interno, responde con error 500
		response.InternalServerError(w, err.Error())
		return
	}

	// Responde con éxito y el cliente actualizado
	response.Success(w, "Cliente actualizado correctamente", resp)
}

// Delete elimina un cliente
func (h *ClienteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Extrae y valida el ID de la URL
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		// Si el ID no es válido, responde con error 400
		response.BadRequest(w, "ID inválido")
		return
	}

	// Llama al servicio para eliminar el cliente
	if err := h.service.DeleteCliente(r.Context(), id); err != nil {
		// Si no se encuentra el cliente, responde con error 404
		response.NotFound(w, "Cliente no encontrado")
		return
	}

	// Responde con código 204 No Content para indicar eliminación exitosa
	w.WriteHeader(http.StatusNoContent)
}
