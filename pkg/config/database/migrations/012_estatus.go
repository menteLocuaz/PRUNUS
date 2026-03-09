package migrations

import "database/sql"

func migrateEstatus(db *sql.DB) error {
	// Habilitar extensión pgcrypto para gen_random_uuid()
	_, err := db.Exec(`CREATE EXTENSION IF NOT EXISTS pgcrypto;`)
	if err != nil {
		return err
	}

	query := `
	CREATE TABLE IF NOT EXISTS estatus (
		id_status       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		std_descripcion VARCHAR(255) NOT NULL,
		stp_tipo_estado VARCHAR(255) NOT NULL,
		mdl_id          INTEGER      NOT NULL,

		created_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at      TIMESTAMP    NULL
	);

	CREATE INDEX IF NOT EXISTS idx_estatus_tipo       ON estatus(stp_tipo_estado);
	CREATE INDEX IF NOT EXISTS idx_estatus_modulo     ON estatus(mdl_id);
	CREATE INDEX IF NOT EXISTS idx_estatus_deleted_at ON estatus(deleted_at);

	-- Insertar valores predeterminados si no existen
	INSERT INTO estatus (id_status, std_descripcion, stp_tipo_estado, mdl_id) VALUES
		('59039503-85CF-E511-80C1-000C29C9E0E0', 'Activo', '1', 8),
		('5A039503-85CF-E511-80C1-000C29C9E0E0', 'Inactivo', '1', 8),
		('99039503-85CF-E511-80C1-000C29C9E0E0', 'Fondo Asignado', '1', 8),
		('9A039503-85CF-E511-80C1-000C29C9E0E0', 'Desmontado', '1', 8),
		('9B039503-85CF-E511-80C1-000C29C9E0E0', 'Ingreso Admin', '1', 8),
		('9C039503-85CF-E511-80C1-000C29C9E0E0', 'Salir Admin', '1', 8),
		('0D4515FE-C907-E611-A6B8-000C29C9E0E0', 'Arqueo', '1', 8),
		('159E3FE6-630E-E611-80C1-000C29C9E0E0', 'Arqueo Retiros', '1', 8),
		('E8297CFA-630E-E611-80C1-000C29C9E0E0', 'Retiro Efectivo', '1', 8),
		('84920103-640E-E611-80C1-000C29C9E0E0', 'Retiro Total', '1', 8),
		('A864475F-0D34-E711-80C1-000C29C9E0E0', 'Fondo Activo', '1', 8),
		('2160B065-0D34-E711-80C1-000C29C9E0E0', 'Fondo Retirado', '1', 8),
		('5E8DD0FB-5550-E711-80C1-000C29C9E0E0', 'Fondo Por Confirmar', '1', 8)
	ON CONFLICT (id_status) DO NOTHING;
	`

	_, err = db.Exec(query)
	return err
}
