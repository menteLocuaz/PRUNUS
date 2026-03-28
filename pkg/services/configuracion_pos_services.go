package services

import (
	"context"

	"github.com/prunus/pkg/dto"
	"github.com/prunus/pkg/store"
)

type ServiceConfiguracion struct {
	store store.StoreConfiguracion
}

func NewServiceConfiguracion(s store.StoreConfiguracion) *ServiceConfiguracion {
	return &ServiceConfiguracion{store: s}
}

func (s *ServiceConfiguracion) ListarCanales(ctx context.Context, chainID int) ([]dto.ConfigResponseDTO, error) {
	canales, err := s.store.GetCanalesImpresionActivos(ctx, chainID)
	if err != nil {
		return nil, err
	}

	var res []dto.ConfigResponseDTO
	for _, c := range canales {
		res = append(res, dto.ConfigResponseDTO{ID: c.IDCanalImpresion, Descripcion: c.Descripcion})
	}
	return res, nil
}

func (s *ServiceConfiguracion) ListarImpresoras(ctx context.Context, restaurantID int) ([]dto.ConfigResponseDTO, error) {
	impresoras, err := s.store.GetImpresorasActivas(ctx, restaurantID)
	if err != nil {
		return nil, err
	}

	var res []dto.ConfigResponseDTO
	for _, i := range impresoras {
		res = append(res, dto.ConfigResponseDTO{ID: i.IDImpresora, Descripcion: i.Nombre})
	}
	return res, nil
}

func (s *ServiceConfiguracion) ListarPuertos(ctx context.Context) ([]dto.ConfigResponseDTO, error) {
	puertos, err := s.store.GetPuertosActivos(ctx)
	if err != nil {
		return nil, err
	}

	var res []dto.ConfigResponseDTO
	for _, p := range puertos {
		res = append(res, dto.ConfigResponseDTO{ID: p.IDPuertos, Descripcion: p.Descripcion})
	}
	return res, nil
}
