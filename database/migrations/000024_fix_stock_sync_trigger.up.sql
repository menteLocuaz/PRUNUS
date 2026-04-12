-- 1. CORRECCIÓN DE FUNCIÓN: fn_actualizar_stock_movimiento
-- Se elimina la referencia a producto.stock (columna deprecada en 000017)
CREATE OR REPLACE FUNCTION fn_actualizar_stock_movimiento()
RETURNS TRIGGER AS $$
DECLARE
    v_stock_actual NUMERIC(12,2) := 0;
    v_signo INTEGER := 1;
BEGIN
    -- A. Determinar el signo según el tipo de movimiento
    IF NEW.tipo_movimiento IN ('VENTA', 'AJUSTE_SALIDA', 'SALIDA') THEN
        v_signo := -1;
    ELSIF NEW.tipo_movimiento IN ('COMPRA', 'DEVOLUCION', 'ENTRADA', 'AJUSTE_ENTRADA') THEN
        v_signo := 1;
    ELSE
        -- AJUSTE genérico: asume entrada a menos que la cantidad sea negativa
        IF NEW.cantidad < 0 THEN
            v_signo := 1;
        END IF;
    END IF;

    -- B. Obtener stock actual o inicializar registro en inventario
    SELECT stock_actual INTO v_stock_actual
    FROM inventario
    WHERE id_producto = NEW.id_producto AND id_sucursal = NEW.id_sucursal;

    IF NOT FOUND THEN
        v_stock_actual := 0;
        INSERT INTO inventario (id_producto, id_sucursal, stock_actual)
        VALUES (NEW.id_producto, NEW.id_sucursal, 0);
    END IF;

    -- C. Registrar valores de auditoría en el movimiento
    NEW.stock_anterior  := v_stock_actual;
    NEW.stock_posterior := v_stock_actual + (NEW.cantidad * v_signo);

    -- D. Actualizar stock en la sucursal correspondiente
    UPDATE inventario
    SET stock_actual = NEW.stock_posterior,
        updated_at   = CURRENT_TIMESTAMP
    WHERE id_producto = NEW.id_producto AND id_sucursal = NEW.id_sucursal;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
