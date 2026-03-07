package transport

import (
	"encoding/json"
	"net/http"

	"github.com/prunus/pkg/helper"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/services"
	"github.com/prunus/pkg/utils/response" // Importa el paquete response para respuestas estandarizadas
)

// AuthHandler maneja las peticiones relacionadas con autenticación
type AuthHandler struct {
	usuarioService *services.ServiceUsuario
}

// NewAuthHandler crea una nueva instancia del handler de autenticación
func NewAuthHandler(usuarioService *services.ServiceUsuario) *AuthHandler {
	return &AuthHandler{
		usuarioService: usuarioService,
	}
}

// Login maneja la petición de inicio de sesión
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Decodificar el body de la petición
	var loginReq models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		// Si el formato es inválido, responde con error 400
		response.BadRequest(w, "Formato de petición inválido")
		return
	}

	// Autenticar al usuario
	usuario, err := h.usuarioService.AuthenticateUsuario(loginReq.Email, loginReq.Password)
	if err != nil {
		// Retornar 401 para errores de autenticación
		response.Unauthorized(w, err.Error())
		return
	}

	// Generar JWT token
	token, expiresAt, err := helper.GenerateToken(usuario)
	if err != nil {
		// Error al generar token, responde con error 500
		response.InternalServerError(w, "Error al generar token")
		return
	}

	// Preparar respuesta
	responseData := models.LoginResponse{
		Token:     token,
		Usuario:   usuario,
		ExpiresAt: expiresAt,
	}

	// Retornar respuesta exitosa con token
	response.Success(w, "Inicio de sesión exitoso", responseData)
}

// Logout maneja la petición de cierre de sesión
// En un sistema JWT stateless, el logout se maneja en el cliente eliminando el token
// Este endpoint puede usarse para logging o futura implementación de blacklist
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Obtener el token del header (opcional, para logging)
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		// Aquí podrías agregar el token a una blacklist si implementas esa funcionalidad
		// Por ahora, solo confirmamos el logout
	}

	// Preparar respuesta de cierre de sesión
	logoutResponse := models.LogoutResponse{
		Message: "Sesión cerrada exitosamente",
	}

	// Retornar respuesta exitosa
	response.Success(w, "Cierre de sesión exitoso", logoutResponse)
}

// GetMe retorna la información del usuario autenticado actual
// Este endpoint requiere autenticación (middleware)
func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	// Obtener los claims del contexto (agregados por el middleware de auth)
	claims, ok := r.Context().Value("claims").(*models.JWTClaims)
	if !ok {
		// Si no se pueden obtener los claims, responde con error 500
		response.InternalServerError(w, "No se pudo obtener información del usuario")
		return
	}

	// Obtener información completa del usuario desde la BD
	usuario, err := h.usuarioService.GetUsuarioByID(claims.IDUsuario)
	if err != nil {
		// Si no se encuentra el usuario, responde con error 404
		response.NotFound(w, "Usuario no encontrado")
		return
	}

	// Limpiar password antes de enviar
	usuario.UsuPassword = ""

	// Retornar usuario autenticado
	response.Success(w, "Información de usuario obtenida correctamente", usuario)
}

// RefreshToken genera un nuevo token basado en el token actual
// Útil para renovar la sesión antes de que expire
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// Obtener el token del header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		// Si no se proporciona token, responde con error 401
		response.Unauthorized(w, "Token no proporcionado")
		return
	}

	// Extraer token del header
	token, err := helper.ExtractTokenFromHeader(authHeader)
	if err != nil {
		// Si hay error al extraer el token, responde con error 401
		response.Unauthorized(w, err.Error())
		return
	}

	// Refrescar el token
	newToken, expiresAt, err := helper.RefreshToken(token)
	if err != nil {
		// Si hay error al refrescar el token, responde con error 401
		response.Unauthorized(w, "Error al refrescar token")
		return
	}

	// Preparar respuesta con nuevo token
	tokenResponse := map[string]interface{}{
		"token":      newToken,
		"expires_at": expiresAt,
	}

	// Retornar nuevo token
	response.Success(w, "Token refrescado correctamente", tokenResponse)
}
