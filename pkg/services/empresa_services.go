// Package services contiene la lógica de negocio y reglas de validación
// para las entidades del sistema. En este caso, maneja las operaciones
// relacionadas con Empresa.
package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
	zaplogger "github.com/prunus/pkg/utils/logger"
	"go.uber.org/zap"
)

// ServiceEmpresa define un servicio que encapsula la lógica de negocio
// para la entidad Empresa. Se apoya en un StoreEmpresa para interactuar
// con la capa de persistencia (base de datos).
type ServiceEmpresa struct {
	store  store.StoreEmpresa
	logger *zap.Logger
}

// NewServiceEmpresa crea una nueva instancia del servicio de Empresa,
// recibiendo como dependencia un StoreEmpresa. Esto facilita la inyección
// de dependencias y el testeo.
func NewServiceEmpresa(s store.StoreEmpresa, logger *zap.Logger) *ServiceEmpresa {
	return &ServiceEmpresa{
		store:  s,
		logger: logger,
	}
}

// GetAllEmpresa devuelve todas las empresas registradas en el sistema.
func (s *ServiceEmpresa) GetAllEmpresa(ctx context.Context) ([]*models.Empresa, error) {
	return s.store.GetAllEmpresa(ctx)
}

// GetByIDEmpresa busca una empresa por su ID único.
func (s *ServiceEmpresa) GetByIDEmpresa(ctx context.Context, id uuid.UUID) (*models.Empresa, error) {
	if id == uuid.Nil {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de obtener empresa con ID nulo")
		return nil, errors.New("el ID de la empresa es requerido")
	}
	return s.store.GetByIdEmpresa(ctx, id)
}

// CrearEmpresa valida y crea una nueva empresa en el sistema.
func (s *ServiceEmpresa) CrearEmpresa(ctx context.Context, empresa models.Empresa) (*models.Empresa, error) {
	if empresa.Nombre == "" {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de crear empresa con nombre vacío")
		return nil, errors.New("Falta el nombre de la empresa")
	}
	if empresa.IDStatus == uuid.Nil {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de crear empresa sin estatus", zap.String("nombre", empresa.Nombre))
		return nil, errors.New("Falta el ID de estatus")
	}
	return s.store.CreateEmpresa(ctx, &empresa)
}

// UpdateEmpresa valida y actualiza una empresa existente.
func (s *ServiceEmpresa) UpdateEmpresa(ctx context.Context, id uuid.UUID, empresa models.Empresa) (*models.Empresa, error) {
	if id == uuid.Nil {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de actualizar empresa con ID nulo")
		return nil, errors.New("el ID de la empresa es requerido")
	}
	if empresa.Nombre == "" {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de actualizar empresa con nombre vacío", zap.String("id", id.String()))
		return nil, errors.New("Falta el nombre de la empresa")
	}
	return s.store.UpdateEmpresa(ctx, id, &empresa)
}

// ElimminarEmpresa elimina una empresa por su ID.
func (s *ServiceEmpresa) ElimminarEmpresa(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de eliminar empresa con ID nulo")
		return errors.New("el ID de la empresa es requerido")
	}
	return s.store.DeleteEmpresa(ctx, id)
}
