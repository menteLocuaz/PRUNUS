package migrations

import "database/sql"

func migrateAddUsernameToUsuario(db *sql.DB) error {
	query := `
	ALTER TABLE usuario ADD COLUMN IF NOT EXISTS username VARCHAR(50);
	UPDATE usuario SET username = SPLIT_PART(email, '@', 1) WHERE username IS NULL;
	ALTER TABLE usuario ALTER COLUMN username SET NOT NULL;
	ALTER TABLE usuario ADD CONSTRAINT uk_usuario_username UNIQUE (username);
	`
	_, err := db.Exec(query)
	return err
}
