package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/prunus/pkg/helper"
	"github.com/prunus/pkg/utils/tenancy"
)

// RequireAuth es un middleware que valida la presencia y validez del JWT token
// Extrae el token del header Authorization y valida su firma y expiración
func RequireAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Obtener el header Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Token de autenticación requerido", http.StatusUnauthorized)
				return
			}

			// Extraer el token del header (formato: "Bearer <token>")
			token, err := helper.ExtractTokenFromHeader(authHeader)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			// Validar el token
			claims, err := helper.ValidateToken(token)
			if err != nil {
				// Manejar errores específicos
				if err == helper.ErrExpiredToken {
					http.Error(w, "Token expirado", http.StatusUnauthorized)
					return
				}
				http.Error(w, "Token inválido", http.StatusUnauthorized)
				return
			}

			// Agregar los claims al contexto de la petición
			// Los handlers pueden acceder a esta información
			ctx := context.WithValue(r.Context(), "claims", claims)
			ctx = context.WithValue(ctx, "user_id", claims.IDUsuario)
			ctx = context.WithValue(ctx, "user_email", claims.Email)
			ctx = context.WithValue(ctx, "user_rol", claims.RolNombre)
			ctx = context.WithValue(ctx, "user_sucursal", claims.IDSucursal)

			// Inyectar tenancy usando claves tipadas para que los stores
			// puedan filtrar datos por sucursal/empresa de forma automática
			ctx = tenancy.WithSucursal(ctx, claims.IDSucursal)
			ctx = tenancy.WithEmpresa(ctx, claims.IDEmpresa)

			// Continuar con el siguiente handler
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole es un middleware que verifica que el usuario tenga uno de los roles permitidos
// Debe usarse DESPUÉS de RequireAuth()
func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Obtener el rol del contexto
			userRol, ok := r.Context().Value("user_rol").(string)
			if !ok {
				http.Error(w, "No se pudo verificar el rol del usuario", http.StatusForbidden)
				return
			}

			// Verificar si el rol del usuario está en la lista de roles permitidos
			roleAllowed := false
			for _, role := range allowedRoles {
				if strings.EqualFold(userRol, role) {
					roleAllowed = true
					break
				}
			}

			if !roleAllowed {
				http.Error(w, "No tienes permisos para acceder a este recurso", http.StatusForbidden)
				return
			}

			// Continuar con el siguiente handler
			next.ServeHTTP(w, r)
		})
	}
}

// OptionalAuth es un middleware que intenta autenticar pero no falla si no hay token
// Útil para endpoints que funcionan diferente si el usuario está autenticado
func OptionalAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Obtener el header Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				// No hay token, continuar sin autenticación
				next.ServeHTTP(w, r)
				return
			}

			// Intentar extraer el token
			token, err := helper.ExtractTokenFromHeader(authHeader)
			if err != nil {
				// Token mal formado, continuar sin autenticación
				next.ServeHTTP(w, r)
				return
			}

			// Intentar validar el token
			claims, err := helper.ValidateToken(token)
			if err != nil {
				// Token inválido o expirado, continuar sin autenticación
				next.ServeHTTP(w, r)
				return
			}

			// Token válido, agregar al contexto
			ctx := context.WithValue(r.Context(), "claims", claims)
			ctx = context.WithValue(ctx, "user_id", claims.IDUsuario)
			ctx = context.WithValue(ctx, "user_email", claims.Email)
			ctx = context.WithValue(ctx, "user_rol", claims.RolNombre)
			ctx = context.WithValue(ctx, "user_sucursal", claims.IDSucursal)
			ctx = tenancy.WithSucursal(ctx, claims.IDSucursal)
			ctx = tenancy.WithEmpresa(ctx, claims.IDEmpresa)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
