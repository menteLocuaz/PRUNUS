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

// StoreCaja define la interfaz para el almacenamiento de cajas y sesiones.
type StoreCaja interface {
	// Cajas
	GetAllCajas(ctx context.Context) ([]*models.Caja, error)
	GetCajaByID(ctx context.Context, id uuid.UUID) (*models.Caja, error)
	GetCajasBySucursal(ctx context.Context, idSucursal uuid.UUID) ([]*models.Caja, error)
	CreateCaja(ctx context.Context, c *models.Caja) (*models.Caja, error)
	UpdateCaja(ctx context.Context, id uuid.UUID, c *models.Caja) (*models.Caja, error)
	DeleteCaja(ctx context.Context, id uuid.UUID) error

	// Sesiones
	CreateSesion(ctx context.Context, s *models.SesionCaja) (*models.SesionCaja, error)
	GetSesionByID(ctx context.Context, id uuid.UUID) (*models.SesionCaja, error)
	GetSesionActivaByUsuario(ctx context.Context, idUsuario uuid.UUID) (*models.SesionCaja, error)
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

// NewCajaStore crea una nueva instancia del almacén de cajas.
func NewCajaStore(db *sql.DB) StoreCaja {
	return &storeCaja{db: db}
}

// --- Helpers de Escaneo ---

const cajaSelectFields = `id_caja, nombre, id_sucursal, estado, created_at, updated_at`

func (s *storeCaja) scanRowCaja(scanner interface{ Scan(dest ...any) error }, c *models.Caja) error {
	return scanner.Scan(&c.IDCaja, &c.Nombre, &c.IDSucursal, &c.Estado, &c.CreatedAt, &c.UpdatedAt)
}

const sesionCajaSelectFields = `id_sesion, id_caja, id_usuario, monto_apertura, monto_cierre, fecha_apertura, fecha_cierre, estado, created_at, updated_at`

func (s *storeCaja) scanRowSesion(scanner interface{ Scan(dest ...any) error }, sesion *models.SesionCaja) error {
	return scanner.Scan(
		&sesion.IDSesion, &sesion.IDCaja, &sesion.IDUsuario, &sesion.MontoApertura, &sesion.MontoCierre,
		&sesion.FechaApertura, &sesion.FechaCierre, &sesion.Estado, &sesion.CreatedAt, &sesion.UpdatedAt,
	)
}

const movimientoCajaSelectFields = `id_movimiento, id_sesion, tipo, monto, motivo, created_at, updated_at`

func (s *storeCaja) scanRowMovimiento(scanner interface{ Scan(dest ...any) error }, m *models.MovimientoCaja) error {
	return scanner.Scan(&m.IDMovimiento, &m.IDSesion, &m.Tipo, &m.Monto, &m.Motivo, &m.CreatedAt, &m.UpdatedAt)
}

// --- Implementación Cajas ---

func (s *storeCaja) GetAllCajas(ctx context.Context) ([]*models.Caja, error) {
	defer performance.Trace(ctx, "store", "GetAllCajas", performance.DbThreshold, time.Now())
	query := fmt.Sprintf(`SELECT %s FROM caja WHERE deleted_at IS NULL`, cajaSelectFields)
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener cajas: %w", err)
	}
	defer rows.Close()

	var cajas []*models.Caja
	for rows.Next() {
		c := &models.Caja{}
		if err := s.scanRowCaja(rows, c); err != nil {
			return nil, fmt.Errorf("error al escanear caja: %w", err)
		}
		cajas = append(cajas, c)
	}
	return cajas, nil
}

func (s *storeCaja) GetCajaByID(ctx context.Context, id uuid.UUID) (*models.Caja, error) {
	defer performance.Trace(ctx, "store", "GetCajaByID", performance.DbThreshold, time.Now())
	query := fmt.Sprintf(`SELECT %s FROM caja WHERE id_caja = $1 AND deleted_at IS NULL`, cajaSelectFields)
	c := &models.Caja{}
	err := s.scanRowCaja(s.db.QueryRowContext(ctx, query, id), c)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("caja no encontrada")
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener caja: %w", err)
	}
	return c, nil
}

