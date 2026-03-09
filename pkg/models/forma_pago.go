package models

import (
	"time"

	"github.com/google/uuid"
)

type FormaPago struct {
	IDFormaPago    uuid.UUID  `json:"id_forma_pago"`
	FmpCodigo      string     `json:"fmp_codigo"`
	FmpDescripcion string     `json:"fmp_descripcion"`
	IDStatus       uuid.UUID  `json:"id_status"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}
