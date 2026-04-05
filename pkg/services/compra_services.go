package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/store"
)

type ServiceCompra struct {
	store      store.StoreCompra
	invService *ServiceInventario
}

func NewServiceCompra(s store.StoreCompra, inv *ServiceInventario) *ServiceCompra {
	return &ServiceCompra{
		store:      s,
		invService: inv,
	}
}

func (s *ServiceCompra) CreateOrden(ctx context.Context, o *models.OrdenCompra) (*models.OrdenCompra, error) {
	return s.store.CreateOrden(ctx, o)
}

func (s *ServiceCompra) GetAllOrdenes(ctx context.Context) ([]*models.OrdenCompra, error) {
	return s.store.GetAllOrdenes(ctx)
}

func (s *ServiceCompra) GetOrdenByID(ctx context.Context, id uuid.UUID) (*models.OrdenCompra, error) {
	return s.store.GetOrdenByID(ctx, id)
}

func (s *ServiceCompra) ProcesarRecepcion(ctx context.Context, idCompra uuid.UUID, statusID uuid.UUID, items []*models.DetalleOrdenCompra, userID uuid.UUID) error {
	// 1. Obtener la orden para validar sucursal y existencia
	orden, err := s.store.GetOrdenByID(ctx, idCompra)
	if err != nil {
		return fmt.Errorf("orden de compra no encontrada: %w", err)
	}

	now := time.Now()
	// 2. Por cada item recibido, registrar movimiento de inventario
	for _, item := range items {
		if item.CantidadRecibida <= 0 {
			continue
		}

		// Actualizar cantidad recibida en el detalle de la compra
		if err := s.store.UpdateDetalleRecepcion(ctx, item.IDDetalleCompra, item.CantidadRecibida); err != nil {
			return fmt.Errorf("error al actualizar detalle de recepción: %w", err)
		}

		// Registrar el abastecimiento en el inventario
		mov := models.MovimientoInventario{
			IDProducto:     item.IDProducto,
			IDSucursal:     orden.IDSucursal,
			TipoMovimiento: "ENTRADA", // Abastecimiento por compra
			Cantidad:       item.CantidadRecibida,
			IDUsuario:      userID,
			Referencia:     fmt.Sprintf("Recepción OC: %s", orden.NumeroOrden),
		}

		if _, err := s.invService.RegistrarMovimiento(ctx, mov); err != nil {
			return fmt.Errorf("error al registrar abastecimiento de producto %s: %w", item.IDProducto, err)
		}

		// CREAR LOTE PARA TRAZABILIDAD
		lote := models.Lote{
			IDProducto:      item.IDProducto,
			IDSucursal:      orden.IDSucursal,
			CodigoLote:      fmt.Sprintf("%s-%s", orden.NumeroOrden, item.IDProducto.String()[:8]),
			CantidadInicial: item.CantidadRecibida,
			CantidadActual:  item.CantidadRecibida,
			CostoCompra:     item.PrecioUnitario,
			FechaRecepcion:  now,
		}

		if _, err := s.invService.CreateLote(ctx, lote); err != nil {
			return fmt.Errorf("error al crear lote para producto %s: %w", item.IDProducto, err)
		}
	}

	// 3. Actualizar estado de la orden a "RECIBIDO"
	return s.store.UpdateStatus(ctx, idCompra, statusID, &now)
}
