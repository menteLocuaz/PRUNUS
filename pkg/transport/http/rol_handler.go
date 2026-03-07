package transport

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/prunus/pkg/dto"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/services"
	"github.com/prunus/pkg/utils/response" // Importa el paquete response para respuestas estandarizadas
)

// RolHandler maneja las peticiones HTTP relacionadas con roles
type RolHandler struct {
	service *services.ServiceRol
}

// NewRolHandler crea una nueva instancia del handler de rol
func NewRolHandler(s *services.ServiceRol) *RolHandler {
	return &RolHandler{service: s}
}

// GetAll maneja la petición GET para obtener todos los roles
func (h *RolHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Llama al servicio para obtener todos los roles
	roles, err := h.service.GetAllRoles()
	if err != nil {
		// En caso de error, responde con error 500
		response.InternalServerError(w, err.Error())
		return
	}

	// Responde con éxito y los datos obtenidos
	response.Success(w, "Roles obtenidos correctamente", roles)
}

// GetByID maneja la petición GET para obtener un rol por ID
func (h *RolHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID del parámetro de la URL
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		// Si el ID no es válido, responde con error 400
		response.BadRequest(w, "ID inválido")
		return
	}

	// Llama al servicio para obtener el rol por ID
	rol, err := h.service.GetRolByID(uint(id))
	if err != nil {
		// Si no se encuentra el rol, responde con error 404
		response.NotFound(w, err.Error())
		return
	}

	// Responde con éxito y el rol encontrado
	response.Success(w, "Rol obtenido correctamente", rol)
}

// Create maneja la petición POST para crear un nuevo rol
func (h *RolHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Decodificar el cuerpo de la petición
	var req dto.RolCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Si el JSON es inválido, responde con error 400
		response.BadRequest(w, "JSON inválido")
		return
	}

	// Convertir DTO a modelo
	rol := models.Rol{
		RolNombre:  req.RolNombre,
		IDSucursal: req.IDSucursal,
		Estado:     req.Estado,
	}

	// Crear el rol usando el servicio
	resp, err := h.service.CreateRol(rol)
	if err != nil {
		// Si hay error en la creación, responde con error 400
		response.BadRequest(w, err.Error())
		return
	}

	// Responde con código 201 y el rol creado
	response.Created(w, "Rol creado correctamente", resp)
}

// Update maneja la petición PUT para actualizar un rol existente
func (h *RolHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID del parámetro de la URL
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		// Si el ID no es válido, responde con error 400
		response.BadRequest(w, "ID inválido")
		return
	}

	// Decodificar el cuerpo de la petición
	var req dto.RolUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Si el JSON es inválido, responde con error 400
		response.BadRequest(w, "JSON inválido")
		return
	}

	// Convertir DTO a modelo
	rol := models.Rol{
		RolNombre:  req.RolNombre,
		IDSucursal: req.IDSucursal,
		Estado:     req.Estado,
	}

	// Actualizar el rol usando el servicio
	resp, err := h.service.UpdateRol(uint(id), rol)
	if err != nil {
		// Si hay error en la actualización, responde con error 400
		response.BadRequest(w, err.Error())
		return
	}

	// Responde con éxito y el rol actualizado
	response.Success(w, "Rol actualizado correctamente", resp)
}

// Delete maneja la petición DELETE para eliminar un rol (soft delete)
func (h *RolHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID del parámetro de la URL
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		// Si el ID no es válido, responde con error 400
		response.BadRequest(w, "ID inválido")
		return
	}

	// Eliminar el rol usando el servicio
	if err := h.service.DeleteRol(uint(id)); err != nil {
		// Si hay error en la eliminación, responde con error 400
		response.BadRequest(w, err.Error())
		return
	}

	// Responde con código 204 No Content para indicar eliminación exitosa
	w.WriteHeader(http.StatusNoContent)
}
