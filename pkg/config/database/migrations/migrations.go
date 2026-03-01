package migrations

import (
	"database/sql"
	"fmt"
)

func RunMigrations(db *sql.DB) error {
	steps := []struct {
		name string
		fn   func(*sql.DB) error
	}{
		{"001_empresa", migrateEmpresa},
		{"002_rol", migrateRol},
		{"003_sucursal", migrateSucursal},
		{"004_usuario", migrateUsuario},
		{"005_categoria", migrateCategoria},
		{"006_unidad", migrateUnidad},
		{"007_moneda", migrateMoneda},
		{"008_cliente", migrateCliente},
		{"009_proveedor", migrateProveedor},
		{"010_producto", migrateProducto},
	}

	for _, s := range steps {
		if err := s.fn(db); err != nil {
			return fmt.Errorf("migración %s: %w", s.name, err)
		}
	}

	return nil
}
