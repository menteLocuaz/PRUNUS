package services

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/prunus/pkg/dto"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
)

type ServicePOS struct {
	store store.StorePOS
}

func NewServicePOS(s store.StorePOS) *ServicePOS {
	return &ServicePOS{store: s}
}

// AbrirCaja realiza la apertura de una estación de POS
func (s *ServicePOS) AbrirCaja(input dto.AbrirCajaDTO, idUsuario uuid.UUID) (*models.ControlEstacion, error) {
	// 1. Validar que la estación existe
	_, err := s.store.GetEstacionByID(input.IDEstacion)
	if err != nil {
		return nil, fmt.Errorf("estación no encontrada: %w", err)
	}

	// 2. Validar que no haya una sesión activa en esta estación
	activo, err := s.store.GetActiveControlByEstacion(input.IDEstacion)
	if err != nil {
		return nil, err
	}
	if activo != nil {
		return nil, errors.New("la estación ya tiene una sesión abierta")
	}

	// 3. Validar que exista un periodo activo
	periodo, err := s.store.GetActivePeriodo()
	if err != nil {
		return nil, fmt.Errorf("error al validar periodo: %w", err)
	}

	// 4. Crear el registro de control de estación
	control := &models.ControlEstacion{
		IDEstacion:      input.IDEstacion,
		FondoBase:       input.FondoBase,
		UsuarioAsignado: idUsuario, // Usuario que realiza la apertura (ej. Admin o Gerente)
		IDStatus:        models.EstatusFondoAsignado,
		IDUserPos:       input.IDUserPos, // Usuario que operará el POS
		IDPeriodo:       periodo.IDPeriodo,
	}

	result, err := s.store.CreateControlEstacion(control)
	if err != nil {
		return nil, err
	}

	// 5. Actualizar estatus de la estación a 'Fondo Asignado' o similar si aplica
	err = s.store.UpdateEstacionStatus(input.IDEstacion, models.EstatusFondoAsignado)
	if err != nil {
		// Log error but don't fail the whole operation?
		// In a production app, we might want to use a transaction.
		fmt.Printf("Error al actualizar estatus de estación: %v\n", err)
	}

	return result, nil
}

// GetEstadoCaja obtiene el estado actual de una estación
func (s *ServicePOS) GetEstadoCaja(idEstacion uuid.UUID) (*dto.EstadoCajaDTO, error) {
	estacion, err := s.store.GetEstacionByID(idEstacion)
	if err != nil {
		return nil, err
	}

	control, err := s.store.GetActiveControlByEstacion(idEstacion)
	if err != nil {
		return nil, err
	}

	if control == nil {
		return &dto.EstadoCajaDTO{
			NombreEstacion:    estacion.Nombre,
			StatusDescripcion: "Cerrada",
		}, nil
	}

	return &dto.EstadoCajaDTO{
		IDControlEstacion: control.IDControlEstacion,
		NombreEstacion:    estacion.Nombre,
		FondoBase:         control.FondoBase,
		IDStatus:          control.IDStatus,
		FechaInicio:       control.FechaInicio,
	}, nil
}
