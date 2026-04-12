package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
)

// ServiceCliente coordina la lógica de negocio para la gestión de clientes.
type ServiceCliente struct {
	store  store.StoreCliente
	logger *slog.Logger
}

// NewServiceCliente crea una nueva instancia del servicio de clientes.
func NewServiceCliente(s store.StoreCliente, logger *slog.Logger) *ServiceCliente {
	return &ServiceCliente{
		store:  s,
		logger: logger,
	}
}

// GetAllClientes recupera la lista completa de clientes desde el store.
func (s *ServiceCliente) GetAllClientes(ctx context.Context) ([]*models.Cliente, error) {
	return s.store.GetAllClientes(ctx)
}

// GetClienteByID busca un cliente específico validando que el ID sea correcto.
func (s *ServiceCliente) GetClienteByID(ctx context.Context, id uuid.UUID) (*models.Cliente, error) {
	if id == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de obtener cliente con ID nulo")
		return nil, errors.New("el ID del cliente es requerido")
	}
	return s.store.GetClienteByID(ctx, id)
}

// CreateCliente valida los datos básicos antes de persistir un nuevo cliente.
func (s *ServiceCliente) CreateCliente(ctx context.Context, cliente models.Cliente) (*models.Cliente, error) {
	if cliente.NombreCompleto == "" {
		s.logger.WarnContext(ctx, "Intento de crear cliente sin nombre completo")
		return nil, errors.New("falta el nombre completo del cliente")
	}
	if cliente.IDStatus == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de crear cliente sin estatus", slog.String("nombre", cliente.NombreCompleto))
		return nil, errors.New("falta el ID de estatus")
	}
	return s.store.CreateCliente(ctx, &cliente)
}

// UpdateCliente valida la existencia del ID y la integridad de los datos antes de actualizar.
func (s *ServiceCliente) UpdateCliente(ctx context.Context, id uuid.UUID, cliente models.Cliente) (*models.Cliente, error) {
	if id == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de actualizar cliente con ID nulo")
		return nil, errors.New("el ID del cliente es requerido")
	}
	if cliente.NombreCompleto == "" {
		s.logger.WarnContext(ctx, "Intento de actualizar cliente con nombre vacío", slog.String("id", id.String()))
		return nil, errors.New("falta el nombre completo del cliente")
	}
	return s.store.UpdateCliente(ctx, id, &cliente)
}

// DeleteCliente gestiona la eliminación lógica de un cliente.
func (s *ServiceCliente) DeleteCliente(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de eliminar cliente con ID nulo")
		return errors.New("el ID del cliente es requerido")
	}
	return s.store.DeleteCliente(ctx, id)
}
