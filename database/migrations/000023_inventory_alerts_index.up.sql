-- 1. ÍNDICE DE ALERTAS DE STOCK (PostgreSQL Optimization)
-- Acelera consultas del tipo: "Mostrar productos con stock bajo en la sucursal X"
-- Al ser un índice parcial, su tamaño es mínimo y su velocidad máxima.
CREATE INDEX IF NOT EXISTS idx_inventario_alertas_stock 
ON inventario(id_sucursal, stock_actual, stock_minimo) 
WHERE deleted_at IS NULL AND stock_actual <= stock_minimo;
