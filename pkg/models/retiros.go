package models

import (
	"time"

	"github.com/google/uuid"
)

type Retiro struct {
	IDRetiro               uuid.UUID  `json:"id_retiro"`
	ArcValor               float64    `json:"arc_valor"`
	ArcNumeroTransacciones int        `json:"arc_numero_transacciones"`
	IDControlEstacion      uuid.UUID  `json:"id_control_estacion"`
	IDFormaPago            uuid.UUID  `json:"id_forma_pago"`
	IDUserPos              uuid.UUID  `json:"id_user_pos"`
	UsuarioInicia          uint       `json:"usuario_inicia"`
	UsuarioFinaliza        *uint      `json:"usuario_finaliza,omitempty"`
	FechaInicio            time.Time  `json:"fecha_inicio"`
	FechaFinaliza          *time.Time `json:"fecha_finaliza,omitempty"`
	IDStatus               uuid.UUID  `json:"id_status"`
	PosCalculado           float64    `json:"pos_calculado"`
	DiferenciaValor        float64    `json:"diferencia_valor"`
	RetiroValor            float64    `json:"retiro_valor"`
	TPEnvID                int        `json:"tpenv_id"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
	DeletedAt              *time.Time `json:"deleted_at,omitempty"`
}
