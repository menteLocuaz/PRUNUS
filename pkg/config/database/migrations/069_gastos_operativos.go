package migrations

import "database/sql"

func migrateGastosOperativos(db *sql.DB) error {
	query := `
	-- Tabla de gastos operativos (Fijos y Variables)
	CREATE TABLE IF NOT EXISTS gastos_operativos (
		id_gasto    UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
		id_sucursal UUID        NOT NULL REFERENCES sucursal(id_sucursal),
		id_usuario  UUID        NOT NULL REFERENCES usuario(id_usuario),
		descripcion VARCHAR(255) NOT NULL,
		monto       NUMERIC(18,4) NOT NULL DEFAULT 0,
		frecuencia  VARCHAR(20) NOT NULL DEFAULT 'MENSUAL', -- MENSUAL, ANUAL, UNICO
		fecha_gasto DATE        NOT NULL DEFAULT CURRENT_DATE,
		created_at  TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at  TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at  TIMESTAMP   NULL
	);

	CREATE INDEX IF NOT EXISTS idx_gastos_sucursal_fecha ON gastos_operativos(id_sucursal, fecha_gasto);

	-- Función auxiliar para obtener el total de gastos operativos mensuales de una sucursal
	CREATE OR REPLACE FUNCTION fn_get_gastos_mensuales(p_sucursal_id UUID, p_fecha DATE)
	RETURNS NUMERIC AS $$
	DECLARE
		v_total NUMERIC;
	BEGIN
		SELECT COALESCE(SUM(
			CASE 
				WHEN frecuencia = 'MENSUAL' THEN monto
				WHEN frecuencia = 'ANUAL' THEN monto / 12
				WHEN frecuencia = 'UNICO' AND date_trunc('month', fecha_gasto) = date_trunc('month', p_fecha) THEN monto
				ELSE 0
			END
		), 0) INTO v_total
		FROM gastos_operativos
		WHERE id_sucursal = p_sucursal_id 
		  AND deleted_at IS NULL;
		
		RETURN v_total;
	END;
	$$ LANGUAGE plpgsql;
	`
	_, err := db.Exec(query)
	return err
}
