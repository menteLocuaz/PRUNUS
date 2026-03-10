package store

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/prunus/pkg/models"
)

type StorePOS interface {
	// Control Estación
	GetActiveControlByEstacion(idEstacion uuid.UUID) (*models.ControlEstacion, error)
	CreateControlEstacion(control *models.ControlEstacion) (*models.ControlEstacion, error)
	UpdateControlEstacion(control *models.ControlEstacion) error

	// Estaciones
	GetEstacionByID(id uuid.UUID) (*models.EstacionPos, error)
	UpdateEstacionStatus(id uuid.UUID, idStatus uuid.UUID) error

	// Periodos
	GetActivePeriodo() (*models.Periodo, error)
}

type storePOS struct {
	db *sql.DB
}

func NewPOSStore(db *sql.DB) StorePOS {
	return &storePOS{db: db}
}

func (s *storePOS) GetActiveControlByEstacion(idEstacion uuid.UUID) (*models.ControlEstacion, error) {
	query := `
		SELECT id_control_estacion, id_estacion, fecha_inicio, fecha_salida, fondo_base, 
		       usuario_asignado, id_status, id_user_pos, id_periodo, created_at, updated_at
		FROM control_estacion
		WHERE id_estacion = $1 AND fecha_salida IS NULL AND deleted_at IS NULL
		LIMIT 1
	`
	c := &models.ControlEstacion{}
	err := s.db.QueryRow(query, idEstacion).Scan(
		&c.IDControlEstacion, &c.IDEstacion, &c.FechaInicio, &c.FechaSalida, &c.FondoBase,
		&c.UsuarioAsignado, &c.IDStatus, &c.IDUserPos, &c.IDPeriodo, &c.CreatedAt, &c.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error al buscar control activo: %w", err)
	}
	return c, nil
}

func (s *storePOS) CreateControlEstacion(control *models.ControlEstacion) (*models.ControlEstacion, error) {
	query := `
		INSERT INTO control_estacion (id_estacion, fondo_base, usuario_asignado, id_status, id_user_pos, id_periodo)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id_control_estacion, fecha_inicio, created_at, updated_at
	`
	err := s.db.QueryRow(query,
		control.IDEstacion, control.FondoBase, control.UsuarioAsignado,
		control.IDStatus, control.IDUserPos, control.IDPeriodo,
	).Scan(&control.IDControlEstacion, &control.FechaInicio, &control.CreatedAt, &control.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("error al crear control estacion: %w", err)
	}
	return control, nil
}

func (s *storePOS) UpdateControlEstacion(control *models.ControlEstacion) error {
	query := `
		UPDATE control_estacion
		SET fecha_salida = $1, fondo_retirado = $2, usuario_retiro_fondo = $3, 
		    id_status = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id_control_estacion = $5 AND deleted_at IS NULL
	`
	_, err := s.db.Exec(query,
		control.FechaSalida, control.FondoRetirado, control.UsuarioRetiroFondo,
		control.IDStatus, control.IDControlEstacion,
	)
	return err
}

func (s *storePOS) GetEstacionByID(id uuid.UUID) (*models.EstacionPos, error) {
	query := `SELECT id_estacion, codigo, nombre, ip, id_sucursal, id_status FROM estaciones_pos WHERE id_estacion = $1 AND deleted_at IS NULL`
	e := &models.EstacionPos{}
	err := s.db.QueryRow(query, id).Scan(&e.IDEstacion, &e.Codigo, &e.Nombre, &e.IP, &e.IDSucursal, &e.IDStatus)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (s *storePOS) UpdateEstacionStatus(id uuid.UUID, idStatus uuid.UUID) error {
	query := `UPDATE estaciones_pos SET id_status = $1, updated_at = CURRENT_TIMESTAMP WHERE id_estacion = $2`
	_, err := s.db.Exec(query, idStatus, id)
	return err
}

func (s *storePOS) GetActivePeriodo() (*models.Periodo, error) {
	query := `SELECT id_periodo, prd_fecha_apertura, id_status FROM periodo WHERE prd_fecha_cierre IS NULL AND deleted_at IS NULL LIMIT 1`
	p := &models.Periodo{}
	err := s.db.QueryRow(query).Scan(&p.IDPeriodo, &p.PrdFechaApertura, &p.IDStatus)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no hay un periodo activo abierto")
	}
	return p, err
}
