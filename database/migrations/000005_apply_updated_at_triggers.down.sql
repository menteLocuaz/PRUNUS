-- Revertir: Eliminar los triggers de updated_at creados masivamente
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
        EXECUTE format('DROP TRIGGER IF EXISTS trg_update_updated_at_%I ON %I', t, t);
    END LOOP;
END $$;
