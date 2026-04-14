package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/prunus/pkg/dto"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/utils/performance"
)

// StoreUsuario interfaz que define las operaciones de acceso a datos para usuario
type StoreUsuario interface {
	GetAllUsuarios(ctx context.Context, params dto.PaginationParams) ([]*models.Usuario, error)
	GetUsuarioByID(ctx context.Context, id uuid.UUID) (*models.Usuario, error)
	GetUsuarioByEmail(ctx context.Context, email string) (*models.Usuario, error)
	GetUsuarioByUsername(ctx context.Context, username string) (*models.Usuario, error)
	GetUsuarioByPin(ctx context.Context, pin string) (*models.Usuario, error)
	CreateUsuario(ctx context.Context, usuario *models.Usuario) (*models.Usuario, error)
	UpdateUsuario(ctx context.Context, id uuid.UUID, usuario *models.Usuario) (*models.Usuario, error)
	DeleteUsuario(ctx context.Context, id uuid.UUID) error

	// Accesos Multi-Sucursal
	AssignSucursales(ctx context.Context, userID uuid.UUID, sucursales []uuid.UUID) error
	GetSucursalesAcceso(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)

	// Permisos y Módulos
	GetPermisosByRol(ctx context.Context, rolID uuid.UUID) ([]string, error)

	// Turnos
	UpdateTurnoStatus(ctx context.Context, id uuid.UUID, enTurno bool) error
}

// storeUsuario implementación de la interfaz StoreUsuario
type storeUsuario struct {
	db *sql.DB
}

// Campos base para SELECT de usuario con sus joins de Rol y Sucursal.
// Se incluye su.id_empresa para que el JWT pueda transportar el tenant de empresa.
const usuarioSelectFields = `
	u.id_usuario, u.id_sucursal, u.id_rol, u.username, u.email, u.usu_nombre, u.usu_dni,
	COALESCE(u.usu_telefono, ''), COALESCE(u.usu_tarjeta_nfc, ''), COALESCE(u.usu_pin_pos, ''),
	COALESCE(u.nombre_ticket, ''), u.password, u.id_status, u.metadata, u.created_at, u.updated_at, u.deleted_at, u.en_turno,

	r.id_rol, r.nombre_rol, r.id_status,
	su.id_sucursal, su.nombre_sucursal, su.id_status, su.id_empresa
`

// scanRowUsuario es un helper centralizado para escanear las columnas definidas en usuarioSelectFields
func (s *storeUsuario) scanRowUsuario(scanner interface{ Scan(dest ...any) error }, u *models.Usuario) error {
	if u.Rol == nil {
		u.Rol = &models.Rol{}
	}
	if u.Sucursal == nil {
		u.Sucursal = &models.Sucursal{}
	}
	var metadataJSON []byte
	err := scanner.Scan(
		&u.IDUsuario, &u.IDSucursal, &u.IDRol, &u.Username, &u.Email, &u.UsuNombre, &u.UsuDNI,
		&u.UsuTelefono, &u.UsuTarjetaNFC, &u.UsuPinPOS, &u.NombreTicket, &u.Password, &u.IDStatus, &metadataJSON,
		&u.CreatedAt, &u.UpdatedAt, &u.DeletedAt, &u.EnTurno,
		&u.Rol.IDRol, &u.Rol.RolNombre, &u.Rol.IDStatus,
		&u.Sucursal.IDSucursal, &u.Sucursal.NombreSucursal, &u.Sucursal.IDStatus, &u.Sucursal.IDEmpresa,
	)
	if err != nil {
		return err
	}
	if len(metadataJSON) > 0 {
		json.Unmarshal(metadataJSON, &u.Metadata)
	}
	return nil
}

// NewUsuario crea una nueva instancia del store de usuario
func NewUsuario(db *sql.DB) StoreUsuario {
	return &storeUsuario{db: db}
}

// GetAllUsuarios obtiene todos los usuarios activos con soporte para paginación.
func (s *storeUsuario) GetAllUsuarios(ctx context.Context, params dto.PaginationParams) ([]*models.Usuario, error) {
	defer performance.Trace(ctx, "store", "GetAllUsuarios", performance.DbThreshold, time.Now())
	
	if params.Limit <= 0 {
		params.Limit = dto.DefaultLimit
	}

	query := fmt.Sprintf(`
		SELECT %s
		FROM usuario u
		LEFT JOIN rol r ON u.id_rol = r.id_rol
		LEFT JOIN sucursal su ON su.id_sucursal = u.id_sucursal
		WHERE u.deleted_at IS NULL
	`, usuarioSelectFields)

	var args []interface{}
	if params.LastDate != nil {
		query += " AND u.created_at < $1"
		args = append(args, params.LastDate)
	}

	query += " ORDER BY u.created_at DESC LIMIT $" + fmt.Sprint(len(args)+1)
	args = append(args, params.Limit)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error al obtener usuarios: %w", err)
	}
	defer rows.Close()

	var usuarios []*models.Usuario
	for rows.Next() {
		usuario := &models.Usuario{}
		if err := s.scanRowUsuario(rows, usuario); err != nil {
			return nil, fmt.Errorf("error al escanear usuario: %w", err)
		}
		usuarios = append(usuarios, usuario)
	}

	return usuarios, nil
}

