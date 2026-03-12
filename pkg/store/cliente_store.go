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

func NewCliente(db *sql.DB) StoreCliente {
	return &storeCliente{db: db}
}

func (s *storeCliente) GetAllClientes(ctx context.Context) ([]*models.Cliente, error) {
	defer performance.Trace(ctx, "store", "GetAllClientes", performance.DbThreshold, time.Now())
	query := `
	SELECT
		id_cliente,
		empresa_cliente,
		nombre,
		ruc,
		direccion,
		telefono,
		email,
		id_status,
		created_at,
		updated_at
	FROM cliente
	WHERE deleted_at IS NULL
	ORDER BY id_cliente
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener clientes: %w", err)
	}
	defer rows.Close()

	var clientes []*models.Cliente

	for rows.Next() {
		c := &models.Cliente{}

		if err := rows.Scan(
			&c.IDCliente,
			&c.EmpresaCliente,
			&c.Nombre,
			&c.RUC,
			&c.Direccion,
			&c.Telefono,
			&c.Email,
			&c.IDStatus,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error al escanear cliente: %w", err)
		}

		clientes = append(clientes, c)
	}

	return clientes, nil
}

func (s *storeCliente) GetClienteByID(ctx context.Context, id uuid.UUID) (*models.Cliente, error) {
	defer performance.Trace(ctx, "store", "GetClienteByID", performance.DbThreshold, time.Now())
	query := `
	SELECT
		id_cliente,
		empresa_cliente,
		nombre,
		ruc,
		direccion,
		telefono,
		email,
		id_status,
		created_at,
		updated_at
	FROM cliente
	WHERE id_cliente = $1
	  AND deleted_at IS NULL
	`

	c := &models.Cliente{}

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&c.IDCliente,
		&c.EmpresaCliente,
		&c.Nombre,
		&c.RUC,
		&c.Direccion,
		&c.Telefono,
		&c.Email,
		&c.IDStatus,
		&c.CreatedAt,
		&c.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("cliente con ID %s no encontrado", id)
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener cliente: %w", err)
	}

	return c, nil
}

func (s *storeCliente) CreateCliente(ctx context.Context, cliente *models.Cliente) (*models.Cliente, error) {
	defer performance.Trace(ctx, "store", "CreateCliente", performance.DbThreshold, time.Now())
	query := `
		INSERT INTO cliente (empresa_cliente, nombre, ruc, direccion, telefono, email, id_status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id_cliente
	`

	var id uuid.UUID
	err := s.db.QueryRowContext(ctx, query,
		cliente.EmpresaCliente,
		cliente.Nombre,
		cliente.RUC,
		cliente.Direccion,
		cliente.Telefono,
		cliente.Email,
		cliente.IDStatus,
	).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("error al crear cliente: %w", err)
	}

	cliente.IDCliente = id
	return cliente, nil
}

func (s *storeCliente) UpdateCliente(ctx context.Context, id uuid.UUID, cliente *models.Cliente) (*models.Cliente, error) {
	defer performance.Trace(ctx, "store", "UpdateCliente", performance.DbThreshold, time.Now())
	query := `
		UPDATE cliente
		SET
			empresa_cliente = $1,
			nombre = $2,
			ruc = $3,
			direccion = $4,
			telefono = $5,
			email = $6,
			id_status = $7,
			updated_at = CURRENT_TIMESTAMP
		WHERE id_cliente = $8
		  AND deleted_at IS NULL
		RETURNING
			id_cliente,
			empresa_cliente,
			nombre,
			ruc,
			direccion,
			telefono,
			email,
			id_status,
			created_at,
			updated_at
	`

	err := s.db.QueryRowContext(ctx, query,
		cliente.EmpresaCliente,
		cliente.Nombre,
		cliente.RUC,
		cliente.Direccion,
		cliente.Telefono,
		cliente.Email,
		cliente.IDStatus,
		id,
	).Scan(
		&cliente.IDCliente,
		&cliente.EmpresaCliente,
		&cliente.Nombre,
		&cliente.RUC,
		&cliente.Direccion,
		&cliente.Telefono,
		&cliente.Email,
		&cliente.IDStatus,
		&cliente.CreatedAt,
		&cliente.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("cliente con ID %s no encontrado", id)
	}
	if err != nil {
		return nil, fmt.Errorf("error al actualizar cliente: %w", err)
	}

	return cliente, nil
}

func (s *storeCliente) DeleteCliente(ctx context.Context, id uuid.UUID) error {
	defer performance.Trace(ctx, "store", "DeleteCliente", performance.DbThreshold, time.Now())
	query := `UPDATE cliente SET deleted_at = $1 WHERE id_cliente = $2 AND deleted_at IS NULL`

	result, err := s.db.ExecContext(ctx, query, time.Now(), id)

	if err != nil {
		return fmt.Errorf("error al eliminar cliente: %w", err)
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
