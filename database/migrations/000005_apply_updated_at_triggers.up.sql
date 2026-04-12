-- 1. Asegurar que la función genérica exista
CREATE OR REPLACE FUNCTION fn_update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 2. Aplicar trigger a todas las tablas que tengan la columna 'updated_at'
-- Esto automatiza el mantenimiento para tablas existentes y futuras que se creen antes de esta migración.
DO $$ 
DECLARE 
    t text;
BEGIN
    FOR t IN 
        SELECT table_name 
        FROM information_schema.columns 
        WHERE column_name = 'updated_at' 
        AND table_schema = 'public'
    LOOP
        -- Eliminar trigger si ya existe para evitar duplicados
        EXECUTE format('DROP TRIGGER IF EXISTS trg_update_updated_at_%I ON %I', t, t);
        
        -- Crear el nuevo trigger
        EXECUTE format('CREATE TRIGGER trg_update_updated_at_%I 
                        BEFORE UPDATE ON %I 
                        FOR EACH ROW 
                        EXECUTE FUNCTION fn_update_updated_at_column()', t, t);
    END LOOP;
END $$;
