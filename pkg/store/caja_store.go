package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/utils/performance"
)

type StoreCaja interface {
	GetAllCajas(ctx context.Context) ([]*models.Caja, error)
	GetCajaByID(ctx context.Context, id uuid.UUID) (*models.Caja, error)
	CreateCaja(ctx context.Context, c *models.Caja) (*models.Caja, error)
	UpdateCaja(ctx context.Context, id uuid.UUID, c *models.Caja) (*models.Caja, error)
	DeleteCaja(ctx context.Context, id uuid.UUID) error

	// Sesiones
	CreateSesion(ctx context.Context, s *models.SesionCaja) (*models.SesionCaja, error)
	GetSesionByID(ctx context.Context, id uuid.UUID) (*models.SesionCaja, error)
	UpdateSesion(ctx context.Context, id uuid.UUID, s *models.SesionCaja) (*models.SesionCaja, error)

	// Movimientos
	CreateMovimiento(ctx context.Context, m *models.MovimientoCaja) (*models.MovimientoCaja, error)
	GetMovimientosBySesion(ctx context.Context, sesionID uuid.UUID) ([]*models.MovimientoCaja, error)

	// Arqueo
	GetVentasEfectivoBySesion(ctx context.Context, sesionID uuid.UUID) (float64, error)
}

type storeCaja struct {
	db *sql.DB
}

func NewCaja(db *sql.DB) StoreCaja {
	return &storeCaja{db: db}
}

