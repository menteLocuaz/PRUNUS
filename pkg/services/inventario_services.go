package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/prunus/pkg/dto"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
	"go.uber.org/zap"
)

type ServiceInventario struct {
	store  store.StoreInventario
	logger *zap.Logger
}

func NewServiceInventario(s store.StoreInventario, logger *zap.Logger) *ServiceInventario {
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

func (s *ServiceInventario) GetInventarioBySucursal(ctx context.Context, idSucursal uuid.UUID, params dto.PaginationParams) ([]*models.Inventario, error) {
	return s.store.GetInventarioBySucursal(ctx, idSucursal, params)
}

func (s *ServiceInventario) CreateInventario(ctx context.Context, inventario models.Inventario) (*models.Inventario, error) {
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
	return s.store.GetAnalisisRotacion(ctx, sucursalID)
}

func (s *ServiceInventario) CreateLote(ctx context.Context, lote models.Lote) (*models.Lote, error) {
	return s.store.CreateLote(ctx, &lote)
}

func (s *ServiceInventario) AdjustStock(ctx context.Context, idInventario uuid.UUID, delta float64) error {
	return s.store.AdjustStock(ctx, idInventario, delta)
}

func (s *ServiceInventario) AdjustLoteStock(ctx context.Context, idLote uuid.UUID, delta float64) error {
	return s.store.AdjustLoteStock(ctx, idLote, delta)
}

func (s *ServiceInventario) GetRotacionDetalle(ctx context.Context, sucursalID uuid.UUID, params dto.RotacionFiltroParams) ([]*dto.RotacionProductoResponse, error) {
	if params.FechaFin.Before(params.FechaInicio) {
		return nil, errors.New("fecha_fin debe ser posterior a fecha_inicio")
	}
	return s.store.GetRotacionDetalle(ctx, sucursalID, params)
}

func (s *ServiceInventario) GetComposicionCategoria(ctx context.Context, sucursalID uuid.UUID) ([]*dto.ComposicionCategoriaResponse, error) {
	return s.store.GetComposicionCategoria(ctx, sucursalID)
}

func (s *ServiceInventario) GetAlertasStockDetalle(ctx context.Context, sucursalID uuid.UUID) ([]*dto.AlertaStockResponse, error) {
	return s.store.GetAlertasStockDetalle(ctx, sucursalID)
}

func (s *ServiceInventario) CapturarSnapshotInventario(ctx context.Context, sucursalID uuid.UUID) error {
	return s.store.CapturarSnapshotInventario(ctx, sucursalID)
}

func (s *ServiceInventario) GetValorHistorico(ctx context.Context, sucursalID uuid.UUID, params dto.RotacionFiltroParams) ([]*dto.ValorHistoricoResponse, error) {
	if params.FechaFin.Before(params.FechaInicio) {
		return nil, errors.New("fecha_fin debe ser posterior a fecha_inicio")
	}
	return s.store.GetValorHistorico(ctx, sucursalID, params)
}

func (s *ServiceInventario) GetPerdidas(ctx context.Context, sucursalID uuid.UUID, params dto.RotacionFiltroParams) ([]*dto.PerdidaResponse, error) {
	if params.FechaFin.Before(params.FechaInicio) {
		return nil, errors.New("fecha_fin debe ser posterior a fecha_inicio")
	}
	return s.store.GetPerdidas(ctx, sucursalID, params)
}

func (s *ServiceInventario) GetMargenGanancia(ctx context.Context, sucursalID uuid.UUID, params dto.RotacionFiltroParams) ([]*dto.MargenProductoResponse, error) {
	if params.FechaFin.Before(params.FechaInicio) {
		return nil, errors.New("fecha_fin debe ser posterior a fecha_inicio")
	}
	return s.store.GetMargenGanancia(ctx, sucursalID, params)
}
