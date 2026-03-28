package migrations

import "database/sql"

func addTPEnvIDToRetiros(db *sql.DB) error {
	query := `ALTER TABLE retiros ADD COLUMN IF NOT EXISTS tpenv_id INTEGER NOT NULL DEFAULT -1;`
	_, err := db.Exec(query)
	return err
}
