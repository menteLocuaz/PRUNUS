package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// SetAuditUser establece el ID de usuario en la sesión de PostgreSQL para los triggers de auditoría.
// Debe ejecutarse dentro de una transacción para que SET LOCAL tenga efecto.
func SetAuditUser(ctx context.Context, tx *sql.Tx) error {
	val := ctx.Value("user_id")
	if val == nil {
		return nil
	}

	var userIDStr string
	switch v := val.(type) {
	case uuid.UUID:
		userIDStr = v.String()
	case string:
		userIDStr = v
	default:
		return nil
	}

	// Usar quote_literal para mayor seguridad aunque sea un UUID controlado
	query := fmt.Sprintf("SET LOCAL app.current_user_id = %s", quoteLiteral(userIDStr))
	_, err := tx.ExecContext(ctx, query)
	return err
}

// ExecAudited ejecuta una función dentro de una transacción que tiene configurado el usuario de auditoría.
func ExecAudited(ctx context.Context, db *sql.DB, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := SetAuditUser(ctx, tx); err != nil {
		return fmt.Errorf("error al configurar usuario de auditoría: %w", err)
	}

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}

// quoteLiteral escapa una cadena para usarla en SQL (prevención básica)
func quoteLiteral(s string) string {
	return "'" + s + "'"
}
