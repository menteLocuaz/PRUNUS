package store

import (
	"context"
	"database/sql"

	"github.com/prunus/pkg/dto"

	"github.com/google/uuid"
)

type StoreDashboard interface {
	GetResumen(ctx context.Context, sucursalID uuid.UUID) (*dto.DashboardResumen, error)
	GetStockBajo(ctx context.Context, sucursalID uuid.UUID) ([]dto.TopProductoDTO, error)
	GetVentasVsCompras(ctx context.Context, sucursalID uuid.UUID) ([]dto.VentasComprasDTO, error)
	GetRentabilidadTop(ctx context.Context, sucursalID uuid.UUID) ([]dto.TopProductoDTO, error)
	GetAntiguedadDeuda(ctx context.Context, sucursalID uuid.UUID) ([]dto.AntiguedadDeudaDTO, error)
}

type dashboardStore struct {
	db *sql.DB
}

func NewDashboardStore(db *sql.DB) StoreDashboard {
	return &dashboardStore{db: db}
}

func (s *dashboardStore) GetResumen(ctx context.Context, sucursalID uuid.UUID) (*dto.DashboardResumen, error) {
	resumen := &dto.DashboardResumen{}

	// 1. Valor Inventario Total
	queryValor := `
		SELECT COALESCE(SUM(stock_actual * precio_compra), 0)
		FROM inventario
		WHERE id_sucursal = $1 AND deleted_at IS NULL`
	err := s.db.QueryRowContext(ctx, queryValor, sucursalID).Scan(&resumen.ValorInventarioTotal)
	if err != nil {
		return nil, err
	}

	// 2. Productos bajo stock
	queryBajoStock := `
		SELECT COUNT(*)
		FROM inventario
		WHERE id_sucursal = $1 AND stock_actual <= stock_minimo AND deleted_at IS NULL`
	err = s.db.QueryRowContext(ctx, queryBajoStock, sucursalID).Scan(&resumen.ProductosBajoStock)
	if err != nil {
		return nil, err
	}

	// 3. Ventas Mes Actual
	queryVentas := `
		SELECT COALESCE(SUM(cfac_total), 0)
		FROM factura f
		JOIN control_estacion ce ON f.id_control_estacion = ce.id_control_estacion
		WHERE ce.id_sucursal = $1 
		  AND f.id_status = '0f447fd7-9849-4a68-b82f-c69297e7a924' -- Pagada
		  AND f.cfac_fecha_creacion >= date_trunc('month', current_date)
		  AND f.deleted_at IS NULL`
	err = s.db.QueryRowContext(ctx, queryVentas, sucursalID).Scan(&resumen.VentasMesActual)
	if err != nil {
		return nil, err
	}

	// 4. Cuentas por Cobrar (Facturas Pendientes)
	queryCxC := `
		SELECT COALESCE(SUM(cfac_total), 0)
		FROM factura f
		JOIN control_estacion ce ON f.id_control_estacion = ce.id_control_estacion
		WHERE ce.id_sucursal = $1 
		  AND f.id_status = '892340e0-4328-491d-9102-80550bb6aac4' -- Pendiente de Pago
		  AND f.deleted_at IS NULL`
	err = s.db.QueryRowContext(ctx, queryCxC, sucursalID).Scan(&resumen.CuentasPorCobrar)
	if err != nil {
		return nil, err
	}

	// 5. Cuentas por Pagar (Ordenes de Compra Recibidas pero no Pagadas)
	// Nota: Asumimos que ID_STATUS en orden_compra refleja el estado de pago/recepción
	queryCxP := `
		SELECT COALESCE(SUM(total), 0)
		FROM orden_compra
		WHERE id_sucursal = $1 
		  AND id_status IN (
			  '00363491-8508-4220-9661-e99f05b0d545', -- Recibida Parcialmente
			  '00363491-8508-4220-9661-e99f05b00001'  -- Recibida Completa
		  )
		  AND deleted_at IS NULL`
	err = s.db.QueryRowContext(ctx, queryCxP, sucursalID).Scan(&resumen.CuentasPorPagar)
	if err != nil {
		return nil, err
	}

	// 6. Gastos Mensuales (Usando la nueva función de Postgres)
	queryGastos := `SELECT fn_get_gastos_mensuales($1, current_date)`
	err = s.db.QueryRowContext(ctx, queryGastos, sucursalID).Scan(&resumen.GastosMensuales)
	if err != nil {
		return nil, err
	}

	// 7. Punto de Equilibrio (Ventas necesarias para cubrir gastos)
	// PE = Gastos Fijos / Margen de Contribución Promedio
	// Simplificado: Ventas donde (Ventas - Costos) = Gastos
	queryPE := `
		WITH MargenPromedio AS (
			SELECT 
				COALESCE(SUM((m.precio_unitario - m.costo_unitario) * m.cantidad) / NULLIF(SUM(m.precio_unitario * m.cantidad), 0), 0.2) as margen
			FROM movimientos_inventario m
			WHERE m.id_sucursal = $1 
			  AND m.tipo_movimiento = 'VENTA'
			  AND m.fecha >= date_trunc('month', current_date)
		)
		SELECT COALESCE($2 / NULLIF(margen, 0), 0) FROM MargenPromedio`
	err = s.db.QueryRowContext(ctx, queryPE, sucursalID, resumen.GastosMensuales).Scan(&resumen.PuntoEquilibrio)
	if err != nil {
		// Si no hay ventas, no podemos calcular margen real, usamos 0 o un default
		resumen.PuntoEquilibrio = 0
	}

	return resumen, nil
}

