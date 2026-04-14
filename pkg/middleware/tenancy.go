package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/utils/tenancy"
)

// Tenancy extrae el sucursal_id y empresa_id de los claims JWT (inyectados
// previamente por RequireAuth) y los almacena en el contexto usando claves
// tipadas del paquete tenancy.
//
// Los stores y servicios consumen estos valores mediante tenancy.MustSucursalID(ctx)
// y tenancy.MustEmpresaID(ctx), garantizando que cada operación quede aislada
// al tenant del usuario autenticado sin depender de parámetros manuales.
//
// Comportamiento:
//   - Si no hay claims en el contexto (ruta pública o sin RequireAuth previo),
//     actúa como no-op para no bloquear endpoints de autenticación.
//   - Solo inyecta un valor si el UUID no es Nil, evitando tenants inválidos.
func Tenancy() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value("claims").(*models.JWTClaims)
			if !ok || claims == nil {
				// Ruta pública — continuar sin inyectar tenancy
				next.ServeHTTP(w, r)
				return
			}

			ctx := r.Context()

			if claims.IDSucursal != uuid.Nil {
				ctx = tenancy.WithSucursal(ctx, claims.IDSucursal)
			}
			if claims.IDEmpresa != uuid.Nil {
				ctx = tenancy.WithEmpresa(ctx, claims.IDEmpresa)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
