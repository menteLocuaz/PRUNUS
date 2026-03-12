// Package services contiene la lógica de negocio y reglas de validación
// para las entidades del sistema. En este caso, maneja las operaciones
// relacionadas con Empresa.
package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
)

// ServiceEmpresa define un servicio que encapsula la lógica de negocio
// para la entidad Empresa. Se apoya en un StoreEmpresa para interactuar
// con la capa de persistencia (base de datos).
type ServiceEmpresa struct {
	store  store.StoreEmpresa
	logger *slog.Logger
}

// NewServiceEmpresa crea una nueva instancia del servicio de Empresa,
// recibiendo como dependencia un StoreEmpresa. Esto facilita la inyección
// de dependencias y el testeo.
func NewServiceEmpresa(s store.StoreEmpresa, logger *slog.Logger) *ServiceEmpresa {
	return &ServiceEmpresa{
		store:  s,
		logger: logger,
	}
}

// GetAllEmpresa devuelve todas las empresas registradas en el sistema.
// Retorna un slice de punteros a Empresa y un posible error.
func (s *ServiceEmpresa) GetAllEmpresa(ctx context.Context) ([]*models.Empresa, error) {
	return s.store.GetAllEmpresa(ctx)
}

// GetByIDEmpresa busca una empresa por su ID único.
// Retorna la empresa encontrada o un error si no existe.
func (s *ServiceEmpresa) GetByIDEmpresa(ctx context.Context, id uuid.UUID) (*models.Empresa, error) {
	if id == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de obtener empresa con ID nulo")
		return nil, errors.New("el ID de la empresa es requerido")
	}
	return s.store.GetByIdEmpresa(ctx, id)
}

// CrearEmpresa valida y crea una nueva empresa en el sistema.
// Si el nombre está vacío, retorna un error de validación.
// En caso contrario, delega la creación al Store.
func (s *ServiceEmpresa) CrearEmpresa(ctx context.Context, empresa models.Empresa) (*models.Empresa, error) {
	if empresa.Nombre == "" {
		s.logger.WarnContext(ctx, "Intento de crear empresa con nombre vacío")
		return nil, errors.New("Falta el nombre de la empresa")
	}
	if empresa.IDStatus == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de crear empresa sin estatus", slog.String("nombre", empresa.Nombre))
		return nil, errors.New("Falta el ID de estatus")
	}
	return s.store.CreateEmpresa(ctx, &empresa)
}

// UpdateEmpresa valida y actualiza una empresa existente.
// Si el nombre está vacío, retorna un error de validación.
// En caso contrario, delega la actualización al Store.
func (s *ServiceEmpresa) UpdateEmpresa(ctx context.Context, id uuid.UUID, empresa models.Empresa) (*models.Empresa, error) {
	if id == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de actualizar empresa con ID nulo")
		return nil, errors.New("el ID de la empresa es requerido")
	}
	if empresa.Nombre == "" {
		s.logger.WarnContext(ctx, "Intento de actualizar empresa con nombre vacío", slog.String("id", id.String()))
		return nil, errors.New("Falta el nombre de la empresa")
	}
	return s.store.UpdateEmpresa(ctx, id, &empresa)
}

// ElimminarEmpresa elimina una empresa por su ID.
// Retorna un error si la operación falla.
func (s *ServiceEmpresa) ElimminarEmpresa(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		s.logger.WarnContext(ctx, "Intento de eliminar empresa con ID nulo")
		return errors.New("el ID de la empresa es requerido")
	}
	return s.store.DeleteEmpresa(ctx, id)
}
