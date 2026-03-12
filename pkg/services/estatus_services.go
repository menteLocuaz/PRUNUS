package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
)

type ServiceEstatus struct {
	store  store.StoreEstatus
	cache  models.CacheStore
	logger *slog.Logger
}

func NewServiceEstatus(s store.StoreEstatus, c models.CacheStore, logger *slog.Logger) *ServiceEstatus {
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
	cacheExpirationEstatus  = 24 * time.Hour // Los estatus no cambian frecuentemente
)

// Nombres de módulos (puedes mover esto a un archivo de constantes global si es necesario)
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
	var catalog map[int]interface{}

	// Intentar caché
	err := s.cache.Get(ctx, cacheKeyEstatusMaster, &catalog)
	if err == nil {
		return catalog, nil
	}

	// Obtener todos
	all, err := s.store.GetAllEstatus(ctx)
	if err != nil {
		return nil, err
	}

	// Agrupar por módulo
	catalog = make(map[int]interface{})
	for _, e := range all {
		if _, ok := catalog[e.MdlID]; !ok {
			name := "Módulo Desconocido"
			if n, exists := moduloNames[e.MdlID]; exists {
				name = n
			}
			catalog[e.MdlID] = struct {
				Modulo string           `json:"modulo"`
				Items  []*models.Estatus `json:"items"`
			}{
				Modulo: name,
				Items:  []*models.Estatus{},
			}
		}

		// Hack para añadir al slice de una interfaz{} (o podrías usar structs tipados en el servicio)
		group := catalog[e.MdlID].(struct {
			Modulo string           `json:"modulo"`
			Items  []*models.Estatus `json:"items"`
		})
		group.Items = append(group.Items, e)
		catalog[e.MdlID] = group
	}

	// Guardar en caché
	_ = s.cache.Set(ctx, cacheKeyEstatusMaster, catalog, cacheExpirationEstatus)

	return catalog, nil
}

func (s *ServiceEstatus) GetAllEstatus(ctx context.Context) ([]*models.Estatus, error) {
	var estatusList []*models.Estatus

	// Intentar obtener del caché
	if err := s.cache.Get(ctx, cacheKeyEstatusAll, &estatusList); err == nil {
		return estatusList, nil
	}

	// Si no hay caché, ir a la base de datos
	estatusList, err := s.store.GetAllEstatus(ctx)
	if err != nil {
		return nil, err
	}

	// Guardar en caché (no bloqueante — ignoramos el error de cacheo)
	_ = s.cache.Set(ctx, cacheKeyEstatusAll, estatusList, cacheExpirationEstatus)

	return estatusList, nil
}

func (s *ServiceEstatus) GetEstatusByID(ctx context.Context, id uuid.UUID) (*models.Estatus, error) {
	var estatus *models.Estatus
	key := fmt.Sprintf(cacheKeyEstatusID, id.String())

	// Intentar obtener del caché
	if err := s.cache.Get(ctx, key, &estatus); err == nil {
		return estatus, nil
	}

	// Si no hay caché, ir a la base de datos
	estatus, err := s.store.GetEstatusByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Guardar en caché
	_ = s.cache.Set(ctx, key, estatus, cacheExpirationEstatus)

	return estatus, nil
}

func (s *ServiceEstatus) GetEstatusByTipo(ctx context.Context, tipo string) ([]*models.Estatus, error) {
	var estatusList []*models.Estatus
	key := fmt.Sprintf(cacheKeyEstatusTipo, tipo)

	// Intentar obtener del caché
	if err := s.cache.Get(ctx, key, &estatusList); err == nil {
		return estatusList, nil
	}

	// Si no hay caché, ir a la base de datos
	estatusList, err := s.store.GetEstatusByTipo(ctx, tipo)
	if err != nil {
		return nil, err
	}

	// Guardar en caché
	_ = s.cache.Set(ctx, key, estatusList, cacheExpirationEstatus)

	return estatusList, nil
}

func (s *ServiceEstatus) GetEstatusByModulo(ctx context.Context, moduloID int) ([]*models.Estatus, error) {
	var estatusList []*models.Estatus
	key := fmt.Sprintf(cacheKeyEstatusModuloID, moduloID)

	// Intentar obtener del caché
	if err := s.cache.Get(ctx, key, &estatusList); err == nil {
		return estatusList, nil
	}

	// Si no hay caché, ir a la base de datos
	estatusList, err := s.store.GetEstatusByModulo(ctx, moduloID)
	if err != nil {
		return nil, err
	}

	// Guardar en caché
	_ = s.cache.Set(ctx, key, estatusList, cacheExpirationEstatus)

	return estatusList, nil
}

func (s *ServiceEstatus) CreateEstatus(ctx context.Context, estatus models.Estatus) (*models.Estatus, error) {
	if estatus.StdDescripcion == "" {
		s.logger.WarnContext(ctx, "Intento de creación de estatus con descripción vacía")
		return nil, errors.New("la descripción es obligatoria")
	}
	if estatus.StpTipoEstado == "" {
		s.logger.WarnContext(ctx, "Intento de creación de estatus sin tipo de estado", slog.String("descripcion", estatus.StdDescripcion))
		return nil, errors.New("el tipo de estado es obligatorio")
	}
	if estatus.MdlID == 0 {
		s.logger.WarnContext(ctx, "Intento de creación de estatus sin ID de módulo", slog.String("descripcion", estatus.StdDescripcion))
		return nil, errors.New("el módulo (mdl_id) es obligatorio")
	}

	result, err := s.store.CreateEstatus(ctx, &estatus)
	if err != nil {
		return nil, err
	}

	// Invalidar caché
	s.invalidateCache(ctx, result)

	return result, nil
}

func (s *ServiceEstatus) UpdateEstatus(ctx context.Context, id uuid.UUID, estatus models.Estatus) (*models.Estatus, error) {
	if estatus.StdDescripcion == "" {
		s.logger.WarnContext(ctx, "Intento de actualización de estatus con descripción vacía", slog.String("id", id.String()))
		return nil, errors.New("la descripción es obligatoria")
	}

	result, err := s.store.UpdateEstatus(ctx, id, &estatus)
	if err != nil {
		return nil, err
	}

	// Invalidar caché
	s.invalidateCache(ctx, result)

	return result, nil
}

func (s *ServiceEstatus) DeleteEstatus(ctx context.Context, id uuid.UUID) error {
	// Obtener el registro antes de borrarlo para saber qué claves invalidar en el caché
	estatus, err := s.store.GetEstatusByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.store.DeleteEstatus(ctx, id); err != nil {
		return err
	}

	// Invalidar caché
	s.invalidateCache(ctx, estatus)

	return nil
}

// invalidateCache invalida las claves relacionadas con el estatus proporcionado.
func (s *ServiceEstatus) invalidateCache(ctx context.Context, e *models.Estatus) {
	if e == nil {
		return
	}
	_ = s.cache.Delete(ctx, cacheKeyEstatusAll)
	_ = s.cache.Delete(ctx, cacheKeyEstatusMaster)
	_ = s.cache.Delete(ctx, fmt.Sprintf(cacheKeyEstatusID, e.IDStatus.String()))
	_ = s.cache.Delete(ctx, fmt.Sprintf(cacheKeyEstatusTipo, e.StpTipoEstado))
	_ = s.cache.Delete(ctx, fmt.Sprintf(cacheKeyEstatusModuloID, e.MdlID))
}
