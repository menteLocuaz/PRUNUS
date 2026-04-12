-- Revertir: Eliminar motor de consistencia de estados
DO $$ 
DECLARE 
    r RECORD;
BEGIN
    FOR r IN SELECT table_name FROM config.trigger_module_map LOOP
        EXECUTE format('DROP TRIGGER IF EXISTS %I ON %I', 'tr_status_check_' || r.table_name, r.table_name);
    END LOOP;
END $$;

DROP PROCEDURE IF EXISTS config.sp_apply_status_rule(TEXT);
DROP FUNCTION IF EXISTS config.fn_validate_status_module();
DROP TABLE IF EXISTS config.trigger_module_map;