// GetUsuarioByID obtiene un usuario por su ID con su rol
func (s *storeUsuario) GetUsuarioByID(ctx context.Context, id uuid.UUID) (*models.Usuario, error) {
	defer performance.Trace(ctx, "store", "GetUsuarioByID", performance.DbThreshold, time.Now())
	query := fmt.Sprintf(`
		SELECT %s
		FROM usuario u
		LEFT JOIN rol r ON u.id_rol = r.id_rol
		LEFT JOIN sucursal su ON su.id_sucursal = u.id_sucursal
		WHERE u.id_usuario = $1 AND u.deleted_at IS NULL
	`, usuarioSelectFields)

	usuario := &models.Usuario{}
	err := s.scanRowUsuario(s.db.QueryRowContext(ctx, query, id), usuario)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("usuario con ID %s no encontrado", id)
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener usuario: %w", err)
	}

	return usuario, nil
}

// GetUsuarioByEmail obtiene un usuario por su email con su rol
func (s *storeUsuario) GetUsuarioByEmail(ctx context.Context, email string) (*models.Usuario, error) {
	defer performance.Trace(ctx, "store", "GetUsuarioByEmail", performance.DbThreshold, time.Now())
	query := fmt.Sprintf(`
		SELECT %s
		FROM usuario u
		LEFT JOIN rol r ON u.id_rol = r.id_rol
		LEFT JOIN sucursal su ON su.id_sucursal = u.id_sucursal
		WHERE u.email = $1 AND u.deleted_at IS NULL
	`, usuarioSelectFields)

	usuario := &models.Usuario{}
	err := s.scanRowUsuario(s.db.QueryRowContext(ctx, query, email), usuario)

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

	err := ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		metadataJSON, _ := json.Marshal(usuario.Metadata)
		query := `
			INSERT INTO usuario (id_sucursal, id_rol, username, email, usu_nombre, usu_dni, usu_telefono, usu_tarjeta_nfc, usu_pin_pos, nombre_ticket, password, id_status, metadata)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
			RETURNING id_usuario, created_at, updated_at
		`
		return tx.QueryRowContext(
			ctx, query,
			usuario.IDSucursal, usuario.IDRol, usuario.Username, usuario.Email, usuario.UsuNombre,
			usuario.UsuDNI, usuario.UsuTelefono, usuario.UsuTarjetaNFC, usuario.UsuPinPOS,
			usuario.NombreTicket, usuario.Password, usuario.IDStatus, metadataJSON,
		).Scan(&usuario.IDUsuario, &usuario.CreatedAt, &usuario.UpdatedAt)
	})

	if err != nil {
		return nil, fmt.Errorf("error al crear usuario: %w", err)
	}

	return usuario, nil
}

// UpdateUsuario actualiza un usuario existente en la base de datos
func (s *storeUsuario) UpdateUsuario(ctx context.Context, id uuid.UUID, usuario *models.Usuario) (*models.Usuario, error) {
	defer performance.Trace(ctx, "store", "UpdateUsuario", performance.DbThreshold, time.Now())

	err := ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		metadataJSON, _ := json.Marshal(usuario.Metadata)
		query := `
			UPDATE usuario
			SET
				id_sucursal = $1, id_rol = $2, username = $3, email = $4, usu_nombre = $5,
				usu_dni = $6, usu_telefono = $7, usu_tarjeta_nfc = $8, usu_pin_pos = $9,
				nombre_ticket = $10, password = CASE WHEN $11 <> '' THEN $11 ELSE password END,
				id_status = $12, metadata = $13, updated_at = CURRENT_TIMESTAMP
			WHERE id_usuario = $14 AND deleted_at IS NULL
			RETURNING
				id_usuario, id_sucursal, id_rol, username, email, usu_nombre, usu_dni,
				COALESCE(usu_telefono, ''), COALESCE(usu_tarjeta_nfc, ''), COALESCE(usu_pin_pos, ''),
				COALESCE(nombre_ticket, ''), password, id_status, created_at, updated_at
		`
		// Nota: El RETURNING en el Scan original no incluía metadata, voy a ajustarlo para que sea consistente pero simplificar el retorno
		return tx.QueryRowContext(
			ctx, query,
			usuario.IDSucursal, usuario.IDRol, usuario.Username, usuario.Email, usuario.UsuNombre,
			usuario.UsuDNI, usuario.UsuTelefono, usuario.UsuTarjetaNFC, usuario.UsuPinPOS,
			usuario.NombreTicket, usuario.Password, usuario.IDStatus, metadataJSON, id,
		).Scan(
			&usuario.IDUsuario, &usuario.IDSucursal, &usuario.IDRol, &usuario.Username, &usuario.Email,
			&usuario.UsuNombre, &usuario.UsuDNI, &usuario.UsuTelefono, &usuario.UsuTarjetaNFC,
			&usuario.UsuPinPOS, &usuario.NombreTicket, &usuario.Password, &usuario.IDStatus,
			&usuario.CreatedAt, &usuario.UpdatedAt,
		)
	})

	if err != nil {
		return nil, fmt.Errorf("error al actualizar usuario: %w", err)
	}

	return usuario, nil
}

