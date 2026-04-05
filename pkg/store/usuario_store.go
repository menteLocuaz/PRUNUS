package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/utils/performance"
)

// StoreUsuario interfaz que define las operaciones de acceso a datos para usuario
type StoreUsuario interface {
	GetAllUsuarios(ctx context.Context) ([]*models.Usuario, error)
	GetUsuarioByID(ctx context.Context, id uuid.UUID) (*models.Usuario, error)
	GetUsuarioByEmail(ctx context.Context, email string) (*models.Usuario, error)
	CreateUsuario(ctx context.Context, usuario *models.Usuario) (*models.Usuario, error)
	UpdateUsuario(ctx context.Context, id uuid.UUID, usuario *models.Usuario) (*models.Usuario, error)
	DeleteUsuario(ctx context.Context, id uuid.UUID) error

	// Accesos Multi-Sucursal
	AssignSucursales(ctx context.Context, userID uuid.UUID, sucursales []uuid.UUID) error
	GetSucursalesAcceso(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
}

// storeUsuario implementación de la interfaz StoreUsuario
type storeUsuario struct {
	db *sql.DB
}

// NewUsuario crea una nueva instancia del store de usuario
func NewUsuario(db *sql.DB) StoreUsuario {
	return &storeUsuario{db: db}
}

// GetAllUsuarios obtiene todos los usuarios activos (no eliminados) con su rol
func (s *storeUsuario) GetAllUsuarios(ctx context.Context) ([]*models.Usuario, error) {
	defer performance.Trace(ctx, "store", "GetAllUsuarios", performance.DbThreshold, time.Now())
	query := `
		SELECT
			u.id_usuario,
			u.id_sucursal,
			u.id_rol,
			u.username,
			u.email,
			u.usu_nombre,
			u.usu_dni,
			COALESCE(u.usu_telefono, ''),
			COALESCE(u.usu_tarjeta_nfc, ''),
			COALESCE(u.usu_pin_pos, ''),
			COALESCE(u.nombre_ticket, ''),
			u.password,
			u.id_status,
			u.created_at,
			u.updated_at,
			u.deleted_at,
			
			r.id_rol,
			r.nombre_rol,
			r.id_status,

			su.id_sucursal,
			su.nombre_sucursal,
			su.id_status
		FROM usuario u
		LEFT JOIN rol r ON u.id_rol = r.id_rol
		LEFT JOIN sucursal su ON su.id_sucursal = u.id_sucursal
		WHERE u.deleted_at IS NULL
		ORDER BY u.created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener usuarios: %w", err)
	}
	defer rows.Close()

	var usuarios []*models.Usuario
	for rows.Next() {
		usuario := &models.Usuario{
			Rol:      &models.Rol{},
			Sucursal: &models.Sucursal{},
		}

		err := rows.Scan(
			// Usuario
			&usuario.IDUsuario,
			&usuario.IDSucursal,
			&usuario.IDRol,
			&usuario.Username,
			&usuario.Email,
			&usuario.UsuNombre,
			&usuario.UsuDNI,
			&usuario.UsuTelefono,
			&usuario.UsuTarjetaNFC,
			&usuario.UsuPinPOS,
			&usuario.NombreTicket,
			&usuario.Password,
			&usuario.IDStatus,
			&usuario.CreatedAt,
			&usuario.UpdatedAt,
			&usuario.DeletedAt,

			// Rol
			&usuario.Rol.IDRol,
			&usuario.Rol.RolNombre,
			&usuario.Rol.IDStatus,

			// Sucursal
			&usuario.Sucursal.IDSucursal,
			&usuario.Sucursal.NombreSucursal,
			&usuario.Sucursal.IDStatus,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear usuario: %w", err)
		}

		usuarios = append(usuarios, usuario)
	}

	return usuarios, nil
}

// GetUsuarioByID obtiene un usuario por su ID con su rol
func (s *storeUsuario) GetUsuarioByID(ctx context.Context, id uuid.UUID) (*models.Usuario, error) {
	defer performance.Trace(ctx, "store", "GetUsuarioByID", performance.DbThreshold, time.Now())
	query := `
		SELECT
			u.id_usuario,
			u.id_sucursal,
			u.id_rol,
			u.username,
			u.email,
			u.usu_nombre,
			u.usu_dni,
			COALESCE(u.usu_telefono, ''),
			COALESCE(u.usu_tarjeta_nfc, ''),
			COALESCE(u.usu_pin_pos, ''),
			COALESCE(u.nombre_ticket, ''),
			u.password,
			u.id_status,
			u.created_at,
			u.updated_at,
			u.deleted_at,

			r.id_rol,
			r.nombre_rol,
			r.id_status,

			su.id_sucursal,
			su.nombre_sucursal,
			su.id_status
		FROM usuario u
		LEFT JOIN rol r ON u.id_rol = r.id_rol
		LEFT JOIN sucursal su ON su.id_sucursal = u.id_sucursal
		WHERE u.id_usuario = $1 AND u.deleted_at IS NULL
	`

	usuario := &models.Usuario{
		Rol:      &models.Rol{},
		Sucursal: &models.Sucursal{},
	}

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		// Usuario
		&usuario.IDUsuario,
		&usuario.IDSucursal,
		&usuario.IDRol,
		&usuario.Username,
		&usuario.Email,
		&usuario.UsuNombre,
		&usuario.UsuDNI,
		&usuario.UsuTelefono,
		&usuario.UsuTarjetaNFC,
		&usuario.UsuPinPOS,
		&usuario.NombreTicket,
		&usuario.Password,
		&usuario.IDStatus,
		&usuario.CreatedAt,
		&usuario.UpdatedAt,
		&usuario.DeletedAt,

		// Rol
		&usuario.Rol.IDRol,
		&usuario.Rol.RolNombre,
		&usuario.Rol.IDStatus,

		// Sucursal
		&usuario.Sucursal.IDSucursal,
		&usuario.Sucursal.NombreSucursal,
		&usuario.Sucursal.IDStatus,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("usuario con ID %s no encontrado", id)
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener usuario: %w", err)
	}

	return usuario, nil
}

// GetUsuarioByEmail obtiene un usuario por su email con su rol
// Este método incluye el password hasheado para validación de autenticación
func (s *storeUsuario) GetUsuarioByEmail(ctx context.Context, email string) (*models.Usuario, error) {
	defer performance.Trace(ctx, "store", "GetUsuarioByEmail", performance.DbThreshold, time.Now())
	query := `
		SELECT
			u.id_usuario,
			u.id_sucursal,
			u.id_rol,
			u.username,
			u.email,
			u.usu_nombre,
			u.usu_dni,
			COALESCE(u.usu_telefono, ''),
			COALESCE(u.usu_tarjeta_nfc, ''),
			COALESCE(u.usu_pin_pos, ''),
			COALESCE(u.nombre_ticket, ''),
			u.password,
			u.id_status,
			u.created_at,
			u.updated_at,
			u.deleted_at,

			r.id_rol,
			r.nombre_rol,
			r.id_status,

			su.id_sucursal,
			su.nombre_sucursal,
			su.id_status
		FROM usuario u
		LEFT JOIN rol r ON u.id_rol = r.id_rol
		LEFT JOIN sucursal su ON su.id_sucursal = u.id_sucursal
		WHERE u.email = $1 AND u.deleted_at IS NULL
	`

	usuario := &models.Usuario{
		Rol:      &models.Rol{},
		Sucursal: &models.Sucursal{},
	}

	err := s.db.QueryRowContext(ctx, query, email).Scan(
		// Usuario
		&usuario.IDUsuario,
		&usuario.IDSucursal,
		&usuario.IDRol,
		&usuario.Username,
		&usuario.Email,
		&usuario.UsuNombre,
		&usuario.UsuDNI,
		&usuario.UsuTelefono,
		&usuario.UsuTarjetaNFC,
		&usuario.UsuPinPOS,
		&usuario.NombreTicket,
		&usuario.Password,
		&usuario.IDStatus,
		&usuario.CreatedAt,
		&usuario.UpdatedAt,
		&usuario.DeletedAt,

		// Rol
		&usuario.Rol.IDRol,
		&usuario.Rol.RolNombre,
		&usuario.Rol.IDStatus,

		// Sucursal
		&usuario.Sucursal.IDSucursal,
		&usuario.Sucursal.NombreSucursal,
		&usuario.Sucursal.IDStatus,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("usuario con email %s no encontrado", email)
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener usuario por email: %w", err)
	}

	return usuario, nil
}

// CreateUsuario crea un nuevo usuario en la base de datos
func (s *storeUsuario) CreateUsuario(ctx context.Context, usuario *models.Usuario) (*models.Usuario, error) {
	defer performance.Trace(ctx, "store", "CreateUsuario", performance.DbThreshold, time.Now())
	query := `
		INSERT INTO usuario (id_sucursal, id_rol, username, email, usu_nombre, usu_dni, usu_telefono, usu_tarjeta_nfc, usu_pin_pos, nombre_ticket, password, id_status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id_usuario, created_at, updated_at
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		usuario.IDSucursal,
		usuario.IDRol,
		usuario.Username,
		usuario.Email,
		usuario.UsuNombre,
		usuario.UsuDNI,
		usuario.UsuTelefono,
		usuario.UsuTarjetaNFC,
		usuario.UsuPinPOS,
		usuario.NombreTicket,
		usuario.Password,
		usuario.IDStatus,
	).Scan(&usuario.IDUsuario, &usuario.CreatedAt, &usuario.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("error al crear usuario: %w", err)
	}

	return usuario, nil
}

// UpdateUsuario actualiza un usuario existente en la base de datos
func (s *storeUsuario) UpdateUsuario(ctx context.Context, id uuid.UUID, usuario *models.Usuario) (*models.Usuario, error) {
	defer performance.Trace(ctx, "store", "UpdateUsuario", performance.DbThreshold, time.Now())
	query := `
		UPDATE usuario
		SET
			id_sucursal = $1,
			id_rol = $2,
			username = $3,
			email = $4,
			usu_nombre = $5,
			usu_dni = $6,
			usu_telefono = $7,
			usu_tarjeta_nfc = $8,
			usu_pin_pos = $9,
			nombre_ticket = $10,
			password = CASE WHEN $11 <> '' THEN $11 ELSE password END,
			id_status = $12,
			updated_at = CURRENT_TIMESTAMP
		WHERE id_usuario = $13 AND deleted_at IS NULL
		RETURNING
			id_usuario,
			id_sucursal,
			id_rol,
			username,
			email,
			usu_nombre,
			usu_dni,
			COALESCE(usu_telefono, ''),
			COALESCE(usu_tarjeta_nfc, ''),
			COALESCE(usu_pin_pos, ''),
			COALESCE(nombre_ticket, ''),
			password,
			id_status,
			created_at,
			updated_at
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		usuario.IDSucursal,
		usuario.IDRol,
		usuario.Username,
		usuario.Email,
		usuario.UsuNombre,
		usuario.UsuDNI,
		usuario.UsuTelefono,
		usuario.UsuTarjetaNFC,
		usuario.UsuPinPOS,
		usuario.NombreTicket,
		usuario.Password,
		usuario.IDStatus,
		id,
	).Scan(
		&usuario.IDUsuario,
		&usuario.IDSucursal,
		&usuario.IDRol,
		&usuario.Username,
		&usuario.Email,
		&usuario.UsuNombre,
		&usuario.UsuDNI,
		&usuario.UsuTelefono,
		&usuario.UsuTarjetaNFC,
		&usuario.UsuPinPOS,
		&usuario.NombreTicket,
		&usuario.Password,
		&usuario.IDStatus,
		&usuario.CreatedAt,
		&usuario.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("usuario con ID %s no encontrado", id)
	}
	if err != nil {
		return nil, fmt.Errorf("error al actualizar usuario: %w", err)
	}

	return usuario, nil
}

// DeleteUsuario realiza un soft delete del usuario (actualiza deleted_at)
func (s *storeUsuario) DeleteUsuario(ctx context.Context, id uuid.UUID) error {
	defer performance.Trace(ctx, "store", "DeleteUsuario", performance.DbThreshold, time.Now())
	query := `
		UPDATE usuario
		SET deleted_at = $1
		WHERE id_usuario = $2 AND deleted_at IS NULL
	`

	result, err := s.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("error al eliminar usuario: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("usuario con ID %s no encontrado", id)
	}

	return nil
}

func (s *storeUsuario) AssignSucursales(ctx context.Context, userID uuid.UUID, sucursales []uuid.UUID) error {
	defer performance.Trace(ctx, "store", "AssignSucursales", performance.DbThreshold, time.Now())
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Limpiar accesos previos
	queryDelete := `DELETE FROM usuario_sucursal_acceso WHERE id_usuario = $1`
	if _, err := tx.ExecContext(ctx, queryDelete, userID); err != nil {
		return fmt.Errorf("error al eliminar accesos previos: %w", err)
	}

	// 2. Insertar nuevos accesos
	queryInsert := `INSERT INTO usuario_sucursal_acceso (id_usuario, id_sucursal) VALUES ($1, $2)`
	for _, sucursalID := range sucursales {
		if _, err := tx.ExecContext(ctx, queryInsert, userID, sucursalID); err != nil {
			return fmt.Errorf("error al insertar acceso a sucursal %s: %w", sucursalID, err)
		}
	}

	return tx.Commit()
}

func (s *storeUsuario) GetSucursalesAcceso(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	defer performance.Trace(ctx, "store", "GetSucursalesAcceso", performance.DbThreshold, time.Now())
	query := `SELECT id_sucursal FROM usuario_sucursal_acceso WHERE id_usuario = $1`
	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sucursales []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		sucursales = append(sucursales, id)
	}
	return sucursales, nil
}
