package migrations

import "database/sql"

func migrateModulos(db *sql.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS modulo (
			id_modulo       UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
			mdl_id          SERIAL       NOT NULL,
			mdl_descripcion VARCHAR(150) NOT NULL,
			id_status       UUID         NOT NULL,
			replica         INTEGER,
			is_active       BOOLEAN      NOT NULL DEFAULT FALSE,
			abreviatura     VARCHAR(50),
			nivel           INTEGER      NOT NULL DEFAULT 0,

			created_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
			deleted_at      TIMESTAMP    NULL,

			CONSTRAINT fk_modulo_status
				FOREIGN KEY (id_status)
				REFERENCES estatus(id_status)
				ON UPDATE CASCADE
				ON DELETE RESTRICT
		);`,
		// Asegurar id_modulo si la tabla ya existía
		`ALTER TABLE modulo ADD COLUMN IF NOT EXISTS id_modulo UUID DEFAULT gen_random_uuid();`,
		// Si no tiene PK, ponérsela a id_modulo
		`DO $$ 
		BEGIN 
			IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints 
						  WHERE table_name='modulo' AND constraint_type='PRIMARY KEY') THEN
				ALTER TABLE modulo ADD PRIMARY KEY (id_modulo);
			END IF;
		END $$;`,
		`ALTER TABLE modulo ADD COLUMN IF NOT EXISTS mdl_id SERIAL;`,
		`ALTER TABLE modulo ADD COLUMN IF NOT EXISTS mdl_descripcion VARCHAR(150);`,
		`ALTER TABLE modulo ADD COLUMN IF NOT EXISTS id_status UUID;`,
		`ALTER TABLE modulo ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT FALSE;`,
		`ALTER TABLE modulo ADD COLUMN IF NOT EXISTS abreviatura VARCHAR(50);`,
		`ALTER TABLE modulo ADD COLUMN IF NOT EXISTS nivel INTEGER DEFAULT 0;`,
		`ALTER TABLE modulo ADD COLUMN IF NOT EXISTS created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;`,
		`ALTER TABLE modulo ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;`,
		`ALTER TABLE modulo ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP NULL;`,
		
		`CREATE INDEX IF NOT EXISTS idx_modulo_id_status ON modulo(id_status);`,
		`CREATE INDEX IF NOT EXISTS idx_modulo_deleted_at ON modulo(deleted_at);`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_modulo_mdl_id ON modulo(mdl_id);`,
	}

	for _, stmt := range statements {
		if _, err := db.Exec(stmt); err != nil {
			return err
		}
	}

	return seedModulos(db)
}

// seedModulos inserta los módulos predefinidos del sistema solo si la tabla
// está vacía, evitando duplicados en re-ejecuciones de la migración.
func seedModulos(db *sql.DB) error {
	// Verificar si ya existen registros para hacer el seed idempotente
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM modulo`).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil // Ya fue poblada; no insertar de nuevo
	}

	// id_status fijo: el UUID visible en la imagen (59039503-85...)
	// Se asume que este estatus ("Activo") ya existe en la tabla estatus.
	// Ajusta el UUID completo según tu tabla de estatus real.
	const idStatusActivo = "59039503-85CF-E511-80C1-000C29C9E0E0"

	seed := `
	INSERT INTO modulo (mdl_id, mdl_descripcion, id_status, replica, is_active, abreviatura, nivel)
	OVERRIDING SYSTEM VALUE
	VALUES
		(1, 'Tienda',   	 $1, 1, FALSE, 'Tienda',  0),
		(2, 'Anulación',     $1, 1, FALSE, NULL,   0),
		(4, 'Orden Pedido',  $1, 1, FALSE, NULL,   0),
		(5, 'Apertura',      $1, 1, FALSE, NULL,   0),
		(6, 'Estaciones',    $1, 1, FALSE, NULL,   0);`

	if _, err := db.Exec(seed, idStatusActivo); err != nil {
		return err
	}

	// Resincronizar la secuencia del SERIAL en una llamada separada
	setvalQuery := `
	SELECT setval(
		pg_get_serial_sequence('modulo', 'mdl_id'),
		(SELECT MAX(mdl_id) FROM modulo)
	);`

	_, err = db.Exec(setvalQuery)
	return err
}
