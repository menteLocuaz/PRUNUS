package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/prunus/pkg/models"
)

// StoreRol interfaz que define las operaciones de acceso a datos para rol
type StoreRol interface {
	GetAllRoles() ([]*models.Rol, error)
	GetRolByID(id uint) (*models.Rol, error)
	CreateRol(rol *models.Rol) (*models.Rol, error)
	UpdateRol(id uint, rol *models.Rol) (*models.Rol, error)
	DeleteRol(id uint) error
}

// storeRol implementación de la interfaz StoreRol
type storeRol struct {
	db *sql.DB
}

// NewRol crea una nueva instancia del store de rol
func NewRol(db *sql.DB) StoreRol {
	return &storeRol{db: db}
}

// GetAllRoles obtiene todos los roles activos (no eliminados) de la base de datos
func (s *storeRol) GetAllRoles() ([]*models.Rol, error) {
	query := `
	SELECT
		r.id_rol,
		r.nombre_rol,
		r.id_sucursal,
		r.estado,

		s.id_sucursal,
		s.nombre_sucursal,
		s.estado
	FROM rol r
	JOIN sucursal s ON s.id_sucursal = r.id_sucursal
	WHERE r.deleted_at IS NULL
	  AND s.deleted_at IS NULL
	ORDER BY r.created_at DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener roles: %w", err)
	}
	defer rows.Close()

	var roles []*models.Rol

	for rows.Next() {
		rol := &models.Rol{
			Sucursal: &models.Sucursal{},
		}

		err := rows.Scan(
			&rol.IDRol,
			&rol.RolNombre,
			&rol.IDSucursal,
			&rol.Estado,

			&rol.Sucursal.IDSucursal,
			&rol.Sucursal.NombreSucursal,
			&rol.Sucursal.Estado,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear rol: %w", err)
		}

		roles = append(roles, rol)
	}

	return roles, nil
}

// GetRolByID obtiene un rol por su ID
func (s *storeRol) GetRolByID(id uint) (*models.Rol, error) {
	query := `
		SELECT
			r.id_rol,
			r.nombre_rol,
			r.id_sucursal,
			r.estado,
			
			s.id_sucursal,
    		s.nombre_sucursal,
			s.estado
		FROM rol r
		JOIN sucursal s ON s.id_sucursal = r.id_sucursal
		WHERE r.id_rol = $1 AND r.deleted_at IS NULL
	`

	rol := models.Rol{
		Sucursal: &models.Sucursal{},
	}
	err := s.db.QueryRow(query, id).Scan(
		&rol.IDRol,
		&rol.RolNombre,
		&rol.IDSucursal,
		&rol.Estado,

		&rol.Sucursal.IDSucursal,
		&rol.Sucursal.NombreSucursal,
		&rol.Sucursal.Estado,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("rol con ID %d no encontrado", id)
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener rol: %w", err)
	}

	return &rol, nil
}

// CreateRol crea un nuevo rol en la base de datos
func (s *storeRol) CreateRol(rol *models.Rol) (*models.Rol, error) {
	query := `
		INSERT INTO rol (nombre_rol, id_sucursal, estado)
		VALUES ($1, $2, $3)
		RETURNING id_rol, created_at, updated_at
	`

	err := s.db.QueryRow(
		query,
		rol.RolNombre,
		rol.IDSucursal,
		rol.Estado,
	).Scan(&rol.IDRol, &rol.CreatedAt, &rol.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("error al crear rol: %w", err)
	}

	return rol, nil
}

// UpdateRol actualiza un rol existente en la base de datos
func (s *storeRol) UpdateRol(id uint, rol *models.Rol) (*models.Rol, error) {
	query := `
		UPDATE rol
		SET nombre_rol = $1, id_sucursal = $2, estado = $3
		WHERE id_rol = $4 AND deleted_at IS NULL
		RETURNING id_rol, nombre_rol, id_sucursal, estado, created_at, updated_at
	`

	err := s.db.QueryRow(
		query,
		rol.RolNombre,
		rol.IDSucursal,
		rol.Estado,
		id,
	).Scan(
		&rol.IDRol,
		&rol.RolNombre,
		&rol.IDSucursal,
		&rol.Estado,
		&rol.CreatedAt,
		&rol.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("rol con ID %d no encontrado", id)
	}
	if err != nil {
		return nil, fmt.Errorf("error al actualizar rol: %w", err)
	}

	return rol, nil
}

// DeleteRol realiza un soft delete del rol (actualiza deleted_at)
func (s *storeRol) DeleteRol(id uint) error {
	query := `
		UPDATE rol
		SET deleted_at = $1
		WHERE id_rol = $2 AND deleted_at IS NULL
	`

	result, err := s.db.Exec(query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("error al eliminar rol: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("rol con ID %d no encontrado", id)
	}

	return nil
}
