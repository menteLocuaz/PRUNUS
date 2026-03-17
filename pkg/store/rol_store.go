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

// StoreRol interfaz que define las operaciones de acceso a datos para rol
type StoreRol interface {
	GetAllRoles(ctx context.Context) ([]*models.Rol, error)
	GetRolByID(ctx context.Context, id uuid.UUID) (*models.Rol, error)
	CreateRol(ctx context.Context, rol *models.Rol) (*models.Rol, error)
	UpdateRol(ctx context.Context, id uuid.UUID, rol *models.Rol) (*models.Rol, error)
	DeleteRol(ctx context.Context, id uuid.UUID) error
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
func (s *storeRol) GetAllRoles(ctx context.Context) ([]*models.Rol, error) {
	defer performance.Trace(ctx, "store", "GetAllRoles", performance.DbThreshold, time.Now())
	query := `
	SELECT
		r.id_rol,
		r.nombre_rol,
		r.id_sucursal,
		r.id_status,

		su.id_sucursal,
		su.nombre_sucursal,
		su.id_status
	FROM rol r
	JOIN sucursal su ON su.id_sucursal = r.id_sucursal
	WHERE r.deleted_at IS NULL
	  AND su.deleted_at IS NULL
	ORDER BY r.created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query)
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
			&rol.IDStatus,

			&rol.Sucursal.IDSucursal,
			&rol.Sucursal.NombreSucursal,
			&rol.Sucursal.IDStatus,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear rol: %w", err)
		}

		roles = append(roles, rol)
	}

	return roles, nil
}

// GetRolByID obtiene un rol por su ID
func (s *storeRol) GetRolByID(ctx context.Context, id uuid.UUID) (*models.Rol, error) {
	defer performance.Trace(ctx, "store", "GetRolByID", performance.DbThreshold, time.Now())
	query := `
		SELECT
			r.id_rol,
			r.nombre_rol,
			r.id_sucursal,
			r.id_status,
			
			su.id_sucursal,
    		su.nombre_sucursal,
			su.id_status
		FROM rol r
		JOIN sucursal su ON su.id_sucursal = r.id_sucursal
		WHERE r.id_rol = $1 AND r.deleted_at IS NULL
	`

	rol := models.Rol{
		Sucursal: &models.Sucursal{},
	}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&rol.IDRol,
		&rol.RolNombre,
		&rol.IDSucursal,
		&rol.IDStatus,

		&rol.Sucursal.IDSucursal,
		&rol.Sucursal.NombreSucursal,
		&rol.Sucursal.IDStatus,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("rol con ID %s no encontrado", id)
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener rol: %w", err)
	}

	return &rol, nil
}

// CreateRol crea un nuevo rol en la base de datos
func (s *storeRol) CreateRol(ctx context.Context, rol *models.Rol) (*models.Rol, error) {
	defer performance.Trace(ctx, "store", "CreateRol", performance.DbThreshold, time.Now())
	query := `
		INSERT INTO rol (nombre_rol, id_sucursal, id_status)
		VALUES ($1, $2, $3)
		RETURNING id_rol, created_at, updated_at
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		rol.RolNombre,
		rol.IDSucursal,
		rol.IDStatus,
	).Scan(&rol.IDRol, &rol.CreatedAt, &rol.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("error al crear rol: %w", err)
	}

	return rol, nil
}

// UpdateRol actualiza un rol existente en la base de datos
func (s *storeRol) UpdateRol(ctx context.Context, id uuid.UUID, rol *models.Rol) (*models.Rol, error) {
	defer performance.Trace(ctx, "store", "UpdateRol", performance.DbThreshold, time.Now())
	query := `
		UPDATE rol
		SET nombre_rol = $1, id_sucursal = $2, id_status = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id_rol = $4 AND deleted_at IS NULL
		RETURNING id_rol, nombre_rol, id_sucursal, id_status, created_at, updated_at
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		rol.RolNombre,
		rol.IDSucursal,
		rol.IDStatus,
		id,
	).Scan(
		&rol.IDRol,
		&rol.RolNombre,
		&rol.IDSucursal,
		&rol.IDStatus,
		&rol.CreatedAt,
		&rol.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("rol con ID %s no encontrado", id)
	}
	if err != nil {
		return nil, fmt.Errorf("error al actualizar rol: %w", err)
	}

	return rol, nil
}

// DeleteRol realiza un soft delete del rol (actualiza deleted_at)
func (s *storeRol) DeleteRol(ctx context.Context, id uuid.UUID) error {
	defer performance.Trace(ctx, "store", "DeleteRol", performance.DbThreshold, time.Now())
	query := `
		UPDATE rol
		SET deleted_at = $1
		WHERE id_rol = $2 AND deleted_at IS NULL
	`

	result, err := s.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("error al eliminar rol: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("rol con ID %s no encontrado", id)
	}

	return nil
}
