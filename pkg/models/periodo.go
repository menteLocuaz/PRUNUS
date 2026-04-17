package models

import (
	"time"

	"github.com/google/uuid"
)

type Periodo struct {
	IDPeriodo          uuid.UUID  `json:"id_periodo"`
	IDSucursal         uuid.UUID  `json:"id_sucursal"`
	PrdFechaApertura   time.Time  `json:"prd_fecha_apertura"`
	PrdFechaCierre     *time.Time `json:"prd_fecha_cierre,omitempty"`
	PrdUsuarioApertura uuid.UUID  `json:"prd_usuario_apertura"`
	PrdUsuarioCierre   *uuid.UUID `json:"prd_usuario_cierre,omitempty"`
	PrdIPApertura      string     `json:"prd_ip_apertura,omitempty"`
	PrdMotivoApertura  string     `json:"prd_motivo_apertura,omitempty"`
	PrdIPCierre        string     `json:"prd_ip_cierre,omitempty"`
	IDStatus           uuid.UUID  `json:"id_status"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty"`
}
