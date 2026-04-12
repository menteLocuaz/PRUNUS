package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/utils/performance"
)

// StoreCliente define las operaciones de persistencia para el catálogo de clientes.
// Sigue el patrón Repository para desacoplar la lógica de negocio de la base de datos.
type StoreCliente interface {
	GetAllClientes(ctx context.Context) ([]*models.Cliente, error)
	GetClienteByID(ctx context.Context, id uuid.UUID) (*models.Cliente, error)
	CreateCliente(ctx context.Context, cliente *models.Cliente) (*models.Cliente, error)
	UpdateCliente(ctx context.Context, id uuid.UUID, cliente *models.Cliente) (*models.Cliente, error)
	DeleteCliente(ctx context.Context, id uuid.UUID) error
}

type storeCliente struct {
	db *sql.DB
}

// NewCliente crea una nueva instancia del store de clientes inyectando la conexión a la DB.
func NewCliente(db *sql.DB) StoreCliente {
	return &storeCliente{db: db}
}

// Campos base para SELECT de cliente mapeados al esquema DB (Migración 000034).
// Se utiliza COALESCE para manejar campos opcionales que pueden ser NULL en la base de datos.
const clienteSelectFields = `
	c.id_cliente, c.nombre_completo, COALESCE(c.tipo_documento, ''), 
	COALESCE(c.documento, ''), COALESCE(c.email, ''), COALESCE(c.telefono, ''), 
	COALESCE(c.direccion, ''), c.id_status, c.metadata, c.created_at, c.updated_at
`

// scanRowCliente centraliza el escaneo de resultados para mantener consistencia en todas las consultas de lectura.
// Maneja la conversión de JSONB de PostgreSQL al mapa de Metadata en el modelo Go.
func (s *storeCliente) scanRowCliente(scanner interface{ Scan(dest ...any) error }, c *models.Cliente) error {
	var metadataJSON []byte
	err := scanner.Scan(
		&c.IDCliente, &c.NombreCompleto, &c.TipoDocumento,
		&c.Documento, &c.Email, &c.Telefono,
		&c.Direccion, &c.IDStatus, &metadataJSON, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return err
	}

	if len(metadataJSON) > 0 {
		_ = json.Unmarshal(metadataJSON, &c.Metadata)
	}
	return nil
}

// GetAllClientes recupera todos los clientes activos (no eliminados lógicamente).
func (s *storeCliente) GetAllClientes(ctx context.Context) ([]*models.Cliente, error) {
	defer performance.Trace(ctx, "store", "GetAllClientes", performance.DbThreshold, time.Now())
	
	query := fmt.Sprintf(`
		SELECT %s
		FROM cliente c
		WHERE c.deleted_at IS NULL
		ORDER BY c.nombre_completo ASC
	`, clienteSelectFields)

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener clientes: %w", err)
	}
	defer rows.Close()

	var clientes []*models.Cliente
	for rows.Next() {
		c := &models.Cliente{}
		if err := s.scanRowCliente(rows, c); err != nil {
			return nil, fmt.Errorf("error al escanear cliente: %w", err)
		}
		clientes = append(clientes, c)
	}

	return clientes, nil
}

// GetClienteByID busca un cliente por su identificador único.
func (s *storeCliente) GetClienteByID(ctx context.Context, id uuid.UUID) (*models.Cliente, error) {
	defer performance.Trace(ctx, "store", "GetClienteByID", performance.DbThreshold, time.Now())
	
	query := fmt.Sprintf(`
		SELECT %s
		FROM cliente c
		WHERE c.id_cliente = $1 AND c.deleted_at IS NULL
	`, clienteSelectFields)

	c := &models.Cliente{}
	err := s.scanRowCliente(s.db.QueryRowContext(ctx, query, id), c)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("cliente con ID %s no encontrado", id)
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener cliente: %w", err)
	}

	return c, nil
}

// CreateCliente inserta un nuevo registro de cliente y retorna la entidad con su ID y timestamps generados.
// Utiliza ExecAudited para asegurar que el usuario que realiza la operación sea registrado por los triggers.
func (s *storeCliente) CreateCliente(ctx context.Context, cliente *models.Cliente) (*models.Cliente, error) {
	defer performance.Trace(ctx, "store", "CreateCliente", performance.DbThreshold, time.Now())

	// Serializar metadatos a JSON para la base de datos
	metadata, err := json.Marshal(cliente.Metadata)
	if err != nil || metadata == nil {
		metadata = []byte("{}")
	}

	err = ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		query := `
			INSERT INTO cliente (nombre_completo, tipo_documento, documento, email, telefono, direccion, id_status, metadata)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id_cliente, created_at, updated_at
		`
		return tx.QueryRowContext(ctx, query,
			cliente.NombreCompleto, cliente.TipoDocumento, cliente.Documento,
			cliente.Email, cliente.Telefono, cliente.Direccion,
			cliente.IDStatus, metadata,
		).Scan(&cliente.IDCliente, &cliente.CreatedAt, &cliente.UpdatedAt)
	})

	if err != nil {
		return nil, fmt.Errorf("error al crear cliente: %w", err)
	}

	return cliente, nil
}

// UpdateCliente actualiza los datos de un cliente existente.
func (s *storeCliente) UpdateCliente(ctx context.Context, id uuid.UUID, cliente *models.Cliente) (*models.Cliente, error) {
	defer performance.Trace(ctx, "store", "UpdateCliente", performance.DbThreshold, time.Now())

	metadata, err := json.Marshal(cliente.Metadata)
	if err != nil || metadata == nil {
		metadata = []byte("{}")
	}

	err = ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		query := `
			UPDATE cliente
			SET nombre_completo = $1, tipo_documento = $2, documento = $3, email = $4, 
			    telefono = $5, direccion = $6, id_status = $7, metadata = $8, updated_at = CURRENT_TIMESTAMP
			WHERE id_cliente = $9 AND deleted_at IS NULL
			RETURNING created_at, updated_at
		`
		return tx.QueryRowContext(ctx, query,
			cliente.NombreCompleto, cliente.TipoDocumento, cliente.Documento,
			cliente.Email, cliente.Telefono, cliente.Direccion,
			cliente.IDStatus, metadata, id,
		).Scan(&cliente.CreatedAt, &cliente.UpdatedAt)
	})

	if err != nil {
		return nil, fmt.Errorf("error al actualizar cliente: %w", err)
	}

	cliente.IDCliente = id
	return cliente, nil
}

// DeleteCliente realiza una eliminación lógica del cliente seteando el campo deleted_at.
func (s *storeCliente) DeleteCliente(ctx context.Context, id uuid.UUID) error {
	defer performance.Trace(ctx, "store", "DeleteCliente", performance.DbThreshold, time.Now())

	err := ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		query := `UPDATE cliente SET deleted_at = CURRENT_TIMESTAMP WHERE id_cliente = $1 AND deleted_at IS NULL`
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
		return fmt.Errorf("error al eliminar cliente: %w", err)
	}

	return nil
}