func (s *storeCaja) GetCajasBySucursal(ctx context.Context, idSucursal uuid.UUID) ([]*models.Caja, error) {
	defer performance.Trace(ctx, "store", "GetCajasBySucursal", performance.DbThreshold, time.Now())
	query := fmt.Sprintf(`SELECT %s FROM caja WHERE id_sucursal = $1 AND deleted_at IS NULL`, cajaSelectFields)
	rows, err := s.db.QueryContext(ctx, query, idSucursal)
	if err != nil {
		return nil, fmt.Errorf("error al obtener cajas por sucursal: %w", err)
	}
	defer rows.Close()

	var cajas []*models.Caja
	for rows.Next() {
		c := &models.Caja{}
		if err := s.scanRowCaja(rows, c); err != nil {
			return nil, fmt.Errorf("error al escanear caja: %w", err)
		}
		cajas = append(cajas, c)
	}
	return cajas, nil
}

func (s *storeCaja) CreateCaja(ctx context.Context, c *models.Caja) (*models.Caja, error) {
	defer performance.Trace(ctx, "store", "CreateCaja", performance.DbThreshold, time.Now())
	err := ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		query := `INSERT INTO caja (nombre, id_sucursal, estado) VALUES ($1, $2, $3) RETURNING id_caja, created_at, updated_at`
		return tx.QueryRowContext(ctx, query, c.Nombre, c.IDSucursal, c.Estado).Scan(&c.IDCaja, &c.CreatedAt, &c.UpdatedAt)
	})
	if err != nil {
		return nil, fmt.Errorf("error al crear caja: %w", err)
	}
	return c, nil
}

func (s *storeCaja) UpdateCaja(ctx context.Context, id uuid.UUID, c *models.Caja) (*models.Caja, error) {
	defer performance.Trace(ctx, "store", "UpdateCaja", performance.DbThreshold, time.Now())
	err := ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		query := `UPDATE caja SET nombre = $1, id_sucursal = $2, estado = $3, updated_at = CURRENT_TIMESTAMP 
		          WHERE id_caja = $4 AND deleted_at IS NULL 
		          RETURNING created_at, updated_at`
		return tx.QueryRowContext(ctx, query, c.Nombre, c.IDSucursal, c.Estado, id).Scan(&c.CreatedAt, &c.UpdatedAt)
	})
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("caja no encontrada")
	}
	if err != nil {
		return nil, fmt.Errorf("error al actualizar caja: %w", err)
	}
	c.IDCaja = id
	return c, nil
}

func (s *storeCaja) DeleteCaja(ctx context.Context, id uuid.UUID) error {
	defer performance.Trace(ctx, "store", "DeleteCaja", performance.DbThreshold, time.Now())
	return ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		query := `UPDATE caja SET deleted_at = CURRENT_TIMESTAMP WHERE id_caja = $1 AND deleted_at IS NULL`
		result, err := tx.ExecContext(ctx, query, id)
		if err != nil {
			return err
		}
		rows, _ := result.RowsAffected()
		if rows == 0 {
			return fmt.Errorf("caja no encontrada")
		}
		return nil
	})
}

// --- Implementación Sesiones ---

func (s *storeCaja) CreateSesion(ctx context.Context, sesion *models.SesionCaja) (*models.SesionCaja, error) {
	defer performance.Trace(ctx, "store", "CreateSesion", performance.DbThreshold, time.Now())
	err := ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		query := `INSERT INTO sesion_caja (id_caja, id_usuario, monto_apertura, fecha_apertura, estado) 
		          VALUES ($1, $2, $3, CURRENT_TIMESTAMP, 'ABIERTA') 
		          RETURNING id_sesion, fecha_apertura, created_at, updated_at`
		err := tx.QueryRowContext(ctx, query, sesion.IDCaja, sesion.IDUsuario, sesion.MontoApertura).Scan(
			&sesion.IDSesion, &sesion.FechaApertura, &sesion.CreatedAt, &sesion.UpdatedAt,
		)
		sesion.Estado = "ABIERTA"
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("error al abrir sesión de caja: %w", err)
	}
	return sesion, nil
}

func (s *storeCaja) GetSesionByID(ctx context.Context, id uuid.UUID) (*models.SesionCaja, error) {
	defer performance.Trace(ctx, "store", "GetSesionByID", performance.DbThreshold, time.Now())
	query := fmt.Sprintf(`SELECT %s FROM sesion_caja WHERE id_sesion = $1 AND deleted_at IS NULL`, sesionCajaSelectFields)
	sesion := &models.SesionCaja{}
	err := s.scanRowSesion(s.db.QueryRowContext(ctx, query, id), sesion)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("sesión de caja no encontrada")
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener sesión de caja: %w", err)
	}
	return sesion, nil
}

