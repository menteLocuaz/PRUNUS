package helper

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/prunus/pkg/config"
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
	secret := config.Get("JWT_SECRET")
	if secret == "" {
		return "", ErrMissingSecret
	}
	return secret, nil
}

// getJWTExpirationHours obtiene las horas de expiración desde variables de entorno
// Por defecto: 24 horas
func getJWTExpirationHours() int {
	hoursStr := config.Get("JWT_EXPIRATION_HOURS")
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

	// Extraer empresa del tenant (disponible cuando la sucursal se cargó con join)
	idEmpresa := uuid.Nil
	if usuario.Sucursal != nil {
		idEmpresa = usuario.Sucursal.IDEmpresa
	}

	// Crear los claims
	claims := &models.JWTClaims{
		IDUsuario:  usuario.IDUsuario,
		Email:      usuario.Email,
		IDRol:      idRol,
		RolNombre:  rolNombre,
		IDSucursal: usuario.IDSucursal,
		IDEmpresa:  idEmpresa,
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

// RefreshToken genera un nuevo token basado en un token existente.
// Solo se permite refrescar dentro de la misma ventana de tiempo que el token original
// (JWT_EXPIRATION_HOURS). Un token expirado hace más tiempo que ese período es rechazado.
func RefreshToken(tokenString string) (string, int64, error) {
	claims, err := ValidateToken(tokenString)
	if err != nil {
		if !errors.Is(err, ErrExpiredToken) {
			return "", 0, err
		}

		// Token expirado: extraer claims sin validar expiración para comprobar la ventana de gracia.
		secret, secretErr := getJWTSecret()
		if secretErr != nil {
			return "", 0, secretErr
		}
		parser := jwt.NewParser(jwt.WithoutClaimsValidation())
		tok, _, parseErr := parser.ParseUnverified(tokenString, &models.JWTClaims{})
		if parseErr != nil {
			return "", 0, ErrInvalidToken
		}
		// Verificar firma manualmente para evitar usar un token adulterado.
		if _, verifyErr := jwt.ParseWithClaims(tokenString, &models.JWTClaims{},
			func(_ *jwt.Token) (any, error) { return []byte(secret), nil },
			jwt.WithoutClaimsValidation(),
		); verifyErr != nil {
			return "", 0, ErrInvalidToken
		}

		expiredClaims, ok := tok.Claims.(*models.JWTClaims)
		if !ok || expiredClaims.ExpiresAt == nil {
			return "", 0, ErrInvalidToken
		}

		gracePeriod := time.Duration(getJWTExpirationHours()) * time.Hour
		if time.Since(expiredClaims.ExpiresAt.Time) > gracePeriod {
			return "", 0, errors.New("token expirado hace demasiado tiempo, inicie sesión nuevamente")
		}
		claims = expiredClaims
	}

	usuario := &models.Usuario{
		IDUsuario:  claims.IDUsuario,
		Email:      claims.Email,
		IDSucursal: claims.IDSucursal,
		Rol: &models.Rol{
			IDRol:     claims.IDRol,
			RolNombre: claims.RolNombre,
		},
	}

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
