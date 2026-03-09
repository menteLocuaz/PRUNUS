package helper

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
)

var (
	// ErrInvalidToken cuando el token es inválido
	ErrInvalidToken = errors.New("token inválido")
	// ErrExpiredToken cuando el token ha expirado
	ErrExpiredToken = errors.New("token expirado")
	// ErrMissingSecret cuando no se encuentra JWT_SECRET
	ErrMissingSecret = errors.New("JWT_SECRET no configurado")
)

// getJWTSecret obtiene el secret desde las variables de entorno
func getJWTSecret() (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", ErrMissingSecret
	}
	return secret, nil
}

// getJWTExpirationHours obtiene las horas de expiración desde variables de entorno
// Por defecto: 24 horas
func getJWTExpirationHours() int {
	hoursStr := os.Getenv("JWT_EXPIRATION_HOURS")
	if hoursStr == "" {
		return 24 // Default: 24 horas
	}

	hours, err := strconv.Atoi(hoursStr)
	if err != nil {
		return 24 // Fallback a 24 horas si hay error
	}

	return hours
}

// GenerateToken genera un JWT token para un usuario autenticado
func GenerateToken(usuario *models.Usuario) (string, int64, error) {
	secret, err := getJWTSecret()
	if err != nil {
		return "", 0, err
	}

	expirationHours := getJWTExpirationHours()
	expirationTime := time.Now().Add(time.Duration(expirationHours) * time.Hour)

	// Extraer nombre del rol (puede ser nil si no se cargó)
	rolNombre := ""
	idRol := uuid.Nil
	if usuario.Rol != nil {
		rolNombre = usuario.Rol.RolNombre
		idRol = usuario.Rol.IDRol
	}

	// Crear los claims
	claims := &models.JWTClaims{
		IDUsuario:  usuario.IDUsuario,
		Email:      usuario.Email,
		IDRol:      idRol,
		RolNombre:  rolNombre,
		IDSucursal: usuario.IDSucursal,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "prunus-api",
			Subject:   usuario.IDUsuario.String(),
		},
	}

	// Crear el token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar el token
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", 0, fmt.Errorf("error al firmar token: %w", err)
	}

	return tokenString, expirationTime.Unix(), nil
}

// ValidateToken valida un JWT token y retorna los claims
func ValidateToken(tokenString string) (*models.JWTClaims, error) {
	secret, err := getJWTSecret()
	if err != nil {
		return nil, err
	}

	// Parsear y validar el token
	token, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verificar que el método de firma sea el esperado
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		// Verificar si el error es por expiración
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	// Extraer los claims
	claims, ok := token.Claims.(*models.JWTClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// RefreshToken genera un nuevo token basado en un token existente válido
// Útil para renovar la sesión del usuario antes de que expire
func RefreshToken(tokenString string) (string, int64, error) {
	// Validar el token actual
	claims, err := ValidateToken(tokenString)
	if err != nil {
		// Si está expirado, permitir refrescar si no ha pasado mucho tiempo
		if !errors.Is(err, ErrExpiredToken) {
			return "", 0, err
		}
	}

	// Crear un usuario temporal con los datos del token
	usuario := &models.Usuario{
		IDUsuario:  claims.IDUsuario,
		Email:      claims.Email,
		IDSucursal: claims.IDSucursal,
		Rol: &models.Rol{
			IDRol:     claims.IDRol,
			RolNombre: claims.RolNombre,
		},
	}

	// Generar nuevo token
	return GenerateToken(usuario)
}

// ExtractTokenFromHeader extrae el token del header Authorization
// Espera formato: "Bearer <token>"
func ExtractTokenFromHeader(authHeader string) (string, error) {
	const bearerPrefix = "Bearer "

	if authHeader == "" {
		return "", errors.New("header Authorization vacío")
	}

	if len(authHeader) < len(bearerPrefix) {
		return "", errors.New("formato de Authorization inválido")
	}

	if authHeader[:len(bearerPrefix)] != bearerPrefix {
		return "", errors.New("debe usar formato 'Bearer <token>'")
	}

	token := authHeader[len(bearerPrefix):]
	if token == "" {
		return "", errors.New("token vacío")
	}

	return token, nil
}
