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

type StoreProveedor interface {
	GetAllProveedores(ctx context.Context) ([]*models.Proveedor, error)
	GetProveedorByID(ctx context.Context, id uuid.UUID) (*models.Proveedor, error)
	CreateProveedor(ctx context.Context, proveedor *models.Proveedor) (*models.Proveedor, error)
	UpdateProveedor(ctx context.Context, id uuid.UUID, proveedor *models.Proveedor) (*models.Proveedor, error)
	DeleteProveedor(ctx context.Context, id uuid.UUID) error
}

type storeProveedor struct {
	db *sql.DB
}

func NewProveedor(db *sql.DB) StoreProveedor {
	return &storeProveedor{db: db}
}

func (s *storeProveedor) GetAllProveedores(ctx context.Context) ([]*models.Proveedor, error) {
	defer performance.Trace(ctx, "store", "GetAllProveedores", performance.DbThreshold, time.Now())
	query := `
	SELECT
		p.id_proveedor,
		p.nombre,
		p.ruc,
		p.telefono,
		p.direccion,
		p.email,
		p.id_status,
		p.id_sucursal,
		p.id_empresa,
		p.created_at,
		p.updated_at,

		su.id_sucursal,
		su.nombre_sucursal,
		su.id_status,

		e.id_empresa,
		e.nombre,
		e.rut,
		e.id_status
	FROM proveedor p
	LEFT JOIN sucursal su ON su.id_sucursal = p.id_sucursal
	LEFT JOIN empresa e ON e.id_empresa = p.id_empresa
	WHERE p.deleted_at IS NULL
	ORDER BY p.id_proveedor
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener proveedores: %w", err)
	}
	defer rows.Close()

	var proveedores []*models.Proveedor

	for rows.Next() {
		p := &models.Proveedor{
			Sucursal: &models.Sucursal{},
			Empresa:  &models.Empresa{},
		}

		if err := rows.Scan(
			&p.IDProveedor,
			&p.Nombre,
			&p.RUC,
			&p.Telefono,
			&p.Direccion,
			&p.Email,
			&p.IDStatus,
			&p.IDSucursal,
			&p.IDEmpresa,
			&p.CreatedAt,
			&p.UpdatedAt,

			&p.Sucursal.IDSucursal,
			&p.Sucursal.NombreSucursal,
			&p.Sucursal.IDStatus,

			&p.Empresa.IDEmpresa,
			&p.Empresa.Nombre,
			&p.Empresa.RUT,
			&p.Empresa.IDStatus,
		); err != nil {
			return nil, fmt.Errorf("error al escanear proveedor: %w", err)
		}

		proveedores = append(proveedores, p)
	}

	return proveedores, nil
}

func (s *storeProveedor) GetProveedorByID(ctx context.Context, id uuid.UUID) (*models.Proveedor, error) {
	defer performance.Trace(ctx, "store", "GetProveedorByID", performance.DbThreshold, time.Now())
	query := `
	SELECT
		p.id_proveedor,
		p.nombre,
		p.ruc,
		p.telefono,
		p.direccion,
		p.email,
		p.id_status,
		p.id_sucursal,
		p.id_empresa,
		p.created_at,
		p.updated_at,

		su.id_sucursal,
		su.nombre_sucursal,
		su.id_status,

		e.id_empresa,
		e.nombre,
		e.rut,
		e.id_status
	FROM proveedor p
	LEFT JOIN sucursal su ON su.id_sucursal = p.id_sucursal
	LEFT JOIN empresa e ON e.id_empresa = p.id_empresa
	WHERE p.id_proveedor = $1
	  AND p.deleted_at IS NULL
	`

	p := &models.Proveedor{
		Sucursal: &models.Sucursal{},
		Empresa:  &models.Empresa{},
	}

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&p.IDProveedor,
		&p.Nombre,
		&p.RUC,
		&p.Telefono,
		&p.Direccion,
		&p.Email,
		&p.IDStatus,
		&p.IDSucursal,
		&p.IDEmpresa,
		&p.CreatedAt,
		&p.UpdatedAt,

		&p.Sucursal.IDSucursal,
		&p.Sucursal.NombreSucursal,
		&p.Sucursal.IDStatus,

		&p.Empresa.IDEmpresa,
		&p.Empresa.Nombre,
		&p.Empresa.RUT,
		&p.Empresa.IDStatus,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("proveedor con ID %s no encontrado", id)
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener proveedor: %w", err)
	}

	return p, nil
}

func (s *storeProveedor) CreateProveedor(ctx context.Context, proveedor *models.Proveedor) (*models.Proveedor, error) {
	defer performance.Trace(ctx, "store", "CreateProveedor", performance.DbThreshold, time.Now())
	query := `
		INSERT INTO proveedor (nombre, ruc, telefono, direccion, email, id_status, id_sucursal, id_empresa)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id_proveedor
	`

	var id uuid.UUID
	err := s.db.QueryRowContext(ctx, query,
		proveedor.Nombre,
		proveedor.RUC,
		proveedor.Telefono,
		proveedor.Direccion,
		proveedor.Email,
		proveedor.IDStatus,
		proveedor.IDSucursal,
		proveedor.IDEmpresa,
	).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("error al crear proveedor: %w", err)
	}

	proveedor.IDProveedor = id
	return proveedor, nil
}

func (s *storeProveedor) UpdateProveedor(ctx context.Context, id uuid.UUID, proveedor *models.Proveedor) (*models.Proveedor, error) {
	defer performance.Trace(ctx, "store", "UpdateProveedor", performance.DbThreshold, time.Now())
	query := `
		UPDATE proveedor
		SET
			nombre = $1,
			ruc = $2,
			telefono = $3,
			direccion = $4,
			email = $5,
			id_status = $6,
			id_sucursal = $7,
			id_empresa = $8,
			updated_at = CURRENT_TIMESTAMP
		WHERE id_proveedor = $9
		  AND deleted_at IS NULL
		RETURNING
			id_proveedor,
			nombre,
			ruc,
			telefono,
			direccion,
			email,
			id_status,
			id_sucursal,
			id_empresa,
			created_at,
			updated_at
	`

	err := s.db.QueryRowContext(ctx, query,
		proveedor.Nombre,
		proveedor.RUC,
		proveedor.Telefono,
		proveedor.Direccion,
		proveedor.Email,
		proveedor.IDStatus,
		proveedor.IDSucursal,
		proveedor.IDEmpresa,
		id,
	).Scan(
		&proveedor.IDProveedor,
		&proveedor.Nombre,
		&proveedor.RUC,
		&proveedor.Telefono,
		&proveedor.Direccion,
		&proveedor.Email,
		&proveedor.IDStatus,
		&proveedor.IDSucursal,
		&proveedor.IDEmpresa,
		&proveedor.CreatedAt,
		&proveedor.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("proveedor con ID %s no encontrado", id)
	}
	if err != nil {
		return nil, fmt.Errorf("error al actualizar proveedor: %w", err)
	}

	return proveedor, nil
}

func (s *storeProveedor) DeleteProveedor(ctx context.Context, id uuid.UUID) error {
	defer performance.Trace(ctx, "store", "DeleteProveedor", performance.DbThreshold, time.Now())
	query := `UPDATE proveedor SET deleted_at = $1 WHERE id_proveedor = $2 AND deleted_at IS NULL`

	result, err := s.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("error al eliminar proveedor: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

