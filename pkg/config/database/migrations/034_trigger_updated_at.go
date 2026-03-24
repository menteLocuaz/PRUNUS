package migrations

import "database/sql"

func migrateTriggerUpdatedAt(db *sql.DB) error {
	query := `
	-- Función genérica para actualizar updated_at
	CREATE OR REPLACE FUNCTION fn_update_updated_at_column()
	RETURNS TRIGGER AS $$
	BEGIN
		NEW.updated_at = CURRENT_TIMESTAMP;
		RETURN NEW;
	END;
	$$ language 'plpgsql';

	-- Aplicar trigger a tablas existentes (ejemplos principales)
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
			EXECUTE format('CREATE TRIGGER trg_update_updated_at_%I 
							BEFORE UPDATE ON %I 
							FOR EACH ROW 
							EXECUTE FUNCTION fn_update_updated_at_column()', t, t);
		END LOOP;
	END $$;
	`
	_, err := db.Exec(query)
	return err
}