// DeleteUsuario realiza un soft delete del usuario
func (s *storeUsuario) DeleteUsuario(ctx context.Context, id uuid.UUID) error {
	defer performance.Trace(ctx, "store", "DeleteUsuario", performance.DbThreshold, time.Now())

	return ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		query := `UPDATE usuario SET deleted_at = $1 WHERE id_usuario = $2 AND deleted_at IS NULL`
		result, err := tx.ExecContext(ctx, query, time.Now(), id)
		if err != nil {
			return err
		}
		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			return fmt.Errorf("usuario con ID %s no encontrado", id)
		}
		return nil
	})
}

// AssignSucursales optimizado: usa un solo query con UNNEST para inserción masiva
func (s *storeUsuario) AssignSucursales(ctx context.Context, userID uuid.UUID, sucursales []uuid.UUID) error {
	defer performance.Trace(ctx, "store", "AssignSucursales", performance.DbThreshold, time.Now())

	return ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		// 1. Limpiar accesos previos
		if _, err := tx.ExecContext(ctx, "DELETE FROM usuario_sucursal_acceso WHERE id_usuario = $1", userID); err != nil {
			return fmt.Errorf("error al eliminar accesos previos: %w", err)
		}

		if len(sucursales) == 0 {
			return nil
		}

		// 2. Inserción masiva eficiente usando UNNEST
		query := `
			INSERT INTO usuario_sucursal_acceso (id_usuario, id_sucursal)
			SELECT $1, unnest($2::uuid[])
		`
		if _, err := tx.ExecContext(ctx, query, userID, sucursales); err != nil {
			return fmt.Errorf("error al insertar accesos masivos: %w", err)
		}
		return nil
	})
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

func (s *storeUsuario) GetPermisosByRol(ctx context.Context, rolID uuid.UUID) ([]string, error) {
	defer performance.Trace(ctx, "store", "GetPermisosByRol", performance.DbThreshold, time.Now())
	query := `
		SELECT m.ruta 
		FROM permiso_rol pr
		JOIN modulo m ON pr.id_modulo = m.id_modulo
		WHERE pr.id_rol = $1 AND pr.can_read = true AND m.is_active = true AND m.deleted_at IS NULL
	`
	rows, err := s.db.QueryContext(ctx, query, rolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permisos []string
	for rows.Next() {
		var ruta string
		if err := rows.Scan(&ruta); err != nil {
			return nil, err
		}
		if ruta != "" {
			permisos = append(permisos, ruta)
		}
	}
	return permisos, nil
}

func (s *storeUsuario) GetUsuarioByUsername(ctx context.Context, username string) (*models.Usuario, error) {
	defer performance.Trace(ctx, "store", "GetUsuarioByUsername", performance.DbThreshold, time.Now())
	query := fmt.Sprintf(`
		SELECT %s
		FROM usuario u
		LEFT JOIN rol r ON u.id_rol = r.id_rol
		LEFT JOIN sucursal su ON su.id_sucursal = u.id_sucursal
		WHERE u.username = $1 AND u.deleted_at IS NULL
	`, usuarioSelectFields)

	usuario := &models.Usuario{}
	err := s.scanRowUsuario(s.db.QueryRowContext(ctx, query, username), usuario)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("usuario %s no encontrado", username)
	}
	return usuario, err
}

func (s *storeUsuario) GetUsuarioByPin(ctx context.Context, pin string) (*models.Usuario, error) {
	defer performance.Trace(ctx, "store", "GetUsuarioByPin", performance.DbThreshold, time.Now())
	query := fmt.Sprintf(`
		SELECT %s
		FROM usuario u
		LEFT JOIN rol r ON u.id_rol = r.id_rol
		LEFT JOIN sucursal su ON su.id_sucursal = u.id_sucursal
		WHERE u.usu_pin_pos = $1 AND u.deleted_at IS NULL
	`, usuarioSelectFields)

	usuario := &models.Usuario{}
	err := s.scanRowUsuario(s.db.QueryRowContext(ctx, query, pin), usuario)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("pin inválido")
	}
	return usuario, err
}

func (s *storeUsuario) UpdateTurnoStatus(ctx context.Context, id uuid.UUID, enTurno bool) error {
	defer performance.Trace(ctx, "store", "UpdateTurnoStatus", performance.DbThreshold, time.Now())
	query := `UPDATE usuario SET en_turno = $1, updated_at = CURRENT_TIMESTAMP WHERE id_usuario = $2`
	_, err := s.db.ExecContext(ctx, query, enTurno, id)
	return err
}
