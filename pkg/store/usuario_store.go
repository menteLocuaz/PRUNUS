package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/prunus/pkg/models"
)

// StoreUsuario interfaz que define las operaciones de acceso a datos para usuario
type StoreUsuario interface {
	GetAllUsuarios() ([]*models.Usuario, error)
	GetUsuarioByID(id uint) (*models.Usuario, error)
	CreateUsuario(usuario *models.Usuario) (*models.Usuario, error)
	UpdateUsuario(id uint, usuario *models.Usuario) (*models.Usuario, error)
	DeleteUsuario(id uint) error
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
			u.email,
			u.usu_nombre,
			u.usu_dni,
			u.usu_telefono,
			u.password,
			u.estado,
			u.created_at,
			u.updated_at,
			u.deleted_at,
			r.id_rol,
			r.nombre_rol,
			r.id_sucursal,
			r.estado
		FROM usuario u
		LEFT JOIN rol r ON u.id_rol = r.id_rol
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
		var usuario models.Usuario
		var rol models.Rol

		err := rows.Scan(
			&usuario.IDUsuario,
			&usuario.IDSucursal,
			&usuario.UsuEmail,
			&usuario.UsuNombre,
			&usuario.UsuDni,
			&usuario.UsuTelefono,
			&usuario.UsuPassword,
			&usuario.Estado,
			&usuario.CreatedAt,
			&usuario.UpdatedAt,
			&usuario.DeletedAt,
			&rol.IDRol,
			&rol.RolNombre,
			&rol.IDSucursal,
			&rol.Estado,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear usuario: %w", err)
		}

		usuario.Rol = &rol
		usuarios = append(usuarios, &usuario)
	}

	return usuarios, nil
}

// GetUsuarioByID obtiene un usuario por su ID con su rol
func (s *storeUsuario) GetUsuarioByID(id uint) (*models.Usuario, error) {
	query := `
		SELECT
			u.id_usuario,
			u.id_sucursal,
			u.email,
			u.usu_nombre,
			u.usu_dni,
			u.usu_telefono,
			u.password,
			u.estado,
			u.created_at,
			u.updated_at,
			u.deleted_at,
			r.id_rol,
			r.nombre_rol,
			r.id_sucursal,
			r.estado
		FROM usuario u
		LEFT JOIN rol r ON u.id_rol = r.id_rol
		WHERE u.id_usuario = $1 AND u.deleted_at IS NULL
	`

	var usuario models.Usuario
	var rol models.Rol

	err := s.db.QueryRow(query, id).Scan(
		&usuario.IDUsuario,
		&usuario.IDSucursal,
		&usuario.UsuEmail,
		&usuario.UsuNombre,
		&usuario.UsuDni,
		&usuario.UsuTelefono,
		&usuario.UsuPassword,
		&usuario.Estado,
		&usuario.CreatedAt,
		&usuario.UpdatedAt,
		&usuario.DeletedAt,
		&rol.IDRol,
		&rol.RolNombre,
		&rol.IDSucursal,
		&rol.Estado,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("usuario con ID %d no encontrado", id)
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener usuario: %w", err)
	}

	usuario.Rol = &rol
	return &usuario, nil
}

// CreateUsuario crea un nuevo usuario en la base de datos
func (s *storeUsuario) CreateUsuario(usuario *models.Usuario) (*models.Usuario, error) {
	query := `
		INSERT INTO usuario (id_sucursal, id_rol, email, usu_nombre, usu_dni, usu_telefono, password, estado)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id_usuario, created_at, updated_at
	`

	// Obtener id_rol del objeto Rol si existe
	var idRol uint
	if usuario.Rol != nil {
		idRol = usuario.Rol.IDRol
	}

	err := s.db.QueryRow(
		query,
		usuario.IDSucursal,
		idRol,
		usuario.UsuEmail,
		usuario.UsuNombre,
		usuario.UsuDni,
		usuario.UsuTelefono,
		usuario.UsuPassword,
		usuario.Estado,
	).Scan(&usuario.IDUsuario, &usuario.CreatedAt, &usuario.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("error al crear usuario: %w", err)
	}

	return usuario, nil
}

// UpdateUsuario actualiza un usuario existente en la base de datos
func (s *storeUsuario) UpdateUsuario(id uint, usuario *models.Usuario) (*models.Usuario, error) {
	query := `
		UPDATE usuario
		SET id_sucursal = $1, id_rol = $2, email = $3, usu_nombre = $4,
		    usu_dni = $5, usu_telefono = $6, password = $7, estado = $8
		WHERE id_usuario = $9 AND deleted_at IS NULL
		RETURNING id_usuario, id_sucursal, email, usu_nombre, usu_dni,
		          usu_telefono, password, estado, created_at, updated_at
	`

	// Obtener id_rol del objeto Rol si existe
	var idRol uint
	if usuario.Rol != nil {
		idRol = usuario.Rol.IDRol
	}

	err := s.db.QueryRow(
		query,
		usuario.IDSucursal,
		idRol,
		usuario.UsuEmail,
		usuario.UsuNombre,
		usuario.UsuDni,
		usuario.UsuTelefono,
		usuario.UsuPassword,
		usuario.Estado,
		id,
	).Scan(
		&usuario.IDUsuario,
		&usuario.IDSucursal,
		&usuario.UsuEmail,
		&usuario.UsuNombre,
		&usuario.UsuDni,
		&usuario.UsuTelefono,
		&usuario.UsuPassword,
		&usuario.Estado,
		&usuario.CreatedAt,
		&usuario.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("usuario con ID %d no encontrado", id)
	}
	if err != nil {
		return nil, fmt.Errorf("error al actualizar usuario: %w", err)
	}

	return usuario, nil
}

// DeleteUsuario realiza un soft delete del usuario (actualiza deleted_at)
func (s *storeUsuario) DeleteUsuario(id uint) error {
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
		return fmt.Errorf("usuario con ID %d no encontrado", id)
	}

	return nil
}
