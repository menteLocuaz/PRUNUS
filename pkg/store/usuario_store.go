package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
)

// StoreUsuario interfaz que define las operaciones de acceso a datos para usuario
type StoreUsuario interface {
	GetAllUsuarios() ([]*models.Usuario, error)
	GetUsuarioByID(id uuid.UUID) (*models.Usuario, error)
	GetUsuarioByEmail(email string) (*models.Usuario, error)
	CreateUsuario(usuario *models.Usuario) (*models.Usuario, error)
	UpdateUsuario(id uuid.UUID, usuario *models.Usuario) (*models.Usuario, error)
	DeleteUsuario(id uuid.UUID) error
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
func (s *storeUsuario) GetAllUsuarios() ([]*models.Usuario, error) {
	query := `
		SELECT
			u.id_usuario,
			u.id_sucursal,
			u.id_rol,
			u.email,
			u.usu_nombre,
			u.usu_dni,
			u.usu_telefono,
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

	rows, err := s.db.Query(query)
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
			&usuario.Email,
			&usuario.UsuNombre,
			&usuario.UsuDNI,
			&usuario.UsuTelefono,
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
func (s *storeUsuario) GetUsuarioByID(id uuid.UUID) (*models.Usuario, error) {
	query := `
		SELECT
			u.id_usuario,
			u.id_sucursal,
			u.id_rol,
			u.email,
			u.usu_nombre,
			u.usu_dni,
			u.usu_telefono,
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

	err := s.db.QueryRow(query, id).Scan(
		// Usuario
		&usuario.IDUsuario,
		&usuario.IDSucursal,
		&usuario.IDRol,
		&usuario.Email,
		&usuario.UsuNombre,
		&usuario.UsuDNI,
		&usuario.UsuTelefono,
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
func (s *storeUsuario) GetUsuarioByEmail(email string) (*models.Usuario, error) {
	query := `
		SELECT
			u.id_usuario,
			u.id_sucursal,
			u.id_rol,
			u.email,
			u.usu_nombre,
			u.usu_dni,
			u.usu_telefono,
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

	err := s.db.QueryRow(query, email).Scan(
		// Usuario
		&usuario.IDUsuario,
		&usuario.IDSucursal,
		&usuario.IDRol,
		&usuario.Email,
		&usuario.UsuNombre,
		&usuario.UsuDNI,
		&usuario.UsuTelefono,
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
func (s *storeUsuario) CreateUsuario(usuario *models.Usuario) (*models.Usuario, error) {
	query := `
		INSERT INTO usuario (id_sucursal, id_rol, email, usu_nombre, usu_dni, usu_telefono, password, id_status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id_usuario, created_at, updated_at
	`

	err := s.db.QueryRow(
		query,
		usuario.IDSucursal,
		usuario.IDRol,
		usuario.Email,
		usuario.UsuNombre,
		usuario.UsuDNI,
		usuario.UsuTelefono,
		usuario.Password,
		usuario.IDStatus,
	).Scan(&usuario.IDUsuario, &usuario.CreatedAt, &usuario.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("error al crear usuario: %w", err)
	}

	return usuario, nil
}

// UpdateUsuario actualiza un usuario existente en la base de datos
func (s *storeUsuario) UpdateUsuario(id uuid.UUID, usuario *models.Usuario) (*models.Usuario, error) {
	query := `
		UPDATE usuario
		SET
			id_sucursal = $1,
			id_rol = $2,
			email = $3,
			usu_nombre = $4,
			usu_dni = $5,
			usu_telefono = $6,
			password = $7,
			id_status = $8,
			updated_at = CURRENT_TIMESTAMP
		WHERE id_usuario = $9 AND deleted_at IS NULL
		RETURNING
			id_usuario,
			id_sucursal,
			id_rol,
			email,
			usu_nombre,
			usu_dni,
			usu_telefono,
			password,
			id_status,
			created_at,
			updated_at
	`

	err := s.db.QueryRow(
		query,
		usuario.IDSucursal,
		usuario.IDRol,
		usuario.Email,
		usuario.UsuNombre,
		usuario.UsuDNI,
		usuario.UsuTelefono,
		usuario.Password,
		usuario.IDStatus,
		id,
	).Scan(
		&usuario.IDUsuario,
		&usuario.IDSucursal,
		&usuario.IDRol,
		&usuario.Email,
		&usuario.UsuNombre,
		&usuario.UsuDNI,
		&usuario.UsuTelefono,
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
func (s *storeUsuario) DeleteUsuario(id uuid.UUID) error {
	query := `
		UPDATE usuario
		SET deleted_at = $1
		WHERE id_usuario = $2 AND deleted_at IS NULL
	`

	result, err := s.db.Exec(query, time.Now(), id)
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
