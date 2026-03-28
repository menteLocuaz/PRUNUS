package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/prunus/pkg/dto"
	"github.com/prunus/pkg/models"
)

type StoreLogs interface {
	CreateLog(ctx context.Context, l *models.LogSistema) error
	GetAllLogs(ctx context.Context, params dto.PaginationParams) ([]*models.LogSistema, error)

	CreateAuditoriaCaja(ctx context.Context, a *models.AuditoriaCaja) error
	GetAuditoriaCaja(ctx context.Context, controlID uuid.UUID) ([]*models.AuditoriaCaja, error)
}

type storeLogs struct {
	db *sql.DB
}

func NewLogs(db *sql.DB) StoreLogs {
	return &storeLogs{db: db}
}

func (s *storeLogs) CreateLog(ctx context.Context, l *models.LogSistema) error {
	query := `INSERT INTO log_sistema (id_usuario, accion, tabla, registro_id, ip, fecha) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := s.db.ExecContext(ctx, query, l.IDUsuario, l.Accion, l.Tabla, l.RegistroID, l.IP, time.Now())
	return err
}

func (s *storeLogs) GetAllLogs(ctx context.Context, params dto.PaginationParams) ([]*models.LogSistema, error) {
	if params.Limit <= 0 {
		params.Limit = 100 // Límite por defecto para logs es mayor
	}

	query := `SELECT id_log, id_usuario, accion, tabla, registro_id, ip, fecha FROM log_sistema`
	
	var args []interface{}
	
	if params.LastDate != nil {
		query += " WHERE fecha < $1"
		args = append(args, params.LastDate)
	}

	query += " ORDER BY fecha DESC LIMIT $" + fmt.Sprint(len(args)+1)
	args = append(args, params.Limit)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*models.LogSistema
	for rows.Next() {
		l := &models.LogSistema{}
		if err := rows.Scan(&l.IDLog, &l.IDUsuario, &l.Accion, &l.Tabla, &l.RegistroID, &l.IP, &l.Fecha); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}

func (s *storeLogs) CreateAuditoriaCaja(ctx context.Context, a *models.AuditoriaCaja) error {
	query := `INSERT INTO auditoria_caja (id_control_estacion, tipo_movimiento, valor, fecha, id_usuario, descripcion) 
	          VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := s.db.ExecContext(ctx, query, a.IDControlEstacion, a.TipoMovimiento, a.Valor, time.Now(), a.IDUsuario, a.Descripcion)
	return err
}

func (s *storeLogs) GetAuditoriaCaja(ctx context.Context, controlID uuid.UUID) ([]*models.AuditoriaCaja, error) {
	query := `SELECT id_auditoria, id_control_estacion, tipo_movimiento, valor, fecha, id_usuario, descripcion FROM auditoria_caja WHERE id_control_estacion = $1`
	rows, err := s.db.QueryContext(ctx, query, controlID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var auditorias []*models.AuditoriaCaja
	for rows.Next() {
		a := &models.AuditoriaCaja{}
		if err := rows.Scan(&a.IDAuditoria, &a.IDControlEstacion, &a.TipoMovimiento, &a.Valor, &a.Fecha, &a.IDUsuario, &a.Descripcion); err != nil {
			return nil, err
		}
		auditorias = append(auditorias, a)
	}
	return auditorias, nil
}
