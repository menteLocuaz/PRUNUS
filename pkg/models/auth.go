package models

import "github.com/golang-jwt/jwt/v5"

// LoginRequest representa la petición de login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse representa la respuesta exitosa de login
type LoginResponse struct {
	Token     string   `json:"token"`
	Usuario   *Usuario `json:"usuario"`
	ExpiresAt int64    `json:"expires_at"`
}

// JWTClaims representa los claims personalizados del JWT
type JWTClaims struct {
	IDUsuario  uint   `json:"id_usuario"`
	Email      string `json:"email"`
	IDRol      uint   `json:"id_rol"`
	RolNombre  string `json:"rol_nombre"`
	IDSucursal uint   `json:"id_sucursal"`
	jwt.RegisteredClaims
}

// UsuarioFromClaims convierte los claims JWT a un objeto Usuario simplificado
func UsuarioFromClaims(claims *JWTClaims) *Usuario {
	return &Usuario{
		IDUsuario:  claims.IDUsuario,
		UsuEmail:   claims.Email,
		IDSucursal: claims.IDSucursal,
		Rol: &Rol{
			IDRol:     claims.IDRol,
			RolNombre: claims.RolNombre,
		},
	}
}

// LogoutResponse representa la respuesta de logout
type LogoutResponse struct {
	Message string `json:"message"`
}
