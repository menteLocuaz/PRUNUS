package dto

import (
	"time"
	"github.com/google/uuid"
)

type CajaDTO struct {
	IDCaja      uuid.UUID `json:"id_caja"`
	Nombre      string    `json:"nombre"`
	IDSucursal  uuid.UUID `json:"id_sucursal"`
	IDStatus    uuid.UUID `json:"id_status"`
	CreatedAt   time.Time `json:"created_at"`
}

type DenominacionDTO struct {
	ValorNominal float64 `json:"valor_nominal"`
	Cantidad     int     `json:"cantidad"`
	Subtotal     float64 `json:"subtotal"`
}

type CierreCajaRequest struct {
	IDControlEstacion uuid.UUID         `json:"id_control_estacion"`
	MontoDeclarado    float64           `json:"monto_declarado"`
	Desglose          []DenominacionDTO `json:"desglose"`
	Observaciones     string            `json:"observaciones"`
}

type ResumenCierreDTO struct {
	FondoInicial      float64 `json:"fondo_inicial"`
	VentasEfectivo    float64 `json:"ventas_efectivo"`
	VentasTarjeta     float64 `json:"ventas_tarjeta"`
	VentasTransfer    float64 `json:"ventas_transferencia"`
	TotalRetiros      float64 `json:"total_retiros"`
	TotalGastos       float64 `json:"total_gastos"`
	SaldoEsperado     float64 `json:"saldo_esperado"`
	SaldoReal         float64 `json:"saldo_real"`
	Diferencia        float64 `json:"diferencia"`
	Resultado         string  `json:"resultado"` // CUADRADO, FALTANTE, SOBRANTE
}