func (s *storeCaja) GetSesionActivaByUsuario(ctx context.Context, idUsuario uuid.UUID) (*models.SesionCaja, error) {
	defer performance.Trace(ctx, "store", "GetSesionActivaByUsuario", performance.DbThreshold, time.Now())
	query := fmt.Sprintf(`SELECT %s FROM sesion_caja WHERE id_usuario = $1 AND estado = 'ABIERTA' AND deleted_at IS NULL LIMIT 1`, sesionCajaSelectFields)
	sesion := &models.SesionCaja{}
	err := s.scanRowSesion(s.db.QueryRowContext(ctx, query, idUsuario), sesion)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error al buscar sesión activa: %w", err)
	}
	return sesion, nil
}

func (s *storeCaja) UpdateSesion(ctx context.Context, id uuid.UUID, sesion *models.SesionCaja) (*models.SesionCaja, error) {
	defer performance.Trace(ctx, "store", "UpdateSesion", performance.DbThreshold, time.Now())
	err := ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		query := `UPDATE sesion_caja SET monto_cierre = $1, fecha_cierre = $2, estado = $3, updated_at = CURRENT_TIMESTAMP 
		          WHERE id_sesion = $4 AND deleted_at IS NULL 
		          RETURNING created_at, updated_at`
		return tx.QueryRowContext(ctx, query, sesion.MontoCierre, sesion.FechaCierre, sesion.Estado, id).Scan(&sesion.CreatedAt, &sesion.UpdatedAt)
	})
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("sesión de caja no encontrada")
	}
	if err != nil {
		return nil, fmt.Errorf("error al actualizar sesión de caja: %w", err)
	}
	sesion.IDSesion = id
	return sesion, nil
}

// --- Implementación Movimientos ---

func (s *storeCaja) CreateMovimiento(ctx context.Context, m *models.MovimientoCaja) (*models.MovimientoCaja, error) {
	defer performance.Trace(ctx, "store", "CreateMovimiento", performance.DbThreshold, time.Now())
	err := ExecAudited(ctx, s.db, func(tx *sql.Tx) error {
		query := `INSERT INTO movimiento_caja (id_sesion, tipo, monto, motivo) VALUES ($1, $2, $3, $4) 
		          RETURNING id_movimiento, created_at, updated_at`
		return tx.QueryRowContext(ctx, query, m.IDSesion, m.Tipo, m.Monto, m.Motivo).Scan(&m.IDMovimiento, &m.CreatedAt, &m.UpdatedAt)
	})
	if err != nil {
		return nil, fmt.Errorf("error al crear movimiento de caja: %w", err)
	}
	return m, nil
}

func (s *storeCaja) GetMovimientosBySesion(ctx context.Context, sesionID uuid.UUID) ([]*models.MovimientoCaja, error) {
	defer performance.Trace(ctx, "store", "GetMovimientosBySesion", performance.DbThreshold, time.Now())
	query := fmt.Sprintf(`SELECT %s FROM movimiento_caja WHERE id_sesion = $1 AND deleted_at IS NULL ORDER BY created_at ASC`, movimientoCajaSelectFields)
	rows, err := s.db.QueryContext(ctx, query, sesionID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener movimientos: %w", err)
	}
	defer rows.Close()

	var movimientos []*models.MovimientoCaja
	for rows.Next() {
		m := &models.MovimientoCaja{}
		if err := s.scanRowMovimiento(rows, m); err != nil {
			return nil, fmt.Errorf("error al escanear movimiento: %w", err)
		}
		movimientos = append(movimientos, m)
	}
	return movimientos, nil
}

// --- Implementación Arqueo ---

func (s *storeCaja) GetVentasEfectivoBySesion(ctx context.Context, sesionID uuid.UUID) (float64, error) {
	defer performance.Trace(ctx, "store", "GetVentasEfectivoBySesion", performance.DbThreshold, time.Now())

	// Sumar pagos en efectivo de facturas asociadas a esta sesión de caja.
	// NOTA: Se asocia vía id_control_estacion que en la lógica de negocio se mapea a la sesión activa.
	query := `
		SELECT COALESCE(SUM(fpf.monto), 0)
		FROM factura f
		JOIN forma_pago_factura fpf ON f.id_factura = fpf.id_factura
		JOIN forma_pago fp ON fpf.metodo_pago = fp.nombre
		WHERE f.id_control_estacion = $1 
		  AND fp.nombre ILIKE '%Efectivo%'
		  AND f.deleted_at IS NULL
	`
	var total float64
	err := s.db.QueryRowContext(ctx, query, sesionID).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("error al calcular ventas en efectivo: %w", err)
	}
	return total, nil
}
