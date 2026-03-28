package migrations

import "database/sql"

func addMotivoDescuadreToControlEstacion(db *sql.DB) error {
	query := `ALTER TABLE control_estacion ADD COLUMN IF NOT EXISTS ctrc_motivo_descuadre TEXT;`
	_, err := db.Exec(query)
	return err
}
