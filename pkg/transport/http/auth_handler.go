package transport

import (
	"encoding/json"
	"net/http"

	"github.com/prunus/pkg/helper"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/services"
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
		http.Error(w, "Formato de petición inválido", http.StatusBadRequest)
		return
	}

	// Autenticar al usuario
	usuario, err := h.usuarioService.AuthenticateUsuario(loginReq.Email, loginReq.Password)
	if err != nil {
		// Retornar 401 para errores de autenticación
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Generar JWT token
	token, expiresAt, err := helper.GenerateToken(usuario)
	if err != nil {
		http.Error(w, "Error al generar token", http.StatusInternalServerError)
		return
	}

	// Preparar respuesta
	response := models.LoginResponse{
		Token:     token,
		Usuario:   usuario,
		ExpiresAt: expiresAt,
	}

	// Retornar respuesta
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
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

	// Respuesta exitosa
	response := models.LogoutResponse{
		Message: "Sesión cerrada exitosamente",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetMe retorna la información del usuario autenticado actual
// Este endpoint requiere autenticación (middleware)
func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	// Obtener los claims del contexto (agregados por el middleware de auth)
	claims, ok := r.Context().Value("claims").(*models.JWTClaims)
	if !ok {
		http.Error(w, "No se pudo obtener información del usuario", http.StatusInternalServerError)
		return
	}

	// Obtener información completa del usuario desde la BD
	usuario, err := h.usuarioService.GetUsuarioByID(claims.IDUsuario)
	if err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	// Limpiar password antes de enviar
	usuario.UsuPassword = ""

	// Retornar usuario
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(usuario)
}

// RefreshToken genera un nuevo token basado en el token actual
// Útil para renovar la sesión antes de que expire
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// Obtener el token del header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Token no proporcionado", http.StatusUnauthorized)
		return
	}

	// Extraer token del header
	token, err := helper.ExtractTokenFromHeader(authHeader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Refrescar el token
	newToken, expiresAt, err := helper.RefreshToken(token)
	if err != nil {
		http.Error(w, "Error al refrescar token", http.StatusUnauthorized)
		return
	}

	// Preparar respuesta
	response := map[string]interface{}{
		"token":      newToken,
		"expires_at": expiresAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
