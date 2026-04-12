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

// StoreCategoria define las operaciones de persistencia para el catálogo de categorías.
type StoreCategoria interface {
	GetAllCategorias(ctx context.Context) ([]*models.Categoria, error)
	GetCategoriaByID(ctx context.Context, id uuid.UUID) (*models.Categoria, error)
	GetCategoriasBySucursal(ctx context.Context, sucursalID uuid.UUID) ([]*models.Categoria, error)
	CreateCategoria(ctx context.Context, categoria *models.Categoria) (*models.Categoria, error)
	UpdateCategoria(ctx context.Context, id uuid.UUID, categoria *models.Categoria) (*models.Categoria, error)
	DeleteCategoria(ctx context.Context, id uuid.UUID) error
}

type storeCategoria struct {
	db *sql.DB
}

// NewCategoria crea una nueva instancia del store de categorías.
func NewCategoria(db *sql.DB) StoreCategoria {
	return &storeCategoria{db: db}
}

// Campos base para SELECT de categoria mapeados al esquema DB.
const categoriaSelectFields = `
	c.id_categoria, c.nombre, c.id_status, c.id_sucursal, c.created_at, c.updated_at,
	s.id_sucursal, s.nombre_sucursal, s.id_status
`

// scanRowCategoria centraliza el escaneo de resultados para mantener consistencia.
func (s *storeCategoria) scanRowCategoria(scanner interface{ Scan(dest ...any) error }, c *models.Categoria) error {
	if c.Sucursal == nil {
		c.Sucursal = &models.Sucursal{}
	}

	return scanner.Scan(
		&c.IDCategoria, &c.Nombre, &c.IDStatus, &c.IDSucursal, &c.CreatedAt, &c.UpdatedAt,
		&c.Sucursal.IDSucursal, &c.Sucursal.NombreSucursal, &c.Sucursal.IDStatus,
	)
}

func (s *storeCategoria) GetAllCategorias(ctx context.Context) ([]*models.Categoria, error) {
	defer performance.Trace(ctx, "store", "GetAllCategorias", performance.DbThreshold, time.Now())
	query := fmt.Sprintf(`
		SELECT %s
		FROM categoria c
		LEFT JOIN sucursal s ON s.id_sucursal = c.id_sucursal
		WHERE c.deleted_at IS NULL
		ORDER BY c.nombre ASC
	`, categoriaSelectFields)

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener categorias: %w", err)
	}
	defer rows.Close()

	var categorias []*models.Categoria
	for rows.Next() {
		c := &models.Categoria{}
		if err := s.scanRowCategoria(rows, c); err != nil {
			return nil, fmt.Errorf("error al escanear categoria: %w", err)
		}
		categorias = append(categorias, c)
	}

	return categorias, nil
}

func (s *storeCategoria) GetCategoriaByID(ctx context.Context, id uuid.UUID) (*models.Categoria, error) {
	defer performance.Trace(ctx, "store", "GetCategoriaByID", performance.DbThreshold, time.Now())
	query := fmt.Sprintf(`
		SELECT %s
		FROM categoria c
		LEFT JOIN sucursal s ON s.id_sucursal = c.id_sucursal
		WHERE c.id_categoria = $1 AND c.deleted_at IS NULL
	`, categoriaSelectFields)

	c := &models.Categoria{}
	err := s.scanRowCategoria(s.db.QueryRowContext(ctx, query, id), c)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("categoria con ID %s no encontrada", id)
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener categoria: %w", err)
	}

	return c, nil
}

func (s *storeCategoria) GetCategoriasBySucursal(ctx context.Context, sucursalID uuid.UUID) ([]*models.Categoria, error) {
	defer performance.Trace(ctx, "store", "GetCategoriasBySucursal", performance.DbThreshold, time.Now())
	query := fmt.Sprintf(`
		SELECT %s
		FROM categoria c
		LEFT JOIN sucursal s ON s.id_sucursal = c.id_sucursal
		WHERE c.id_sucursal = $1 AND c.deleted_at IS NULL
		ORDER BY c.nombre ASC
	`, categoriaSelectFields)

	rows, err := s.db.QueryContext(ctx, query, sucursalID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener categorias por sucursal: %w", err)
	}
	defer rows.Close()

	var categorias []*models.Categoria
	for rows.Next() {
		c := &models.Categoria{}
		if err := s.scanRowCategoria(rows, c); err != nil {
			return nil, fmt.Errorf("error al escanear categoria: %w", err)
		}
		categorias = append(categorias, c)
	}

	return categorias, nil
}

func (s *storeCategoria) CreateCategoria(ctx context.Context, categoria *models.Categoria) (*models.Categoria, error) {
	defer performance.Trace(ctx, "store", "CreateCategoria", performance.DbThreshold, time.Now())

	// Validación de seguridad antes de insertar
	if categoria.IDStatus == uuid.Nil {
		return nil, fmt.Errorf("error de integridad: id_status es requerido para crear una categoría")
	}

	err := ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		query := `
			INSERT INTO categoria (nombre, id_status, id_sucursal)
			VALUES ($1, $2, $3)
			RETURNING id_categoria, created_at, updated_at
		`
		return tx.QueryRowContext(ctx, query,
			categoria.Nombre,
			categoria.IDStatus,
			categoria.IDSucursal,
		).Scan(&categoria.IDCategoria, &categoria.CreatedAt, &categoria.UpdatedAt)
	})

	if err != nil {
		return nil, fmt.Errorf("error al crear categoria en DB: %w", err)
	}

	return categoria, nil
}

func (s *storeCategoria) UpdateCategoria(ctx context.Context, id uuid.UUID, categoria *models.Categoria) (*models.Categoria, error) {
	defer performance.Trace(ctx, "store", "UpdateCategoria", performance.DbThreshold, time.Now())

	err := ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		query := `
			UPDATE categoria
			SET nombre = $1, id_status = $2, id_sucursal = $3, updated_at = CURRENT_TIMESTAMP
			WHERE id_categoria = $4 AND deleted_at IS NULL
			RETURNING created_at, updated_at
		`
		return tx.QueryRowContext(ctx, query,
			categoria.Nombre, categoria.IDStatus, categoria.IDSucursal, id,
		).Scan(&categoria.CreatedAt, &categoria.UpdatedAt)
	})

	if err != nil {
		return nil, fmt.Errorf("error al actualizar categoria: %w", err)
	}

	categoria.IDCategoria = id
	return categoria, nil
}

func (s *storeCategoria) DeleteCategoria(ctx context.Context, id uuid.UUID) error {
	defer performance.Trace(ctx, "store", "DeleteCategoria", performance.DbThreshold, time.Now())

	err := ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		query := `UPDATE categoria SET deleted_at = CURRENT_TIMESTAMP WHERE id_categoria = $1 AND deleted_at IS NULL`
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
		return fmt.Errorf("error al eliminar categoria: %w", err)
	}

	return nil
}
