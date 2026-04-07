package services

import (
	"context"

	"github.com/prunus/pkg/dto"
	"github.com/prunus/pkg/store"

	"github.com/google/uuid"
)

type ServiceDashboard interface {
	GetDashboardData(ctx context.Context, sucursalID uuid.UUID) (*dto.DashboardResumen, error)
	GetAntiguedadDeuda(ctx context.Context, sucursalID uuid.UUID) ([]dto.AntiguedadDeudaDTO, error)
}

type dashboardService struct {
	store store.StoreDashboard
}

func NewDashboardService(s store.StoreDashboard) ServiceDashboard {
	return &dashboardService{store: s}
}

func (s *dashboardService) GetDashboardData(ctx context.Context, sucursalID uuid.UUID) (*dto.DashboardResumen, error) {
	resumen, err := s.store.GetResumen(ctx, sucursalID)
	if err != nil {
		return nil, err
	}

	topRentables, err := s.store.GetRentabilidadTop(ctx, sucursalID)
	if err == nil {
		resumen.TopProductos = topRentables
	}

	ventasVsCompras, err := s.store.GetVentasVsCompras(ctx, sucursalID)
	if err == nil {
		resumen.VentasVsCompras = ventasVsCompras
	}

	return resumen, nil
}

func (s *dashboardService) GetAntiguedadDeuda(ctx context.Context, sucursalID uuid.UUID) ([]dto.AntiguedadDeudaDTO, error) {
	return s.store.GetAntiguedadDeuda(ctx, sucursalID)
}
