-- 1. SINCRONIZACIÓN DE DETALLE CON MOVIMIENTO DE INVENTARIO
CREATE OR REPLACE FUNCTION fn_sincronizar_detalle_movimiento()
RETURNS TRIGGER AS $$
BEGIN
    -- CASO INSERT: Registrar venta en inventario
    IF (TG_OP = 'INSERT') THEN
        INSERT INTO movimientos_inventario (
            id_producto, id_sucursal, id_usuario, tipo_movimiento, cantidad, id_referencia, observacion
        )
        SELECT 
            NEW.id_producto, f.id_sucursal, f.id_usuario, 'VENTA', NEW.cantidad, f.id_factura, 'Factura #' || f.fac_numero
        FROM factura f
        WHERE f.id_factura = NEW.id_factura;
        
        RETURN NEW;
    END IF;

    -- CASO UPDATE (Soft Delete): Marcar movimiento como eliminado
    IF (TG_OP = 'UPDATE') THEN
        IF (OLD.deleted_at IS NULL AND NEW.deleted_at IS NOT NULL) THEN
            UPDATE movimientos_inventario 
            SET deleted_at = NEW.deleted_at 
            WHERE id_producto = NEW.id_producto 
              AND id_referencia = NEW.id_factura;
        END IF;
        RETURN NEW;
    END IF;

    -- CASO DELETE (Físico): Eliminar movimiento
    IF (TG_OP = 'DELETE') THEN
        DELETE FROM movimientos_inventario 
        WHERE id_producto = OLD.id_producto 
          AND id_referencia = OLD.id_factura;
        RETURN OLD;
    END IF;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_registrar_movimiento_venta ON detalle_factura;
DROP TRIGGER IF EXISTS trg_sincronizar_detalle_movimiento ON detalle_factura;

CREATE TRIGGER trg_sincronizar_detalle_movimiento
AFTER INSERT OR UPDATE OR DELETE ON detalle_factura
FOR EACH ROW EXECUTE FUNCTION fn_sincronizar_detalle_movimiento();

-- 2. CASCADA DE BORRADO LÓGICO PARA FACTURA
CREATE OR REPLACE FUNCTION fn_factura_soft_delete_cascade()
RETURNS TRIGGER AS $$
BEGIN
    IF (OLD.deleted_at IS NULL AND NEW.deleted_at IS NOT NULL) THEN
        -- Cascada a detalles (esto activará fn_sincronizar_detalle_movimiento)
        UPDATE detalle_factura SET deleted_at = NEW.deleted_at WHERE id_factura = NEW.id_factura AND deleted_at IS NULL;
        -- Cascada a formas de pago
        UPDATE forma_pago_factura SET deleted_at = NEW.deleted_at WHERE id_factura = NEW.id_factura AND deleted_at IS NULL;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_factura_soft_delete_cascade ON factura;
CREATE TRIGGER trg_factura_soft_delete_cascade
AFTER UPDATE ON factura
FOR EACH ROW EXECUTE FUNCTION fn_factura_soft_delete_cascade();