func (s *storeCaja) GetAllCajas(ctx context.Context) ([]*models.Caja, error) {
	defer performance.Trace(ctx, "store", "GetAllCajas", performance.DbThreshold, time.Now())
	query := `SELECT id_caja, nombre, id_sucursal, estado, created_at, updated_at FROM caja WHERE deleted_at IS NULL`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cajas []*models.Caja
	for rows.Next() {
		c := &models.Caja{}
		if err := rows.Scan(&c.IDCaja, &c.Nombre, &c.IDSucursal, &c.Estado, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		cajas = append(cajas, c)
	}
	return cajas, nil
}

func (s *storeCaja) GetCajaByID(ctx context.Context, id uuid.UUID) (*models.Caja, error) {
	query := `SELECT id_caja, nombre, id_sucursal, estado, created_at, updated_at FROM caja WHERE id_caja = $1 AND deleted_at IS NULL`
	c := &models.Caja{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(&c.IDCaja, &c.Nombre, &c.IDSucursal, &c.Estado, &c.CreatedAt, &c.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("caja no encontrada")
	}
	return c, err
}

func (s *storeCaja) CreateCaja(ctx context.Context, c *models.Caja) (*models.Caja, error) {
	query := `INSERT INTO caja (nombre, id_sucursal, estado) VALUES ($1, $2, $3) RETURNING id_caja, created_at, updated_at`
	err := s.db.QueryRowContext(ctx, query, c.Nombre, c.IDSucursal, c.Estado).Scan(&c.IDCaja, &c.CreatedAt, &c.UpdatedAt)
	return c, err
}

func (s *storeCaja) UpdateCaja(ctx context.Context, id uuid.UUID, c *models.Caja) (*models.Caja, error) {
	query := `UPDATE caja SET nombre = $1, id_sucursal = $2, estado = $3, updated_at = CURRENT_TIMESTAMP WHERE id_caja = $4 AND deleted_at IS NULL RETURNING id_caja, created_at, updated_at`
	err := s.db.QueryRowContext(ctx, query, c.Nombre, c.IDSucursal, c.Estado, id).Scan(&c.IDCaja, &c.CreatedAt, &c.UpdatedAt)
	return c, err
}

func (s *storeCaja) DeleteCaja(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE caja SET deleted_at = CURRENT_TIMESTAMP WHERE id_caja = $1 AND deleted_at IS NULL`
	_, err := s.db.ExecContext(ctx, query, id)
	return err
}

func (s *storeCaja) CreateSesion(ctx context.Context, sesion *models.SesionCaja) (*models.SesionCaja, error) {
	query := `INSERT INTO sesion_caja (id_caja, id_usuario, monto_apertura, fecha_apertura, estado) VALUES ($1, $2, $3, $4, $5) RETURNING id_sesion, created_at, updated_at`
	err := s.db.QueryRowContext(ctx, query, sesion.IDCaja, sesion.IDUsuario, sesion.MontoApertura, time.Now(), "ABIERTA").Scan(&sesion.IDSesion, &sesion.CreatedAt, &sesion.UpdatedAt)
	sesion.Estado = "ABIERTA"
	sesion.FechaApertura = time.Now()
	return sesion, err
}

func (s *storeCaja) GetSesionByID(ctx context.Context, id uuid.UUID) (*models.SesionCaja, error) {
	query := `SELECT id_sesion, id_caja, id_usuario, monto_apertura, monto_cierre, fecha_apertura, fecha_cierre, estado, created_at, updated_at FROM sesion_caja WHERE id_sesion = $1 AND deleted_at IS NULL`
	sesion := &models.SesionCaja{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&sesion.IDSesion, &sesion.IDCaja, &sesion.IDUsuario, &sesion.MontoApertura, &sesion.MontoCierre,
		&sesion.FechaApertura, &sesion.FechaCierre, &sesion.Estado, &sesion.CreatedAt, &sesion.UpdatedAt,
	)
	return sesion, err
}

func (s *storeCaja) UpdateSesion(ctx context.Context, id uuid.UUID, sesion *models.SesionCaja) (*models.SesionCaja, error) {
	query := `UPDATE sesion_caja SET monto_cierre = $1, fecha_cierre = $2, estado = $3, updated_at = CURRENT_TIMESTAMP WHERE id_sesion = $4 AND deleted_at IS NULL RETURNING id_sesion, created_at, updated_at`
	err := s.db.QueryRowContext(ctx, query, sesion.MontoCierre, sesion.FechaCierre, sesion.Estado, id).Scan(&sesion.IDSesion, &sesion.CreatedAt, &sesion.UpdatedAt)
	return sesion, err
}

func (s *storeCaja) CreateMovimiento(ctx context.Context, m *models.MovimientoCaja) (*models.MovimientoCaja, error) {
	query := `INSERT INTO movimiento_caja (id_sesion, tipo, monto, motivo) VALUES ($1, $2, $3, $4) RETURNING id_movimiento, created_at, updated_at`
	err := s.db.QueryRowContext(ctx, query, m.IDSesion, m.Tipo, m.Monto, m.Motivo).Scan(&m.IDMovimiento, &m.CreatedAt, &m.UpdatedAt)
	return m, err
}

func (s *storeCaja) GetMovimientosBySesion(ctx context.Context, sesionID uuid.UUID) ([]*models.MovimientoCaja, error) {
	query := `SELECT id_movimiento, id_sesion, tipo, monto, motivo, created_at, updated_at FROM movimiento_caja WHERE id_sesion = $1 AND deleted_at IS NULL`
	rows, err := s.db.QueryContext(ctx, query, sesionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movimientos []*models.MovimientoCaja
	for rows.Next() {
		m := &models.MovimientoCaja{}
		if err := rows.Scan(&m.IDMovimiento, &m.IDSesion, &m.Tipo, &m.Monto, &m.Motivo, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		movimientos = append(movimientos, m)
	}
	return movimientos, nil
}

func (s *storeCaja) GetVentasEfectivoBySesion(ctx context.Context, sesionID uuid.UUID) (float64, error) {
	defer performance.Trace(ctx, "store", "GetVentasEfectivoBySesion", performance.DbThreshold, time.Now())
	
	// Sumar pagos en efectivo de facturas asociadas a esta sesión de caja
	// Unimos factura con forma_pago_factura y forma_pago
	query := `
		SELECT COALESCE(SUM(fpf.valor_billete), 0)
		FROM factura f
		JOIN forma_pago_factura fpf ON f.id_factura = fpf.id_factura
		JOIN forma_pago fp ON fpf.id_forma_pago = fp.id_forma_pago
		WHERE f.id_control_estacion = $1 
		  AND fp.fmp_codigo = 'EF'
		  AND f.deleted_at IS NULL
	`
	var total float64
	err := s.db.QueryRowContext(ctx, query, sesionID).Scan(&total)
	return total, err
}
