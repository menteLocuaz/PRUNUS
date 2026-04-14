package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
	"github.com/prunus/pkg/utils"
	zaplogger "github.com/prunus/pkg/utils/logger"
	"go.uber.org/zap"
)

type ServiceEstatus struct {
	store  store.StoreEstatus
	cache  *utils.CacheManager
	logger *zap.Logger
}

func NewServiceEstatus(s store.StoreEstatus, c *utils.CacheManager, logger *zap.Logger) *ServiceEstatus {
	return &ServiceEstatus{
		store:  s,
		cache:  c,
		logger: logger,
	}
}

const (
	cacheKeyEstatusAll      = "estatus:all"
	cacheKeyEstatusMaster   = "estatus:master_catalog"
	cacheKeyEstatusID       = "estatus:id:%s"
	cacheKeyEstatusTipo     = "estatus:tipo:%s"
	cacheKeyEstatusModuloID = "estatus:modulo:%d"
	cacheExpirationEstatus  = 24 * time.Hour
)

var moduloNames = map[int]string{
	1: "Empresa",
	2: "Sucursal",
	3: "Usuario",
	4: "Producto",
	5: "Venta",
	6: "Compra",
	7: "Finanzas",
	8: "Caja/POS",
}

func (s *ServiceEstatus) GetMasterCatalog(ctx context.Context) (map[int]interface{}, error) {
	return utils.GetOrSet(ctx, s.cache, cacheKeyEstatusMaster, cacheExpirationEstatus, func() (map[int]interface{}, error) {
		all, err := s.store.GetAllEstatus(ctx)
		if err != nil {
			return nil, err
		}

		catalog := make(map[int]interface{})
		for _, e := range all {
			if _, ok := catalog[e.MdlID]; !ok {
				name := "Módulo Desconocido"
				if n, exists := moduloNames[e.MdlID]; exists {
					name = n
				}
				catalog[e.MdlID] = struct {
					Modulo string            `json:"modulo"`
					Items  []*models.Estatus `json:"items"`
				}{
					Modulo: name,
					Items:  []*models.Estatus{},
				}
			}

			group := catalog[e.MdlID].(struct {
				Modulo string            `json:"modulo"`
				Items  []*models.Estatus `json:"items"`
			})
			group.Items = append(group.Items, e)
			catalog[e.MdlID] = group
		}
		return catalog, nil
	})
}

func (s *ServiceEstatus) GetAllEstatus(ctx context.Context) ([]*models.Estatus, error) {
	return utils.GetOrSet(ctx, s.cache, cacheKeyEstatusAll, cacheExpirationEstatus, func() ([]*models.Estatus, error) {
		return s.store.GetAllEstatus(ctx)
	})
}

func (s *ServiceEstatus) GetEstatusByID(ctx context.Context, id uuid.UUID) (*models.Estatus, error) {
	key := fmt.Sprintf(cacheKeyEstatusID, id.String())
	return utils.GetOrSet(ctx, s.cache, key, cacheExpirationEstatus, func() (*models.Estatus, error) {
		return s.store.GetEstatusByID(ctx, id)
	})
}

func (s *ServiceEstatus) GetEstatusByTipo(ctx context.Context, tipo string) ([]*models.Estatus, error) {
	key := fmt.Sprintf(cacheKeyEstatusTipo, tipo)
	return utils.GetOrSet(ctx, s.cache, key, cacheExpirationEstatus, func() ([]*models.Estatus, error) {
		return s.store.GetEstatusByTipo(ctx, tipo)
	})
}

func (s *ServiceEstatus) GetEstatusByModulo(ctx context.Context, moduloID int) ([]*models.Estatus, error) {
	key := fmt.Sprintf(cacheKeyEstatusModuloID, moduloID)
	return utils.GetOrSet(ctx, s.cache, key, cacheExpirationEstatus, func() ([]*models.Estatus, error) {
		return s.store.GetEstatusByModulo(ctx, moduloID)
	})
}

func (s *ServiceEstatus) CreateEstatus(ctx context.Context, estatus models.Estatus) (*models.Estatus, error) {
	if estatus.StdDescripcion == "" {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de creación de estatus con descripción vacía")
		return nil, errors.New("la descripción es obligatoria")
	}
	if estatus.StdTipoEstado == "" {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de creación de estatus sin tipo de estado", zap.String("descripcion", estatus.StdDescripcion))
		return nil, errors.New("el tipo de estado es obligatorio")
	}
	if estatus.MdlID == 0 {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de creación de estatus sin ID de módulo", zap.String("descripcion", estatus.StdDescripcion))
		return nil, errors.New("el módulo (mdl_id) es obligatorio")
	}

	result, err := s.store.CreateEstatus(ctx, &estatus)
	if err != nil {
		return nil, err
	}

	s.cache.Invalidate(ctx, "estatus:")

	return result, nil
}

func (s *ServiceEstatus) UpdateEstatus(ctx context.Context, id uuid.UUID, estatus models.Estatus) (*models.Estatus, error) {
	if estatus.StdDescripcion == "" {
		zaplogger.WithContext(ctx, s.logger).Warn("Intento de actualización de estatus con descripción vacía", zap.String("id", id.String()))
		return nil, errors.New("la descripción es obligatoria")
	}

	result, err := s.store.UpdateEstatus(ctx, id, &estatus)
	if err != nil {
		return nil, err
	}

	s.cache.Invalidate(ctx, "estatus:")

	return result, nil
}

func (s *ServiceEstatus) DeleteEstatus(ctx context.Context, id uuid.UUID) error {
	if err := s.store.DeleteEstatus(ctx, id); err != nil {
		return err
	}

	s.cache.Invalidate(ctx, "estatus:")

	return nil
}
