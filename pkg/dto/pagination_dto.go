package dto

import "time"

// PaginationParams define los parámetros para la paginación por cursor (Keyset)
type PaginationParams struct {
	LastID   string     `json:"last_id"`
	LastDate *time.Time `json:"last_date"`
	Limit    int        `json:"limit"`
}

// DefaultLimit es el límite por defecto si no se especifica uno.
const DefaultLimit = 20

// MaxLimit es el límite máximo permitido para evitar respuestas sin acotamiento.
const MaxLimit = 200
