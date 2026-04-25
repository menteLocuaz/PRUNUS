package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
)

type StorePeriodo interface {
	CreatePeriodo(ctx context.Context, p *models.Periodo) (*models.Periodo, error)
	GetPeriodoByID(ctx context.Context, id uuid.UUID) (*models.Periodo, error)
	GetAllPeriodos(ctx context.Context) ([]*models.Periodo, error)
	UpdatePeriodo(ctx context.Context, id uuid.UUID, p *models.Periodo) (*models.Periodo, error)
	DeletePeriodo(ctx context.Context, id uuid.UUID) error

	// Metodos operativos
	GetActivePeriodo(ctx context.Context, sucursalID uuid.UUID) (*models.Periodo, error)
	CerrarPeriodo(ctx context.Context, id uuid.UUID, idUsuarioCierre uuid.UUID, ipCierre string) error

	// Auditoría y Cierre
	GenerarSnapshotPeriodo(ctx context.Context, idPeriodo uuid.UUID) (*models.PeriodoSnapshot, error)
	GuardarSnapshot(ctx context.Context, s *models.PeriodoSnapshot) error

	// Utils
	GetStatusIDByDesc(ctx context.Context, desc string) (uuid.UUID, error)
}

type PeriodoStore struct {
	db *sql.DB
}

func NewPeriodoStore(db *sql.DB) *PeriodoStore {
	return &PeriodoStore{db: db}
}

func (s *PeriodoStore) CreatePeriodo(ctx context.Context, p *models.Periodo) (*models.Periodo, error) {
	query := `INSERT INTO periodo (
				id_periodo, nombre, id_sucursal, prd_fecha_apertura, prd_usuario_apertura, 
				prd_ip_apertura, prd_motivo_apertura, id_status, created_at, updated_at
			  ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW()) RETURNING id_periodo`

	if p.IDPeriodo == uuid.Nil {
		p.IDPeriodo = uuid.New()
	}

	err := s.db.QueryRowContext(ctx, query,
		p.IDPeriodo, p.Nombre, p.IDSucursal, p.PrdFechaApertura, p.PrdUsuarioApertura,
		p.PrdIPApertura, p.PrdMotivoApertura, p.IDStatus,
	).Scan(&p.IDPeriodo)

	if err != nil {
		return nil, fmt.Errorf("error al crear periodo (posible duplicidad): %w", err)
	}
	return p, nil
}

