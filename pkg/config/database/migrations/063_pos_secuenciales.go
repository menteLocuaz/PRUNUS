package migrations

import "database/sql"

func migratePosSecuenciales(db *sql.DB) error {
	query := `
	-- 1. Tabla para controlar los correlativos por estación
	CREATE TABLE IF NOT EXISTS pos_secuenciales (
		id_secuencial   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		id_estacion     UUID NOT NULL,
		tipo_documento  VARCHAR(20) NOT NULL DEFAULT 'FACTURA', -- FACTURA, NOTA_CREDITO, etc.
		prefijo         VARCHAR(10), -- Ej: '001-001'
		ultimo_numero   BIGINT NOT NULL DEFAULT 0,
		longitud        INTEGER NOT NULL DEFAULT 9,
		created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

		CONSTRAINT fk_secuencial_estacion FOREIGN KEY (id_estacion) REFERENCES estaciones_pos(id_estacion),
		CONSTRAINT uk_secuencial_estacion_tipo UNIQUE (id_estacion, tipo_documento)
	);

	-- 2. Función para obtener y aumentar el siguiente número (Thread-safe)
	CREATE OR REPLACE FUNCTION fn_get_next_secuencial(
		p_id_estacion    UUID,
		p_tipo_documento VARCHAR(20)
	)
	RETURNS VARCHAR(50)
	LANGUAGE plpgsql
	AS $$
	DECLARE
		v_prefijo       VARCHAR(10);
		v_ultimo        BIGINT;
		v_longitud      INTEGER;
		v_nuevo_numero  VARCHAR(50);
	BEGIN
		-- Bloquear la fila para evitar duplicados en concurrencia
		SELECT prefijo, ultimo_numero, longitud 
		INTO v_prefijo, v_ultimo, v_longitud
		FROM pos_secuenciales
		WHERE id_estacion = p_id_estacion AND tipo_documento = p_tipo_documento
		FOR UPDATE;

		-- Si no existe configuración para esta estación, crear una por defecto
		IF NOT FOUND THEN
			-- Intentar construir un prefijo básico basado en el código de la estación
			SELECT SUBSTRING(codigo, 1, 7) INTO v_prefijo FROM estaciones_pos WHERE id_estacion = p_id_estacion;
			v_prefijo := COALESCE(v_prefijo, '001-001');
			v_ultimo := 0;
			v_longitud := 9;

			INSERT INTO pos_secuenciales (id_estacion, tipo_documento, prefijo, ultimo_numero, longitud)
			VALUES (p_id_estacion, p_tipo_documento, v_prefijo, v_ultimo, v_longitud);
			
			-- Re-bloquear
			SELECT prefijo, ultimo_numero, longitud INTO v_prefijo, v_ultimo, v_longitud
			FROM pos_secuenciales WHERE id_estacion = p_id_estacion AND tipo_documento = p_tipo_documento FOR UPDATE;
		END IF;

		-- Incrementar
		v_ultimo := v_ultimo + 1;
		
		-- Formatear (Ej: 001-001-000000001)
		v_nuevo_numero := v_prefijo || '-' || LPAD(v_ultimo::TEXT, v_longitud, '0');

		-- Actualizar tabla
		UPDATE pos_secuenciales 
		SET ultimo_numero = v_ultimo, updated_at = CURRENT_TIMESTAMP
		WHERE id_estacion = p_id_estacion AND tipo_documento = p_tipo_documento;

		RETURN v_nuevo_numero;
	END;
	$$;
	`
	_, err := db.Exec(query)
	return err
}
