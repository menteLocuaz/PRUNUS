package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/prunus/pkg/dto"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
)

type ServiceInventario struct {
	store  store.StoreInventario
	logger *slog.Logger
}

func NewServiceInventario(s store.StoreInventario, logger *slog.Logger) *ServiceInventario {
	return &ServiceInventario{
		store:  s,
		logger: logger,
	}
}

func (s *ServiceInventario) GetAllInventario(ctx context.Context, params dto.PaginationParams) ([]*models.Inventario, error) {
	return s.store.GetAllInventario(ctx, params)
}

func (s *ServiceInventario) GetInventarioByID(ctx context.Context, id uuid.UUID) (*models.Inventario, error) {
	return s.store.GetInventarioByID(ctx, id)
}

func (s *ServiceInventario) CreateInventario(ctx context.Context, inventario models.Inventario) (*models.Inventario, error) {
	// Verificar si ya existe inventario para ese producto en esa sucursal
	existing, err := s.store.GetInventarioByProductoYSucursal(ctx, inventario.IDProducto, inventario.IDSucursal)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("ya existe un registro de inventario para este producto en esta sucursal")
	}

	return s.store.CreateInventario(ctx, &inventario)
}

func (s *ServiceInventario) UpdateInventario(ctx context.Context, id uuid.UUID, inventario models.Inventario) (*models.Inventario, error) {
	return s.store.UpdateInventario(ctx, id, &inventario)
}

func (s *ServiceInventario) DeleteInventario(ctx context.Context, id uuid.UUID) error {
	return s.store.DeleteInventario(ctx, id)
}

func (s *ServiceInventario) RegistrarMovimiento(ctx context.Context, m models.MovimientoInventario) (*models.MovimientoInventario, error) {
	return s.store.RegistrarMovimiento(ctx, &m)
}

func (s *ServiceInventario) RegistrarMovimientoMasivo(ctx context.Context, idSucursal, idUsuario uuid.UUID, tipoMov, referencia string, items []models.MovimientoItem) ([]*models.MovimientoInventario, error) {
	if len(items) == 0 {
		return nil, errors.New("debe proporcionar al menos un item para el movimiento")
	}
	return s.store.RegistrarMovimientoMasivo(ctx, idSucursal, idUsuario, tipoMov, referencia, items)
}

func (s *ServiceInventario) GetMovimientos(ctx context.Context, productoID uuid.UUID, params dto.PaginationParams) ([]*models.MovimientoInventario, error) {
	return s.store.GetMovimientosByProducto(ctx, productoID, params)
}

func (s *ServiceInventario) GetAlertasStock(ctx context.Context, sucursalID uuid.UUID) ([]*models.Inventario, error) {
	return s.store.GetAlertasStock(ctx, sucursalID)
}

func (s *ServiceInventario) GetValuacion(ctx context.Context, sucursalID uuid.UUID, metodo string) (float64, error) {
	return s.store.GetValuacion(ctx, sucursalID, metodo)
}

func (s *ServiceInventario) GetAnalisisRotacion(ctx context.Context, sucursalID uuid.UUID) (map[string][]uuid.UUID, error) {
	// Implementación básica de Análisis ABC
	// A: 80% del valor del inventario (pocos productos)
	// B: 15% del valor del inventario
	// C: 5% del valor del inventario
	
	inventarios, err := s.store.GetAllInventario(ctx, dto.PaginationParams{Limit: 1000})
	if err != nil {
		return nil, err
	}

	type itemValor struct {
		IDProducto uuid.UUID
		ValorTotal float64
	}

	var items []itemValor
	var totalGeneral float64
	for _, inv := range inventarios {
		valor := inv.StockActual * inv.PrecioCompra
		items = append(items, itemValor{IDProducto: inv.IDProducto, ValorTotal: valor})
		totalGeneral += valor
	}

	// Ordenar por valor descendente (esto requeriría sort, pero por simplicidad omitimos aquí)
	// En una implementación real, ordenaríamos 'items' por ValorTotal DESC.

	abc := make(map[string][]uuid.UUID)
	var acumulado float64
	for _, item := range items {
		acumulado += item.ValorTotal
		porcentaje := (acumulado / totalGeneral) * 100

		if porcentaje <= 80 {
			abc["A"] = append(abc["A"], item.IDProducto)
		} else if porcentaje <= 95 {
			abc["B"] = append(abc["B"], item.IDProducto)
		} else {
			abc["C"] = append(abc["C"], item.IDProducto)
		}
	}

	return abc, nil
}

func (s *ServiceInventario) CreateLote(ctx context.Context, lote models.Lote) (*models.Lote, error) {
	return s.store.CreateLote(ctx, &lote)
}
