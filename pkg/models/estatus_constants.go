package models

import "github.com/google/uuid"

var (
	// Estatus Generales
	EstatusActivo   = uuid.MustParse("59039503-85CF-E511-80C1-000C29C9E0E0")
	EstatusInactivo = uuid.MustParse("5A039503-85CF-E511-80C1-000C29C9E0E0")

	// Estatus POS / Control Estación
	EstatusFondoAsignado     = uuid.MustParse("99039503-85CF-E511-80C1-000C29C9E0E0")
	EstatusDesmontado        = uuid.MustParse("9A039503-85CF-E511-80C1-000C29C9E0E0")
	EstatusIngresoAdmin      = uuid.MustParse("9B039503-85CF-E511-80C1-000C29C9E0E0")
	EstatusSalirAdmin        = uuid.MustParse("9C039503-85CF-E511-80C1-000C29C9E0E0")
	EstatusArqueo            = uuid.MustParse("0D4515FE-C907-E611-A6B8-000C29C9E0E0")
	EstatusArqueoRetiros     = uuid.MustParse("159E3FE6-630E-E611-80C1-000C29C9E0E0")
	EstatusRetiroEfectivo    = uuid.MustParse("E8297CFA-630E-E611-80C1-000C29C9E0E0")
	EstatusRetiroTotal       = uuid.MustParse("84920103-640E-E611-80C1-000C29C9E0E0")
	EstatusFondoActivo       = uuid.MustParse("A864475F-0D34-E711-80C1-000C29C9E0E0")
	EstatusFondoRetirado     = uuid.MustParse("2160B065-0D34-E711-80C1-000C29C9E0E0")
	EstatusFondoPorConfirmar = uuid.MustParse("5E8DD0FB-5550-E711-80C1-000C29C9E0E0")
)
