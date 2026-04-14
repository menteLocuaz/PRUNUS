package models

import (
	"time"

	"github.com/google/uuid"
)

type Impuesto struct {
	IDImpuesto uuid.UUID  `json:"id_impuesto"`
	Nombre     string     `json:"nombre"`
	Porcentaje float64    `json:"porcentaje"`
	IDStatus   uuid.UUID  `json:"id_status"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}