func (s *PeriodoStore) GetActivePeriodo(ctx context.Context, sucursalID uuid.UUID) (*models.Periodo, error) {
	query := `SELECT id_periodo, nombre, id_sucursal, prd_fecha_apertura, prd_usuario_apertura, prd_ip_apertura, prd_motivo_apertura, id_status 
			  FROM periodo 
			  WHERE id_sucursal = $1 AND prd_fecha_cierre IS NULL AND deleted_at IS NULL 
			  LIMIT 1`

	var p models.Periodo
	err := s.db.QueryRowContext(ctx, query, sucursalID).Scan(
		&p.IDPeriodo, &p.Nombre, &p.IDSucursal, &p.PrdFechaApertura, &p.PrdUsuarioApertura,
		&p.PrdIPApertura, &p.PrdMotivoApertura, &p.IDStatus,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *PeriodoStore) CerrarPeriodo(ctx context.Context, id uuid.UUID, idUsuarioCierre uuid.UUID, ipCierre string) error {
	query := `UPDATE periodo SET prd_fecha_cierre = NOW(), prd_usuario_cierre = $1, prd_ip_cierre = $2, updated_at = NOW() 
			  WHERE id_periodo = $3 AND prd_fecha_cierre IS NULL`

	result, err := s.db.ExecContext(ctx, query, idUsuarioCierre, ipCierre, id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("periodo no encontrado o ya cerrado")
	}
	return nil
}

func (s *PeriodoStore) GenerarSnapshotPeriodo(ctx context.Context, idPeriodo uuid.UUID) (*models.PeriodoSnapshot, error) {
	query := `
		WITH ventas_metodos AS (
			SELECT 
				COALESCE(SUM(CASE WHEN fp.metodo_pago ILIKE '%Efectivo%' THEN fp.monto ELSE 0 END), 0) as efectivo,
				COALESCE(SUM(CASE WHEN fp.metodo_pago ILIKE '%Tarjeta%' THEN fp.monto ELSE 0 END), 0) as tarjeta,
				COALESCE(SUM(CASE WHEN fp.metodo_pago NOT ILIKE '%Efectivo%' AND fp.metodo_pago NOT ILIKE '%Tarjeta%' THEN fp.monto ELSE 0 END), 0) as otros,
				COUNT(DISTINCT f.id_factura) as conteo
			FROM factura f
			JOIN forma_pago_factura fp ON f.id_factura = fp.id_factura
			WHERE f.id_periodo = $1 AND f.deleted_at IS NULL
		),
		arqueos AS (
			SELECT 
				COALESCE(SUM(monto_cierre - (monto_apertura + COALESCE((SELECT SUM(monto) FROM movimiento_caja WHERE id_sesion = sc.id_sesion AND tipo = 'VENTA'), 0))), 0) as diferencia
			FROM sesion_caja sc
			WHERE EXISTS (SELECT 1 FROM factura WHERE id_control_estacion = sc.id_sesion AND id_periodo = $1)
		)
		SELECT 
			v.efectivo + v.tarjeta + v.otros as total_ventas,
			v.efectivo,
			v.tarjeta,
			v.otros,
			v.conteo,
			a.diferencia
		FROM ventas_metodos v, arqueos a
	`

	snapshot := &models.PeriodoSnapshot{IDPeriodo: idPeriodo}
	err := s.db.QueryRowContext(ctx, query, idPeriodo).Scan(
		&snapshot.TotalVentas,
		&snapshot.TotalEfectivo,
		&snapshot.TotalTarjeta,
		&snapshot.TotalOtros,
		&snapshot.TotalOperaciones,
		&snapshot.TotalDiferencias,
	)
	if err != nil {
		return nil, fmt.Errorf("error al consolidar snapshot de periodo: %w", err)
	}

	return snapshot, nil
}

func (s *PeriodoStore) GuardarSnapshot(ctx context.Context, snapshot *models.PeriodoSnapshot) error {
	dataJSON, _ := json.Marshal(snapshot.DataJSON)

	query := `INSERT INTO periodo_snapshot (
				id_periodo, total_ventas, total_efectivo, total_tarjeta, total_otros, 
				total_diferencias, total_operaciones, data_json, id_usuario_cierre
			  ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := s.db.ExecContext(ctx, query,
		snapshot.IDPeriodo, snapshot.TotalVentas, snapshot.TotalEfectivo, 
		snapshot.TotalTarjeta, snapshot.TotalOtros, snapshot.TotalDiferencias, 
		snapshot.TotalOperaciones, dataJSON, snapshot.IDUsuarioCierre,
	)
	return err
}

func (s *PeriodoStore) GetStatusIDByDesc(ctx context.Context, desc string) (uuid.UUID, error) {
	var id uuid.UUID
	query := `SELECT id_status FROM estatus WHERE std_descripcion ILIKE $1 LIMIT 1`
	err := s.db.QueryRowContext(ctx, query, desc).Scan(&id)
	return id, err
}
func (s *PeriodoStore) GetPeriodoByID(ctx context.Context, id uuid.UUID) (*models.Periodo, error) {
	return nil, nil
}
func (s *PeriodoStore) GetAllPeriodos(ctx context.Context) ([]*models.Periodo, error) {
	return nil, nil
}
func (s *PeriodoStore) UpdatePeriodo(ctx context.Context, id uuid.UUID, p *models.Periodo) (*models.Periodo, error) {
	return nil, nil
}
func (s *PeriodoStore) DeletePeriodo(ctx context.Context, id uuid.UUID) error { return nil }
