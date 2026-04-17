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

// StoreProveedor define las operaciones de persistencia para el catálogo de proveedores.
type StoreProveedor interface {
	GetAllProveedores(ctx context.Context, params dto.PaginationParams) ([]*models.Proveedor, error)
	GetProveedorByID(ctx context.Context, id uuid.UUID) (*models.Proveedor, error)
	CreateProveedor(ctx context.Context, proveedor *models.Proveedor) (*models.Proveedor, error)
	UpdateProveedor(ctx context.Context, id uuid.UUID, proveedor *models.Proveedor) (*models.Proveedor, error)
	DeleteProveedor(ctx context.Context, id uuid.UUID) error
}

type storeProveedor struct {
	db *sql.DB
}

// NewProveedor crea una nueva instancia del store de proveedores.
func NewProveedor(db *sql.DB) StoreProveedor {
	return &storeProveedor{db: db}
}

// Campos base para SELECT de proveedor mapeados al esquema DB (Migración 000034).
const proveedorSelectFields = `
	p.id_proveedor, p.razon_social, p.nit_rut, COALESCE(p.contacto_nombre, ''),
	COALESCE(p.telefono, ''), COALESCE(p.direccion, ''), COALESCE(p.email, ''),
	p.id_status, p.metadata, p.created_at, p.updated_at
`

// scanRowProveedor centraliza el escaneo de resultados para mantener consistencia.
func (s *storeProveedor) scanRowProveedor(scanner interface{ Scan(dest ...any) error }, p *models.Proveedor) error {
	var metadataJSON []byte
	err := scanner.Scan(
		&p.IDProveedor, &p.RazonSocial, &p.NitRut, &p.ContactoNombre,
		&p.Telefono, &p.Direccion, &p.Email,
		&p.IDStatus, &metadataJSON, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return err
	}

	if len(metadataJSON) > 0 {
		_ = json.Unmarshal(metadataJSON, &p.Metadata)
	}
	return nil
}

func (s *storeProveedor) GetAllProveedores(ctx context.Context, params dto.PaginationParams) ([]*models.Proveedor, error) {
	defer performance.Trace(ctx, "store", "GetAllProveedores", performance.DbThreshold, time.Now())

	if params.Limit <= 0 {
		params.Limit = dto.DefaultLimit
	}

	query := fmt.Sprintf(`
		SELECT %s
		FROM proveedor p
		WHERE p.deleted_at IS NULL
	`, proveedorSelectFields)

	var args []interface{}
	if params.LastDate != nil {
		query += " AND p.created_at < $1"
		args = append(args, params.LastDate)
	}

	query += " ORDER BY p.created_at DESC LIMIT $" + fmt.Sprint(len(args)+1)
	args = append(args, params.Limit)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error al obtener proveedores: %w", err)
	}
	defer rows.Close()

	var proveedores []*models.Proveedor
	for rows.Next() {
		p := &models.Proveedor{}
		if err := s.scanRowProveedor(rows, p); err != nil {
			return nil, fmt.Errorf("error al escanear proveedor: %w", err)
		}
		proveedores = append(proveedores, p)
	}

	return proveedores, nil
}

func (s *storeProveedor) GetProveedorByID(ctx context.Context, id uuid.UUID) (*models.Proveedor, error) {
	defer performance.Trace(ctx, "store", "GetProveedorByID", performance.DbThreshold, time.Now())
	query := fmt.Sprintf(`
		SELECT %s
		FROM proveedor p
		WHERE p.id_proveedor = $1 AND p.deleted_at IS NULL
	`, proveedorSelectFields)

	p := &models.Proveedor{}
	err := s.scanRowProveedor(s.db.QueryRowContext(ctx, query, id), p)

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

	metadata, _ := json.Marshal(proveedor.Metadata)
	if metadata == nil {
		metadata = []byte("{}")
	}

	err := ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		query := `
			INSERT INTO proveedor (razon_social, nit_rut, contacto_nombre, telefono, direccion, email, id_status, metadata)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id_proveedor, created_at, updated_at
		`
		return tx.QueryRowContext(ctx, query,
			proveedor.RazonSocial, proveedor.NitRut, proveedor.ContactoNombre,
			proveedor.Telefono, proveedor.Direccion, proveedor.Email,
			proveedor.IDStatus, metadata,
		).Scan(&proveedor.IDProveedor, &proveedor.CreatedAt, &proveedor.UpdatedAt)
	})

	if err != nil {
		return nil, fmt.Errorf("error al crear proveedor: %w", err)
	}

	return proveedor, nil
}

func (s *storeProveedor) UpdateProveedor(ctx context.Context, id uuid.UUID, proveedor *models.Proveedor) (*models.Proveedor, error) {
	defer performance.Trace(ctx, "store", "UpdateProveedor", performance.DbThreshold, time.Now())

	metadata, _ := json.Marshal(proveedor.Metadata)
	if metadata == nil {
		metadata = []byte("{}")
	}

	err := ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		query := `
			UPDATE proveedor
			SET razon_social = $1, nit_rut = $2, contacto_nombre = $3, telefono = $4, 
			    direccion = $5, email = $6, id_status = $7, metadata = $8, updated_at = CURRENT_TIMESTAMP
			WHERE id_proveedor = $9 AND deleted_at IS NULL
			RETURNING created_at, updated_at
		`
		return tx.QueryRowContext(ctx, query,
			proveedor.RazonSocial, proveedor.NitRut, proveedor.ContactoNombre,
			proveedor.Telefono, proveedor.Direccion, proveedor.Email,
			proveedor.IDStatus, metadata, id,
		).Scan(&proveedor.CreatedAt, &proveedor.UpdatedAt)
	})

	if err != nil {
		return nil, fmt.Errorf("error al actualizar proveedor: %w", err)
	}

	proveedor.IDProveedor = id
	return proveedor, nil
}

func (s *storeProveedor) DeleteProveedor(ctx context.Context, id uuid.UUID) error {
	defer performance.Trace(ctx, "store", "DeleteProveedor", performance.DbThreshold, time.Now())

	err := ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		query := `UPDATE proveedor SET deleted_at = CURRENT_TIMESTAMP WHERE id_proveedor = $1 AND deleted_at IS NULL`
		result, err := tx.ExecContext(ctx, query, id)
		if err != nil {
			return err
		}
		rows, _ := result.RowsAffected()
		if rows == 0 {
			return sql.ErrNoRows
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error al eliminar proveedor: %w", err)
	}

	return nil
}
