package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
)

type ServiceCliente struct {
	store  store.StoreCliente
	logger *slog.Logger
}

func NewServiceCliente(s store.StoreCliente, logger *slog.Logger) *ServiceCliente {
	return &ServiceCliente{
		store:  s,
		logger: logger,
	}
}

func (s *ServiceCliente) GetAllClientes(ctx context.Context) ([]*models.Cliente, error) {
	return s.store.GetAllClientes(ctx)
}

func (s *ServiceCliente) GetClienteByID(ctx context.Context, id uuid.UUID) (*models.Cliente, error) {
	if id == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de obtener cliente con ID nulo")
		return nil, errors.New("el ID del cliente es requerido")
	}
	return s.store.GetClienteByID(ctx, id)
}

func (s *ServiceCliente) CreateCliente(ctx context.Context, cliente models.Cliente) (*models.Cliente, error) {
	if cliente.Nombre == "" {
		s.logger.WarnContext(ctx, "Intento de crear cliente sin nombre")
		return nil, errors.New("falta el nombre del cliente")
	}
	if cliente.EmpresaCliente == "" {
		s.logger.WarnContext(ctx, "Intento de crear cliente sin empresa", slog.String("nombre", cliente.Nombre))
		return nil, errors.New("falta el nombre de la empresa del cliente")
	}
	if cliente.IDStatus == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de crear cliente sin estatus", slog.String("nombre", cliente.Nombre))
		return nil, errors.New("falta el ID de estatus")
	}
	return s.store.CreateCliente(ctx, &cliente)
}

func (s *ServiceCliente) UpdateCliente(ctx context.Context, id uuid.UUID, cliente models.Cliente) (*models.Cliente, error) {
	if id == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de actualizar cliente con ID nulo")
		return nil, errors.New("el ID del cliente es requerido")
	}
	if cliente.Nombre == "" {
		s.logger.WarnContext(ctx, "Intento de actualizar cliente con nombre vacío", slog.String("id", id.String()))
		return nil, errors.New("falta el nombre del cliente")
	}
	return s.store.UpdateCliente(ctx, id, &cliente)
}

func (s *ServiceCliente) DeleteCliente(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de eliminar cliente con ID nulo")
		return errors.New("el ID del cliente es requerido")
	}
	return s.store.DeleteCliente(ctx, id)
}