func (s *dashboardStore) GetStockBajo(ctx context.Context, sucursalID uuid.UUID) ([]dto.TopProductoDTO, error) {
	query := `
		SELECT p.nombre, i.stock_actual, 0 as rentabilidad
		FROM inventario i
		JOIN producto p ON i.id_producto = p.id_producto
		WHERE i.id_sucursal = $1 
		  AND i.stock_actual <= i.stock_minimo 
		  AND i.deleted_at IS NULL
		ORDER BY i.stock_actual ASC
		LIMIT 10`

	rows, err := s.db.QueryContext(ctx, query, sucursalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []dto.TopProductoDTO
	for rows.Next() {
		var item dto.TopProductoDTO
		if err := rows.Scan(&item.Nombre, &item.Cantidad, &item.Rentabilidad); err != nil {
			return nil, err
		}
		results = append(results, item)
	}
	return results, nil
}

func (s *dashboardStore) GetVentasVsCompras(ctx context.Context, sucursalID uuid.UUID) ([]dto.VentasComprasDTO, error) {
	query := `
		WITH Meses AS (
			SELECT generate_series(
				date_trunc('month', current_date - interval '5 months'),
				date_trunc('month', current_date),
				interval '1 month'
			)::date as mes
		),
		Ventas AS (
			SELECT date_trunc('month', f.cfac_fecha_creacion)::date as mes, SUM(f.cfac_total) as total
			FROM factura f
			JOIN control_estacion ce ON f.id_control_estacion = ce.id_control_estacion
			WHERE ce.id_sucursal = $1 AND f.id_status = '0f447fd7-9849-4a68-b82f-c69297e7a924'
			GROUP BY 1
		),
		Compras AS (
			SELECT date_trunc('month', oc.fecha_emision)::date as mes, SUM(oc.total) as total
			FROM orden_compra oc
			WHERE oc.id_sucursal = $1 AND oc.deleted_at IS NULL
			GROUP BY 1
		)
		SELECT 
			to_char(m.mes, 'TMMonth') as mes_nombre,
			COALESCE(v.total, 0) as ventas,
			COALESCE(c.total, 0) as compras
		FROM Meses m
		LEFT JOIN Ventas v ON m.mes = v.mes
		LEFT JOIN Compras c ON m.mes = c.mes
		ORDER BY m.mes ASC`

	rows, err := s.db.QueryContext(ctx, query, sucursalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []dto.VentasComprasDTO
	for rows.Next() {
		var item dto.VentasComprasDTO
		if err := rows.Scan(&item.Mes, &item.Ventas, &item.Compras); err != nil {
			return nil, err
		}
		results = append(results, item)
	}
	return results, nil
}

func (s *dashboardStore) GetRentabilidadTop(ctx context.Context, sucursalID uuid.UUID) ([]dto.TopProductoDTO, error) {
	// Cálculo de rentabilidad: (Precio Venta - Costo Unitario) * Cantidad
	// Usamos movimientos_inventario para obtener el costo real al momento de la venta
	query := `
		SELECT 
			p.nombre,
			SUM(m.cantidad) as total_cantidad,
			SUM((m.precio_unitario - m.costo_unitario) * m.cantidad) as total_rentabilidad
		FROM movimientos_inventario m
		JOIN producto p ON m.id_producto = p.id_producto
		WHERE m.id_sucursal = $1 
		  AND m.tipo_movimiento = 'VENTA'
		  AND m.fecha >= date_trunc('month', current_date - interval '1 month')
		GROUP BY p.id_producto, p.nombre
		HAVING SUM((m.precio_unitario - m.costo_unitario) * m.cantidad) > 0
		ORDER BY total_rentabilidad DESC
		LIMIT 10`

	rows, err := s.db.QueryContext(ctx, query, sucursalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []dto.TopProductoDTO
	for rows.Next() {
		var item dto.TopProductoDTO
		if err := rows.Scan(&item.Nombre, &item.Cantidad, &item.Rentabilidad); err != nil {
			return nil, err
		}
		results = append(results, item)
	}
	return results, nil
}

func (s *dashboardStore) GetAntiguedadDeuda(ctx context.Context, sucursalID uuid.UUID) ([]dto.AntiguedadDeudaDTO, error) {
	query := `
		SELECT 
			CASE 
				WHEN current_date - fecha_vencimiento::date <= 30 THEN '0-30 días'
				WHEN current_date - fecha_vencimiento::date <= 60 THEN '31-60 días'
				WHEN current_date - fecha_vencimiento::date <= 90 THEN '61-90 días'
				ELSE '90+ días'
			END as rango,
			SUM(cfac_total) as monto
		FROM factura f
		JOIN control_estacion ce ON f.id_control_estacion = ce.id_control_estacion
		WHERE ce.id_sucursal = $1 
		  AND f.id_status = '892340e0-4328-491d-9102-80550bb6aac4' -- Pendiente de Pago
		  AND f.fecha_vencimiento IS NOT NULL
		  AND f.deleted_at IS NULL
		GROUP BY 1
		ORDER BY 1`

	rows, err := s.db.QueryContext(ctx, query, sucursalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []dto.AntiguedadDeudaDTO
	for rows.Next() {
		var item dto.AntiguedadDeudaDTO
		if err := rows.Scan(&item.Rango, &item.Monto); err != nil {
			return nil, err
		}
		results = append(results, item)
	}
	return results, nil
}
