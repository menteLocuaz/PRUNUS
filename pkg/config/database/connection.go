package database

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/prunus/pkg/config"
	"github.com/redis/go-redis/v9"
)

// InitDB abre la conexión a la base de datos usando la configuración de Viper.
func InitDB() (*sql.DB, error) {
	host := config.Get("DB_HOST")
	user := config.Get("DB_USER")
	password := config.Get("DB_PASSWORD")
	dbname := config.Get("DB_NAME")
	port := config.GetDefault("DB_PORT", "5432")
	sslmode := config.GetDefault("DB_SSLMODE", "disable")

	connStr := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode,
	)

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("error al abrir la base de datos: %w", err)
	}

	maxOpenConns, _ := strconv.Atoi(config.GetDefault("DB_POOL_MAX_OPEN", "25"))
	maxIdleConns, _ := strconv.Atoi(config.GetDefault("DB_POOL_MAX_IDLE", "15"))
	connMaxLifetimeMin, _ := strconv.Atoi(config.GetDefault("DB_POOL_CONN_MAX_LIFETIME_MIN", "60"))

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(time.Duration(connMaxLifetimeMin) * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error al conectar (ping) a PostgreSQL: %w", err)
	}

	// Hotfix: Sincronizar esquema de estatus si el migrador falló
	if err := hotfixSyncSchema(db); err != nil {
		fmt.Printf("⚠️ Aviso: Fallo al sincronizar esquema (hotfix): %v\n", err)
	}

	return db, nil
}

// hotfixSyncSchema asegura que las columnas críticas existan y tengan los nombres correctos.
// Esto es necesario cuando el migrador reporta que está actualizado pero faltan columnas o renombrados.
func hotfixSyncSchema(db *sql.DB) error {
	queries := []string{
		// estatus: columnas nuevas
		`ALTER TABLE estatus
		 ADD COLUMN IF NOT EXISTS std_tipo_estado VARCHAR(100),
		 ADD COLUMN IF NOT EXISTS factor VARCHAR(100),
		 ADD COLUMN IF NOT EXISTS nivel INTEGER DEFAULT 0;`,

		// control_estacion: renombrar columnas si aún tienen los nombres viejos
		`DO $$ BEGIN
		   IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='control_estacion' AND column_name='id_control') THEN
		     ALTER TABLE control_estacion RENAME COLUMN id_control TO id_control_estacion;
		   END IF;
		 END $$;`,
		`DO $$ BEGIN
		   IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='control_estacion' AND column_name='id_usuario') THEN
		     ALTER TABLE control_estacion RENAME COLUMN id_usuario TO usuario_asignado;
		   END IF;
		 END $$;`,
		`DO $$ BEGIN
		   IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='control_estacion' AND column_name='fecha_apertura') THEN
		     ALTER TABLE control_estacion RENAME COLUMN fecha_apertura TO fecha_inicio;
		   END IF;
		 END $$;`,
		`DO $$ BEGIN
		   IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='control_estacion' AND column_name='fecha_cierre') THEN
		     ALTER TABLE control_estacion RENAME COLUMN fecha_cierre TO fecha_salida;
		   END IF;
		 END $$;`,
		`DO $$ BEGIN
		   IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='control_estacion' AND column_name='monto_apertura') THEN
		     ALTER TABLE control_estacion RENAME COLUMN monto_apertura TO fondo_base;
		   END IF;
		 END $$;`,

		// control_estacion: columnas nuevas
		`ALTER TABLE control_estacion
		 ADD COLUMN IF NOT EXISTS id_user_pos          UUID,
		 ADD COLUMN IF NOT EXISTS id_periodo           UUID,
		 ADD COLUMN IF NOT EXISTS fondo_retirado       NUMERIC(18,2),
		 ADD COLUMN IF NOT EXISTS usuario_retiro_fondo UUID,
		 ADD COLUMN IF NOT EXISTS ctrc_motivo_descuadre TEXT;`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return fmt.Errorf("hotfix schema: %w", err)
		}
	}
	return nil
}

// InitRedis abre la conexión a Redis usando la configuración de Viper.
func InitRedis() (*redis.Client, error) {
	host := config.GetDefault("REDIS_HOST", "localhost")
	port := config.GetDefault("REDIS_PORT", "6379")
	password := config.Get("REDIS_PASSWORD")
	db := config.GetInt("REDIS_DB")

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return rdb, rdb.Ping(ctx).Err()
}
