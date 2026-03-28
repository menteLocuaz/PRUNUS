package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/prunus/pkg/models"
)

// StoreConfiguracion define las operaciones de base de datos para la configuración de impresión.
type StoreConfiguracion interface {
	GetCanalesImpresionActivos(ctx context.Context, chainID int) ([]models.CanalImpresion, error)
	GetImpresorasActivas(ctx context.Context, restaurantID int) ([]models.Impresora, error)
	GetPuertosActivos(ctx context.Context) ([]models.Puerto, error)
}

type storeConfiguracion struct {
	db *sql.DB
}

// NewConfiguracion crea una nueva instancia de StoreConfiguracion.
func NewConfiguracion(db *sql.DB) StoreConfiguracion {
	return &storeConfiguracion{db: db}
}

// GetCanalesImpresionActivos obtiene los canales de impresión activos para una cadena.
func (s *storeConfiguracion) GetCanalesImpresionActivos(ctx context.Context, chainID int) ([]models.CanalImpresion, error) {
	query := `
            SELECT id_canal_impresion, descripcion
            FROM canal_impresion
            WHERE cdn_id = $1
              AND deleted_at IS NULL
              AND id_status = (SELECT config.fn_estado('Canales Impresion', 'Activo'))`

	rows, err := s.db.QueryContext(ctx, query, chainID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo canales: %w", err)
	}
	defer rows.Close()

	var canales []models.CanalImpresion
	for rows.Next() {
		var c models.CanalImpresion
		if err := rows.Scan(&c.IDCanalImpresion, &c.Descripcion); err != nil {
			return nil, err
		}
		canales = append(canales, c)
	}
	return canales, nil
}

// GetImpresorasActivas obtiene las impresoras activas para un restaurante.
func (s *storeConfiguracion) GetImpresorasActivas(ctx context.Context, restaurantID int) ([]models.Impresora, error) {
	query := `
            SELECT id_impresora, nombre
            FROM impresora
            WHERE rst_id = $1
              AND deleted_at IS NULL
              AND id_status = (SELECT config.fn_estado('Canales Impresion', 'Activo'))
           ORDER BY nombre`

	rows, err := s.db.QueryContext(ctx, query, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo impresoras: %w", err)
	}
	defer rows.Close()

	var impresoras []models.Impresora
	for rows.Next() {
		var i models.Impresora
		if err := rows.Scan(&i.IDImpresora, &i.Nombre); err != nil {
			return nil, err
		}
		impresoras = append(impresoras, i)
	}
	return impresoras, nil
}

// GetPuertosActivos obtiene todos los puertos activos disponibles.
func (s *storeConfiguracion) GetPuertosActivos(ctx context.Context) ([]models.Puerto, error) {
	query := `
            SELECT id_puertos, descripcion
            FROM puertos
            WHERE deleted_at IS NULL
              AND id_status = (SELECT config.fn_estado('Canales Impresion', 'Activo'))
            ORDER BY
                SUBSTRING(descripcion, 1, 3),
                CAST(SUBSTRING(descripcion FROM 4) AS INT)`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo puertos: %w", err)
	}
	defer rows.Close()

	var puertos []models.Puerto
	for rows.Next() {
		var p models.Puerto
		if err := rows.Scan(&p.IDPuertos, &p.Descripcion); err != nil {
			return nil, err
		}
		puertos = append(puertos, p)
	}
	return puertos, nil
}
