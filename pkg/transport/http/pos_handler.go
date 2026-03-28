package transport

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/prunus/pkg/dto"
	"github.com/prunus/pkg/services"
	"github.com/prunus/pkg/utils/response"
	"github.com/prunus/pkg/utils/validator"
)

type POSHandler struct {
	service *services.ServicePOS
}

func NewPOSHandler(s *services.ServicePOS) *POSHandler {
	return &POSHandler{service: s}
}

// AbrirCajaHandler maneja la petición de apertura de caja
func (h *POSHandler) AbrirCajaHandler(w http.ResponseWriter, r *http.Request) {
	var input dto.AbrirCajaDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.BadRequest(w, "Formato de petición inválido")
		return
	}

	if err := validator.Validate.Struct(input); err != nil {
		response.ValidationError(w, validator.FormatErrors(err))
		return
	}

	idUsuario, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Unauthorized(w, "Usuario no autenticado o ID inválido")
		return
	}

	result, err := h.service.AbrirCaja(r.Context(), input, idUsuario)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Created(w, "Caja abierta exitosamente", result)
}

// DesmontarCajeroHandler maneja la petición de desmontar/cerrar sesión de un cajero (Back Office)
func (h *POSHandler) DesmontarCajeroHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		IDControlEstacion uuid.UUID `json:"id_control_estacion" validate:"required"`
		IDRestaurante     string    `json:"id_restaurante" validate:"required"`
		MotivoDescuadre   string    `json:"motivo_descuadre"`
		AccionInt         int       `json:"accion_int"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.BadRequest(w, "Formato de petición inválido")
		return
	}

	if err := validator.Validate.Struct(input); err != nil {
		response.ValidationError(w, validator.FormatErrors(err))
		return
	}

	idUsuario, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Unauthorized(w, "Usuario no autenticado")
		return
	}

	// Valor por defecto si no se envía
	if input.AccionInt == 0 {
		input.AccionInt = 1
	}

	err := h.service.DesmontarCajero(r.Context(), input.IDControlEstacion, idUsuario, input.IDRestaurante, input.MotivoDescuadre, input.AccionInt)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, "Cajero desmontado correctamente", nil)
}

// ActualizarValoresDeclaradosHandler maneja la actualización del arqueo por forma de pago
func (h *POSHandler) ActualizarValoresDeclaradosHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		IDControlEstacion uuid.UUID `json:"id_control_estacion" validate:"required"`
		IDFormaPago       uuid.UUID `json:"id_forma_pago" validate:"required"`
		Valor             float64   `json:"valor"`
		TPEnvID           int       `json:"tpenv_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.BadRequest(w, "Formato de petición inválido")
		return
	}

	if err := validator.Validate.Struct(input); err != nil {
		response.ValidationError(w, validator.FormatErrors(err))
		return
	}

	idUsuario, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Unauthorized(w, "Usuario no autenticado")
		return
	}

	// Si no se envía tpenv_id, se asume efectivo (-1)
	if input.TPEnvID == 0 {
		input.TPEnvID = -1
	}

	err := h.service.ActualizarValoresDeclarados(r.Context(), input.IDControlEstacion, input.IDFormaPago, idUsuario, input.Valor, input.TPEnvID)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, "Valores declarados actualizados correctamente", nil)
}

// GetEstadoCajaHandler obtiene el estado actual de una estación
func (h *POSHandler) GetEstadoCajaHandler(w http.ResponseWriter, r *http.Request) {
	idEstacionStr := chi.URLParam(r, "id")
	idEstacion, err := uuid.Parse(idEstacionStr)
	if err != nil {
		response.BadRequest(w, "ID de estación inválido")
		return
	}

	result, err := h.service.GetEstadoCaja(r.Context(), idEstacion)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.Success(w, "Estado de caja obtenido correctamente", result)
}
