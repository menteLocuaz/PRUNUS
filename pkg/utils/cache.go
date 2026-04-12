package utils

import (
	"context"
	"time"

	"github.com/prunus/pkg/models"
)

// CacheManager centraliza la abstracción de caché para los servicios.
// Proporciona métodos genéricos para reducir el código repetitivo (boilerplate)
// y estandarizar la invalidación por etiquetas/prefijos.
type CacheManager struct {
	store models.CacheStore
}

// NewCacheManager crea una nueva instancia del gestor de caché.
func NewCacheManager(store models.CacheStore) *CacheManager {
	return &CacheManager{store: store}
}

// GetOrSet encapsula el patrón "Look-aside":
// 1. Intenta obtener el valor del caché (si el store está disponible).
// 2. Si no existe (Cache Miss), ejecuta la función 'fetch'.
// 3. Si 'fetch' es exitoso, guarda el resultado en caché (si el store está disponible) y lo retorna.
func GetOrSet[T any](ctx context.Context, cm *CacheManager, key string, ttl time.Duration, fetch func() (T, error)) (T, error) {
	var result T

	// Si no hay store de caché, simplemente ejecutar fetch y retornar
	if cm == nil || cm.store == nil {
		return fetch()
	}

	// Intentar recuperar de caché
	if err := cm.store.Get(ctx, key, &result); err == nil {
		return result, nil
	}

	// Si falla, obtener de la fuente original
	data, err := fetch()
	if err != nil {
		return data, err
	}

	// Guardar en caché para futuras consultas
	_ = cm.store.Set(ctx, key, data, ttl)

	return data, nil
}

// Invalidate limpia el caché basado en etiquetas o prefijos.
// Esto permite invalidar grupos de datos relacionados (ej: todos los roles).
func (cm *CacheManager) Invalidate(ctx context.Context, tags ...string) {
	if cm == nil || cm.store == nil {
		return
	}
	for _, tag := range tags {
		// Usamos el prefijo como etiqueta para invalidación masiva
		_ = cm.store.DeleteByPrefix(ctx, tag)
	}
}

// Delete elimina una llave específica del caché.
func (cm *CacheManager) Delete(ctx context.Context, key string) {
	_ = cm.store.Delete(ctx, key)
}
