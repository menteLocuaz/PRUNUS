-- 1. ASEGURAR COLUMNA STOCK EN PRODUCTO
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='producto' AND column_name='stock') THEN
        ALTER TABLE producto ADD COLUMN stock NUMERIC(12,2) DEFAULT 0;
    END IF;
END $$;

-- 2. TABLA DE MOVIMIENTOS DE INVENTARIO
CREATE TABLE IF NOT EXISTS movimientos_inventario (
    id_movimiento    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_producto      UUID NOT NULL REFERENCES producto(id_producto),
    id_sucursal      UUID NOT NULL REFERENCES sucursal(id_sucursal),
    tipo_movimiento  VARCHAR(50) NOT NULL, -- VENTA, COMPRA, AJUSTE, etc.
    cantidad         NUMERIC(12,2) NOT NULL,
    stock_anterior   NUMERIC(12,2),
    stock_posterior  NUMERIC(12,2),
    id_referencia    UUID, -- ID de factura o compra relacionada
    observacion      TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at       TIMESTAMPTZ NULL
);

-- 3. FUNCIÓN DE SINCRONIZACIÓN DE STOCK (Lógica de Negocio)
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
        -- Si es AJUSTE genérico, intentamos inferir del valor
        IF NEW.cantidad < 0 THEN
            v_signo := 1; -- Si ya es negativo, no invertimos
        END IF;
    END IF;

    -- B. Obtener stock actual o inicializar (Upsert automático en inventario)
    SELECT stock_actual INTO v_stock_actual 
    FROM inventario 
    WHERE id_producto = NEW.id_producto AND id_sucursal = NEW.id_sucursal;

    IF NOT FOUND THEN
        v_stock_actual := 0;
        INSERT INTO inventario (id_producto, id_sucursal, stock_actual)
        VALUES (NEW.id_producto, NEW.id_sucursal, 0);
    END IF;

    -- C. Asignar valores de auditoría al movimiento
    NEW.stock_anterior := v_stock_actual;
    NEW.stock_posterior := v_stock_actual + (NEW.cantidad * v_signo);

    -- D. Actualizar stock en la sucursal (inventario)
    UPDATE inventario 
    SET stock_actual = NEW.stock_posterior,
        updated_at = CURRENT_TIMESTAMP
    WHERE id_producto = NEW.id_producto AND id_sucursal = NEW.id_sucursal;

    -- E. Sincronizar stock GLOBAL en la tabla producto (PostgreSQL Optimization: SUM)
    UPDATE producto 
    SET stock = (SELECT COALESCE(SUM(stock_actual), 0) FROM inventario WHERE id_producto = NEW.id_producto AND deleted_at IS NULL),
        updated_at = CURRENT_TIMESTAMP
    WHERE id_producto = NEW.id_producto;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 4. CREACIÓN DEL TRIGGER
DROP TRIGGER IF EXISTS trg_actualizar_stock_movimiento ON movimientos_inventario;
CREATE TRIGGER trg_actualizar_stock_movimiento
BEFORE INSERT ON movimientos_inventario
FOR EACH ROW
EXECUTE FUNCTION fn_actualizar_stock_movimiento();

-- 5. APLICAR FRAMEWORK CORE
CALL sp_core_setup_table('movimientos_inventario');
