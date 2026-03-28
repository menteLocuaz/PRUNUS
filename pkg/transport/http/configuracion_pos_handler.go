// Package http contiene los handlers HTTP para la capa de transporte,
// responsables de recibir solicitudes, validar parámetros y delegar
// la lógica de negocio al servicio correspondiente.
package transport

import (
	"net/http"
	"strconv"

	"github.com/prunus/pkg/services"
	"github.com/prunus/pkg/utils/response"

	"github.com/go-chi/chi/v5"
)

// ConfiguracionHandler expone los endpoints HTTP relacionados con la
// configuración de impresión (canales, impresoras y puertos).
// Depende de ServiceConfiguracion para ejecutar la lógica de negocio.
type ConfiguracionHandler struct {
	service *services.ServiceConfiguracion
}

// NewConfiguracionHandler crea e inicializa una nueva instancia de ConfiguracionHandler
// inyectando la dependencia del servicio de configuración.
//
// Parámetros:
//   - s: puntero al servicio de configuración.
//
// Retorna:
//   - *ConfiguracionHandler: instancia lista para registrar sus rutas.
func NewConfiguracionHandler(s *services.ServiceConfiguracion) *ConfiguracionHandler {
	return &ConfiguracionHandler{service: s}
}

// GetCanales maneja la solicitud GET para obtener los canales de impresión
// activos asociados a una cadena específica.
//
// Ruta esperada: GET /canales/{chainId}
//
// Parámetros de URL:
//   - chainId (int): identificador único de la cadena.
//
// Respuestas:
//   - 200 OK:          lista de canales obtenida exitosamente.
//   - 400 Bad Request: el parámetro chainId es inválido o igual a cero.
//   - 500 Internal:    error al consultar los canales en el servicio.
func (h *ConfiguracionHandler) GetCanales(w http.ResponseWriter, r *http.Request) {
	// Extraer y convertir el parámetro de ruta "chainId" a entero.
	chainID, _ := strconv.Atoi(chi.URLParam(r, "chainId"))

	// Validar que el ID sea un valor positivo válido.
	if chainID == 0 {
		response.Error(w, http.StatusBadRequest, "ID de cadena inválido")
		return
	}

	// Delegar la consulta al servicio, propagando el contexto de la solicitud.
	data, err := h.service.ListarCanales(r.Context(), chainID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Error al listar canales: "+err.Error())
		return
	}

	response.Success(w, "Canales obtenidos", data)
}

// GetImpresoras maneja la solicitud GET para obtener las impresoras activas
// asociadas a un restaurante específico.
//
// Ruta esperada: GET /impresoras/{restId}
//
// Parámetros de URL:
//   - restId (int): identificador único del restaurante.
//
// Respuestas:
//   - 200 OK:          lista de impresoras obtenida exitosamente.
//   - 400 Bad Request: el parámetro restId es inválido o igual a cero.
//   - 500 Internal:    error al consultar las impresoras en el servicio.
func (h *ConfiguracionHandler) GetImpresoras(w http.ResponseWriter, r *http.Request) {
	// Extraer y convertir el parámetro de ruta "restId" a entero.
	restID, _ := strconv.Atoi(chi.URLParam(r, "restId"))

	// Validar que el ID sea un valor positivo válido.
	if restID == 0 {
		response.Error(w, http.StatusBadRequest, "ID de restaurante inválido")
		return
	}

	// Delegar la consulta al servicio, propagando el contexto de la solicitud.
	data, err := h.service.ListarImpresoras(r.Context(), restID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Error al listar impresoras: "+err.Error())
		return
	}

	response.Success(w, "Impresoras obtenidas", data)
}

// GetPuertos maneja la solicitud GET para obtener todos los puertos de
// comunicación activos disponibles en el sistema.
//
// Ruta esperada: GET /puertos
//
// No requiere parámetros de entrada.
//
// Respuestas:
//   - 200 OK:       lista de puertos obtenida exitosamente.
//   - 500 Internal: error al consultar los puertos en el servicio.
func (h *ConfiguracionHandler) GetPuertos(w http.ResponseWriter, r *http.Request) {
	// Delegar la consulta al servicio, propagando el contexto de la solicitud.
	data, err := h.service.ListarPuertos(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Error al listar puertos: "+err.Error())
		return
	}

	response.Success(w, "Puertos obtenidos", data)
}
