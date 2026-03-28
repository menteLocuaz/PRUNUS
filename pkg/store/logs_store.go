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

	CreateFacturaAudit(ctx context.Context, a *models.FacturaAudit) error
	GetFacturaAudit(ctx context.Context, facturaID uuid.UUID) ([]*models.FacturaAudit, error)

	GetHistorialPrecios(ctx context.Context, productoID uuid.UUID, sucursalID uuid.UUID) ([]*models.HistorialPrecios, error)
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

func (s *storeLogs) CreateFacturaAudit(ctx context.Context, a *models.FacturaAudit) error {
	query := `INSERT INTO factura_audit (id_factura, id_usuario, accion, estado_anterior, estado_nuevo, observaciones, ip_address) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := s.db.ExecContext(ctx, query, a.IDFactura, a.IDUsuario, a.Accion, a.EstadoAnterior, a.EstadoNuevo, a.Observaciones, a.IPAddress)
	return err
}

func (s *storeLogs) GetFacturaAudit(ctx context.Context, facturaID uuid.UUID) ([]*models.FacturaAudit, error) {
	query := `SELECT id_audit, id_factura, id_usuario, accion, estado_anterior, estado_nuevo, observaciones, fecha, ip_address 
	          FROM factura_audit WHERE id_factura = $1 ORDER BY fecha DESC`
	rows, err := s.db.QueryContext(ctx, query, facturaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var audits []*models.FacturaAudit
	for rows.Next() {
		a := &models.FacturaAudit{}
		if err := rows.Scan(&a.IDAudit, &a.IDFactura, &a.IDUsuario, &a.Accion, &a.EstadoAnterior, &a.EstadoNuevo, &a.Observaciones, &a.Fecha, &a.IPAddress); err != nil {
			return nil, err
		}
		audits = append(audits, a)
	}
	return audits, nil
}

func (s *storeLogs) GetHistorialPrecios(ctx context.Context, productoID uuid.UUID, sucursalID uuid.UUID) ([]*models.HistorialPrecios, error) {
	query := `SELECT id_historial, id_producto, id_sucursal, precio_anterior, precio_nuevo, tipo_precio, id_usuario, fecha 
	          FROM historial_precios WHERE id_producto = $1 AND id_sucursal = $2 ORDER BY fecha DESC`
	rows, err := s.db.QueryContext(ctx, query, productoID, sucursalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var historial []*models.HistorialPrecios
	for rows.Next() {
		h := &models.HistorialPrecios{}
		if err := rows.Scan(&h.IDHistorial, &h.IDProducto, &h.IDSucursal, &h.PrecioAnterior, &h.PrecioNuevo, &h.TipoPrecio, &h.IDUsuario, &h.Fecha); err != nil {
			return nil, err
		}
		historial = append(historial, h)
	}
	return historial, nil
}
