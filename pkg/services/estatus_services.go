package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
)

type ServiceEstatus struct {
	store store.StoreEstatus
	cache models.CacheStore
}

func NewServiceEstatus(s store.StoreEstatus, c models.CacheStore) *ServiceEstatus {
	return &ServiceEstatus{store: s, cache: c}
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

func (s *ServiceEstatus) GetMasterCatalog() (map[int]interface{}, error) {
	ctx := context.Background()
	var catalog map[int]interface{}

	// Intentar caché
	err := s.cache.Get(ctx, cacheKeyEstatusMaster, &catalog)
	if err == nil {
		return catalog, nil
	}

	// Obtener todos
	all, err := s.store.GetAllEstatus()
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

func (s *ServiceEstatus) GetAllEstatus() ([]*models.Estatus, error) {
	ctx := context.Background()
	var estatusList []*models.Estatus

	// Intentar obtener del caché
	if err := s.cache.Get(ctx, cacheKeyEstatusAll, &estatusList); err == nil {
		return estatusList, nil
	}

	// Si no hay caché, ir a la base de datos
	estatusList, err := s.store.GetAllEstatus()
	if err != nil {
		return nil, err
	}

	// Guardar en caché (no bloqueante — ignoramos el error de cacheo)
	_ = s.cache.Set(ctx, cacheKeyEstatusAll, estatusList, cacheExpirationEstatus)

	return estatusList, nil
}

func (s *ServiceEstatus) GetEstatusByID(id uuid.UUID) (*models.Estatus, error) {
	ctx := context.Background()
	var estatus *models.Estatus
	key := fmt.Sprintf(cacheKeyEstatusID, id.String())

	// Intentar obtener del caché
	if err := s.cache.Get(ctx, key, &estatus); err == nil {
		return estatus, nil
	}

	// Si no hay caché, ir a la base de datos
	estatus, err := s.store.GetEstatusByID(id)
	if err != nil {
		return nil, err
	}

	// Guardar en caché
	_ = s.cache.Set(ctx, key, estatus, cacheExpirationEstatus)

	return estatus, nil
}

func (s *ServiceEstatus) GetEstatusByTipo(tipo string) ([]*models.Estatus, error) {
	ctx := context.Background()
	var estatusList []*models.Estatus
	key := fmt.Sprintf(cacheKeyEstatusTipo, tipo)

	// Intentar obtener del caché
	if err := s.cache.Get(ctx, key, &estatusList); err == nil {
		return estatusList, nil
	}

	// Si no hay caché, ir a la base de datos
	estatusList, err := s.store.GetEstatusByTipo(tipo)
	if err != nil {
		return nil, err
	}

	// Guardar en caché
	_ = s.cache.Set(ctx, key, estatusList, cacheExpirationEstatus)

	return estatusList, nil
}

func (s *ServiceEstatus) GetEstatusByModulo(moduloID int) ([]*models.Estatus, error) {
	ctx := context.Background()
	var estatusList []*models.Estatus
	key := fmt.Sprintf(cacheKeyEstatusModuloID, moduloID)

	// Intentar obtener del caché
	if err := s.cache.Get(ctx, key, &estatusList); err == nil {
		return estatusList, nil
	}

	// Si no hay caché, ir a la base de datos
	estatusList, err := s.store.GetEstatusByModulo(moduloID)
	if err != nil {
		return nil, err
	}

	// Guardar en caché
	_ = s.cache.Set(ctx, key, estatusList, cacheExpirationEstatus)

	return estatusList, nil
}

func (s *ServiceEstatus) CreateEstatus(estatus models.Estatus) (*models.Estatus, error) {
	if estatus.StdDescripcion == "" {
		return nil, errors.New("la descripción es obligatoria")
	}
	if estatus.StpTipoEstado == "" {
		return nil, errors.New("el tipo de estado es obligatorio")
	}
	if estatus.MdlID == 0 {
		return nil, errors.New("el módulo (mdl_id) es obligatorio")
	}

	result, err := s.store.CreateEstatus(&estatus)
	if err != nil {
		return nil, err
	}

	// Invalidar caché
	s.invalidateCache(context.Background(), result)

	return result, nil
}

func (s *ServiceEstatus) UpdateEstatus(id uuid.UUID, estatus models.Estatus) (*models.Estatus, error) {
	if estatus.StdDescripcion == "" {
		return nil, errors.New("la descripción es obligatoria")
	}

	result, err := s.store.UpdateEstatus(id, &estatus)
	if err != nil {
		return nil, err
	}

	// Invalidar caché
	s.invalidateCache(context.Background(), result)

	return result, nil
}

func (s *ServiceEstatus) DeleteEstatus(id uuid.UUID) error {
	// Obtener el registro antes de borrarlo para saber qué claves invalidar en el caché
	estatus, err := s.store.GetEstatusByID(id)
	if err != nil {
		return err
	}

	if err := s.store.DeleteEstatus(id); err != nil {
		return err
	}

	// Invalidar caché
	s.invalidateCache(context.Background(), estatus)

	return nil
}

// invalidateCache invalida las claves relacionadas con el estatus proporcionado.
// Nota: el primer parámetro es del tipo context.Context (no context.Background()).
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
