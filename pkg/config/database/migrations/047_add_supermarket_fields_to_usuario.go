package migrations

import "database/sql"

func addSupermarketFieldsToUsuario(db *sql.DB) error {
	query := `
	ALTER TABLE usuario 
	ADD COLUMN IF NOT EXISTS usu_tarjeta_nfc VARCHAR(100) UNIQUE,
	ADD COLUMN IF NOT EXISTS usu_pin_pos     VARCHAR(100),
	ADD COLUMN IF NOT EXISTS nombre_ticket   VARCHAR(50);
	`
	_, err := db.Exec(query)
	return err
}
