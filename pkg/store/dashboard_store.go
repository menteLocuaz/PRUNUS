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
	GetComposicionCategoria(ctx context.Context, sucursalID uuid.UUID) ([]dto.InventarioCategoriaDTO, error)
	GetMermas(ctx context.Context, sucursalID uuid.UUID) ([]dto.TopProductoDTO, error)
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
		SELECT COALESCE(SUM(f.cfac_total), 0)
		FROM factura f
		JOIN estaciones_pos ep ON f.id_estacion = ep.id_estacion
		WHERE ep.id_sucursal = $1 
		  AND f.id_status = '0f447fd7-9849-4a68-b82f-c69297e7a924' -- Pagada
		  AND f.cfac_fecha_creacion >= date_trunc('month', current_date)
		  AND f.deleted_at IS NULL`
	err = s.db.QueryRowContext(ctx, queryVentas, sucursalID).Scan(&resumen.VentasMesActual)
	if err != nil {
		return nil, err
	}

	// 4. Cuentas por Cobrar (Facturas Pendientes)
	queryCxC := `
		SELECT COALESCE(SUM(f.cfac_total), 0)
		FROM factura f
		JOIN estaciones_pos ep ON f.id_estacion = ep.id_estacion
		WHERE ep.id_sucursal = $1 
		  AND f.id_status = '892340e0-4328-491d-9102-80550bb6aac4' -- Pendiente de Pago
		  AND f.deleted_at IS NULL`
	err = s.db.QueryRowContext(ctx, queryCxC, sucursalID).Scan(&resumen.CuentasPorCobrar)
	if err != nil {
		return nil, err
	}

	// 5. Cuentas por Pagar (Ordenes de Compra Recibidas pero no Pagadas)
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

	// 7. Punto de Equilibrio
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
		resumen.PuntoEquilibrio = 0
	}

	// 8. Ciclo de Conversión de Efectivo
	queryCCC := `
		WITH Periodo AS (
			SELECT 90 as dias
		),
		Metricas AS (
			SELECT COALESCE(SUM(m.costo_unitario * m.cantidad), 0) as cogs
			FROM movimientos_inventario m
			WHERE m.id_sucursal = $1 
			  AND m.tipo_movimiento = 'VENTA' 
			  AND m.fecha >= current_date - interval '90 days'
		),
		InventarioPromedio AS (
			SELECT COALESCE(AVG(valor_total), 0) as avg_inv
			FROM inventario_historico
			WHERE id_sucursal = $1 AND fecha_snapshot >= current_date - interval '90 days'
		),
		VentasCredito AS (
			SELECT COALESCE(SUM(f.cfac_total), 0) as total_ventas,
			       COALESCE(AVG(f.cfac_total), 0) as avg_ar
			FROM factura f
			JOIN estaciones_pos ep ON f.id_estacion = ep.id_estacion
			WHERE ep.id_sucursal = $1 
			  AND f.cfac_fecha_creacion >= current_date - interval '90 days'
			  AND f.fecha_vencimiento IS NOT NULL
		),
		ComprasCredito AS (
			SELECT COALESCE(SUM(total), 0) as total_compras,
			       COALESCE(AVG(total), 0) as avg_ap
			FROM orden_compra
			WHERE id_sucursal = $1 
			  AND fecha_emision >= current_date - interval '90 days'
			  AND fecha_vencimiento IS NOT NULL
		)
		SELECT 
			COALESCE((i.avg_inv / NULLIF(m.cogs, 0)) * p.dias, 0) as dio,
			COALESCE((v.avg_ar / NULLIF(v.total_ventas, 0)) * p.dias, 0) as dso,
			COALESCE((c.avg_ap / NULLIF(c.total_compras, 0)) * p.dias, 0) as dpo
		FROM Periodo p, Metricas m, InventarioPromedio i, VentasCredito v, ComprasCredito c`

	err = s.db.QueryRowContext(ctx, queryCCC, sucursalID).Scan(&resumen.DIO, &resumen.DSO, &resumen.DPO)
	if err == nil {
		resumen.CicloConversionEfectivo = resumen.DIO + resumen.DSO - resumen.DPO
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
			JOIN estaciones_pos ep ON f.id_estacion = ep.id_estacion
			WHERE ep.id_sucursal = $1 AND f.id_status = '0f447fd7-9849-4a68-b82f-c69297e7a924'
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
			SUM(f.cfac_total) as monto
		FROM factura f
		JOIN estaciones_pos ep ON f.id_estacion = ep.id_estacion
		WHERE ep.id_sucursal = $1 
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

func (s *dashboardStore) GetComposicionCategoria(ctx context.Context, sucursalID uuid.UUID) ([]dto.InventarioCategoriaDTO, error) {
	query := `
		WITH StockValores AS (
			SELECT 
				c.nombre as categoria,
				SUM(i.stock_actual * i.precio_compra) as valor
			FROM inventario i
			JOIN producto p ON i.id_producto = p.id_producto
			JOIN categoria c ON p.id_categoria = c.id_categoria
			WHERE i.id_sucursal = $1 AND i.deleted_at IS NULL
			GROUP BY c.id_categoria, c.nombre
		),
		Total AS (
			SELECT SUM(valor) as gran_total FROM StockValores
		)
		SELECT 
			categoria, 
			valor, 
			COALESCE((valor / NULLIF(gran_total, 0)) * 100, 0) as porcentaje
		FROM StockValores, Total
		ORDER BY valor DESC`

	rows, err := s.db.QueryContext(ctx, query, sucursalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []dto.InventarioCategoriaDTO
	for rows.Next() {
		var item dto.InventarioCategoriaDTO
		if err := rows.Scan(&item.Categoria, &item.Valor, &item.Porcentaje); err != nil {
			return nil, err
		}
		results = append(results, item)
	}
	return results, nil
}

func (s *dashboardStore) GetMermas(ctx context.Context, sucursalID uuid.UUID) ([]dto.TopProductoDTO, error) {
	query := `
		SELECT 
			p.nombre,
			SUM(m.cantidad) as total_cantidad,
			SUM(m.costo_unitario * m.cantidad) as total_perdida
		FROM movimientos_inventario m
		JOIN producto p ON m.id_producto = p.id_producto
		WHERE m.id_sucursal = $1 
		  AND m.tipo_movimiento IN ('MERMA', 'CADUCADO')
		  AND m.fecha >= date_trunc('month', current_date)
		GROUP BY p.id_producto, p.nombre
		ORDER BY total_perdida DESC
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
