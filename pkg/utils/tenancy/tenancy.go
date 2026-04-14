package tenancy

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// contextKey es un tipo privado para evitar colisiones con otras claves de contexto
// de paquetes externos o del propio proyecto.
type contextKey int

const (
	sucursalKey contextKey = iota
	empresaKey
)

var (
	// ErrNoSucursal se retorna cuando el sucursal_id no está disponible en el contexto.
	ErrNoSucursal = errors.New("sucursal_id no disponible en el contexto de tenancy")
	// ErrNoEmpresa se retorna cuando el empresa_id no está disponible en el contexto.
	ErrNoEmpresa = errors.New("empresa_id no disponible en el contexto de tenancy")
)

// WithSucursal retorna un nuevo contexto con el sucursal_id inyectado.
// Lo llama el middleware Tenancy() después de validar el JWT.
func WithSucursal(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, sucursalKey, id)
}

// WithEmpresa retorna un nuevo contexto con el empresa_id inyectado.
func WithEmpresa(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, empresaKey, id)
}

// SucursalID extrae el sucursal_id del contexto.
// Retorna (uuid.Nil, false) si no está presente o es uuid.Nil.
func SucursalID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(sucursalKey).(uuid.UUID)
	return id, ok && id != uuid.Nil
}

// EmpresaID extrae el empresa_id del contexto.
// Retorna (uuid.Nil, false) si no está presente o es uuid.Nil.
func EmpresaID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(empresaKey).(uuid.UUID)
	return id, ok && id != uuid.Nil
}

// MustSucursalID retorna el sucursal_id o un error si no está disponible.
// Los stores deben usar esta función cuando el filtro por sucursal es obligatorio.
//
//	idSucursal, err := tenancy.MustSucursalID(ctx)
//	if err != nil {
//	    return nil, err
//	}
func MustSucursalID(ctx context.Context) (uuid.UUID, error) {
	id, ok := SucursalID(ctx)
	if !ok {
		return uuid.Nil, ErrNoSucursal
	}
	return id, nil
}

// MustEmpresaID retorna el empresa_id o un error si no está disponible.
func MustEmpresaID(ctx context.Context) (uuid.UUID, error) {
	id, ok := EmpresaID(ctx)
	if !ok {
		return uuid.Nil, ErrNoEmpresa
	}
	return id, nil
}
