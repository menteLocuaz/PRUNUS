package transport

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/prunus/pkg/dto"
	"github.com/prunus/pkg/services"
	"github.com/prunus/pkg/utils/response"
	"github.com/prunus/pkg/utils/validator"
	"go.uber.org/zap"
)

type POSHandler struct {
	service *services.ServicePOS
}

func NewPOSHandler(s *services.ServicePOS) *POSHandler {
	return &POSHandler{service: s}
}

// AbrirCajaHandler gestiona el inicio de un turno para un cajero en una estación específica.
// Es una operación crítica que vincula al usuario con un fondo inicial y un periodo contable.
func (h *POSHandler) AbrirCajaHandler(w http.ResponseWriter, r *http.Request) {
	var input dto.AbrirCajaDTO

	// 1. Decodificación y validación inicial del JSON
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.service.Logger().Error("Error al decodificar JSON de apertura", zap.Error(err))
		response.BadRequest(w, "Cuerpo de solicitud inválido: verifique el formato JSON y los tipos de datos.")
		return
	}

	// 2. Validación de campos obligatorios (Struct Tags)
	if err := validator.Validate.Struct(input); err != nil {
		response.ValidationError(w, validator.FormatErrors(err))
		return
	}

	// 3. Recuperar identidad del usuario desde el contexto de seguridad
	idUsuario, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Unauthorized(w, "Sesión no válida o expirada. Por favor, inicie sesión nuevamente.")
		return
	}

	// 4. Ejecución de lógica de negocio (Apertura o Recuperación de sesión activa)
	result, err := h.service.AbrirCaja(r.Context(), input, idUsuario)
	if err != nil {
		// El servicio ahora maneja la idempotencia, por lo que los errores aquí son fallos reales
		response.BadRequest(w, err.Error())
		return
	}

	// 5. Respuesta exitosa (201 Created si es nueva, 200 OK si ya existía)
	// Nota: Por simplicidad devolvemos 201 en ambos casos si el frontend lo espera así
	response.Created(w, "Sesión de caja procesada correctamente", result)
}

// DesmontarCajeroHandler gestiona el cierre forzado o administrativo de una estación
func (h *POSHandler) DesmontarCajeroHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		IDControlEstacion uuid.UUID `json:"id_control_estacion" validate:"required"`
		IDRestaurante     string    `json:"id_restaurante" validate:"required"`
		MotivoDescuadre   string    `json:"motivo_descuadre"`
		AccionInt         int       `json:"accion_int"` // 1: BackOffice, 2: Admin, 3: Descuadre
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.BadRequest(w, "Cuerpo de solicitud mal formado")
		return
	}

	if err := validator.Validate.Struct(input); err != nil {
		response.ValidationError(w, validator.FormatErrors(err))
		return
	}

	idUsuario, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Unauthorized(w, "Acceso denegado: usuario no identificado")
		return
	}

	// Normalización de acción por defecto
	if input.AccionInt == 0 {
		input.AccionInt = 1
	}

	err := h.service.DesmontarCajero(r.Context(), input.IDControlEstacion, idUsuario, input.IDRestaurante, input.MotivoDescuadre, input.AccionInt)
	if err != nil {
		response.InternalServerError(w, "Error interno al procesar el desmontado: "+err.Error())
		return
	}

	response.Success(w, "Cajero desmontado de la estación exitosamente", nil)
}

// ActualizarValoresDeclaradosHandler procesa el arqueo o declaración de valores por forma de pago
func (h *POSHandler) ActualizarValoresDeclaradosHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		IDControlEstacion uuid.UUID `json:"id_control_estacion" validate:"required"`
		IDFormaPago       uuid.UUID `json:"id_forma_pago" validate:"required"`
		Valor             float64   `json:"valor" validate:"min=0"`
		TPEnvID           int       `json:"tpenv_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.BadRequest(w, "Datos de arqueo inválidos")
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

	// -1 representa efectivo por convención interna si no se especifica tpenv_id
	if input.TPEnvID == 0 {
		input.TPEnvID = -1
	}

	err := h.service.ActualizarValoresDeclarados(r.Context(), input.IDControlEstacion, input.IDFormaPago, idUsuario, input.Valor, input.TPEnvID)
	if err != nil {
		response.InternalServerError(w, "No se pudo actualizar el arqueo: "+err.Error())
		return
	}

	response.Success(w, "Declaración de valores registrada", nil)
}

// GetEstadoCajaHandler devuelve la situación actual de una estación de POS
func (h *POSHandler) GetEstadoCajaHandler(w http.ResponseWriter, r *http.Request) {
	idEstacionStr := chi.URLParam(r, "id")
	idEstacion, err := uuid.Parse(idEstacionStr)
	if err != nil {
		response.BadRequest(w, "El ID de estación proporcionado no es un UUID válido")
		return
	}

	result, err := h.service.GetEstadoCaja(r.Context(), idEstacion)
	if err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(w, "No se encontró información para la estación solicitada")
			return
		}
		response.InternalServerError(w, "Error al consultar el estado de la caja")
		return
	}

	response.Success(w, "Estado de caja consultado con éxito", result)
}
